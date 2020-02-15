package restful

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/driver"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/onvif"
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
func (h *handler) getPresetPosition(w http.ResponseWriter, r *http.Request) {
	type responce struct {
		PresetPosition string `json:"presetposition"`
	}
	deviceName := "deviceName"
	// TODO ,whether the deviceName is obtained from the request
	if jOnvif := driver.CurrentDriver.JDevices[deviceName].Onvif; jOnvif == nil {
		h.respFailed(fmt.Errorf("this %s devicee not support onvif", deviceName), w)
		return
	} else {
		resp := responce{
			PresetPosition: jOnvif.GetPresets(),
		}
		h.respSuccess(resp, w)

	}
}

// 移动
func (h *handler) postOnvifMove(w http.ResponseWriter, r *http.Request) {
	type vector2D struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
	type request struct {
		PanTiltSpeed vector2D `json:"pantiltspeed"`
		Zoom         float64  `json:"root"`
		TimeOut      int      `json:"timeout"`
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.respFailed(err, w)
		return
	}

	req := request{
		PanTiltSpeed: vector2D{},
	}
	err = json.Unmarshal(data, &req)
	if err != nil {
		h.respFailed(err, w)
	}

	deviceName := "deviceName"
	// TODO ,whether the deviceName is obtained from the request
	move := onvif.Move{
		PanTiltSpeed: onvif.Vector2D{
			X: req.PanTiltSpeed.X,
			Y: req.PanTiltSpeed.Y,
		},
		Zoom: req.Zoom,
	}

	h.checkOnvifAndDo(
		deviceName,
		driver.CurrentDriver.JDevices[deviceName].Onvif,
		w,
		func(jOnvif jdevice.Onvif) error {
			return jOnvif.ContinuousMove(time.Duration(req.TimeOut)*time.Second, move)
		})
}

//停止
func (h *handler) postOnvifStop(w http.ResponseWriter, r *http.Request) {
	deviceName := "deviceName"
	// TODO ,whether the deviceName is obtained from the request
	h.checkOnvifAndDo(
		deviceName,
		driver.CurrentDriver.JDevices[deviceName].Onvif,
		w,
		func(jOnvif jdevice.Onvif) error {
			return jOnvif.Stop()
		})
}

//设置零点
func (h *handler) postSetHomePosition(w http.ResponseWriter, r *http.Request) {
	deviceName := "deviceName"
	// TODO ,whether the deviceName is obtained from the request
	h.checkOnvifAndDo(
		deviceName,
		driver.CurrentDriver.JDevices[deviceName].Onvif,
		w,
		func(jOnvif jdevice.Onvif) error {
			return jOnvif.SetHomePosition()
		})
}

//回到零点
func (h *handler) postResetPosition(w http.ResponseWriter, r *http.Request) {
	deviceName := "deviceName"
	// TODO ,whether the deviceName is obtained from the request
	h.checkOnvifAndDo(
		deviceName,
		driver.CurrentDriver.JDevices[deviceName].Onvif,
		w,
		func(jOnvif jdevice.Onvif) error {
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
	deviceName := "deviceName"
	h.checkOnvifAndDo(
		deviceName,
		driver.CurrentDriver.JDevices[deviceName].Onvif,
		w,
		func(jOnvif jdevice.Onvif) error {
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
	deviceName := "deviceName"
	h.checkOnvifAndDo(
		deviceName,
		// presetNumber,
		driver.CurrentDriver.JDevices[deviceName].Onvif,
		w,
		func(jOnvif jdevice.Onvif) error {
			return jOnvif.GotoPreset(presetNumber)
		})
}

func (h *handler) checkOnvifAndDo(deviceName string, jOnvif jdevice.Onvif, w http.ResponseWriter, toDo func(jOnvif jdevice.Onvif) error) {
	// var number int64 = 0
	// if v, ok := presetNumber.(int64); ok {
	// 	number = v
	// }
	if jOnvif == nil {
		h.respFailed(fmt.Errorf("this %s devicee not support onvif", deviceName), w)
		return
	} else {
		err := toDo(jOnvif)
		if err != nil {
			h.respFailed(err, w)
			return
		}
		h.respSuccess(nil, w)
	}
}
