package jasmine

import (
	"github.com/emersion/go-webdav/caldav"
	"net/http"
	"time"
)

type Conf struct {
	Listen string `json:"listen"`
}

func NewHttpServer(conf Conf, wd string) *http.Server {
	h := &caldav.Handler{
		Backend: &Backend{},
	}

	return &http.Server{
		Addr:              conf.Listen,
		Handler:           h,
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      time.Minute,
		IdleTimeout:       time.Minute,
		MaxHeaderBytes:    2500,
	}
}
