package restful

import (
	"net/http"

	"github.com/gorilla/mux"
)

func appendOnvifRoute(r *mux.Route) {}

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
