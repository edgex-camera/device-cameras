package jdevice

import (
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice/camera"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice/onvif"
)

type JDevice struct {
	Camera camera.Camera
	Onvif  onvif.Onvif
}
