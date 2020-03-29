package restful

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/edgex-camera/device-cameras/internal/driver"
	"github.com/edgex-camera/device-cameras/internal/jdevice"
	"github.com/edgex-camera/device-cameras/internal/lib/onvif"
	"github.com/gorilla/mux"
)

func appendOnvifRoute(r *mux.Router, h *handler) {
	prefix := "/control"
	subRouter := r.PathPrefix(prefix).Subrouter()
	subRouter.Use(h.checkDeviceMiddvare, h.checkOnvifMiddvare)

	subRouter.Path("/{camera_name}/config").HandlerFunc(h.setConfig).Methods(http.MethodPost)
	subRouter.Path("/{camera_name}/config").HandlerFunc(h.getConfig).Methods(http.MethodGet)
	subRouter.Path("/{camera_name}/presets").HandlerFunc(h.getPresetPosition).Methods(http.MethodGet)
	subRouter.Path("/{camera_name}/continuous_move").HandlerFunc(h.postOnvifMove).Methods(http.MethodPost)
	subRouter.Path("/{camera_name}/stop").HandlerFunc(h.postOnvifStop).Methods(http.MethodPost)
	subRouter.Path("/{camera_name}/set_home_position").HandlerFunc(h.postSetHomePosition).Methods(http.MethodPost)
	subRouter.Path("/{camera_name}/reset_position").HandlerFunc(h.postResetPosition).Methods(http.MethodPost)
	subRouter.Path("/{camera_name}/set_preset/{preset_number}").HandlerFunc(h.postSetPresetPosition).Methods(http.MethodPost)
	subRouter.Path("/{camera_name}/goto_preset/{preset_number}").HandlerFunc(h.postGotoPresetPosition).Methods(http.MethodPost)

}

func (h *handler) setConfig(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.respFailed(err, w)
		return
	}

	driver.CurrentDriver.JDevices[getCameraName(r)].Control.PutConfig(data)
	h.respSuccess(nil, w)
}
func (h *handler) getConfig(w http.ResponseWriter, r *http.Request) {
	data := driver.CurrentDriver.JDevices[getCameraName(r)].Control.GetConfigure()

	h.respSuccess(data, w)
}

// 获取预置位信息
func (h *handler) getPresetPosition(w http.ResponseWriter, r *http.Request) {
	type responce struct {
		PresetPosition string `json:"presetposition"`
	}

	resp := responce{
		PresetPosition: driver.CurrentDriver.JDevices[getCameraName(r)].Control.GetPresets(),
	}
	h.respSuccess(resp, w)
}

// 移动
func (h *handler) postOnvifMove(w http.ResponseWriter, r *http.Request) {

	type request struct {
		X       float64 `json:"x"`
		Y       float64 `json:"y"`
		Zoom    float64 `json:"zoom"`
		TimeOut int     `json:"timeout"`
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.respFailed(err, w)
		return
	}

	req := request{}
	err = json.Unmarshal(data, &req)
	if err != nil {
		h.respFailed(err, w)
	}

	move := onvif.Move{
		PanTiltSpeed: onvif.Vector2D{
			X: req.X,
			Y: req.Y,
		},
		Zoom: req.Zoom,
	}

	h.onvifDo(
		driver.CurrentDriver.JDevices[getCameraName(r)].Control,
		w,
		func(jOnvif jdevice.Control) error {
			return jOnvif.ContinuousMove(time.Duration(req.TimeOut)*time.Millisecond, move)
		})
}

//停止
func (h *handler) postOnvifStop(w http.ResponseWriter, r *http.Request) {
	h.onvifDo(
		driver.CurrentDriver.JDevices[getCameraName(r)].Control,
		w,
		func(jOnvif jdevice.Control) error {
			return jOnvif.Stop()
		})
}

//设置零点
func (h *handler) postSetHomePosition(w http.ResponseWriter, r *http.Request) {
	h.onvifDo(
		driver.CurrentDriver.JDevices[getCameraName(r)].Control,
		w,
		func(jOnvif jdevice.Control) error {
			return jOnvif.SetHomePosition()
		})
}

//回到零点
func (h *handler) postResetPosition(w http.ResponseWriter, r *http.Request) {
	h.onvifDo(
		driver.CurrentDriver.JDevices[getCameraName(r)].Control,
		w,
		func(jOnvif jdevice.Control) error {
			return jOnvif.Reset()
		})

}

//设置预留位置
func (h *handler) postSetPresetPosition(w http.ResponseWriter, r *http.Request) {
	presetNumber, err := getPresetNumber(r)
	if err != nil {
		h.respFailed(fmt.Errorf("err posittion number %v", presetNumber), w)
		return
	}
	h.onvifDo(
		driver.CurrentDriver.JDevices[getCameraName(r)].Control,
		w,
		func(jOnvif jdevice.Control) error {
			return jOnvif.SetPreset(presetNumber)
		})
}

//回到预留位置
func (h *handler) postGotoPresetPosition(w http.ResponseWriter, r *http.Request) {
	presetNumber, err := getPresetNumber(r)
	if err != nil {
		h.respFailed(fmt.Errorf("err posittion number %v", presetNumber), w)
		return
	}
	h.onvifDo(
		driver.CurrentDriver.JDevices[getCameraName(r)].Control,
		w,
		func(jOnvif jdevice.Control) error {
			return jOnvif.GotoPreset(presetNumber)
		})
}

func (h *handler) onvifDo(onvif jdevice.Control, w http.ResponseWriter, toDo func(jOnvif jdevice.Control) error) {

	err := toDo(onvif)
	if err != nil {
		h.respFailed(err, w)
		return
	}
	h.respSuccess(nil, w)
}
