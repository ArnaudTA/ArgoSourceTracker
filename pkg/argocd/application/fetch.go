package application

import (
	"argocd-watcher/pkg/argocd"
	"context"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var InstanceLabel string = "argocd.argoproj.io/instance"

func GetApplications() (*v1alpha1.ApplicationList, error) {
	argoClient := argocd.GetArgoCDClient()
	return argoClient.ArgoprojV1alpha1().Applications(argocd.ArgocdNs).List(context.TODO(), metav1.ListOptions{})
}

func GetApplication(name string) (*v1alpha1.Application, error) {
	argoClient := argocd.GetArgoCDClient()
	return argoClient.ArgoprojV1alpha1().Applications(argocd.ArgocdNs).Get(context.TODO(), name, metav1.GetOptions{})
}

func GetApplicationTrack(name string) []*v1alpha1.Application {
	track := []*v1alpha1.Application{}
	for i := 0; i < 10; i++ {
		app, err := GetApplication(name)
		if err != nil {
			break
		}
		track = append(track, app)
		if instance, ok := app.Labels[InstanceLabel]; ok {
			name = instance
		} else {
			break
		}
	}
	return track
}
