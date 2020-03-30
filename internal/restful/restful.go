package restful

import (
	"log"
	"net/http"

	"github.com/edgex-camera/device-sdk-go"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/gorilla/mux"
)

type handler struct {
	*mux.Router
	lc logger.LoggingClient
}

const APIv1Prefix = "/api/v1"

func InitRestRoutes(lc logger.LoggingClient) http.Handler {
	h := &handler{
		Router: mux.NewRouter(),
		lc:     lc,
	}
	h.Router.Use(logMiddeware)
	r := h.Router.PathPrefix(APIv1Prefix).Subrouter()
	appendDeviceManageRoute(r, h)
	appendCameraResourcesRoute(r, h)
	appendOnvifRoute(r, h)
	appendConfigRoute(r, h)
	appendStaticRoute(r, h)
	return h
}

func (h *handler) getTemplateOutput(w http.ResponseWriter, r *http.Request) {
	out := device.DriverConfigs()["TemplateOutput"]
	_, err := w.Write([]byte(out))
	if err != nil {
		h.lc.Error(err.Error())
	}
}
func logMiddeware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.String())
		h.ServeHTTP(w, r)
	})
}
