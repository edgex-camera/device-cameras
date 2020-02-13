package jdevice

import (
	"os"
	"time"

	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/onvif"
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

	// Json config
	MergeConfig(configPatch []byte) error
	GetConfigure() []byte

	// For channel manage
	AddChannel() error
	RemoveChannel(channelId string) error
	ListChannels() []string
}

type Onvif interface {
	MergeConfig(configPatch []byte) error

	ContinuousMove(time time.Duration, move onvif.Move) error
	Stop() error

	SetHomePosition() error
	Reset() error

	GetPresets() string
	SetPreset(number int64) error
	GotoPreset(number int64) error

	SyncTime() error
}
