package graph

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/RipperAcskt/innotaxiorder/internal/service"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	s   *service.Service
	Cfg *config.Config
	log *zap.Logger
}

func New(s *service.Service, cfg *config.Config, log *zap.Logger) (*Resolver, error) {
	return &Resolver{s, cfg, log}, nil
}

func (r *Resolver) InitRouters() *http.ServeMux {
	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: r}))

	mux := http.DefaultServeMux

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	return mux
}
