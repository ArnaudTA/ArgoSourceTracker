package types

import "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"

type ApplicationStatus string

const (
	UpToDate ApplicationStatus = "Up-to-date"
	Ignored  ApplicationStatus = "Ignored"
	Outdated ApplicationStatus = "Outdated"
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
	Error    string            `json:"error"`
	RepoURL  string            `json:"repoURL"`
	Status   ApplicationStatus `json:"status"`
	Revision string            `json:"revision"`
	NewTags  []string          `json:"newTags,omitempty"`
	Protocol string            `json:"type"`
	Chart    string            `json:"chart" binding:"required"`
}

type ApplicationSourceWithRevision struct {
	Source   v1alpha1.ApplicationSource
	Revision string
}

type Parent struct {
	Kind           string `json:"kind" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Namespace      string `json:"namespace" binding:"required"`
	ApplicationUrl string `json:"applicationUrl,omitempty"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
}

type AppStats map[ApplicationStatus]int
type ListApplicationRes struct {
	Items []Summary `json:"items" binding:"required"`
	Stats AppStats `json:"stats" binding:"required"`
}