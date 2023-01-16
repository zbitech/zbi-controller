package logger

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/pkg/rctx"
	"time"
)

var (
	log    = logrus.New()
	fields logrus.Fields
	txid   = uuid.New().String()
)

type Log struct {
	startTime time.Time
	log       *logrus.Entry
}

func Init() {
	log.SetLevel(logrus.DebugLevel)
	//	log.SetOutput(os.Stdout)
	//log.SetFormatter(&logrus.TextFormatter{
	//	DisableColors:   true,
	//	FullTimestamp:   true,
	//	TimestampFormat: "2006-01-02 15:04:05",
	//})
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:00",
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		PrettyPrint:       false,
	})

	fields = logrus.Fields{
		rctx.TXID:    txid,
		rctx.USERID:  "zbi",
		rctx.SERVICE: "main",
	}
}

func GetLogger(ctx context.Context) *logrus.Entry {
	mylog := ctx.Value(rctx.LOGGER)

	if mylog == nil {
		return log.WithFields(fields)
	}

	return mylog.(*logrus.Entry)
}

func GetServiceLogger(ctx context.Context, service string) *logrus.Entry {
	mylog := ctx.Value(rctx.LOGGER)
	if mylog == nil {
		mylog = log.WithFields(fields)
	}
	stime := time.Now()

	return mylog.(*logrus.Entry).WithFields(logrus.Fields{rctx.SERVICE: service, rctx.START_TIME: stime})
}

func LogServiceTime(log *logrus.Entry) {

	sTime, ok := log.Data[rctx.START_TIME]
	if ok {
		eTime := time.Since(sTime.(time.Time)).Round(time.Millisecond).String()
		log.Data[rctx.ELAPSED_TIME] = eTime
		delete(log.Data, rctx.START_TIME)
	}

	log.Info("service running time")
}
