package fake_zbi

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/rctx"
)

func InitContext() context.Context {
	logger.Init()
	return context.WithValue(context.Background(), rctx.LOGGER, logrus.WithFields(logrus.Fields{}))
	// return context.Background()
}
