package restful

import (
	"io/ioutil"
	"net/http"

	"github.com/edgex-camera/device-cameras/internal/driver"
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
	JDevice := driver.CurrentDriver.JDevices[getCameraName(r)]
	cameraConfigure := JDevice.Camera.GetConfigure()
	h.respSuccess(cameraConfigure, w)
}

//修改camera配置
func (h *handler) postModifyCameraConfig(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.respFailed(err, w)
	}
	JDevice := driver.CurrentDriver.JDevices[getCameraName(r)]
	err = JDevice.Camera.PutConfig(data)
	if err != nil {
		h.respFailed(err, w)
	}
	newConfigure := JDevice.Camera.GetConfigure()
	h.respSuccess(newConfigure, w)
}
