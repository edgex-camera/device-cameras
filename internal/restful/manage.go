package restful

import (
	"net/http"

	"github.com/gorilla/mux"
)

func appendDeviceManageRoute(r *mux.Router, h *handler) {
	prefix := "/device"
	subRouter := r.PathPrefix(prefix).Subrouter()

	subRouter.Path("/").HandlerFunc(h.getAllCameraDevice).Methods(http.MethodGet)
	subRouter.Path("/devices").HandlerFunc(h.deleteRemoveAllCamera).Methods(http.MethodDelete)
	subRouter.Path("/{camera_name}").HandlerFunc(h.getCameraInfo).Methods(http.MethodGet)
	subRouter.Path("/{camera_name}").HandlerFunc(h.postAddCamera).Methods(http.MethodPost)
	subRouter.Path("/{camera_name}").HandlerFunc(h.deleteRemoveCamera).Methods(http.MethodDelete)

}

//获取所有的的camera设备
func (h *handler) getAllCameraDevice(w http.ResponseWriter, r *http.Request) {

}

//获取camera信息
func (h *handler) getCameraInfo(w http.ResponseWriter, r *http.Request) {

}

//添加camera
func (h *handler) postAddCamera(w http.ResponseWriter, r *http.Request) {}

//删除camera
func (h *handler) deleteRemoveCamera(w http.ResponseWriter, r *http.Request) {}

//删除所有摄像头
func (h *handler) deleteRemoveAllCamera(w http.ResponseWriter, r *http.Request) {}
