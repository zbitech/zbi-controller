package zbi

import (
	"github.com/stretchr/testify/assert"
	fake_zbi "github.com/zbitech/controller/fake-zbi"
	"github.com/zbitech/controller/internal/klient/client"
	"github.com/zbitech/controller/internal/manager"
	"github.com/zbitech/controller/internal/vars"
	"testing"
)

func TestZBIClient_CreateProject(t *testing.T) {

	ctx := fake_zbi.InitContext()
	vars.ManagerFactory = manager.NewResourceManagerFactory()
	vars.ManagerFactory.Init(ctx)

	k, err := client.NewKlient(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, k)

	z := NewZBIClient(ctx, k)
	assert.NotNil(t, z)

	//resources, err := z.CreateProject(ctx, &data.Project1)
	//assert.NoError(t, err)
	//assert.NotNil(t, resources)
	//
	//k.DeleteResources(ctx, resources)
}
