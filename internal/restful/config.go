package restful

import (
	"net/http"

	"github.com/gorilla/mux"
)

func appendConfigRoute(r *mux.Router, h *handler) {
	prefix := "/config"
	subRouter := r.PathPrefix(prefix).Subrouter()

	subRouter.Path("/{camera_name}").HandlerFunc(h.getCameraConfig).Methods(http.MethodGet)
	subRouter.Path("/{camera_name}").HandlerFunc(h.postModifyCameraConfig).Methods(http.MethodPost)

}

//获取camera配置
func (h *handler) getCameraConfig(w http.ResponseWriter, r *http.Request) {

}

//修改camera配置
func (h *handler) postModifyCameraConfig(w http.ResponseWriter, r *http.Request) {

}
