package helper

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/pkg/file_template"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
	"github.com/zbitech/controller/pkg/object"
	"os"
)

type TemplateConfig struct {
	settingsConfig   *object.Settings
	policyConfig     *object.PolicyConfig
	instanceConfig   map[model.InstanceType]object.InstanceConfig
	instanceTemplate map[model.InstanceType]file_template.FileTemplate
	appTemplate      *file_template.FileTemplate
	projectTemplate  *file_template.FileTemplate
}

var CONFIG_DIR = GetEnv("ZBI_CONFIG_DIRECTORY", "/etc/zbi/conf/")
var TEMPLATE_DIR = GetEnv("ZBI_TEMPLATE_DIRECTORY", "/etc/zbi/templates/")
var CONFIG_FILE = "zbi-conf.json"

var Config *TemplateConfig = newTemplateConfig()

var (
	NAMESPACE       = "NAMESPACE"
	ZCASH_CONF      = "ZCASH_CONF"
	ENVOY_CONF      = "ENVOY_CONF"
	LWD_CONF        = "LWD_CONF"
	CREDENTIALS     = "CREDENTIALS"
	DEPLOYMENT      = "DEPLOYMENT"
	SERVICE         = "SERVICE"
	INGRESS         = "INGRESS"
	INGRESS_INCLUDE = "INGRESS_INCLUDE"
	VOLUME          = "VOLUME"
	SNAPSHOT        = "SNAPSHOT"
	VOLUME_SNAPSHOT = "VOLUME_SNAPSHOT"
	INSTANCE_LIST   = "INSTANCE_LIST"
)

func newTemplateConfig() *TemplateConfig {
	return &TemplateConfig{
		instanceConfig:   map[model.InstanceType]object.InstanceConfig{},
		instanceTemplate: map[model.InstanceType]file_template.FileTemplate{},
	}
}

func (tc *TemplateConfig) SetConfig(app object.AppConfig) {
	tc.settingsConfig = &app.Settings
	tc.policyConfig = &app.Policy
	for _, ic := range app.Instances {
		tc.instanceConfig[ic.InstanceType] = ic
	}
}

func (tc *TemplateConfig) LoadConfig(ctx context.Context) {
	content, err := ReadFile(CONFIG_DIR + CONFIG_FILE)
	if err != nil {
		panic(err)
	}

	log := logger.GetLogger(ctx)
	log.WithFields(logrus.Fields{"content": string(content)}).Infof("setting config")
	var appConfig object.AppConfig
	if err = json.Unmarshal(content, &appConfig); err != nil {
		panic(err)
	}

	tc.SetConfig(appConfig)

	settings, _ := json.Marshal(tc)
	log.WithFields(logrus.Fields{"content": string(settings)}).Infof("config settings")
}

func (tc *TemplateConfig) LoadTemplates(ctx context.Context) {

	log := logger.GetLogger(ctx)

	log.Infof("creating app template from %s", TEMPLATE_DIR+tc.settingsConfig.Templates["app"])
	tc.appTemplate = file_template.CreateFilePathTemplate("app", TEMPLATE_DIR+tc.settingsConfig.Templates["app"], file_template.NO_FUNCS)

	log.Infof("creating project template from %s", TEMPLATE_DIR+tc.settingsConfig.Templates["project"])
	tc.projectTemplate = file_template.CreateFilePathTemplate("project", TEMPLATE_DIR+tc.settingsConfig.Templates["project"], file_template.NO_FUNCS)

	for instanceType := range tc.instanceConfig {
		name := string(instanceType)
		path := TEMPLATE_DIR + tc.settingsConfig.Templates[name]

		log.Infof("creating instance %s template from %s", instanceType, path)
		tc.instanceTemplate[instanceType] = *file_template.CreateFilePathTemplate(name, path, file_template.FUNCTIONS)
	}

}

func (tc *TemplateConfig) GetSettings() *object.Settings {
	return tc.settingsConfig
}

func (tc *TemplateConfig) GetPolicyConfig() *object.PolicyConfig {
	return tc.policyConfig
}

func (tc *TemplateConfig) GetInstanceConfig(iType model.InstanceType) (*object.InstanceConfig, error) {
	ic, ok := tc.instanceConfig[iType]
	if !ok {
		return nil, errors.New("instance config not found for ")
	}

	return &ic, nil
}

func (tc *TemplateConfig) GetAppTemplate() *file_template.FileTemplate {
	return tc.appTemplate
}

func (tc *TemplateConfig) GetProjectTemplate() *file_template.FileTemplate {
	return tc.projectTemplate
}

func (tc *TemplateConfig) GetInstanceTemplate(iType model.InstanceType) (*file_template.FileTemplate, error) {
	it, ok := tc.instanceTemplate[iType]
	if !ok {
		return nil, errors.New("instance config not found for ")
	}

	return &it, nil
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func ReadFile(path string) ([]byte, error) {
	//name := strings.Split(filepath.Base(path), ".")[0]
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}
