package camera

import (
	"os/exec"
)

// CameraCmder produce command for stream and image
type CameraCmder interface {
	GetCmdProducers(cc CameraConfig) []func() *exec.Cmd
}

type CameraConfig struct {
	Enabled       bool   `json:"enabled"`
	InputAddr     string `json:"input_addr"`
	Auth          `json:"auth"`
	StreamConfig  `json:"stream"`
	CaptureConfig `json:"capture"`
	VideoConfig   `json:"video"`
	QualityConfig `json:"quality"`
}
