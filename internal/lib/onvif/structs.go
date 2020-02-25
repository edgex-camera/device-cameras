package onvif

import (
	"sync"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/yakovlevdmv/goonvif/xsd/onvif"
)

type Vector2D struct {
	X float64
	Y float64
}

type Move struct {
	PanTiltSpeed Vector2D // -1~1
	Zoom         float64  // -1~1
}

type OnvifCamera struct {
	Name         string
	device       OnvifDevice
	Lc           logger.LoggingClient
	OnvifConfig  OnvifConfig
	stopTimer    *time.Timer
	mutex        *sync.Mutex
	profileToken onvif.ReferenceToken
}

type OnvifConfig struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}
