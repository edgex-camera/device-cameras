package camera

import (
	"os/exec"
)

// CameraCmder produce command for stream and image
type CameraCmder interface {
	GetCmdProducers(cc CameraConfig) []func() *exec.Cmd
}
