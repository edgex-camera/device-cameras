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
	subRouter.Path("/{camera_name}/presets").HandlerFunc(h.getPresetPosition).Methods(http.MethodGet)
	subRouter.Path("/{camera_name}/continuous_move").HandlerFunc(h.postOnvifMove).Methods(http.MethodPost)

	subRouter.Path("/{camera_name}/stop").HandlerFunc(h.postOnvifStop).Methods(http.MethodPost)
	subRouter.Path("/{camera_name}/set_home_position").HandlerFunc(h.postSetHomePosition).Methods(http.MethodPost)
	subRouter.Path("/{camera_name}/reset_position").HandlerFunc(h.postResetPosition).Methods(http.MethodPost)
	subRouter.Path("/{camera_name}/set_preset/{preset-number}").HandlerFunc(h.postSetPresetPosition).Methods(http.MethodPost)
	subRouter.Path("/{camera_name}/goto_preset/{preset-number}").HandlerFunc(h.postGotoPresetPosition).Methods(http.MethodPost)

}

// 获取预置位信息
func (h *handler) getPresetPosition(w http.ResponseWriter, r *http.Request) {
	type responce struct {
		PresetPosition string `json:"presetposition"`
	}
	cameraName := getCameraName(r)
	// TODO ,whether the deviceName is obtained from the request
	if jOnvif := driver.CurrentDriver.JDevices[cameraName].Onvif; jOnvif == nil {
		h.respFailed(fmt.Errorf("this %s devicee not support onvif", cameraName), w)
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

	move := onvif.Move{
		PanTiltSpeed: onvif.Vector2D{
			X: req.PanTiltSpeed.X,
			Y: req.PanTiltSpeed.Y,
		},
		Zoom: req.Zoom,
	}

	cameraName := getCameraName(r)
	h.checkOnvifAndDo(
		cameraName,
		driver.CurrentDriver.JDevices,
		w,
		func(jOnvif jdevice.Onvif) error {
			return jOnvif.ContinuousMove(time.Duration(req.TimeOut)*time.Second, move)
		})
}

//停止
func (h *handler) postOnvifStop(w http.ResponseWriter, r *http.Request) {
	cameraName := getCameraName(r)
	h.checkOnvifAndDo(
		cameraName,
		driver.CurrentDriver.JDevices,
		w,
		func(jOnvif jdevice.Onvif) error {
			return jOnvif.Stop()
		})
}

//设置零点
func (h *handler) postSetHomePosition(w http.ResponseWriter, r *http.Request) {
	cameraName := getCameraName(r)
	h.checkOnvifAndDo(
		cameraName,
		driver.CurrentDriver.JDevices,
		w,
		func(jOnvif jdevice.Onvif) error {
			return jOnvif.SetHomePosition()
		})
}

//回到零点
func (h *handler) postResetPosition(w http.ResponseWriter, r *http.Request) {
	cameraName := getCameraName(r)
	h.checkOnvifAndDo(
		cameraName,
		driver.CurrentDriver.JDevices,
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
	cameraName := getCameraName(r)
	h.checkOnvifAndDo(
		cameraName,
		driver.CurrentDriver.JDevices,
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
	cameraName := getCameraName(r)
	h.checkOnvifAndDo(
		cameraName,
		driver.CurrentDriver.JDevices,
		w,
		func(jOnvif jdevice.Onvif) error {
			return jOnvif.GotoPreset(presetNumber)
		})
}

func (h *handler) checkOnvifAndDo(deviceName string, jDevices map[string]jdevice.JDevice, w http.ResponseWriter, toDo func(jOnvif jdevice.Onvif) error) {
	device, ok := jDevices[deviceName]
	if !ok {
		h.respFailed(fmt.Errorf("has not device %s", deviceName), w)
		return
	}

	if device.Onvif == nil {
		h.respFailed(fmt.Errorf("this %s devicee not support onvif", deviceName), w)
		return
	}

	err := toDo(device.Onvif)
	if err != nil {
		h.respFailed(err, w)
		return
	}
	h.respSuccess(nil, w)
}

// type operation struct {
// 	DeviceName string
// 	Resource   interface{}
// 	w          http.ResponseWriter
// 	h          handler
// }

// func (o *operation) checkAndDo(deviceName string, resource interface{}, middleware ...func(deviceName string, resource interface{}) (interface{}, error)) {

// 	for _, do := range middleware {
// 		var err error
// 		o.Resource, err = do(o.DeviceName, o.Resource)
// 		if err != nil {
// 			o.h.respFailed(err, o.w)
// 		}

// 	}
// }

// func checkDeviceMiddleware(deviceName string, resource interface{}) (interface{}, error) {
// 	v, ok := resource.(map[string]jdevice.JDevice)
// 	if !ok {
// 		return nil, fmt.Errorf("invaild interface: not map[string]jdevice.JDevice")
// 	}

// 	return v[deviceName], nil
// }
// func checkCameraMiddleware(deviceName string, resource interface{}) (interface{}, error) {
// 	v, ok := resource.(jdevice.Camera)
// 	if !ok {
// 		return nil, fmt.Errorf("invaild interface: not jdevice.Camera")
// 	}
// 	return v, nil

// }

// func checkOnvifMiddleware(deviceName string, resource interface{}) (interface{}, error) {
// 	v, ok := resource.(jdevice.Onvif)
// 	if !ok {
// 		return nil, fmt.Errorf("invaild interface: not jdevice.Onvif")
// 	}
// 	return v, nil
// }
