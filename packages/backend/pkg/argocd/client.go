package argocd

import (
	cfg "argocd-watcher/pkg/config"
	"context"
	"fmt"
	"log"
	"os"

	argoclientset "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var argoClient *argoclientset.Clientset
var ArgocdNs string

func InitClient() {
	// Charger le namespace depuis le FS
	data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err == nil {
		ArgocdNs = string(data)
	}

	if cfg.Global.Argocd.Namespace != "" {
		ArgocdNs = cfg.Global.Argocd.Namespace
	}

	// Charger la config Kube
	config, err := clientcmd.BuildConfigFromFlags("", cfg.Global.Kubeconfig)
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

func LoadArgoConf() {
	// Charger le namespace depuis le FS
	data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err == nil {
		ArgocdNs = string(data)
	}

	if cfg.Global.Argocd.Namespace != "" {
		ArgocdNs = cfg.Global.Argocd.Namespace
	}

	// Charger la config Kube
	config, err := clientcmd.BuildConfigFromFlags("", cfg.Global.Kubeconfig)
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la config K8s: %v", err)
	}

	// Créer le client ArgoCD
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Erreur lors de la création du client ArgoCD: %v", err)
	}
	if cfg.Global.Argocd.Url != "" {
		fmt.Printf("ArgoCD set by env or flag: %s\n", cfg.Global.Argocd.Url)
		return
	}
	fmt.Println("Discovering ArgoCD url...")
	configmaps, err := clientset.CoreV1().ConfigMaps(cfg.Global.Argocd.Namespace).List(context.Background(), v1.ListOptions{LabelSelector: "app.kubernetes.io/component=server,app.kubernetes.io/name=argocd-cm"})
	if err != nil {
		fmt.Printf("Can't List configmaps in namespace: %s", cfg.Global.Argocd.Namespace)
		fmt.Println("ArgoCD url not found and not set !")
		return
	}
	for _, cm := range configmaps.Items {
		if url, ok := cm.Data["url"]; ok {
			cfg.Global.Argocd.Url = url
			fmt.Printf("ArgoCD domain found in configmap %s: %s\n", cm.Name, cfg.Global.Argocd.Url)
			return
		}
	}
}

func GetArgoCDClient() *argoclientset.Clientset {
	if argoClient == nil {
		panic("No Argocd client Found")
	}
	return argoClient
}
