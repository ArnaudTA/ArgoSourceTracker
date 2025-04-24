package application

import (
	"fmt"
	"sync"
	"time"

	"github.com/cableship/chart-sentinel/pkg/config"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	appinformers "github.com/argoproj/argo-cd/v2/pkg/client/informers/externalversions"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/cache"
)

var AppCache sync.Map // key = "namespace/name", value = *v1alpha1.Application

func Watch(appClient versioned.Interface) {
	// Création de l'informer factory pour ArgoCD
	informerFactory := appinformers.NewFilteredSharedInformerFactory(appClient, time.Minute*10, config.Global.Argocd.Namespace, nil)

	// Informer typé pour les Applications
	appInformer := informerFactory.Argoproj().V1alpha1().Applications()

	// Handlers avec cache
	appInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			resource := obj.(*v1alpha1.Application)
			key := fmt.Sprintf("%s/%s", resource.Namespace, resource.Name)
			logrus.Debugf("[ADD] Application %s\n", key)

			app := &Application{Resource: resource}
			AppCache.Store(key, app)
			go app.Parse()
		},
		UpdateFunc: func(_, newObj interface{}) {
			resource := newObj.(*v1alpha1.Application)
			key := fmt.Sprintf("%s/%s", resource.Namespace, resource.Name)
			logrus.Debugf("[UPDATE] Application %s\n", key)

			app := &Application{Resource: resource}
			AppCache.Store(key, app)
			go app.Parse()
		},
		DeleteFunc: func(obj interface{}) {
			app := obj.(*v1alpha1.Application)
			key := fmt.Sprintf("%s/%s", app.Namespace, app.Name)
			logrus.Debugf("[DELETE] Application %s\n", key)

			AppCache.Delete(key)
		},
	})

	// Démarrer l'informer
	stopCh := make(chan struct{})
	defer close(stopCh)

	logrus.Println("🚀 Démarrage de l'informer ArgoCD (Application)...")
	informerFactory.Start(stopCh)

	// Attendre que le cache soit prêt
	informerFactory.WaitForCacheSync(stopCh)
	<-stopCh
}
