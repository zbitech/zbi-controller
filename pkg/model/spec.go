package model

type ProjectSpec struct {
	Name         string            `json:"name"`
	Network      NetworkType       `json:"network"`
	Owner        string            `json:"owner"`
	TeamId       string            `json:"team"`
	Namespace    string            `json:"namespace"`
	Instances    string            `json:"instances"`
	Labels       map[string]string `json:"labels"`
	InstancesMap string            `json:"instanceMap"`
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
