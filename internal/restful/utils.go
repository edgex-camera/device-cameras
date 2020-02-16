package restful

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	resp := baseResponse{Data: err, Message: "failed"}
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