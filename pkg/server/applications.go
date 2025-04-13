package server

import (
	"argocd-watcher/pkg/parser"
	"context"
	"log"
	"net/http"
	"os"

	argoclientset "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func getArgoCDClient() *argoclientset.Clientset {
	// Charger la config Kube
	kubeConfig := os.Getenv("KUBECONFIG")
	if kubeConfig == "" {
		kubeConfig = os.Getenv("HOME") + "/.kube/config"
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la config K8s: %v", err)
	}

	// Créer le client ArgoCD
	argoClient, err := argoclientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("Erreur lors de la création du client ArgoCD: %v", err)
	}
	return argoClient
}

func fetchApplications(c *gin.Context) {
	// Liste des applications ArgoCD
	argoClient := getArgoCDClient()
	applications, err := argoClient.ArgoprojV1alpha1().Applications("argocd").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Créer une liste simplifiée des applications
	var result []parser.ChartVersion
	appList := applications.Items
	for _, app := range appList {
		result = append(result, parser.ParseApplication(app)...)
	}

	c.JSON(http.StatusOK, result)
}
