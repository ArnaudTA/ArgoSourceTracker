package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	argoclientset "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Charger la config Kube
	kubeConfig := os.Getenv("HOME") + "/.kube/config" // Remplace par ton chemin de kubeconfig si nécessaire
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la config K8s: %v", err)
	}

	// Créer le client ArgoCD
	argoClient, err := argoclientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("Erreur lors de la création du client ArgoCD: %v", err)
	}

	// Serveur web avec Gin
	r := gin.Default()

	// Endpoint pour récupérer les applications ArgoCD dans un namespace donné
	r.GET("/argocd/apps", func(c *gin.Context) {
		// Liste des applications ArgoCD
		applications, err := argoClient.ArgoprojV1alpha1().Applications("argocd").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Créer une liste simplifiée des applications
		var result []gin.H
		appList := applications.Items
		for _, app := range appList {
			if app.Spec.Source == nil {
				continue
			}
			repoURL := app.Spec.Source.RepoURL
			if strings.HasSuffix(repoURL, ".git") {
				continue
			}
			applicationSummary := gin.H{
				"name": app.Name,
				"repo": repoURL,
			}
			if strings.HasPrefix(repoURL, "https://") {
				applicationSummary["protocol"] = "https"
			} else if strings.HasPrefix(repoURL, "8gears") {
				applicationSummary["protocol"] = "oci"
				// tags, err := getTagByIndex(repoURL, app.Spec.Source.Chart)
				// if err == nil {
				// 	applicationSummary["tags"] = tags
				// }
			}
			result = append(result, applicationSummary)
		}

		c.JSON(http.StatusOK, result)
	})

	// Lancer le serveur
	r.Run(":8080")
}

// func getTagByIndex(repository, chart string) (string, error) {
// 	tags := []string{}
// 	// Appel HTTP pour récupérer le fichier index.yaml
// 	resp, err := http.Get(repository + "/index.yaml")
// 	if err != nil {
// 		log.Fatalf("Erreur HTTP : %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		log.Fatalf("Code HTTP inattendu : %d", resp.StatusCode)
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalf("Erreur de lecture du body : %v", err)
// 	}

// 	// Parsing YAML au format Helm
// 	index := &repo.IndexFile{}
// 	if err := yaml.Unmarshal(body, index); err != nil {
// 		log.Fatalf("Erreur de parsing YAML : %v", err)
// 	}

// 	// Optionnel : vérifier et trier l'index
// 	index.SortEntries()

// 	// Afficher les charts trouvés
// 	for name, versions := range index.Entries {
// 		fmt.Printf("Chart: %s (%d versions)\n", name, len(versions))
// 	}
// 	for entry, versions := range index.Entries {
// 		if entry != chart {
// 			continue
// 		}
// 		for _, version := range versions {
// 			tags = append(tags, version.Version)
// 		}
// 	}
// 	return strings.Join(tags, " "), nil
// }
