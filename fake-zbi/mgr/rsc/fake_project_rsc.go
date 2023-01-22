package rsc

import (
	"context"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type FakeProjectResourceManager struct {
	FakeCreateProjectResource          func(ctx context.Context, project *model.Project) ([]unstructured.Unstructured, error)
	FakeCreateProjectIngressResource   func(ctx context.Context, appIngress *unstructured.Unstructured, project *model.Project, action model.EventAction) ([]unstructured.Unstructured, error)
	FakeCreateInstanceResource         func(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([]unstructured.Unstructured, [][]unstructured.Unstructured, error)
	FakeCreateUpdateResource           func(ctx context.Context, project *model.Project, instance *model.Instance, peers ...model.Instance) ([][]unstructured.Unstructured, error)
	FakeCreateStartResource            func(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error)
	FakeCreateStopResource             func(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]model.KubernetesResource, []unstructured.Unstructured, error)
	FakeCreateRepairResource           func(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([]unstructured.Unstructured, error)
	FakeCreateIngressResource          func(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, action model.EventAction) (*unstructured.Unstructured, error)
	FakeCreateSnapshotResource         func(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error)
	FakeCreateSnapshotScheduleResource func(ctx context.Context, project *model.Project, instance *model.Instance, schedule model.SnapshotScheduleType) ([]unstructured.Unstructured, error)
	FakeCreateRotationResource         func(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error)
	FakeCreateDeleteResource           func(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, resources []model.KubernetesResource) ([]model.KubernetesResource, []unstructured.Unstructured, error)
}

func NewFakeProjectResourceManager() interfaces.ProjectResourceManagerIF {
	return &FakeProjectResourceManager{}
}

func (f FakeProjectResourceManager) CreateProjectResource(ctx context.Context, project *model.Project) ([]unstructured.Unstructured, error) {
	return f.FakeCreateProjectResource(ctx, project)
}

func (f FakeProjectResourceManager) CreateProjectIngressResource(ctx context.Context, appIngress *unstructured.Unstructured, project *model.Project, action model.EventAction) ([]unstructured.Unstructured, error) {
	return f.FakeCreateProjectIngressResource(ctx, appIngress, project, action)
}

func (f FakeProjectResourceManager) CreateInstanceResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([]unstructured.Unstructured, [][]unstructured.Unstructured, error) {
	return f.FakeCreateInstanceResource(ctx, projIngress, project, instance, peers...)
}

func (f FakeProjectResourceManager) CreateUpdateResource(ctx context.Context, project *model.Project, instance *model.Instance, peers ...model.Instance) ([][]unstructured.Unstructured, error) {
	return f.FakeCreateUpdateResource(ctx, project, instance, peers...)
}

func (f FakeProjectResourceManager) CreateStartResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error) {
	return f.FakeCreateStartResource(ctx, projIngress, project, instance)
}

func (f FakeProjectResourceManager) CreateStopResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]model.KubernetesResource, []unstructured.Unstructured, error) {
	return f.FakeCreateStopResource(ctx, projIngress, project, instance)
}

func (f FakeProjectResourceManager) CreateRepairResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([]unstructured.Unstructured, error) {
	return f.FakeCreateRepairResource(ctx, projIngress, project, instance, peers...)
}

func (f FakeProjectResourceManager) CreateIngressResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, action model.EventAction) (*unstructured.Unstructured, error) {
	return f.FakeCreateIngressResource(ctx, projIngress, project, instance, action)
}

func (f FakeProjectResourceManager) CreateSnapshotResource(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error) {
	return f.FakeCreateSnapshotResource(ctx, project, instance)
}

func (f FakeProjectResourceManager) CreateSnapshotScheduleResource(ctx context.Context, project *model.Project, instance *model.Instance, schedule model.SnapshotScheduleType) ([]unstructured.Unstructured, error) {
	return f.FakeCreateSnapshotScheduleResource(ctx, project, instance, schedule)
}

func (f FakeProjectResourceManager) CreateRotationResource(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error) {
	return f.FakeCreateRotationResource(ctx, project, instance)
}

func (f FakeProjectResourceManager) CreateDeleteResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, resources []model.KubernetesResource) ([]model.KubernetesResource, []unstructured.Unstructured, error) {
	return f.FakeCreateDeleteResource(ctx, projIngress, project, instance, resources)
}
