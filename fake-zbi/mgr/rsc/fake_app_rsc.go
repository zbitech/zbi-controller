package rsc

import (
	"context"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type FakeAppResourceManager struct {
	FakeCreateSnapshotResource         func(ctx context.Context, project, instance string, req *model.SnapshotRequest) ([]unstructured.Unstructured, error)
	FakeCreateSnapshotScheduleResource func(ctx context.Context, project, instance string, req *model.SnapshotScheduleRequest) ([]unstructured.Unstructured, error)
	FakeCreateVolumeResource           func(ctx context.Context, project, instance string, volumes ...model.VolumeSpec) ([]unstructured.Unstructured, error)
}

func NewFakeIngressResourceManager() interfaces.AppResourceManagerIF {
	return &FakeAppResourceManager{}
}

func (f FakeAppResourceManager) CreateSnapshotResource(ctx context.Context, project, instance string, req *model.SnapshotRequest) ([]unstructured.Unstructured, error) {
	return f.FakeCreateSnapshotResource(ctx, project, instance, req)
}

func (f FakeAppResourceManager) CreateSnapshotScheduleResource(ctx context.Context, project, instance string, req *model.SnapshotScheduleRequest) ([]unstructured.Unstructured, error) {
	return f.FakeCreateSnapshotScheduleResource(ctx, project, instance, req)
}

func (f FakeAppResourceManager) CreateVolumeResource(ctx context.Context, project, instance string, volumes ...model.VolumeSpec) ([]unstructured.Unstructured, error) {
	return f.FakeCreateVolumeResource(ctx, project, instance, volumes...)
}
