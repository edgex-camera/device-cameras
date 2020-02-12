package normalcam

import (
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/camera"
)

type NormalCamera struct {
	ChannelId string
	Camera    camera.RawCamera
}

func NewCamera(channelId string, lc logger.LoggingClient, cmder camera.CameraCmder, cc camera.CameraConfig) jdevice.Camera {
	ca := NormalCamera{}
	ca.Camera = camera.RawCamera{Lc: lc, Cmder: cmder, CameraConfig: cc}
	ca.ChannelId = channelId
	return &ca
}
