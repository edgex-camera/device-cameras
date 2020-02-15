package restful

import (
	"net/http"

	"github.com/gorilla/mux"
)

func appendCameraResourcesRoute(r *mux.Router, h *handler) {
	prefix := "/source"
	subRouter := r.PathPrefix(prefix).Subrouter()

	subRouter.Path("/video_paths").HandlerFunc(h.getVideoPaths).Methods(http.MethodGet)
	subRouter.Path("/images").HandlerFunc(h.getImageURls).Methods(http.MethodGet)
	subRouter.Path("/vidoes").HandlerFunc(h.getImageURls).Methods(http.MethodGet)
}

//获取视频存储路径
func (h *handler) getVideoPaths(w http.ResponseWriter, r *http.Request) {}

//获取图片链接列表
func (h *handler) getImageURls(w http.ResponseWriter, r *http.Request) {}

//获取video链接列表
func (h *handler) getVideoURLs(w http.ResponseWriter, r *http.Request) {}
