package restful

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/driver"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice"
)

func appendCameraResourcesRoute(r *mux.Router, h *handler) {
	prefix := "/resource"
	subRouter := r.PathPrefix(prefix).Subrouter()
	subRouter.Use(h.checkDeviceMiddvare, h.checkCameraMiddvare)

	subRouter.Path("/{camera_name}/video_paths").HandlerFunc(h.getVideoPaths).Methods(http.MethodGet)
	subRouter.Path("/{camera_name}/images").HandlerFunc(h.getImageURls).Methods(http.MethodGet)
	subRouter.Path("/{camera_name}/vidoes").HandlerFunc(h.getVideoURLs).Methods(http.MethodGet)
}

//获取视频存储路径
func (h *handler) getVideoPaths(w http.ResponseWriter, r *http.Request) {
	type responce struct {
		VideoPaths []string `json:"videopaths"`
	}
	h.DoResp(
		driver.CurrentDriver.JDevices[getCameraName(r)].Camera,
		w,
		func(c jdevice.Camera) interface{} {
			return responce{
				VideoPaths: c.GetVideoPaths(),
			}
		})

}

//获取图片链接列表
func (h *handler) getImageURls(w http.ResponseWriter, r *http.Request) {
	type responce struct {
		ImageURLs []string `json:"imageurls"`
		Total     int      `json:"total"`
	}
	offset, limit := h.getPageInfo(r)
	h.DoResp(
		driver.CurrentDriver.JDevices[getCameraName(r)].Camera,
		w,
		func(c jdevice.Camera) interface{} {

			pageInfo, err := h.currentPage(offset, limit, c.GetVideoPaths())
			if err != nil {
				return err
			}
			return responce{
				ImageURLs: pageInfo,
				Total:     len(c.GetImagePaths()),
			}
		})
}

//获取video链接列表
func (h *handler) getVideoURLs(w http.ResponseWriter, r *http.Request) {
	type responce struct {
		VideoURLs []string `json:"videourls"`
		Total     int      `json:"total"`
	}
	offset, limit := h.getPageInfo(r)

	h.DoResp(
		driver.CurrentDriver.JDevices[getCameraName(r)].Camera,
		w,
		func(c jdevice.Camera) interface{} {
			pageInfo, err := h.currentPage(offset, limit, c.GetVideoPaths())
			if err != nil {
				return err
			}
			return responce{
				VideoURLs: pageInfo,
				Total:     len(c.GetVideoPaths()),
			}
		})

}

func (h *handler) DoResp(camera jdevice.Camera, w http.ResponseWriter, getResp func(c jdevice.Camera) interface{}) {
	resp := getResp(camera)
	if v, ok := resp.(error); ok {
		h.respFailed(v, w)
		return
	}
	h.respSuccess(resp, w)
}
