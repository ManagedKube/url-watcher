package controller

import (
	"managedkube.com/url-watcher/pkg/controller/watcher"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, watcher.Add)
}
