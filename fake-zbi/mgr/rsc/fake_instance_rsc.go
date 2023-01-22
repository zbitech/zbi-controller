package rsc

import (
	"context"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type FakeInstanceResourceManager struct {
	FakeCreateInstanceResource         func(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([][]unstructured.Unstructured, error)
	FakeCreateUpdateResource           func(ctx context.Context, project *model.Project, instance *model.Instance, peers ...model.Instance) ([][]unstructured.Unstructured, error)
	FakeCreateIngressResource          func(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, action model.EventAction) (*unstructured.Unstructured, error)
	FakeCreateStartResource            func(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error)
	FakeCreateStopResource             func(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]model.KubernetesResource, []unstructured.Unstructured, error)
	FakeCreateRepairResource           func(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([]unstructured.Unstructured, error)
	FakeCreateSnapshotResource         func(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error)
	FakeCreateSnapshotScheduleResource func(ctx context.Context, project *model.Project, instance *model.Instance, scheduleType model.SnapshotScheduleType) ([]unstructured.Unstructured, error)
	FakeCreateRotationResource         func(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error)
	FakeCreateDeleteResource           func(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, resources []model.KubernetesResource) ([]model.KubernetesResource, []unstructured.Unstructured, error)
}

func NewFakeInstanceResourceManager() interfaces.InstanceResourceManagerIF {
	return &FakeInstanceResourceManager{}
}

func (f FakeInstanceResourceManager) CreateInstanceResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([][]unstructured.Unstructured, error) {
	return f.FakeCreateInstanceResource(ctx, projIngress, project, instance, peers...)
}

func (f FakeInstanceResourceManager) CreateUpdateResource(ctx context.Context, project *model.Project, instance *model.Instance, peers ...model.Instance) ([][]unstructured.Unstructured, error) {
	return f.FakeCreateUpdateResource(ctx, project, instance, peers...)
}

func (f FakeInstanceResourceManager) CreateIngressResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, action model.EventAction) (*unstructured.Unstructured, error) {
	return f.FakeCreateIngressResource(ctx, projIngress, project, instance, action)
}

func (f FakeInstanceResourceManager) CreateStartResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error) {
	return f.FakeCreateStartResource(ctx, projIngress, project, instance)
}

func (f FakeInstanceResourceManager) CreateStopResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]model.KubernetesResource, []unstructured.Unstructured, error) {
	return f.FakeCreateStopResource(ctx, projIngress, project, instance)
}

func (f FakeInstanceResourceManager) CreateRepairResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([]unstructured.Unstructured, error) {
	return f.FakeCreateRepairResource(ctx, projIngress, project, instance, peers...)
}

func (f FakeInstanceResourceManager) CreateSnapshotResource(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error) {
	return f.FakeCreateSnapshotResource(ctx, project, instance)
}

func (f FakeInstanceResourceManager) CreateSnapshotScheduleResource(ctx context.Context, project *model.Project, instance *model.Instance, scheduleType model.SnapshotScheduleType) ([]unstructured.Unstructured, error) {
	return f.FakeCreateSnapshotScheduleResource(ctx, project, instance, scheduleType)
}

func (f FakeInstanceResourceManager) CreateRotationResource(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error) {
	return f.FakeCreateRotationResource(ctx, project, instance)
}

func (f FakeInstanceResourceManager) CreateDeleteResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, resources []model.KubernetesResource) ([]model.KubernetesResource, []unstructured.Unstructured, error) {
	return f.FakeCreateDeleteResource(ctx, projIngress, project, instance, resources)
}
