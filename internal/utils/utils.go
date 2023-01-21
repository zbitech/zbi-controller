package utils

import (
	b64 "encoding/base64"
	"encoding/json"
	"github.com/sethvargo/go-password/password"
	"github.com/zbitech/controller/pkg/model"
	"os"
	"strconv"
	"strings"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func GetIntEnv(key string, fallback int) int {
	val := GetEnv(key, strconv.Itoa(fallback))
	y, e := strconv.Atoi(val)
	if e != nil {
		return fallback
	}
	return y
}

func Base64EncodeString(value string) string {
	return string([]byte(b64.StdEncoding.EncodeToString([]byte(value))))
}

func Base64DecodeString(value string) string {
	b, err := b64.StdEncoding.DecodeString(value)
	if err != nil {
		return ""
	}

	return string(b)
}

func MarshalObject(obj interface{}) string {
	if obj != nil {
		c, err := json.Marshal(obj)
		if err != nil {
			return err.Error()
		}

		return string(c)
	}
	return ""
}

func UnMarshalObject(data string, obj interface{}) error {
	if err := json.Unmarshal([]byte(data), obj); err != nil {
		return err
	}

	return nil
}

func ReadFile(path string) ([]byte, error) {
	//name := strings.Split(filepath.Base(path), ".")[0]
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func MarshalIndentObject(obj interface{}) string {
	if obj != nil {
		c, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return err.Error()
		}

		return string(c)
	}
	return ""
}

func GenerateRandomString(length int, noupper bool) string {
	text, err := password.Generate(length, 1, 0, noupper, false)
	if err != nil {
		return "password"
	}

	return text
}

func GenerateSecurePassword() string {
	passwd, err := password.Generate(12, 4, 0, false, false)
	if err != nil {
		return "password"
	}

	return passwd
}

func MergeMap(map1, map2 map[string]interface{}) map[string]interface{} {
	for k, v := range map2 {
		map1[k] = v
	}
	return map1
}

func ResourceObjectType(value string) model.ResourceObjectType {
	if strings.ToLower(value) == strings.ToLower(string(model.ResourceNamespace)) {
		return model.ResourceNamespace
	} else if strings.ToLower(value) == strings.ToLower(string(model.ResourceConfigMap)) {
		return model.ResourceConfigMap
	} else if strings.ToLower(value) == strings.ToLower(string(model.ResourceSecret)) {
		return model.ResourceSecret
	} else if strings.ToLower(value) == strings.ToLower(string(model.ResourcePersistentVolumeClaim)) {
		return model.ResourcePersistentVolumeClaim
	} else if strings.ToLower(value) == strings.ToLower(string(model.ResourceDeployment)) {
		return model.ResourceDeployment
	} else if strings.ToLower(value) == strings.ToLower(string(model.ResourceService)) {
		return model.ResourceService
	} else if strings.ToLower(value) == strings.ToLower(string(model.ResourceHTTPProxy)) {
		return model.ResourceHTTPProxy
	}

	return ""
}
