package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/internal/utils"
	"github.com/zbitech/controller/pkg/file_template"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	DATA_PATH = utils.GetEnv("DATA_PATH", "./fake-zbi/data/")

	projectTemplate *file_template.FileTemplate = file_template.CreateFilePathTemplate("project", DATA_PATH+"manifest/project/project_templates_v1.tmpl", file_template.NO_FUNCS)
	zcashTemplate   *file_template.FileTemplate = file_template.CreateFilePathTemplate("zcash", DATA_PATH+"manifest/zcash/zcash_templates_v1.tmpl", file_template.FUNCTIONS)
	lwdTemplate     *file_template.FileTemplate = file_template.CreateFilePathTemplate("lwd", DATA_PATH+"manifest/lwd/lwd_templates_v1.tmpl", file_template.FUNCTIONS)

	ingress = `
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  annotations:
    external-dns.alpha.kubernetes.io/hostname: apps.zbitech.local
  labels:
    app.kubernetes.io/instance: controller
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/name: zbi
    app.kubernetes.io/version: 1.16.0
    helm.sh/chart: zbi-0.1.0
  name: zbi-proxy
  namespace: zbi
spec:
  routes:
  - conditions:
    - prefix: /
    pathRewritePolicy:
      replacePrefix:
      - replacement: /status
    services:
    - name: controller-zbi-svc
      port: 8080
  virtualhost:
    fqdn: api.zbitech.local
    tls:
      secretName: zbi-api-tls
`
	secret = `
apiVersion: v1
kind: Secret
metadata:
  name: zbi-api-tls
  namespace: zbi
type: kubernetes.io/tls
data:
  ca.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURBVENDQWVtZ0F3SUJBZ0lSQUl2L2hOSndmNGNVR2Qrc1FmT2NGS1V3RFFZSktvWklodmNOQVFFTEJRQXcKSERFYU1CZ0dBMVVFQXhNUllYQnBMbnBpYVhSbFkyZ3ViRzlqWVd3d0hoY05Nakl4TVRFd01UWTFNakV5V2hjTgpNak14TVRFd01UWTFNakV5V2pBY01Sb3dHQVlEVlFRREV4RmhjR2t1ZW1KcGRHVmphQzVzYjJOaGJEQ0NBU0l3CkRRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFPbXhpTFdpMVcrZUtTd2JVeDc2RzM4OHVQcncKa3EzTXhOZDhQa0g0cVlIRmI1TmxhSTQrcWFpL1RhbjFVZzhIN1Y2ZU9RZUVCVFlwVWg3SlQ0QlR3UTgweUdudgo1eGhUSmQxcjRCN2hRSkJ6OGk2MjJTWjIxLzdyZjl3QTRYK2czT1Zpd0Y3QlRpR3gzZCtKNndzR1FKYzdubTdhCldvVFc1cURHc2sxRnArRUMxMVR4cnduNVU2RE83eW9JcEVXRU43WWM5QUJnYjJvU1RsZnJUaVBNanhvakRWSFUKNnA1dXdEckY2TEZPeVV4ZkVHUzk4Q2Qzd3MwSFVvVFdUSDJTUkJuQWFCaDRiT2w4TmFkMytpajMxTkJ6eHVHagpPTEpOOFY3VC81ak5kdTZzZFM1c3BMZGhKRVlIYTIwWm95VTJOY3pTWjhNaXVUc1BtU05vcEh3R3FHa0NBd0VBCkFhTStNRHd3RGdZRFZSMFBBUUgvQkFRREFnV2dNQXdHQTFVZEV3RUIvd1FDTUFBd0hBWURWUjBSQkJVd0U0SVIKWVhCcExucGlhWFJsWTJndWJHOWpZV3d3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUp5dXU4dlF3bURJRnR6VgpXYWo4RFRESFJUSkhEVUhvRzY0akVkRW45eElqRUV0c0ZnelljNmFSQUpQeFlNa05KOE1NTVFhbml4a2VlV20vCngyb1FjRkY5d3N3c09BNW9qemVLQ1Jia1c2YjAyNTgwdGp3RUJldTJFRmY3Z3FBUjZ6dUJuMzY3UEp3VDQzWnMKTEdnemNRZFltWjZUWmY4eWd1RU5ST20vODB6NGgzQm1OU0xnUDYzNTJTVVJxL05JZDRpN3ZLaExZQmdSRDNaZQpDS1Iyd2dDYklqSDliZHdLcndmNnc5L0FRdEVqalZJTGtJanhxUmRnK2RMMUpGQVFJdDR2ckt3ejFPWlh2MTZCCkJVNUVRZ096UnRjTXkwWHJjVnh3N0NSVU5BRERGd1h0ZmVUd3Z5L01hcTF4emFSWDAwSjZjRjZDcWY3emY4UHkKSDk2K0ZrOD0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURBVENDQWVtZ0F3SUJBZ0lSQUl2L2hOSndmNGNVR2Qrc1FmT2NGS1V3RFFZSktvWklodmNOQVFFTEJRQXcKSERFYU1CZ0dBMVVFQXhNUllYQnBMbnBpYVhSbFkyZ3ViRzlqWVd3d0hoY05Nakl4TVRFd01UWTFNakV5V2hjTgpNak14TVRFd01UWTFNakV5V2pBY01Sb3dHQVlEVlFRREV4RmhjR2t1ZW1KcGRHVmphQzVzYjJOaGJEQ0NBU0l3CkRRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFPbXhpTFdpMVcrZUtTd2JVeDc2RzM4OHVQcncKa3EzTXhOZDhQa0g0cVlIRmI1TmxhSTQrcWFpL1RhbjFVZzhIN1Y2ZU9RZUVCVFlwVWg3SlQ0QlR3UTgweUdudgo1eGhUSmQxcjRCN2hRSkJ6OGk2MjJTWjIxLzdyZjl3QTRYK2czT1Zpd0Y3QlRpR3gzZCtKNndzR1FKYzdubTdhCldvVFc1cURHc2sxRnArRUMxMVR4cnduNVU2RE83eW9JcEVXRU43WWM5QUJnYjJvU1RsZnJUaVBNanhvakRWSFUKNnA1dXdEckY2TEZPeVV4ZkVHUzk4Q2Qzd3MwSFVvVFdUSDJTUkJuQWFCaDRiT2w4TmFkMytpajMxTkJ6eHVHagpPTEpOOFY3VC81ak5kdTZzZFM1c3BMZGhKRVlIYTIwWm95VTJOY3pTWjhNaXVUc1BtU05vcEh3R3FHa0NBd0VBCkFhTStNRHd3RGdZRFZSMFBBUUgvQkFRREFnV2dNQXdHQTFVZEV3RUIvd1FDTUFBd0hBWURWUjBSQkJVd0U0SVIKWVhCcExucGlhWFJsWTJndWJHOWpZV3d3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUp5dXU4dlF3bURJRnR6VgpXYWo4RFRESFJUSkhEVUhvRzY0akVkRW45eElqRUV0c0ZnelljNmFSQUpQeFlNa05KOE1NTVFhbml4a2VlV20vCngyb1FjRkY5d3N3c09BNW9qemVLQ1Jia1c2YjAyNTgwdGp3RUJldTJFRmY3Z3FBUjZ6dUJuMzY3UEp3VDQzWnMKTEdnemNRZFltWjZUWmY4eWd1RU5ST20vODB6NGgzQm1OU0xnUDYzNTJTVVJxL05JZDRpN3ZLaExZQmdSRDNaZQpDS1Iyd2dDYklqSDliZHdLcndmNnc5L0FRdEVqalZJTGtJanhxUmRnK2RMMUpGQVFJdDR2ckt3ejFPWlh2MTZCCkJVNUVRZ096UnRjTXkwWHJjVnh3N0NSVU5BRERGd1h0ZmVUd3Z5L01hcTF4emFSWDAwSjZjRjZDcWY3emY4UHkKSDk2K0ZrOD0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  tls.key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBNmJHSXRhTFZiNTRwTEJ0VEh2b2Jmenk0K3ZDU3JjekUxM3crUWZpcGdjVnZrMlZvCmpqNnBxTDlOcWZWU0R3ZnRYcDQ1QjRRRk5pbFNIc2xQZ0ZQQkR6VElhZS9uR0ZNbDNXdmdIdUZBa0hQeUxyYloKSm5iWC91dC8zQURoZjZEYzVXTEFYc0ZPSWJIZDM0bnJDd1pBbHp1ZWJ0cGFoTmJtb01heVRVV240UUxYVlBHdgpDZmxUb003dktnaWtSWVEzdGh6MEFHQnZhaEpPVit0T0k4eVBHaU1OVWRUcW5tN0FPc1hvc1U3SlRGOFFaTDN3CkozZkN6UWRTaE5aTWZaSkVHY0JvR0hoczZYdzFwM2Y2S1BmVTBIUEc0YU00c2szeFh0UC9tTTEyN3F4MUxteWsKdDJFa1JnZHJiUm1qSlRZMXpOSm53eUs1T3crWkkyaWtmQWFvYVFJREFRQUJBb0lCQUVmdTh1TGVMWTYvQTNObApNYy9PTXRxV2lXWU0yVW1RUjJNQkJuVHVJdGNrTy9VRitRb1g5Y2RRbzRwV1RoejhWcStTU29HcXZLUHdVaXZSCjBadnhxL0tQVDhWMEtCRlB2czhLWHFHQ3VvbjhkcWEwZCtFa0lkYUJEUWxlYUFzT0xCQ2J0aFUwc1dVangrVUEKSWc1eHJUNGdCdU9lYU5DTkNjNmhlczdZU3hXeVlJTEtHS29NRTRqLzl6M3Y3N2dBaEEzYUVtRmNDU2psZEdQUApLVGppTUYrRjN2aFg2RUZZcGN4czdNeVp2dTkzYkZTWXdvWmpsK0ZRNmJpVkZwRGpvK0lSU1JyRU9JVk1udGtrCkFWM01Pc3p1UmJBNks2ajdSS3BNaldlcWlVcnluSnlXUlJ2SjRvUE5OR0d0NEhTaUtsaGN5bGZMM2l5WHpLdnoKTzZjZVRjVUNnWUVBL1k2WDFhNjdVSGtmTC85TWUvNEdMdXE4L2Z0Y1dUVnhXM2gwUFRGMjQ4YnpKWFBwcXp6UQprejQ0TUl0VTRrcTcvQUNWcm9sNzA0bHpOSytXczlNWndabFFSM3VCTVZLZVB6SkZDdGpxRWRWdlM2Q2R6YVZ1CnZXakJOZVk4KzFtdW1XNW1hNDNQU1daZm1LTVRTZVBJbUNDQnI5L0VIYVgrSzI4Z1ZRKzF5MzhDZ1lFQTYvSHkKWjcwdEVFQVluYUM5T2ttWWptUjdHanRtcUdrYmtoNkhUeTNPL0EwWWpLalRGeW1xK25QL1hzT3U5Ykx3dTlhcwpPSEN3OGRRY09BRzBhSFM4Mlh2NENaNG5qYk9CMFMxTVNKREtyNm42RGcrWU1kMjl4TW1oSTFBcVEzSnN0cEN3CjlpVUFXRTdRRmxDNmNsbXlXL0dXN0tXTmd0NWF2OU5zN1RWK29CY0NnWUVBdnNldXFPSWJJSmF5QjZ4QlFUNUcKS3NFR3lOZDdpY2Z6Ymc2NDcxNHJoWUVwYS9IR1RNaXFhMCt5ZVp3c2wwUUNJNy9RNEEya05PdEQyczJQUitpNwpoWGEwOThRTzFpekMwdXdoRk9OWFkybkRueFRRQjI3Rlh4RFY1NWRBSlNNNmcwbVZHTElQMkx1RmpGU1BhOVpQCkZWL0lGS3Y2WlJDRHFPeXBXRGRFNDBNQ2dZQUV3bXFyWVF1SnFtRlV2S3RVZzZ1S0k4aS83TGJUYXR0ZGhUUWgKaXNFRUlKZUFMdCtqTmZuMjkyUU5XMUVxTDZQZmhpTVBPR2E1V2hmL29Ua3NhajVzL0sweU5IaUR5VDB6SlFERwo3ZlRJdWxzSzhaR1dYK3kyRFlNc25TOWRFTy9VZHZLNjVHQXZaOWVXdTZZbkxGd0dzc2JpbXl4Um1YNm5JL0tzCmprbXhuUUtCZ1FEbURNaElvQ1gvU09wTzJUQjFEazNGNXdjN0dGdTl0QlpLc3dMYzdjRVg3czgvRWUwZlJkODAKSnZHV05YUDFTTDBZSkZWQzNCdHJ2UHZDRG1ZWWNGY2dqOE5ORC9iVEVZYTQ4ekR5bVc1UFFyMEViek8ybkpFRQp6U0pteFplUzBIcHY0Q0hxLzg4ckxHN2t3MERJQThOYWhiOGRqaS84VVZncjVpaVYwUU9GZlE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
`
)

