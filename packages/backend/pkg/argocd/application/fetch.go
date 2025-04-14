package application

import (
	"argocd-watcher/pkg/argocd"
	"argocd-watcher/pkg/argocd/applicationset"
	"context"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var InstanceLabel string = "argocd.argoproj.io/instance"

func getApplications() ([]v1alpha1.Application, error) {
	argoClient := argocd.GetArgoCDClient()
	apps, err := argoClient.ArgoprojV1alpha1().Applications(argocd.ArgocdNs).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return apps.Items, nil
}

func getApplication(name string) (*v1alpha1.Application, error) {
	argoClient := argocd.GetArgoCDClient()
	return argoClient.ArgoprojV1alpha1().Applications(argocd.ArgocdNs).Get(context.TODO(), name, metav1.GetOptions{})
}

type TrackRecord struct {
	Kind string `json:"kind" binding:"required"`
	Name string `json:"name" binding:"required"`
	Sources []v1alpha1.ApplicationSource `json:"sources" binding:"required"`
}

func GetApplicationTrack(name string) []TrackRecord {
	track := []TrackRecord{}
	application, err := StoreGet(name)
	if err != nil {
		return track
	}
	metadata := application.ObjectMeta
	for i := 0; i < 10; i++ {
		previousResource := getPreviousResource(metadata)
		switch previousResource.Kind {
		case "Application":
			app, err := StoreGet(previousResource.Name)
			if err != nil {
				break
			}
			track = append(track, TrackRecord{
				Kind: previousResource.Kind,
				Name: app.Name,
				Sources: append(app.Status.Sync.ComparedTo.Sources, app.Status.Sync.ComparedTo.Source),
			})
			metadata = app.ObjectMeta
		case "ApplicationSet":
			app, err := applicationset.StoreGet(previousResource.Name)
			if err != nil {
				break
			}
			track = append(track, TrackRecord{
				Kind: previousResource.Kind,
				Name: app.Name,
			})
			metadata = app.ObjectMeta
		}
	}
	return track
}

type PreviousResource struct {
	Kind string
	Name string
}

func getPreviousResource(metadata metav1.ObjectMeta) PreviousResource {
	if instance, ok := metadata.Labels[InstanceLabel]; ok {
		return PreviousResource{
			Kind: "Application",
			Name: instance,
		}
	}
	for _, owner := range metadata.OwnerReferences {
		if owner.Kind == "ApplicationSet" {
			return PreviousResource{
				Kind: "ApplicationSet",
				Name: owner.Name,
			}
		}
	}
	return PreviousResource{}
}
