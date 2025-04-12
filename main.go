package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Charger la config Kube
	kubeConfig := os.Getenv("HOME")+"/.kube/config" // Remplace par ton chemin de kubeconfig si nécessaire
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la config K8s: %v", err)
	}

	// Créer le client Kubernetes
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Erreur lors de la création du client Kubernetes: %v", err)
	}

	// Serveur web avec Gin
	r := gin.Default()

	// Endpoint pour récupérer les nodes dans un namespace donné
	r.GET("/nodes/", func(c *gin.Context) {
		// Obtenir la liste des Pods dans le namespace
		nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Créer une liste simplifiée des applications
		var result []gin.H
		for _, pod := range nodes.Items {
			result = append(result, gin.H{
				"name": pod.Name,
				"status": pod.Status.Addresses,
			})
		}

		c.JSON(http.StatusOK, result)
	})

	// Lancer le
	r.Run(":8080")
}