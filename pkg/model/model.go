package model

import (
	"time"
)

var (
	TESTNET = []string{"testnet.z.cash"}
	MAINNET = []string{"mainnet.z.cash"}

	ZCASH_MAX_CONNECTIONS   = "6"
	ZCASH_RPCCLIENT_TIMEOUT = "30"
	ZCASH_SOLVER            = "tromp"
	ZCASH_MAINNET_SOLVER    = "default"
	ZCASH_TESTNET_SOLVER    = "tromp"
)

type SnapshotRequest struct {
	Version       string `json:"version" validate:"required"`
	VolumeName    string `json:"volumeName"`
	Namespace     string `json:"namespace"`
	SnapshotClass string `json:"snapshotClass"`
	Labels        map[string]string
}

type SnapshotScheduleRequest struct {
	Version          string               `json:"version" validate:"required"`
	Schedule         SnapshotScheduleType `json:"schedule" validate:"required"`
	VolumeName       string               `json:"volumeName"`
	Namespace        string               `json:"namespace"`
	SnapshotClass    string
	BackupExpiration string
	MaxBackupCount   int
	Labels           map[string]string
}

type ProjectSpec struct {
	Name      string            `json:"name"`
	Network   NetworkType       `json:"network"`
	Owner     string            `json:"owner"`
	TeamId    string            `json:"team"`
	Namespace string            `json:"namespace"`
	Instances string            `json:"instances"`
	Labels    map[string]string `json:"labels"`
}

type InstanceSpec struct {
	Name               string                 `json:"name"`
	Project            string                 `json:"project"`
	Namespace          string                 `json:"namespace"`
	ServiceAccountName string                 `json:"serviceAccountName"`
	DataVolumeName     string                 `json:"dataVolumeName"`
	DomainName         string                 `json:"domainName"`
	DomainSecret       string                 `json:"domainSecret"`
	Labels             map[string]string      `json:"labels"`
	Envoy              EnvoySpec              `json:"envoy"`
	Images             map[string]string      `json:"images"`
	Ports              map[string]int32       `json:"ports"`
	Properties         map[string]interface{} `json:"properties"`
	//	Version            string         `json:"version"`
	//VolumeType     VolumeType         `json:"dataVolumeType"`
	//VolumeSourceType     VolumeSourceType         `json:"dataSourceType"`
	//VolumeSourceName         string                 `json:"dataSource"`
}

type VolumeSpec struct {
	VolumeName     string
	StorageClass   string
	Namespace      string
	VolumeDataType string
	DataSourceType DataSourceType
	SourceName     string
	Size           string
	Labels         map[string]string
}

type SnapshotSpec struct {
	SnapshotName  string
	Namespace     string
	Owner         string
	VolumeName    string
	Labels        map[string]string
	SnapshotClass string
}

type SnapshotScheduleSpec struct {
	ScheduleName     string
	Namespace        string
	Schedule         string
	ScheduleType     SnapshotScheduleType
	SnapshotClass    string
	BackupExpiration string
	MaxBackupCount   int
	Labels           map[string]string
	ClaimLabels      map[string]string
	SnapshotLabels   map[string]string
}

type EnvoySpec struct {
	Image                 string
	Command               string
	Port                  int32
	Timeout               float32
	AccessAuthorization   bool
	AuthServerURL         string
	AuthServerPort        int32
	AuthenticationEnabled bool
}

type Project struct {
	Name      string               `json:"name"`
	Network   NetworkType          `json:"network"`
	Owner     string               `json:"owner"`
	TeamId    string               `json:"team"`
	Instances []map[string]string  `json:"instances"`
	Resources []KubernetesResource `json:"resources,omitempty"`
}

type Instance struct {
	Name         string               `json:"name"`
	InstanceType InstanceType         `json:"instanceType"`
	Project      string               `json:"project"`
	Resources    *KubernetesResources `json:"resources"`
}

