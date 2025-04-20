package types

import "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"

type ApplicationStatus string

const (
	UpToDate ApplicationStatus = "Up-to-date"
	Unknown  ApplicationStatus = "Unknown"
	Outdated ApplicationStatus = "Outdated"
	Ignored  ApplicationStatus = "Ignored"
	Checking ApplicationStatus = "Checking"
	Error    ApplicationStatus = "Error"
)

type Summary struct {
	Name           string            `json:"name" binding:"required"`
	Namespace      string            `json:"namespace" binding:"required"`
	Charts         []ChartSummary    `json:"charts" binding:"required"`
	Status         ApplicationStatus `json:"status" binding:"required"`
	Instance       string            `json:"instance,omitempty" binding:"required"`
	ApplicationUrl string            `json:"applicationUrl" binding:"required"`
}

type ChartSummary struct {
	RepoURL  string   `json:"repoURL"`
	Status   string   `json:"status,omitempty"`
	Revision string   `json:"revision"`
	NewTags  []string `json:"newTags,omitempty"`
	Protocol string   `json:"type"`
	Chart    string   `json:"chart" binding:"required"`
}

type ApplicationSourceWithRevision struct {
	v1alpha1.ApplicationSource
	Revision string
}

type Parent struct {
	Kind           string `json:"kind" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Namespace      string `json:"namespace" binding:"required"`
	ApplicationUrl string `json:"applicationUrl,omitempty"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
}
