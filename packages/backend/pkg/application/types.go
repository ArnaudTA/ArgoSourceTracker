package application

import "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"

type ChartSummary struct {
	RepoURL  string   `json:"repoURL"`
	Status   string   `json:"status,omitempty"`
	Revision string   `json:"revision"`
	NewTags  []string `json:"newTags,omitempty"`
	Protocol string   `json:"type"`
	Chart    string   `json:"chart" binding:"required"`
}

type ApplicationSummary struct {
	Name      string         `json:"name" binding:"required"`
	Namespace string         `json:"namespace" binding:"required"`
	Charts    []ChartSummary `json:"charts" binding:"required"`
	Status    string         `json:"status" binding:"required"`
	Instance  string         `json:"instance,omitempty" binding:"required"`
	ApplicationUrl string `json:"applicationUrl" binding:"required"`
}

type ApplicationSourceWithRevision struct {
	v1alpha1.ApplicationSource
	revision string
}

type GenealogyItem struct {
	Kind           string `json:"kind" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Namespace           string `json:"namespace" binding:"required"`
	ApplicationUrl string `json:"applicationUrl,omitempty"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
}
