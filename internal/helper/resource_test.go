package helper

import (
	"github.com/stretchr/testify/assert"
	fake_zbi "github.com/zbitech/controller/fake-zbi"
	"testing"
)

func TestConfig_newTemplateConfig(t *testing.T) {

	config := newTemplateConfig()
	assert.NotNil(t, config)
	assert.NotNil(t, config.instanceConfig)
	assert.NotNil(t, config.instanceTemplate)
}

func TestConfig_LoadConfig(t *testing.T) {

	ctx := fake_zbi.InitContext()
	config := newTemplateConfig()
	assert.NotNil(t, config)

	config.LoadConfig(ctx)
	assert.NotNil(t, config.settingsConfig)
	assert.NotNil(t, config.policyConfig)
	assert.Len(t, config.instanceConfig, 2)
}

func TestConfig_SetConfig(t *testing.T) {

}

func TestConfig_LoadTemplates(t *testing.T) {
	ctx := fake_zbi.InitContext()
	config := newTemplateConfig()
	assert.NotNil(t, config)

	config.LoadConfig(ctx)
	config.LoadTemplates(ctx)
	assert.NotNil(t, config.appTemplate)
	assert.NotNil(t, config.projectTemplate)
	assert.Len(t, config.instanceTemplate, 2)
}

func TestConfig_GetPolicyConfig(t *testing.T) {

}

func TestConfig_GetInstanceConfig(t *testing.T) {

}

func TestConfig_GetAppTemplate(t *testing.T) {

}

func TestConfig_GetProjectTemplate(t *testing.T) {

}

func TestConfig_GetInstanceTemplate(t *testing.T) {

}
