package restful

import (
	"net/http"

	"github.com/gorilla/mux"
)

func appendDeviceManageRoute(r *mux.Route) {
	// prefix := "device"
	// subRouter := r.PathPrefix(prefix).Subrouter()

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
