package argocd

import (
	"argocd-watcher/pkg/config"
	"context"
	"log"
	"os"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	argoclientset "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var client *argoclientset.Clientset
var argocdNs string

func InitClient(cfg config.Config) {
	// Charger le namespace depuis le FS
	data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err == nil {
		argocdNs = string(data)
	}

	if cfg.ArgocdNs != "" {
		argocdNs = cfg.ArgocdNs
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

func getArgoCDClient() *argoclientset.Clientset {
	if client == nil {
		panic("No Argocd client Found")
	}
	return client
}

func GetApplications() (*v1alpha1.ApplicationList, error) {
	argoClient := getArgoCDClient()
	return argoClient.ArgoprojV1alpha1().Applications(argocdNs).List(context.TODO(), metav1.ListOptions{})
}
