package monitoring

import (
	"context"
	"net/http"

	"github.com/codeinuit/shibesbot/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Monitoring interface {
	Run()
	Stop()
}

type HttpMonitor struct {
	log logger.Logger
	srv *http.Server
}

func NewHTTPMonitorServer(l logger.Logger) *HttpMonitor {
	http.Handle("/metrics", promhttp.Handler())
	return &HttpMonitor{
		log: l,
		srv: &http.Server{Addr: ":8080"},
	}
}

func (hm *HttpMonitor) Run() {
	err := hm.srv.ListenAndServe()
	if err != http.ErrServerClosed {
		hm.log.Error(err.Error())
	}
}

func (hm *HttpMonitor) Stop() {
	if hm.srv != nil {
		if err := hm.srv.Shutdown(context.Background()); err != nil {
			hm.log.Error(err.Error())
		}
	}
}
