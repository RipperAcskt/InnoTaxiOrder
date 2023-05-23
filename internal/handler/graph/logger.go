package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const logger key = "logger"

func (r *Resolver) Log(handler http.Handler) http.Handler {
	resp := make(map[string]string)
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log, err := zap.NewProduction(zap.Fields(zap.String("url", r.URL.Path), zap.String("method", r.Method), zap.Any("uuid", uuid.New()), zap.String("request time", time.Now().String())))
		if err != nil {
			rw.WriteHeader(http.StatusForbidden)
			resp["error"] = fmt.Errorf("create logger failed").Error()
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				_, err := rw.Write([]byte(err.Error()))
				if err != nil {
					log.Error("Log", zap.Error(fmt.Errorf("write  failed: %w", err)))
				}
				return
			}
			_, err = rw.Write(jsonResp)
			if err != nil {
				log.Error("Log", zap.Error(fmt.Errorf("write  failed: %w", err)))
			}
			return
		}

		r = r.WithContext(ContextWithLogger(r.Context(), log))

		handler.ServeHTTP(rw, r)

		log.Info("request", zap.String("time", time.Since(start).String()))
	})
}

func ContextWithLogger(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, logger, log)
}

func LoggerFromContext(ctx context.Context) (*zap.Logger, bool) {
	log, ok := ctx.Value(logger).(*zap.Logger)
	return log, ok
}
