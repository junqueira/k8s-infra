/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

// zips package is where all the fun Azure client, cache, throttling, CRUD will go. Right now, it just provides an
// Apply and Delete interface
package zips

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
)

type (
	Resourcer interface {
		ToResource() (Resource, error)
		FromResource(Resource) (runtime.Object, error)
	}

	Applier interface {
		Apply(ctx context.Context, res Resource) (Resource, error)
		Delete(ctx context.Context, resourceID string) error
	}

	Resource struct {
		ResourceGroup     string      `json:"-"` // resource group should not be serialized as part of the resource. This indicates that this should be within a resource group or at a subscription level deployment.
		SubscriptionID    string      `json:"-"`
		ProvisioningState string      `json:"-"`
		DeploymentID      string      `json:"-"`
		ID                string      `json:"id,omitempty"`
		Name              string      `json:"name,omitempty"`
		Location          string      `json:"location,omitempty"`
		Type              string      `json:"type,omitempty"`
		APIVersion        string      `json:"apiVersion,omitempty"`
		Properties        interface{} `json:"properties,omitempty"`
	}
)