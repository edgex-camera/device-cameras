package restful

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/driver"
)

func appendDeviceManageRoute(r *mux.Router, h *handler) {
	prefix := "/device"
	subRouter := r.PathPrefix(prefix).Subrouter()

	subRouter.Path("/").HandlerFunc(h.getAllCameraDevice).Methods(http.MethodGet)
	subRouter.Path("/devices").HandlerFunc(h.deleteRemoveAllCamera).Methods(http.MethodDelete)
	subRouter.Path("/{camera_name}").HandlerFunc(h.getCameraInfo).Methods(http.MethodGet)
	subRouter.Path("/{camera_name}").HandlerFunc(h.deleteRemoveCamera).Methods(http.MethodDelete)
	subRouter.Path("/{camera_name}/{device_type}").HandlerFunc(h.postAddCamera).Methods(http.MethodPost)

}

//获取所有的的camera设备
func (h *handler) getAllCameraDevice(w http.ResponseWriter, r *http.Request) {
	type device struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}
	type response struct {
		Devices []device `json:"devices"`
	}
	resp := response{Devices: []device{}}
	for _, jDevice := range driver.CurrentDriver.JDevices {
		resp.Devices = append(resp.Devices, device{Name: jDevice.Name, ID: jDevice.Id})
	}
	h.respSuccess(resp, w)
}

//获取camera信息
func (h *handler) getCameraInfo(w http.ResponseWriter, r *http.Request) {
	type CameraInfo struct {
		IsEnable    bool     `json:"isenable"`
		Capturepath string   `json:"capturepath"`
		StreamAddr  string   `json:"streamaddr"`
		Channels    []string `json:"channels"`
	}
	cameraName := getCameraName(r)

	jDevice, ok := driver.CurrentDriver.JDevices[cameraName]
	if !ok {
		h.respFailed(fmt.Errorf("can not find %s", cameraName), w)
		return
	}
	resp := CameraInfo{
		IsEnable:    jDevice.Camera.IsEnabled(),
		Capturepath: jDevice.Camera.GetCapturePath(),
		StreamAddr:  jDevice.Camera.GetStreamAddr(),
		Channels:    jDevice.Camera.ListChannels(),
	}
	h.respSuccess(resp, w)

}

//添加camera
func (h *handler) postAddCamera(w http.ResponseWriter, r *http.Request) {
	err := driver.CurrentDriver.AddJdevice(getCameraName(r), getDeviceType(r))
	if err != nil {
		h.respFailed(err, w)
		return
	}
	h.respSuccess(nil, w)
}

//删除camera
func (h *handler) deleteRemoveCamera(w http.ResponseWriter, r *http.Request) {
	err := driver.CurrentDriver.RemoveJdevice(getCameraName(r))
	if err != nil {
		h.respFailed(err, w)
		return
	}
	h.respSuccess(nil, w)
}

//删除所有摄像头
func (h *handler) deleteRemoveAllCamera(w http.ResponseWriter, r *http.Request) {
	for _, jDevice := range driver.CurrentDriver.JDevices {
		err := driver.CurrentDriver.RemoveJdevice(jDevice.Name)
		if err != nil {
			h.respFailed(err, w)
			return
		}
	}
	h.respSuccess(nil, w)
}
