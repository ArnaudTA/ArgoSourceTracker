package applicationset

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

var AppSetCache sync.Map // key = "namespace/name", value = *v1alpha1.Application

func Watch(appClient versioned.Interface) {
	// Création de l'informer factory pour ArgoCD
	informerFactory := appinformers.NewFilteredSharedInformerFactory(appClient, time.Minute*10, config.Global.Argocd.Namespace, nil)

	// Informer typé pour les Applications
	appInformer := informerFactory.Argoproj().V1alpha1().ApplicationSets()

	// Handlers avec cache
	appInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			app := obj.(*v1alpha1.ApplicationSet)
			key := fmt.Sprintf("%s/%s", app.Namespace, app.Name)
			AppSetCache.Store(key, app)
			logrus.Printf("[ADD] ApplicationSet %s\n", key)
		},
		UpdateFunc: func(_, newObj interface{}) {
			app := newObj.(*v1alpha1.ApplicationSet)
			key := fmt.Sprintf("%s/%s", app.Namespace, app.Name)
			AppSetCache.Store(key, app)
			logrus.Printf("[UPDATE] ApplicationSet %s\n", key)
		},
		DeleteFunc: func(obj interface{}) {
			app := obj.(*v1alpha1.ApplicationSet)
			key := fmt.Sprintf("%s/%s", app.Namespace, app.Name)
			AppSetCache.Delete(key)
			logrus.Printf("[DELETE] ApplicationSet %s\n", key)
		},
	})

	// Démarrer l'informer
	stopCh := make(chan struct{})
	defer close(stopCh)

	logrus.Println("🚀 Démarrage de l'informer ArgoCD (ApplicationSet)...")
	informerFactory.Start(stopCh)

	// Attendre que le cache soit prêt
	informerFactory.WaitForCacheSync(stopCh)

	<-stopCh
}
