package onvif

import (
	"time"
)

type Onvif interface {
	MergeConfig(configPatch []byte) error

	ContinuousMove(time time.Duration, move Move) error
	Stop() error

	SetHomePosition() error
	Reset() error

	GetPresets() string
	SetPreset(number int64) error
	GotoPreset(number int64) error

	SyncTime() error
}

type Vector2D struct {
	X float64
	Y float64
}

type Move struct {
	PanTiltSpeed Vector2D // -1~1
	Zoom         float64  // -1~1
}
