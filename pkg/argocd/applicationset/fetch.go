package applicationset

import (
	"argocd-watcher/pkg/argocd"
	"context"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var InstanceLabel string = "argocd.argoproj.io/instance"

func getApplicationsets() ([]v1alpha1.ApplicationSet, error) {
	argoClient := argocd.GetArgoCDClient()
	apps, err := argoClient.ArgoprojV1alpha1().ApplicationSets(argocd.ArgocdNs).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return apps.Items, nil
}

func getApplicationset(name string) (*v1alpha1.ApplicationSet, error) {
	argoClient := argocd.GetArgoCDClient()
	return argoClient.ArgoprojV1alpha1().ApplicationSets(argocd.ArgocdNs).Get(context.TODO(), name, metav1.GetOptions{})
}

