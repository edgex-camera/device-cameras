package normalcam

import (
	"errors"
	"os"

	"github.com/edgex-camera/device-cameras/internal/lib/utils"
	"github.com/edgex-camera/device-sdk-go/pkg/camstartup"
)

// camera functions
func (nc *NormalCamera) Enable() {
	nc.Camera.Enable()
}

func (nc *NormalCamera) Disable(wait bool) {
	nc.Camera.Disable(wait)
}

func (nc *NormalCamera) IsEnabled() bool {
	return nc.Camera.IsEnabled()
}

func (nc *NormalCamera) CapturePhotoJPG() (*os.File, error) {
	return nc.Camera.CapturePhotoJPG()
}

func (nc *NormalCamera) GetCapturePath() string {
	return nc.Camera.GetCapturePath()
}

func (nc *NormalCamera) GetImagePaths() []string {
	return nc.Camera.GetImagePaths()
}

func (nc *NormalCamera) GetVideoPaths() []string {
	return nc.Camera.GetVideoPaths()
}

func (nc *NormalCamera) GetStreamAddr() string {
	return nc.Camera.GetStreamAddr()
}

// configs
func (nc *NormalCamera) MergeConfig(conf map[string]string) error {
	configPatch := conf[nc.Name+".camera."+nc.ChannelId]
	return nc.Camera.MergeConfig([]byte(configPatch))
}

func (nc *NormalCamera) GetConfigure() []byte {
	return nc.Camera.GetConfigure()
}

func (nc *NormalCamera) PutConfig(config []byte) error {
	configName := nc.Name + ".camera." + nc.ChannelId
	return camstartup.PutDriverConfig(configName, config)
}

func (nc *NormalCamera) AddChannel() error {
	if nc.ChannelId != "" {
		return errors.New("A channel already exists.")
	}
	channelId := utils.GenUUID()
	rawcam, err := NewRawCamera(nc.Name, channelId, nc.lc, nc.cmder, nc.cc)
	if err != nil {
		return nil
	}
	nc.Camera = rawcam
	nc.ChannelId = channelId
	return nil
}

func (nc *NormalCamera) RemoveChannel(channelId string) error {
	if channelId != nc.ChannelId {
		return errors.New("No such channel id")
	}
	nc.Camera.Disable(true)
	nc.ChannelId = ""
	nc.Camera.CameraConfig.Enabled = false
	return SetupRawCameraConfig(nc.Camera, nc.Name, channelId)
}

func (nc *NormalCamera) ListChannels() []string {
	var res []string
	res = append(res, nc.ChannelId)
	return res
}
