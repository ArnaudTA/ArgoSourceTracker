package application

import (
	"fmt"
	"sync"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

var Cache = sync.Map{}
// var Store = map[string]v1alpha1.Application{}

func StoreGet(applicationName string) (v1alpha1.Application, error) {
	cachedApplication, ok := Cache.Load(applicationName)
	if ok {
		fmt.Printf("Use Application cache for: %s\n", applicationName)
		return cachedApplication.(v1alpha1.Application), nil
	}

	fmt.Printf("\nFetching application: %s\n", applicationName)

	application, err := getApplication(applicationName)
	if err != nil {
		panic(err)
	}

	Cache.Store(applicationName, *application)

	time.AfterFunc(300*time.Second, func() { StoreDeleteApplication(applicationName) })

	return *application, nil
}

func StoreList() ([]v1alpha1.Application, error) {
	fmt.Printf("Fetching all applications\n")

	applications, err := getApplications()
	if err != nil {
		panic(err)
	}

	for _, app := range applications {
		Cache.Store(app.Name, app)
		time.AfterFunc(300*time.Second, func() { StoreDeleteApplication(app.Name) })
	}

	return applications, nil
}

func StoreDeleteApplication(applicationName string) {
	Cache.Delete(applicationName)
}
