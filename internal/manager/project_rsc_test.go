package manager

import (
	"github.com/stretchr/testify/assert"
	fake_zbi "github.com/zbitech/controller/fake-zbi"
	"github.com/zbitech/controller/fake-zbi/data"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/model"
	"testing"
)

func TestProject_NewProjectResourceManager(t *testing.T) {

	var manager = NewProjectResourceManager(map[model.InstanceType]interfaces.InstanceResourceManagerIF{})

	assert.NotNil(t, manager)
}

func TestProject_CreateProjectResource(t *testing.T) {
	ctx := fake_zbi.InitContext()
	var manager = NewProjectResourceManager(map[model.InstanceType]interfaces.InstanceResourceManagerIF{})

	assert.NotNil(t, manager)

	project := &data.Project1
	resources, err := manager.CreateProjectResource(ctx, project)
	assert.NoError(t, err)
	assert.NotNil(t, resources)
}

func TestProject_CreateProjectIngressResource(t *testing.T) {
	//	ctx := initContext()
	//	var manager = NewProjectResourceManager(map[model.InstanceType]interfaces.InstanceResourceManagerIF{})

	//	assert.NotNil(t, manager)

	//	project := &data.Project1
	//	resources, err := manager.CreateProjectIngressResource(ctx, nil, project, model.EventActionCreate)
	//	assert.NoError(t, err)
	//	assert.NotNil(t, resources)
}

func TestProject_CreateDeploymentResource(t *testing.T) {

}

func TestProject_CreateStartResource(t *testing.T) {

}

func TestProject_CreateIngressResource(t *testing.T) {

}

func TestProject_CreateSnapshotResource(t *testing.T) {

}

func TestProject_CreateSnapshotScheduleResource(t *testing.T) {

}

func TestProject_CreateRotationResource(t *testing.T) {

}
