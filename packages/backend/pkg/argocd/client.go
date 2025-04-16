package argocd

import (
	cfg "argocd-watcher/pkg/config"
	"context"
	"os"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	appclientset "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"

)

func LoadArgoConf() {
	// Charger le namespace depuis le FS
	data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err == nil {
		cfg.Global.Argocd.Namespace = string(data)
	}

	// Charger la config Kube
	config, err := clientcmd.BuildConfigFromFlags("", cfg.Global.Kubeconfig)
	if err != nil {
		logrus.Fatalf("Erreur lors du chargement de la config K8s: %v", err)
	}

	// Créer le client ArgoCD
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("Erreur lors de la création du client ArgoCD: %v", err)
	}
	configmaps, err := clientset.CoreV1().ConfigMaps(cfg.Global.Argocd.Namespace).List(context.Background(), v1.ListOptions{LabelSelector: "app.kubernetes.io/component=server,app.kubernetes.io/instance=" + cfg.Global.Argocd.Instance})

	if err != nil {
		logrus.Fatalf("Can't List configmaps in namespace: %s", cfg.Global.Argocd.Namespace)
		logrus.Fatalln("ArgoCD url not found and not set !")
		return
	}
	if cfg.Global.Argocd.Url != "" {
		logrus.Infof("ArgoCD Url set by env or flag: %s\n", cfg.Global.Argocd.Url)
	} else {
		logrus.Infof("Discovering ArgoCD url...")
	}
	for _, cm := range configmaps.Items {
		if instanceLabelKey, ok := cm.Data["application.instanceLabelKey"]; ok {
			cfg.Global.Argocd.InstanceLabelKey = instanceLabelKey
			logrus.Infof("ArgoCD instanceLabelKey found in configmap %s: %s\n", cm.Name, instanceLabelKey)
		}

		if url, ok := cm.Data["url"]; ok {
			cfg.Global.Argocd.Url = url
			logrus.Infof("ArgoCD url found in configmap %s: %s\n", cm.Name, cfg.Global.Argocd.Url)
		}
	}
	if cfg.Global.Argocd.Url != "" {
		logrus.Warnln("Can't find ArgoCD url")
	}
}

func GetClient() *appclientset.Clientset {
		// Charger la config Kube
		config, err := clientcmd.BuildConfigFromFlags("", cfg.Global.Kubeconfig)
		if err != nil {
			logrus.Fatalf("Erreur lors du chargement de la config K8s: %v", err)
		}
	
		// Création du client ArgoCD typé
		appClient, err := appclientset.NewForConfig(config)
		if err != nil {
			panic(err)
		}
		return appClient
}