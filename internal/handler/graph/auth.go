package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/RipperAcskt/innotaxiorder/internal/service"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

func (res *Resolver) VerifyToken(handler http.Handler) http.Handler {
	resp := make(map[string]string)
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log, ok := LoggerFromContext(r.Context())
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			_, err := rw.Write([]byte(fmt.Errorf("get logger failed").Error()))
			if err != nil {
				log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
			}
			return
		}

		token := strings.Split(r.Header.Get("Authorization"), " ")
		if len(token) < 2 {
			rw.WriteHeader(http.StatusUnauthorized)
			resp["error"] = fmt.Errorf("access token required").Error()
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Error("verefy", zap.Error(fmt.Errorf("json marshal failed: %w", err)))

				rw.WriteHeader(http.StatusInternalServerError)
				_, err := rw.Write([]byte(err.Error()))
				if err != nil {
					log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
				}
				return
			}
			_, err = rw.Write(jsonResp)
			if err != nil {
				log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
			}
			return
		}
		accessToken := token[1]

		id, err := service.Verify(accessToken, res.Cfg)
		if err != nil {
			if errors.Is(err, jwt.ValidationError{Errors: jwt.ValidationErrorExpired}) {
				rw.WriteHeader(http.StatusUnauthorized)
				resp["error"] = fmt.Errorf("verify failed: %w", err).Error()
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Error("verefy", zap.Error(fmt.Errorf("json marshal failed: %w", err)))

					rw.WriteHeader(http.StatusInternalServerError)
					_, err := rw.Write([]byte(err.Error()))
					if err != nil {
						log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
					}
					return
				}
				_, err = rw.Write(jsonResp)
				if err != nil {
					log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
				}
				return
			}
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				rw.WriteHeader(http.StatusForbidden)
				resp["error"] = fmt.Errorf("wrong signature").Error()
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Error("verefy", zap.Error(fmt.Errorf("json marshal failed: %w", err)))

					rw.WriteHeader(http.StatusInternalServerError)
					_, err := rw.Write([]byte(err.Error()))
					if err != nil {
						log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
					}
					return
				}
				_, err = rw.Write(jsonResp)
				if err != nil {
					log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
				}
				return
			}
			if errors.Is(err, service.ErrTokenId) {
				rw.WriteHeader(http.StatusForbidden)
				resp["error"] = fmt.Errorf("id failed").Error()
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Error("verefy", zap.Error(fmt.Errorf("json marshal failed: %w", err)))

					rw.WriteHeader(http.StatusInternalServerError)
					_, err := rw.Write([]byte(err.Error()))
					if err != nil {
						log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
					}
					return
				}
				_, err = rw.Write(jsonResp)
				if err != nil {
					log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
				}
				return
			}

			rw.WriteHeader(http.StatusInternalServerError)
			resp["error"] = fmt.Errorf("verify failed: %w", err).Error()
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Error("verefy", zap.Error(fmt.Errorf("json marshal failed: %w", err)))

				rw.WriteHeader(http.StatusInternalServerError)
				_, err := rw.Write([]byte(err.Error()))
				if err != nil {
					log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
				}
				return
			}
			_, err = rw.Write(jsonResp)
			if err != nil {
				log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
			}
			return
		}
		ctx := ContextWithId(r.Context(), id)
		r = r.WithContext(ctx)
		handler.ServeHTTP(rw, r)

	})
}
