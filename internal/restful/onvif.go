package restful

import (
	"net/http"

	"github.com/gorilla/mux"
)

func appendOnvifRoute(r *mux.Router, h *handler) {
	prefix := "/onvif"
	subRouter := r.PathPrefix(prefix).Subrouter()

	subRouter.Path("/presets").HandlerFunc(h.getPresetPosition).Methods(http.MethodGet)
	subRouter.Path("/continuous_move").HandlerFunc(h.postOnvifMove).Methods(http.MethodPost)

	subRouter.Path("/stop").HandlerFunc(h.postOnvifStop).Methods(http.MethodPost)
	subRouter.Path("/set_home_position").HandlerFunc(h.postSetHomePosition).Methods(http.MethodPost)
	subRouter.Path("/reset_position").HandlerFunc(h.postResetPosition).Methods(http.MethodPost)
	subRouter.Path("/set_preset/{preset-number}").HandlerFunc(h.postSetPresetPosition).Methods(http.MethodPost)
	subRouter.Path("/goto_preset/{preset-number}").HandlerFunc(h.postGotoPresetPosition).Methods(http.MethodPost)

}

// 获取预置位信息
func (h *handler) getPresetPosition(w http.ResponseWriter, r *http.Request) {}

// 移动
func (h *handler) postOnvifMove(w http.ResponseWriter, r *http.Request) {}

//停止
func (h *handler) postOnvifStop(w http.ResponseWriter, r *http.Request) {}

//设置零点
func (h *handler) postSetHomePosition(w http.ResponseWriter, r *http.Request) {}

//回到零点
func (h *handler) postResetPosition(w http.ResponseWriter, r *http.Request) {}

//设置预留位置
func (h *handler) postSetPresetPosition(w http.ResponseWriter, r *http.Request) {}

//回到预留位置
func (h *handler) postGotoPresetPosition(w http.ResponseWriter, r *http.Request) {}
