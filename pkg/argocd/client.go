package argocd

import (
	"argocd-watcher/pkg/config"
	"log"
	"os"

	argoclientset "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	"k8s.io/client-go/tools/clientcmd"
)

var client *argoclientset.Clientset
var ArgocdNs string

func InitClient(cfg config.Config) {
	// Charger le namespace depuis le FS
	data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err == nil {
		ArgocdNs = string(data)
	}

	if cfg.ArgocdNs != "" {
		ArgocdNs = cfg.ArgocdNs
	}

	// Charger la config Kube
	config, err := clientcmd.BuildConfigFromFlags("", cfg.Kubeconfig)
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la config K8s: %v", err)
	}

	// Créer le client ArgoCD
	argoClient, err := argoclientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("Erreur lors de la création du client ArgoCD: %v", err)
	}
	client = argoClient

}

func GetArgoCDClient() *argoclientset.Clientset {
	if client == nil {
		panic("No Argocd client Found")
	}
	return client
}
