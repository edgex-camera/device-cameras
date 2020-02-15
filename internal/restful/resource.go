package restful

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/driver"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice"
)

func appendCameraResourcesRoute(r *mux.Router, h *handler) {
	prefix := "/source"
	subRouter := r.PathPrefix(prefix).Subrouter()

	subRouter.Path("/video_paths").HandlerFunc(h.getVideoPaths).Methods(http.MethodGet)
	subRouter.Path("/images").HandlerFunc(h.getImageURls).Methods(http.MethodGet)
	subRouter.Path("/vidoes").HandlerFunc(h.getVideoURLs).Methods(http.MethodGet)
}

//获取视频存储路径
func (h *handler) getVideoPaths(w http.ResponseWriter, r *http.Request) {
	type responce struct {
		VideoPaths []string `json:"videopaths"`
	}
	deviceName := "deviceName"
	h.checkCamerAndDo(
		deviceName,
		driver.CurrentDriver.JDevices[deviceName].Camera,
		w,
		func(c jdevice.Camera) interface{} {
			resp := responce{
				VideoPaths: c.GetVideoPaths(),
			}
			return resp
		})

}

//获取图片链接列表
func (h *handler) getImageURls(w http.ResponseWriter, r *http.Request) {
	type responce struct {
		ImageURLs []string `json:"imageurls"`
	}
	deviceName := "deviceName"
	h.checkCamerAndDo(
		deviceName,
		driver.CurrentDriver.JDevices[deviceName].Camera,
		w,
		func(c jdevice.Camera) interface{} {
			resp := responce{
				ImageURLs: c.GetImagePaths(),
			}
			return resp
		})
}

//获取video链接列表
func (h *handler) getVideoURLs(w http.ResponseWriter, r *http.Request) {
	type responce struct {
		VideoURLs []string `json:"videourls"`
	}
	deviceName := "deviceName"
	h.checkCamerAndDo(
		deviceName,
		driver.CurrentDriver.JDevices[deviceName].Camera,
		w,
		func(c jdevice.Camera) interface{} {
			resp := responce{
				VideoURLs: c.GetVideoPaths(),
			}
			return resp
		})

}

func (h *handler) checkCamerAndDo(deviceName string, camera jdevice.Camera, w http.ResponseWriter, getResp func(c jdevice.Camera) interface{}) {
	if camera == nil {
		h.respFailed(fmt.Errorf("%s has not camera", deviceName), w)
		return
	}
	resp := getResp(camera)
	h.respSuccess(resp, w)
}
