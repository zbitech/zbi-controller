package vars

import (
	"context"
	"github.com/zbitech/controller/internal/utils"
	"github.com/zbitech/controller/pkg/interfaces"
	"strconv"
)

//const (
//	REQUEST      = "request"
//	SERVICE      = "service"
//	ELAPSED_TIME = "etime"
//	USERID       = "userid"
//	ROLE         = "role"
//	TXID         = "txid"
//	LOGGER       = "logger"
//)

var (
	CTX                      context.Context
	DATABASE_FACTORY         = utils.GetEnv("DATABASE_FACTORY", "mongo")
	DATABASE_URL             = utils.GetEnv("DATABASE_URL", "mongodb://root:wzZxK2YCCn@db.zbitech.local:27017")
	DATABASE_NAME            = utils.GetEnv("DATABASE_NAME", "zbiRepo")
	ASSET_PATH_DIRECTORY     = utils.GetEnv("ASSET_PATH_DIRECTORY", "tests/files/etc/zbi")
	KUBECONFIG               = utils.GetEnv("KUBECONFIG", "cfg/kubeconfig")
	ZBI_LOG_LEVEL            = utils.GetIntEnv("ZBI_LOG_LEVEL", 0)
	USE_KUBERNETES_CONFIG, _ = strconv.ParseBool(utils.GetEnv("USEKUBERNETESCONFIG", "false"))
	CONTROLLER_METRICS, _    = strconv.ParseBool(utils.GetEnv("METRICS", "false"))
	K8S_INCLUSTER            = true
	EXPIRATION_HOURS, _      = strconv.Atoi(utils.GetEnv("EXPIRATION_HOURS", "8760"))
	ADMIN_USER               = utils.GetEnv("ZBI_ADMIN_USER", "admin")
	ADMIN_EMAIL              = utils.GetEnv("ZBI_ADMIN_EMAIL", "admin@alphegasolutions.com")
	ADMIN_PASSWORD           = utils.GetEnv("ZBI_ADMIN_PASSWORD", "password")
	ZBI_NAMESPACE            = utils.GetEnv("ZBI_NAMESPACE", "")
	HOURS_IN_YEAR            = 8760

	KlientFactory     interfaces.KlientFactoryIF
	ManagerFactory    interfaces.ResourceManagerFactoryIF
	RepositoryFactory interfaces.RepositoryServiceFactoryIF

	//	ServiceFactory    interfaces.ServiceFactoryIF
	//	AuthorizerFactory interfaces.AccessAuthorizerFactoryIF
)
