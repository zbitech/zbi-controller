package zklient

import (
	"context"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/model"
)

type FakeZBIClient struct {
	client interfaces.KlientIF

	FakeRunInformer func(ctx context.Context)

	FakeCreateProject             func(ctx context.Context, project *model.Project) error
	FakeGetProjects               func(ctx context.Context) ([]model.Project, error)
	FakeRepairProject             func(ctx context.Context, project *model.Project) error
	FakeDeleteProject             func(ctx context.Context, project *model.Project, instances []model.Instance) error
	FakeGetProjectResources       func(ctx context.Context, project string) ([]model.KubernetesResource, error)
	FakeGetProjectResource        func(ctx context.Context, project, resourceName string, resourceType model.ResourceObjectType) (*model.KubernetesResource, error)
	FakeCreateInstance            func(ctx context.Context, project *model.Project, instance *model.Instance) error
	FakeGetAllInstances           func(ctx context.Context, project *model.Project) ([]model.Instance, error)
	FakeGetInstances              func(ctx context.Context, project *model.Project, instances []string) ([]model.Instance, error)
	FakeGetInstanceResources      func(ctx context.Context, project *model.Project, instance string) (*model.KubernetesResources, error)
	FakeGetInstanceResource       func(ctx context.Context, project *model.Project, instance, resourceName string, resourceType model.ResourceObjectType) (*model.KubernetesResource, error)
	FakeDeleteInstanceResource    func(ctx context.Context, project *model.Project, instance *model.Instance, resourceName string, resourceType model.ResourceObjectType) error
	FakeUpdateInstance            func(ctx context.Context, project *model.Project, instance *model.Instance) error
	FakeDeleteInstance            func(ctx context.Context, project *model.Project, instance *model.Instance) error
	FakeRepairInstance            func(ctx context.Context, project *model.Project, instance *model.Instance) error
	FakeStopInstance              func(ctx context.Context, project *model.Project, instance *model.Instance) error
	FakeStartInstance             func(ctx context.Context, project *model.Project, instance *model.Instance) error
	FakeRotateInstanceCredentials func(ctx context.Context, project *model.Project, instance *model.Instance) error
	FakeCreateSnapshot            func(ctx context.Context, project *model.Project, instance *model.Instance) error
	FakeCreateSnapshotSchedule    func(ctx context.Context, project *model.Project, instance *model.Instance, schedule model.SnapshotScheduleType) error
	FakeGetProject                func(ctx context.Context, project string) (*model.Project, error)
	FakeGetInstance               func(ctx context.Context, project *model.Project, instance string) (*model.Instance, error)
}

func (f FakeZBIClient) GetAllInstances(ctx context.Context, project *model.Project) ([]model.Instance, error) {
	return f.FakeGetAllInstances(ctx, project)
}

func NewFakeZBIClient(client interfaces.KlientIF) interfaces.ZBIClientIF {
	return &FakeZBIClient{client: client}
}

func (f FakeZBIClient) GetProject(ctx context.Context, project string) (*model.Project, error) {
	return f.FakeGetProject(ctx, project)
}

func (f FakeZBIClient) GetInstance(ctx context.Context, project *model.Project, instance string) (*model.Instance, error) {
	return f.FakeGetInstance(ctx, project, instance)
}

func (f FakeZBIClient) RunInformer(ctx context.Context) {
	f.FakeRunInformer(ctx)
}

func (f FakeZBIClient) CreateProject(ctx context.Context, project *model.Project) error {
	return f.FakeCreateProject(ctx, project)
}

func (f FakeZBIClient) GetProjects(ctx context.Context) ([]model.Project, error) {
	return f.FakeGetProjects(ctx)
}

func (f FakeZBIClient) RepairProject(ctx context.Context, project *model.Project) error {
	return f.FakeRepairProject(ctx, project)
}

func (f FakeZBIClient) DeleteProject(ctx context.Context, project *model.Project, instances []model.Instance) error {
	return f.FakeDeleteProject(ctx, project, instances)
}

func (f FakeZBIClient) GetProjectResources(ctx context.Context, project string) ([]model.KubernetesResource, error) {
	return f.FakeGetProjectResources(ctx, project)
}

func (f FakeZBIClient) GetProjectResource(ctx context.Context, project, resourceName string, resourceType model.ResourceObjectType) (*model.KubernetesResource, error) {
	return f.FakeGetProjectResource(ctx, project, resourceName, resourceType)
}

func (f FakeZBIClient) CreateInstance(ctx context.Context, project *model.Project, instance *model.Instance) error {
	return f.FakeCreateInstance(ctx, project, instance)
}

func (f FakeZBIClient) GetInstances(ctx context.Context, project *model.Project, instances []string) ([]model.Instance, error) {
	return f.FakeGetInstances(ctx, project, instances)
}

func (f FakeZBIClient) GetInstanceResources(ctx context.Context, project *model.Project, instance string) (*model.KubernetesResources, error) {
	return f.FakeGetInstanceResources(ctx, project, instance)
}

func (f FakeZBIClient) GetInstanceResource(ctx context.Context, project *model.Project, instance, resourceName string, resourceType model.ResourceObjectType) (*model.KubernetesResource, error) {
	return f.FakeGetInstanceResource(ctx, project, instance, resourceName, resourceType)
}

func (f FakeZBIClient) DeleteInstanceResource(ctx context.Context, project *model.Project, instance *model.Instance, resourceName string, resourceType model.ResourceObjectType) error {
	return f.FakeDeleteInstanceResource(ctx, project, instance, resourceName, resourceType)
}

func (f FakeZBIClient) UpdateInstance(ctx context.Context, project *model.Project, instance *model.Instance) error {
	return f.FakeUpdateInstance(ctx, project, instance)
}

func (f FakeZBIClient) DeleteInstance(ctx context.Context, project *model.Project, instance *model.Instance) error {
	return f.FakeDeleteInstance(ctx, project, instance)
}

func (f FakeZBIClient) RepairInstance(ctx context.Context, project *model.Project, instance *model.Instance) error {
	return f.FakeRepairInstance(ctx, project, instance)
}

func (f FakeZBIClient) StopInstance(ctx context.Context, project *model.Project, instance *model.Instance) error {
	return f.FakeStopInstance(ctx, project, instance)
}

func (f FakeZBIClient) StartInstance(ctx context.Context, project *model.Project, instance *model.Instance) error {
	return f.FakeStartInstance(ctx, project, instance)
}

func (f FakeZBIClient) RotateInstanceCredentials(ctx context.Context, project *model.Project, instance *model.Instance) error {
	return f.FakeRotateInstanceCredentials(ctx, project, instance)
}

func (f FakeZBIClient) CreateSnapshot(ctx context.Context, project *model.Project, instance *model.Instance) error {
	return f.FakeCreateSnapshot(ctx, project, instance)
}

func (f FakeZBIClient) CreateSnapshotSchedule(ctx context.Context, project *model.Project, instance *model.Instance, schedule model.SnapshotScheduleType) error {
	return f.FakeCreateSnapshotSchedule(ctx, project, instance, schedule)
}
