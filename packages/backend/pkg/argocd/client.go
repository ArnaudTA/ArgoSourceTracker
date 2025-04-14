package argocd

import (
	"argocd-watcher/pkg/config"
	"context"
	"fmt"
	"log"
	"os"

	argoclientset "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	core "k8s.io/api/core"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/apis/core"
)

var argoClient *argoclientset.Clientset
var ArgocdNs string

func InitClient(cfg config.Config) {
	// Charger le namespace depuis le FS
	data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err == nil {
		ArgocdNs = string(data)
	}

	if cfg.Argocd.Namespace != "" {
		ArgocdNs = cfg.Argocd.Namespace
	}

	// Charger la config Kube
	config, err := clientcmd.BuildConfigFromFlags("", cfg.Kubeconfig)
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la config K8s: %v", err)
	}

	// Créer le client ArgoCD
	client, err := argoclientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("Erreur lors de la création du client ArgoCD: %v", err)
	}
	argoClient = client
}

func LoadArgoConf(cfg *config.Config) {
	// Charger le namespace depuis le FS
	data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err == nil {
		ArgocdNs = string(data)
	}

	if cfg.Argocd.Namespace != "" {
		ArgocdNs = cfg.Argocd.Namespace
	}

	// Charger la config Kube
	config, err := clientcmd.BuildConfigFromFlags("", cfg.Kubeconfig)
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la config K8s: %v", err)
	}

	// Créer le client ArgoCD
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Erreur lors de la création du client ArgoCD: %v", err)
	}
	argoCDCm := core.ConfigMap
	if cfg.Argocd.Url == "" {
		fmt.Println("Discovering ArgoCD url...")
		clientset.CoreV1().ConfigMaps(cfg.Argocd.Namespace).List(context.Background(), v1.ListOptions{LabelSelector: ""})
	}
	fmt.Printf("ArgoCD domain found in : %s\n", cfg.Argocd.Url)
}

func GetArgoCDClient() *argoclientset.Clientset {
	if argoClient == nil {
		panic("No Argocd client Found")
	}
	return argoClient
}
