/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package zips_test

import (
	"errors"
	"testing"

	"github.com/onsi/gomega"

	"github.com/stretchr/testify/mock"

	"github.com/Azure/k8s-infra/pkg/zips"
)

type (
	EnvMock struct {
		mock.Mock
	}
)

func (e *EnvMock) Getenv(key string) string {
	args := e.Called(key)
	return args.String(0)
}

func TestNewAzureTemplateClient(t *testing.T) {
	cases := []struct {
		Name     string
		EnvSetup func(*EnvMock) *EnvMock
		Expect   func(*gomega.GomegaWithT, *zips.AzureTemplateClient, error)
	}{
		{
			Name: "WithServicePrincipalEnv",
			EnvSetup: func(env *EnvMock) *EnvMock {
				env.On("Getenv", "AZURE_SUBSCRIPTION_ID").Return("foo")
				env.On("Getenv", "AZURE_TENANT_ID").Return("bar")
				env.On("Getenv", "AZURE_CLIENT_ID").Return("buzz")
				env.On("Getenv", "AZURE_CLIENT_SECRET").Return("bazz")
				env.On("Getenv", mock.Anything).Return("")
				return env
			},
			Expect: func(g *gomega.GomegaWithT, client *zips.AzureTemplateClient, err error) {
				g.Expect(err).To(gomega.BeNil())
				g.Expect(client.RawClient).ToNot(gomega.BeNil())
			},
		},
		{
			Name: "WithNoEnv",
			EnvSetup: func(env *EnvMock) *EnvMock {
				env.On("Getenv", mock.Anything).Return("")
				return env
			},
			Expect: func(g *gomega.GomegaWithT, client *zips.AzureTemplateClient, err error) {
				g.Expect(err).To(gomega.MatchError(errors.New("env var \"AZURE_SUBSCRIPTION_ID\" was not set")))
			},
		},
		{
			Name: "WillDefaultToMSIWithoutServicePrincipal",
			EnvSetup: func(env *EnvMock) *EnvMock {
				env.On("Getenv", "AZURE_SUBSCRIPTION_ID").Return("foo")
				env.On("Getenv", mock.Anything).Return("")
				return env
			},
			Expect: func(g *gomega.GomegaWithT, client *zips.AzureTemplateClient, err error) {
				g.Expect(err).To(gomega.BeNil())
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			env := new(EnvMock)
			env = c.EnvSetup(env)
			client, err := zips.NewAzureTemplateClient(zips.WithEnv(env))
			g := gomega.NewGomegaWithT(t)
			c.Expect(g, client, err)
		})
	}
}
