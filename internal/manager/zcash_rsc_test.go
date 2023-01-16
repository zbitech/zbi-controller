package manager

import (
	"github.com/stretchr/testify/assert"
	fake_zbi "github.com/zbitech/controller/fake-zbi"
	"github.com/zbitech/controller/fake-zbi/data"
	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/internal/vars"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func TestZcash_NewZcashInstanceResourceManager(t *testing.T) {

	zcash := NewZcashInstanceResourceManager()
	assert.NotNil(t, zcash)
}

func TestZcash_CreateDeploymentResource(t *testing.T) {

	ctx := fake_zbi.InitContext()
	helper.Config.LoadConfig(ctx)
	helper.Config.LoadTemplates(ctx)

	vars.ManagerFactory = NewResourceManagerFactory()
	vars.ManagerFactory.Init(ctx)

	zcash := NewZcashInstanceResourceManager()
	assert.NotNil(t, zcash)

	var projIngress unstructured.Unstructured
	zcash.CreateInstanceResource(ctx, &projIngress, &data.Project1, data.Instance1)
}

func TestZcash_CreateIngressResource(t *testing.T) {

}

func TestZcash_CreateStartResource(t *testing.T) {

}

func TestZcash_CreateSnapshotResource(t *testing.T) {

}

func TestZcash_CreateSnapshotScheduleResource(t *testing.T) {

}

func TestZcash_CreateRotationResource(t *testing.T) {

}
