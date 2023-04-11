package handler

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/RipperAcskt/innotaxiorder/internal/handler/graph"
)

type Handler struct {
	Cfg *config.Config
	log *zap.Logger
}

func New(cfg *config.Config, log *zap.Logger) (*Handler, error) {
	return &Handler{cfg, log}, nil
}

func (h *Handler) InitRouters() *http.ServeMux {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	mux := http.DefaultServeMux

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	return mux
}
