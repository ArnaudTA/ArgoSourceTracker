package application

import (
	"argocd-watcher/pkg/config"
	"fmt"
	"sync"
	"time"

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
			app := obj.(*v1alpha1.Application)
			key := fmt.Sprintf("%s/%s", app.Namespace, app.Name)
			AppCache.Store(key, app)
			go ParseApplication(app)
			logrus.Debugf("[ADD] Application %s\n", key)
		},
		UpdateFunc: func(_, newObj interface{}) {
			app := newObj.(*v1alpha1.Application)
			key := fmt.Sprintf("%s/%s", app.Namespace, app.Name)
			AppCache.Store(key, app)
			go ParseApplication(app)
			logrus.Debugf("[UPDATE] Application %s\n", key)
		},
		DeleteFunc: func(obj interface{}) {
			app := obj.(*v1alpha1.Application)
			key := fmt.Sprintf("%s/%s", app.Namespace, app.Name)
			AppCache.Delete(key)
			logrus.Debugf("[DELETE] Application %s\n", key)
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
