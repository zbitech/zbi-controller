package manager

//import (
//	"encoding/json"
//	"errors"
//	"github.com/zbitech/controller/internal/utils"
//	"github.com/zbitech/controller/pkg/file_template"
//	"github.com/zbitech/controller/pkg/model"
//	"github.com/zbitech/controller/pkg/object"
//)
//
//type TemplateConfig struct {
//	settingsConfig   *object.Settings
//	policyConfig     *object.PolicyConfig
//	instanceConfig   map[model.InstanceType]object.InstanceConfig
//	instanceTemplate map[model.InstanceType]file_template.FileTemplate
//	appTemplate      *file_template.FileTemplate
//	projectTemplate  *file_template.FileTemplate
//}
//
//var ASSET_PATH_DIR = utils.GetEnv("ASSET_PATH_DIRECTORY", "etc/zbi/")
//var CONFIG_DIR = utils.GetEnv("ZBI_CONFIG_DIRECTORY", "/etc/zbi/conf")
//var TEMPLATE_DIR = utils.GetEnv("ZBI_TEMPLATE_DIRECTORY", "etc/zbi/templates")
//var CONFIG_FILE = "zbi-conf.json"
//
//var Config *TemplateConfig = newTemplateConfig()
//
//var (
//	NAMESPACE       = "NAMESPACE"
//	SERVICE         = "SERVICE"
//	INGRESS         = "INGRESS"
//	INGRESS_INCLUDE = "INGRESS_INCLUDE"
//)
//
//func newTemplateConfig() *TemplateConfig {
//	return &TemplateConfig{
//		instanceConfig:   map[model.InstanceType]object.InstanceConfig{},
//		instanceTemplate: map[model.InstanceType]file_template.FileTemplate{},
//	}
//}
//
//func (tc *TemplateConfig) SetConfig(app object.AppConfig) {
//	tc.settingsConfig = &app.Settings
//	tc.policyConfig = &app.Policy
//	for _, ic := range app.Instances {
//		tc.instanceConfig[ic.InstanceType] = ic
//	}
//}
//
//func (tc *TemplateConfig) LoadConfig() {
//	content, err := utils.ReadFile(CONFIG_DIR + CONFIG_FILE)
//	if err != nil {
//		panic(err)
//	}
//
//	//fmt.Printf("Data2 => %s", string(content))
//	var appConfig object.AppConfig
//	if err = json.Unmarshal(content, &appConfig); err != nil {
//		panic(err)
//	}
//
//	tc.SetConfig(appConfig)
//}
//
//func (tc *TemplateConfig) LoadTemplates() {
//
//	tc.appTemplate = file_template.CreateFilePathTemplate("app", TEMPLATE_DIR+tc.settingsConfig.Templates["app"], file_template.NO_FUNCS)
//	tc.projectTemplate = file_template.CreateFilePathTemplate("project", TEMPLATE_DIR+tc.settingsConfig.Templates["project"], file_template.NO_FUNCS)
//
//	for instanceType := range tc.instanceConfig {
//		name := string(instanceType)
//		path := TEMPLATE_DIR + tc.settingsConfig.Templates[name]
//
//		tc.instanceTemplate[instanceType] = *file_template.CreateFilePathTemplate(name, path, file_template.FUNCTIONS)
//	}
//
//}
//
//func (tc *TemplateConfig) GetPolicyConfig() *object.PolicyConfig {
//	return tc.policyConfig
//}
//
//func (tc *TemplateConfig) GetInstanceConfig(iType model.InstanceType) (*object.InstanceConfig, error) {
//	ic, ok := tc.instanceConfig[iType]
//	if !ok {
//		return nil, errors.New("instance config not found for ")
//	}
//
//	return &ic, nil
//}
//
//func (tc *TemplateConfig) GetAppTemplate() *file_template.FileTemplate {
//	return tc.appTemplate
//}
//
//func (tc *TemplateConfig) GetProjectTemplate() *file_template.FileTemplate {
//	return tc.projectTemplate
//}
//
//func (tc *TemplateConfig) GetInstanceTemplate(iType model.InstanceType) (*file_template.FileTemplate, error) {
//	it, ok := tc.instanceTemplate[iType]
//	if !ok {
//		return nil, errors.New("instance config not found for ")
//	}
//
//	return &it, nil
//}