type ResourceRequest struct {
	VolumeType          DataVolumeType         `json:"volumeType,omitempty"`
	VolumeSize          string                 `json:"volumeSize,omitempty"`
	VolumeSourceType    DataSourceType         `json:"volumeSourceType,omitempty"`
	VolumeSourceName    string                 `json:"volumeSourceName,omitempty"`
	VolumeSourceProject string                 `json:"volumeSourceProject"`
	Cpu                 string                 `json:"cpu,omitempty"`
	Memory              string                 `json:"memory,omitempty"`
	Peers               []string               `json:"peers"`
	Properties          map[string]interface{} `json:"properties,omitempty"`
}

type KubernetesResources struct {
	Resources []KubernetesResource `json:"resources"`
	Snapshots []KubernetesResource `json:"snapshots"`
	Schedules []KubernetesResource `json:"schedule"`
}

type KubernetesResource struct {
	Name       string                 `json:"name,omitempty"`
	Namespace  string                 `json:"namespace,omitempty"`
	Type       ResourceObjectType     `json:"type,omitempty"`
	Status     string                 `json:"status,omitempty"`
	Created    *time.Time             `json:"created,omitempty"`
	Updated    *time.Time             `json:"updated,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type ResourceObjectType string

const (
	ResourceNamespace             ResourceObjectType = "Namespace"
	ResourceDeployment            ResourceObjectType = "Deployment"
	ResourceService               ResourceObjectType = "Service"
	ResourceConfigMap             ResourceObjectType = "ConfigMap"
	ResourceSecret                ResourceObjectType = "Secret"
	ResourcePod                   ResourceObjectType = "Pod"
	ResourcePersistentVolume      ResourceObjectType = "PersistentVolume"
	ResourcePersistentVolumeClaim ResourceObjectType = "PersistentVolumeClaim"
	ResourceVolumeSnapshot        ResourceObjectType = "VolumeSnapshot"
	ResourceVolumeSnapshotClass   ResourceObjectType = "VolumeSnapshotClass"
	ResourceSnapshotSchedule      ResourceObjectType = "SnapshotSchedule"
	ResourceHTTPProxy             ResourceObjectType = "HTTPProxy"
)

type EventAction string

const (
	EventActionCreate         EventAction = "create"
	EventActionDelete         EventAction = "delete"
	EventActionUpdate         EventAction = "update"
	EventActionResource       EventAction = "resource"
	EventActionDeactivate     EventAction = "deactivate"
	EventActionReactivate     EventAction = "reactivate"
	EventActionRepair         EventAction = "repair"
	EventActionSnapshot       EventAction = "snapshot"
	EventActionSchedule       EventAction = "schedule"
	EventActionPurge          EventAction = "purge"
	EventActionStopInstance   EventAction = "stop"
	EventActionStartInstance  EventAction = "start"
	EventActionRotate         EventAction = "rotate"
	EventActionUpdatePolicy   EventAction = "updatepolicy"
	EventActionAddMember      EventAction = "addmember"
	EventActionRemoveMember   EventAction = "removemember"
	EventActionUpdateMember   EventAction = "updatemember"
	EventActionRegister       EventAction = "register"
	EventActionCreateKey      EventAction = "createkey"
	EventActionDeleteKey      EventAction = "deletekey"
	EventActionChangePassword EventAction = "changepassword"
	EventActionChangeEmail    EventAction = "changeemail"
	EventActionUpdateProfile  EventAction = "updateprofile"
	EventActionAcceptInvite   EventAction = "acceptinvite"
	EventActionRejectInvite   EventAction = "rejectinvite"
	EventActionExpireInvite   EventAction = "expireinvite"
)

type NetworkType string

const (
	NetworkTypeMain NetworkType = "mainnet"
	NetworkTypeTest NetworkType = "testnet"
)

type InstanceType string

const (
	InstanceTypeZCASH InstanceType = "zcash"
	InstanceTypeLWD   InstanceType = "lwd"
)

type StatusType string

const (
	StatusNew         StatusType = "new"
	StatusActive      StatusType = "active"
	StatusInActive    StatusType = "inactive"
	StatusFailed      StatusType = "failed"
	StatusPending     StatusType = "pending"
	StatusProgressing StatusType = "progressing"
	StatusBound       StatusType = "bound"
	StatusRunning     StatusType = "running"
	StatusStopped     StatusType = "stopped"
	StatusValid       StatusType = "valid"
	StatusReady       StatusType = "ready"
)

type SnapshotScheduleType string

const (
	HourlySnapshotScheduled SnapshotScheduleType = "hourly"
	DailySnapshotSchedule   SnapshotScheduleType = "daily"
	WeeklySnapshotSchedule  SnapshotScheduleType = "weekly"
	MonthlySnapshotSchedule SnapshotScheduleType = "monthly"
)

type DataSourceType string

const (
	NoDataSource       DataSourceType = "none"
	NewDataSource      DataSourceType = "new"
	VolumeDataSource   DataSourceType = "pvc"
	SnapshotDataSource DataSourceType = "snapshot"
)

type DataVolumeType string

const (
	EphemeralDataVolume  DataVolumeType = "eph"
	PersistentDataVolume DataVolumeType = "pvc"
)

type Metadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Labels            map[string]string `json:"labels,omitempty"`
	ResourceVersion   string            `json:"resourceVersion,omitempty"`
	Uid               string            `json:"uid,omitempty"`
	CreationTimestamp string            `json:"creationTimestamp"`
}

type IngressCondition struct {
	Prefix string `json:"prefix,omitempty"`
}

type IngressService struct {
	Name string `json:"name"`
	Port int32  `json:"port"`
}

type IngressRoute struct {
	Conditions        []IngressCondition       `json:"conditions,omitempty"`
	Services          []IngressService         `json:"services,omitempty"`
	PathRewritePolicy IngressPathRewritePolicy `json:"pathRewritePolicy,omitempty"`
}

type IngressPathRewritePolicy struct {
	ReplacePrefix []struct {
		Replacement string `json:"replacement,omitempty"`
	} `json:"replacePrefix,omitempty"`
}

type IngressInclude struct {
	Name       string             `json:"name,omitempty"`
	Namespace  string             `json:"namespace,omitempty"`
	Conditions []IngressCondition `json:"conditions,omitempty"`
}

type Ingress struct {
	ApiVersion string   `json:"apiVersion"`
	Metadata   Metadata `json:"metadata"`
	Spec       struct {
		Includes    []IngressInclude `json:"includes,omitempty"`
		Virtualhost struct {
			Fqdn string `json:"fqdn"`
			Tls  struct {
				SecretName string `json:"secretName"`
			} `json:"tls"`
		} `json:"virtualhost,omitempty"`
		Routes []IngressRoute `json:"routes,omitempty"`
	} `json:"spec"`
	Status struct {
		Conditions []struct {
			Errors []struct {
				Message string `json:"message,omitempty"`
				Reason  string `json:"reason,omitempty"`
				Status  string `json:"status,omitempty"`
				Type    string `json:"type,omitempty"`
			} `json:"errors,omitempty"`
			LastTransitionTime string `json:"lastTransitionTime,omitempty"`
			Message            string `json:"message,omitempty"`
			ObservedGeneration int    `json:"observedGeneration,omitempty"`
			Reason             string `json:"reason,omitempty"`
			Status             string `json:"status,omitempty"`
			Type               string `json:"type,omitempty"`
		} `json:"conditions,omitempty"`
		CurrentStatus string `json:"currentStatus,omitempty"`
		Description   string `json:"description,omitempty"`
		LoadBalancer  struct {
			Ingress []struct {
				Hostname string `json:"hostname,omitempty"`
			} `json:"ingress,omitempty"`
		} `json:"loadBalancer,omitempty"`
	} `json:"status,omitempty"`
}

func (project *Project) GetInstanceType(name string) InstanceType {
	for _, entry := range project.Instances {
		if entry["name"] == name {
			return InstanceType(entry["type"])
		}
	}

	return ""
}

func appendResource(resources []KubernetesResource, newResource KubernetesResource) []KubernetesResource {

	for index := 0; index < len(resources); index++ {
		if resources[index].Name == newResource.Name && resources[index].Type == newResource.Type {
			resources[index].Properties = newResource.Properties
			resources[index].Updated = newResource.Updated
			resources[index].Status = newResource.Status
			return resources
		}
	}

	return append(resources, newResource)
}

func (instance *Instance) GetResource(resourceName string, resourceType ResourceObjectType) *KubernetesResource {

	var resources []KubernetesResource
	if resourceType == ResourceVolumeSnapshot {
		resources = instance.Resources.Snapshots
	} else if resourceType == ResourceSnapshotSchedule {
		resources = instance.Resources.Schedules
	} else {
		resources = instance.Resources.Resources
	}

	for _, resource := range resources {
		if resource.Name == resourceName && resource.Type == resourceType {
			return &resource
		}
	}

	return nil
}

func (instance *Instance) GetResourceByName(resourceName string) *KubernetesResource {

	for _, resource := range instance.Resources.Resources {
		if resource.Name == resourceName {
			return &resource
		}
	}

	for _, resource := range instance.Resources.Snapshots {
		if resource.Name == resourceName {
			return &resource
		}
	}

	for _, resource := range instance.Resources.Schedules {
		if resource.Name == resourceName {
			return &resource
		}
	}

	return nil
}

func (instance *Instance) GetResourceByType(resourceType ResourceObjectType) *KubernetesResource {

	var resources []KubernetesResource
	if resourceType == ResourceVolumeSnapshot {
		resources = instance.Resources.Snapshots
	} else if resourceType == ResourceSnapshotSchedule {
		resources = instance.Resources.Schedules
	} else {
		resources = instance.Resources.Resources
	}

	for _, resource := range resources {
		if resource.Type == resourceType {
			return &resource
		}
	}

	return nil
}

func (instance *Instance) AddResource(resource KubernetesResource) {

	if resource.Type == ResourceVolumeSnapshot {
		instance.Resources.Snapshots = append(instance.Resources.Snapshots, resource)
	} else if resource.Type == ResourceSnapshotSchedule {
		instance.Resources.Schedules = append(instance.Resources.Schedules, resource)
	} else {
		instance.Resources.Resources = append(instance.Resources.Resources, resource)
	}
}

func (instance *Instance) AddResources(resources ...KubernetesResource) {
	for _, resource := range resources {
		instance.AddResource(resource)
	}
}

func (instance *Instance) HasResources() bool {
	return len(instance.Resources.Resources) > 0 ||
		len(instance.Resources.Snapshots) > 0 ||
		len(instance.Resources.Schedules) > 0
}

func (instance *Instance) GetResources() []KubernetesResource {
	var resources = make([]KubernetesResource, 0)

	resources = append(resources, instance.Resources.Resources...)
	resources = append(resources, instance.Resources.Snapshots...)
	resources = append(resources, instance.Resources.Schedules...)

	return resources
}

func (project *Project) GetNamespace() string {
	return project.Name
}

func (project *Project) GetResource(resourceName string, resourceType ResourceObjectType) *KubernetesResource {
	for _, resource := range project.Resources {
		if resource.Name == resourceName && resource.Type == resourceType {
			return &resource
		}
	}

	return nil
}

func (project *Project) GetResourceByName(resourceName string) *KubernetesResource {
	for _, resource := range project.Resources {
		if resource.Name == resourceName {
			return &resource
		}
	}

	return nil
}

func (project *Project) GetResourceByType(resourceType ResourceObjectType) *KubernetesResource {
	for _, resource := range project.Resources {
		if resource.Type == resourceType {
			return &resource
		}
	}

	return nil
}

func (project *Project) AddResource(resource KubernetesResource) {
	if project.Resources == nil {
		project.Resources = make([]KubernetesResource, 0)
	}
	project.Resources = appendResource(project.Resources, resource)
}

func (project *Project) AddResources(resources ...KubernetesResource) {
	for _, resource := range resources {
		project.AddResource(resource)
	}
}
