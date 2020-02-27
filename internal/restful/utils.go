package restful

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/driver"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/camera"
)

//基本responce 结构
type baseResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (h *handler) respWithStatusCode(resp interface{}, w http.ResponseWriter, statusCode int) {
	data, err := json.Marshal(resp)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(data)
	if err != nil {
		h.lc.Error("error", err)
	}
}

//返回成功
func (h *handler) respSuccess(body interface{}, w http.ResponseWriter) {
	if value, ok := body.([]byte); ok {
		body = string(value)
	}

	resp := baseResponse{Data: body, Message: "success"}
	h.respWithStatusCode(resp, w, 200)
}

func (h *handler) respFailed(err error, w http.ResponseWriter) {
	resp := baseResponse{Data: nil, Message: fmt.Sprintln(err)}
	h.respWithStatusCode(resp, w, 500)
}

func getCameraName(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["camera_name"]
}

func getDeviceType(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["device_type"]
}

func getPresetNumber(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	return strconv.ParseInt(vars["preset_number"], 10, 64)
}

func (h *handler) checkDeviceMiddvare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deviceName := getCameraName(r)
		_, ok := driver.CurrentDriver.JDevices[deviceName]
		if !ok {
			h.respFailed(fmt.Errorf("has not device %s", deviceName), w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *handler) checkOnvifMiddvare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deviceName := getCameraName(r)

		if driver.CurrentDriver.JDevices[deviceName].Control == nil {
			h.respFailed(fmt.Errorf("this %s devicee not support onvif", deviceName), w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *handler) checkCameraMiddvare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deviceName := getCameraName(r)

		if driver.CurrentDriver.JDevices[deviceName].Camera == nil {
			h.respFailed(fmt.Errorf("this %s devicee not support camera", deviceName), w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *handler) getPageInfo(r *http.Request) (int, int) {
	offset := 0
	v := r.URL.Query()["offset"]
	if len(v) == 1 {
		offset, _ = strconv.Atoi(v[0])
	}
	limit := 10
	v = r.URL.Query()["limit"]
	if len(v) == 1 {
		limit, _ = strconv.Atoi(v[0])
	}
	return offset, limit
}

func (h *handler) getTimestampInfo(r *http.Request) (int64, int64) {
	from := 0
	v := r.URL.Query()["from"]
	if len(v) == 1 {
		from, _ = strconv.Atoi(v[0])
	}
	until := 9999999999
	v = r.URL.Query()["until"]
	if len(v) == 1 {
		until, _ = strconv.Atoi(v[0])
	}
	return int64(from), int64(until)
}

func (h *handler) timeFiliter(from, until int64, data []string, model string, tailProcess func(data string) string) ([]string, int, error) {
	length := len(data)
	res := make([]string, 0)
	timestampPreset := strings.Index(model, "%s")
	for _, name := range data {

		timestamp, err := strconv.Atoi(name[timestampPreset : timestampPreset+10])
		if err != nil {
			return nil, 0, err
		}
		log.Println(int(timestamp), from, until)
		if from < int64(timestamp) && int64(timestamp) < until {
			if tailProcess != nil {
				res = append(res, tailProcess(name))
			} else {
				res = append(res, name)
			}
		}
	}

	return res, length, nil
}

func (h *handler) listFilter(offset, limit int, data []string, preProcess func(data string) string) ([]string, int, error) {

	res := make([]string, 0)
	if preProcess != nil {
		for _, name := range data {
			res = append(res, preProcess(name))
		}
	} else {
		res = data
	}

	length := len(res)
	if offset > length {
		return []string{}, length, fmt.Errorf("offset more than total")
	}

	if offset+limit < length {
		return res[offset : offset+limit], length, nil

	}
	return res[offset:], length, nil
}

func (h *handler) GenUrlByName(data []string, mode, camerName string) []string {
	prefix := "/api/v1/driver/device-cameras/api/v1/static"
	res := make([]string, 0)
	for _, name := range data {
		res = append(res, strings.Join([]string{prefix, mode, camerName, name}, "/"))
	}
	return res

}

func getCameraConfigJson(cameraName string) (*camera.CameraConfig, error) {
	JDevice := driver.CurrentDriver.JDevices[cameraName]
	data := JDevice.Camera.GetConfigure()
	CameraConfigure := camera.CameraConfig{}
	err := json.Unmarshal(data, &CameraConfigure)
	return &CameraConfigure, err
}
