package normalcam

import (
	"github.com/edgexfoundry/device-sdk-go"
	"github.com/edgexfoundry/device-sdk-go/pkg/jxstartup"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/camera"
)

type NormalCamera struct {
	Name      string
	ChannelId string
	Camera    camera.RawCamera

	lc    logger.LoggingClient
	cmder camera.CameraCmder
	cc    camera.CameraConfig
}

func NewCamera(name, channelId string, lc logger.LoggingClient, cmder camera.CameraCmder, cc camera.CameraConfig) (jdevice.Camera, error) {
	ca := NormalCamera{Name: name}
	rawcam, err := NewRawCamera(name, channelId, lc, cmder, cc)
	if err != nil {
		return nil, err
	}
	ca.Camera = rawcam
	ca.ChannelId = channelId
	ca.lc = lc
	ca.cmder = cmder
	ca.cc = cc
	return &ca, nil
}

func NewRawCamera(name, channelId string, lc logger.LoggingClient, cmder camera.CameraCmder, cc camera.CameraConfig) (camera.RawCamera, error) {
	rawcam := camera.RawCamera{Lc: lc, Cmder: cmder, CameraConfig: cc}
	err := SetupRawCameraConfig(rawcam, name, channelId)
	if err != nil {
		return rawcam, err
	}
	return rawcam, nil
}

func SetupRawCameraConfig(c camera.RawCamera, name, channelId string) error {
	configName := name + ".camera." + channelId
	if _, ok := device.DriverConfigs()[configName]; !ok {
		return jxstartup.PutDriverConfig(configName, c.GetConfigure())
	}
	return nil
}
