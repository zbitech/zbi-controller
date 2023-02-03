package object

import "github.com/zbitech/controller/pkg/model"

type Settings struct {
	EnableRepository      bool              `json:"enableRepository"`
	Repository            string            `json:"repository"`
	Templates             map[string]string `json:"templates"`
	InformerResync        int               `json:"informerResync"`
	EnableMonitor         bool              `json:"enableMonitor"`
	RequireAuthentication bool              `json:"requireAuthentication"`
}

type PolicyConfig struct {
	TokenExpirationPolicy int32                `json:"tokenExpirationPolicy"`
	StorageClass          string               `json:"storageClass"`
	SnapshotClass         string               `json:"snapshotClass"`
	DomainName            string               `json:"domainName"`
	CertificateName       string               `json:"certificateName"`
	ServiceAccount        string               `json:"serviceAccount"`
	BackupExpiration      string               `json:"backupExpiration"`
	NetworkTypes          []string             `json:"networkTypes"`
	SnapshotTypes         []string             `json:"snapshotTypes"`
	ScheduleTypes         []string             `json:"scheduleTypes"`
	EndpointAccessTypes   []string             `json:"endpointAccessTypes"`
	InstanceTypes         []model.InstanceType `json:"instanceTypes"`
	Envoy                 struct {
		Image                 string   `json:"image"`
		Command               []string `json:"command"`
		Timeout               float32  `json:"timeout"`
		AccessAuthorization   bool     `json:"accessAuthorization"`
		AuthServerURL         string   `json:"authServerURL"`
		AuthServerPort        int32    `json:"authServerPort"`
		AuthenticationEnabled bool     `json:"authenticationEnabled"`
	} `json:"envoy"`
	Limits struct {
		MaxBackupCount int    `json:"maxBackupCount"`
		MaxProjects    int    `json:"maxProjects"`
		MaxInstances   int    `json:"maxInstances"`
		ResourceLimit  string `json:"resourceLimit"`
		MaxCPU         int    `json:"maxCPU"`
		MaxMemroy      string `json:"maxMemroy"`
	} `json:"limits"`
}

type ImageConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Url     string `json:"url"`
}

type KVPair struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type InstanceConfig struct {
	InstanceType model.InstanceType  `json:"instanceType"`
	Name         string              `json:"name"`
	Images       []ImageConfig       `json:"images"`
	Endpoints    map[string][]string `json:"endpoints"`
	Ports        map[string]int32    `json:"ports"`
	Settings     map[string][]KVPair `json:"settings"`
}

type AppConfig struct {
	Settings  Settings         `json:"settings"`
	Policy    PolicyConfig     `json:"policy"`
	Instances []InstanceConfig `json:"instances"`
}

func (ic *InstanceConfig) GetImage(name string) *ImageConfig {
	for _, image := range ic.Images {
		if image.Name == name {
			return &image
		}
	}

	return nil
}

func (ic *InstanceConfig) GetImageRepository(name string) string {
	for _, image := range ic.Images {
		if image.Name == name {
			return image.Url
		}
	}

	return ""
}

func (ic *InstanceConfig) GetPort(name string) int32 {
	port, ok := ic.Ports[name]
	if !ok {
		return -1
	}

	return port
}
