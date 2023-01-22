package interfaces

import (
	"context"
	"github.com/zbitech/controller/pkg/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ResourceManagerFactoryIF interface {
	Init(ctx context.Context) error
	GetAppResourceManager(ctx context.Context) AppResourceManagerIF
	GetProjectDataManager(ctx context.Context) ProjectResourceManagerIF
}

type AppResourceManagerIF interface {
	CreateSnapshotResource(ctx context.Context, project, instance string, req *model.SnapshotRequest) ([]unstructured.Unstructured, error)
	CreateSnapshotScheduleResource(ctx context.Context, project, instance string, req *model.SnapshotScheduleRequest) ([]unstructured.Unstructured, error)
	CreateVolumeResource(ctx context.Context, project, instance string, volumes ...model.VolumeSpec) ([]unstructured.Unstructured, error)
}

type ProjectResourceManagerIF interface {
	CreateProjectResource(ctx context.Context, project *model.Project) ([]unstructured.Unstructured, error)
	CreateProjectIngressResource(ctx context.Context, appIngress *unstructured.Unstructured, project *model.Project, action model.EventAction) ([]unstructured.Unstructured, error)
	CreateInstanceResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([]unstructured.Unstructured, [][]unstructured.Unstructured, error)
	CreateUpdateResource(ctx context.Context, project *model.Project, instance *model.Instance, peers ...model.Instance) ([][]unstructured.Unstructured, error)
	CreateStartResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error)
	CreateStopResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]model.KubernetesResource, []unstructured.Unstructured, error)
	CreateRepairResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([]unstructured.Unstructured, error)
	CreateIngressResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, action model.EventAction) (*unstructured.Unstructured, error)
	CreateSnapshotResource(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error)
	CreateSnapshotScheduleResource(ctx context.Context, project *model.Project, instance *model.Instance, schedule model.SnapshotScheduleType) ([]unstructured.Unstructured, error)
	CreateRotationResource(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error)
	CreateDeleteResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, resources []model.KubernetesResource) ([]model.KubernetesResource, []unstructured.Unstructured, error)
}

type InstanceResourceManagerIF interface {
	CreateInstanceResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([][]unstructured.Unstructured, error)
	CreateUpdateResource(ctx context.Context, project *model.Project, instance *model.Instance, peers ...model.Instance) ([][]unstructured.Unstructured, error)
	CreateIngressResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, action model.EventAction) (*unstructured.Unstructured, error)
	CreateStartResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error)
	CreateStopResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]model.KubernetesResource, []unstructured.Unstructured, error)
	CreateRepairResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([]unstructured.Unstructured, error)
	CreateSnapshotResource(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error)
	CreateSnapshotScheduleResource(ctx context.Context, project *model.Project, instance *model.Instance, scheduleType model.SnapshotScheduleType) ([]unstructured.Unstructured, error)
	CreateRotationResource(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error)
	CreateDeleteResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, resources []model.KubernetesResource) ([]model.KubernetesResource, []unstructured.Unstructured, error)
}
