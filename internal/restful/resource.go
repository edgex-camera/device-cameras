package restful

import (
	"net/http"
	"path"

	"github.com/edgex-camera/device-cameras/internal/driver"
	"github.com/edgex-camera/device-cameras/internal/jdevice"
	"github.com/gorilla/mux"
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
		Total      int      `json:"total"`
	}

	CameraConfigure, err := getCameraConfigJson(getCameraName(r))
	if err != nil {
		h.respFailed(err, w)
		return
	}

	from, until := h.getTimestampInfo(r)

	h.DoResp(
		driver.CurrentDriver.JDevices[getCameraName(r)].Camera,
		w,
		func(c jdevice.Camera) interface{} {
			pageInfo, total, err := h.timeFiliter(from, until, c.GetVideoPaths(), CameraConfigure.VideoConfig.Path, nil)
			if err != nil {
				return err
			}
			return responce{
				VideoPaths: pageInfo,
				Total:      total,
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
			names, total, err := h.listFilter(offset, limit, c.GetImagePaths(), func(data string) string {
				_, name := path.Split(data)
				return name
			})
			if err != nil {
				return err
			}
			pageInfo := h.GenUrlByName(names, "images", getCameraName(r))
			return responce{
				ImageURLs: pageInfo,
				Total:     total,
			}
		})
}

//获取video链接列表
func (h *handler) getVideoURLs(w http.ResponseWriter, r *http.Request) {
	type responce struct {
		VideoURLs []string `json:"videourls"`
		Total     int      `json:"total"`
	}

	CameraConfigure, err := getCameraConfigJson(getCameraName(r))
	if err != nil {
		h.respFailed(err, w)
		return
	}

	from, until := h.getTimestampInfo(r)

	h.DoResp(
		driver.CurrentDriver.JDevices[getCameraName(r)].Camera,
		w,
		func(c jdevice.Camera) interface{} {
			names, total, err := h.timeFiliter(from, until,
				c.GetVideoPaths(),
				CameraConfigure.VideoConfig.Path,
				func(data string) string {
					_, name := path.Split(data)
					return name
				})

			if err != nil {
				return err
			}
			pageInfo := h.GenUrlByName(names, "videos", getCameraName(r))
			return responce{
				VideoURLs: pageInfo,
				Total:     total,
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