func getArrayData(filePath string, obj []unstructured.Unstructured) error {
	ctx := context.Background()
	log := logger.GetLogger(ctx)
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	log.Debugf("File Data: %s", string(file))
	if err = json.Unmarshal(file, &obj); err != nil {
		return err
	}
	log.Debugf("READ Data: %s", utils.MarshalObject(obj))
	return nil
}

func getData(filePath string, obj interface{}) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	//	logger.Debugf(rctx.CTX, "File Data: %s", string(file))
	if err = json.Unmarshal(file, &obj); err != nil {
		return err
	}
	//	logger.Debugf(rctx.CTX, "READ Data: %s", utils.MarshalObject(obj))
	return nil
}

func GetInstanceResources(iType model.InstanceType, rType model.ResourceObjectType, size int) ([]unstructured.Unstructured, error) {

	var filePath string
	switch rType {
	case model.ResourcePersistentVolumeClaim:
		filePath = fmt.Sprintf("%s/data/manifest/%s/pvc.json", DATA_PATH, string(iType)) //"data/manifest/zcash/pvc.json"
	}

	var objects = make([]unstructured.Unstructured, size)

	if err := getArrayData(filePath, objects); err != nil {
		return nil, err
	}

	return objects, nil
}

func GetResource(rType model.ResourceObjectType) ([]unstructured.Unstructured, error) {

	var name []string
	switch rType {
	case model.ResourceNamespace:
		name = []string{helper.NAMESPACE}
	case model.ResourceDeployment:
		name = []string{helper.DEPLOYMENT}
	case model.ResourceConfigMap:
		name = []string{helper.ZCASH_CONF}
	case model.ResourceSecret:
		name = []string{helper.CREDENTIALS}
	case model.ResourceService:
		name = []string{helper.SERVICE}
	case model.ResourcePersistentVolumeClaim:
		name = []string{helper.VOLUME}
	}

	data, err := zcashTemplate.ExecuteTemplates(name, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return helper.CreateYAMLObjects(data)
}

func GetZcashResources() ([]unstructured.Unstructured, error) {
	var names = []string{helper.NAMESPACE, helper.ENVOY_CONF, helper.ZCASH_CONF, helper.CREDENTIALS, helper.VOLUME, helper.DEPLOYMENT, helper.SERVICE}
	data, err := zcashTemplate.ExecuteTemplates(names, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return helper.CreateYAMLObjects(data)
}

func GetLWDResources() ([]unstructured.Unstructured, error) {
	var names = []string{helper.NAMESPACE, helper.LWD_CONF, helper.ZCASH_CONF, helper.CREDENTIALS, helper.VOLUME, helper.DEPLOYMENT, helper.SERVICE}
	data, err := projectTemplate.ExecuteTemplates(names, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return helper.CreateYAMLObjects(data)
}
