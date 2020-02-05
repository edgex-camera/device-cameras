package camera

import (
	"os"
)

// Camera stores config and provide REST interface
type Camera interface {
	Enable()
	Disable(wait bool)
	IsEnabled() bool

	CapturePhotoJPG() (*os.File, error)
	GetCapturePath() string
	GetVideoPaths() []string
	GetStreamAddr() string

	// use Json config
	MergeConfig(configPatch []byte) error
	GetConfigure() []byte
}
