package driver

import (
	"time"

	"github.com/edgex-camera/device-cameras/internal/lib/onvif"
	dsModels "github.com/edgex-camera/device-sdk-go/pkg/models"
)

func (d *Driver) HandleMoveCommand(deviceName string, params []*dsModels.CommandValue) error {
	move := onvif.Move{}
	timeout := 1 * time.Second
	for _, param := range params {
		v, err := param.Float32Value()
		if err != nil {
			return err
		}
		switch param.DeviceResourceName {
		case "pan":
			move.PanTiltSpeed.X = float64(v)
		case "tilt":
			move.PanTiltSpeed.Y = float64(v)
		case "zoom":
			move.Zoom = float64(v)
		case "timeout":
			timeout = time.Duration(int64(float32(time.Second) * v))
		}
	}
	return d.JDevices[deviceName].Control.ContinuousMove(timeout, move)
}
