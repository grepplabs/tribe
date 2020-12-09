package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/grepplabs/tribe/api/v1/server/restapi/healthz"
)

func NewHealthzGetReadyHandler() healthz.GetReadyHandler {
	return &healthzGetReadyHandler{}
}

type healthzGetReadyHandler struct {
}

func (h *healthzGetReadyHandler) Handle(_ healthz.GetReadyParams) middleware.Responder {
	return healthz.NewGetHealthyOK()
}

func NewHealthzGetHealthyHandler() healthz.GetHealthyHandler {
	return &healthzGetHealthyHandler{}
}

type healthzGetHealthyHandler struct {
}

func (h *healthzGetHealthyHandler) Handle(_ healthz.GetHealthyParams) middleware.Responder {
	return healthz.NewGetHealthyOK()
}
