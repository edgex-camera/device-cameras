package jdevice

type JDevice struct {
	Name   string
	Id     string
	Camera Camera
	Onvif  Onvif
}

type JDeviceConfig struct {
	Enabled bool   `json:enabled`
	Name    string `json:"name"`
	Id      string `json:"id"`
	Onvif   bool   `json:"onvif"`
}
