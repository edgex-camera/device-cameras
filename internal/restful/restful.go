package restful

import (
	"net/http"

	"github.com/edgexfoundry/device-sdk-go"
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
	r := h.Router.PathPrefix(APIv1Prefix).Subrouter()
	appendDeviceManageRoute(r, h)
	appendCameraResourcesRoute(r, h)
	appendOnvifRoute(r, h)
	appendConfigRoute(r, h)
	return h
}

func (h *handler) getTemplateOutput(w http.ResponseWriter, r *http.Request) {
	out := device.DriverConfigs()["TemplateOutput"]
	_, err := w.Write([]byte(out))
	if err != nil {
		h.lc.Error(err.Error())
	}
}
