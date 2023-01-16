package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/rctx"
	"net/http"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...mux.MiddlewareFunc) http.Handler {

	if len(middlewares) == 0 {
		return f
	}

	return middlewares[0](Chain(f, middlewares[1:cap(middlewares)]...))

	//for _, m := range middlewares {
	//	f = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	//		m(f)
	//	})
	//}
	//
	//return f
}

func InitRequest(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger(r.Context())
		userid := uuid.New().String()
		role := "admin"
		contextLogger := log.WithFields(logrus.Fields{
			rctx.TXID:   uuid.New().String(),
			rctx.USERID: userid,
			rctx.IP:     r.RemoteAddr,
			rctx.XIP:    r.Header.Get("X-Forwarded-For"),
		})

		ctx := context.WithValue(r.Context(), rctx.LOGGER, contextLogger)
		ctx = context.WithValue(ctx, rctx.USERID, userid)
		ctx = context.WithValue(ctx, rctx.ROLE, role)
		f.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Logging(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger(r.Context())
		start := time.Now()
		defer func() {
			log.WithFields(logrus.Fields{rctx.ELAPSED_TIME: time.Since(start).Round(time.Millisecond).String(), rctx.REQUEST: r.URL.Path}).Info()
		}()

		f.ServeHTTP(w, r)
	})
}
