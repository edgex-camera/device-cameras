package restful

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

var prefix = "/static"

func appendStaticRoute(r *mux.Router, h *handler) {

	subRouter := r.PathPrefix(prefix).Subrouter()
	subRouter.PathPrefix("/images/{camera_name}").HandlerFunc(h.getImageStatic)
	subRouter.PathPrefix("/videos/{camera_name}").HandlerFunc(h.getVideoStatic)

}

func (h *handler) getImageStatic(w http.ResponseWriter, r *http.Request) {
	cameraName := getCameraName(r)
	CameraConfigure, err := getCameraConfigJson(cameraName)
	dir, _ := path.Split(CameraConfigure.CaptureConfig.Path)
	if err != nil {
		h.respFailed(err, w)
		return
	}
	http.StripPrefix(APIv1Prefix+prefix+"/images/"+cameraName+"/", http.FileServer(http.Dir(dir))).ServeHTTP(w, r)
}

func (h *handler) getVideoStatic(w http.ResponseWriter, r *http.Request) {
	cameraName := getCameraName(r)
	CameraConfigure, err := getCameraConfigJson(cameraName)
	dir, _ := path.Split(CameraConfigure.VideoConfig.Path)
	if err != nil {
		h.respFailed(err, w)
		return
	}
	http.StripPrefix(APIv1Prefix+prefix+"/videos/"+cameraName+"/", http.FileServer(http.Dir(dir))).ServeHTTP(w, r)
}
