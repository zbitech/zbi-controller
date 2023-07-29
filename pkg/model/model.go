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

type Project struct {
	Id        string               `json:"id"`
	Name      string               `json:"name"`
	Network   NetworkType          `json:"network"`
	Owner     string               `json:"owner"`
	TeamId    string               `json:"team"`
	Instances []Instance           `json:"instances"`
	Resources []KubernetesResource `json:"resources,omitempty"`
}

type Instance struct {
	Id           string               `json:"id"`
	Name         string               `json:"name"`
	InstanceType InstanceType         `json:"instanceType"`
	Project      string               `json:"project,omitempty"`
	Request      *ResourceRequest     `json:"request"`
	Resources    *KubernetesResources `json:"resources,omitempty"`
}

type ResourceRequest struct {
	VolumeType          DataVolumeType         `json:"volumeType,omitempty"`
	VolumeSize          string                 `json:"volumeSize,omitempty"`
	VolumeSourceType    DataSourceType         `json:"volumeSourceType,omitempty"`
	VolumeSourceName    string                 `json:"volumeSourceName,omitempty"`
	VolumeSourceProject string                 `json:"volumeSourceProject,omitempty"`
	Cpu                 string                 `json:"cpu,omitempty"`
	Memory              string                 `json:"memory,omitempty"`
	Peers               []string               `json:"peers,omitempty"`
	Properties          map[string]interface{} `json:"properties,omitempty"`
}

type KubernetesResources struct {
	Resources []KubernetesResource `json:"resources"`
	Snapshots []KubernetesResource `json:"snapshots"`
	Schedules []KubernetesResource `json:"schedule"`
}

type KubernetesResource struct {
	//	Id         string                 `json:"id"`
	Name       string                 `json:"name,omitempty"`
	Namespace  string                 `json:"namespace,omitempty"`
	Type       ResourceObjectType     `json:"type,omitempty"`
	Status     string                 `json:"status,omitempty"`
	Created    *time.Time             `json:"created,omitempty"`
	Updated    *time.Time             `json:"updated,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

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
