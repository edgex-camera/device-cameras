package normalcam

import "os"

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

func (nc *NormalCamera) GetVideoPaths() []string {
	return nc.Camera.GetVideoPaths()
}

func (nc *NormalCamera) GetStreamAddr() string {
	return nc.Camera.GetStreamAddr()
}

// configs
func (nc *NormalCamera) MergeConfig(configPatch []byte) error {
	return nc.Camera.MergeConfig(configPatch)
}

func (nc *NormalCamera) GetConfigure() []byte {
	return nc.Camera.GetConfigure()
}

// TODO: channel management
func (nc *NormalCamera) AddChannel(channelId string) error {
	return nil
}

func (nc *NormalCamera) RemoveChannel(channelId string) error {
	return nil
}

func (nc *NormalCamera) ListChannels() []byte {
	res := make([]byte, 10)
	return res
}
