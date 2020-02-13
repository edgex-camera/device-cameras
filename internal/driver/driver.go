package driver

import (
	dsModels "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice"
)

const CAMERA_FACTORY = "camera-factory"

type Driver struct {
	lc            logger.LoggingClient
	asyncCh       chan<- *dsModels.AsyncValues
	JDevices      map[string]jdevice.JDevice
	ProcessMethod string
}

var CurrentDriver Driver

func (d *Driver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues) error {
	d.lc = lc
	d.asyncCh = asyncCh
	d.JDevices = make(map[string]jdevice.JDevice)

	return nil
}

func (d *Driver) HandleReadCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {
	return nil, err
}

func (d *Driver) HandleWriteCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest, params []*dsModels.CommandValue) error {
	if deviceName == CAMERA_FACTORY {
		for _, param := range params {
			switch param.DeviceResourceName {
			case "add_device", "device_type":
				var name, deviceType string
				switch param.DeviceResourceName {
				case "add_device":
					v, err := param.StringValue()
					if err != nil {
						return err
					}
					name = v
				case "device_type":
					v, err := param.StringValue()
					if err != nil {
						return err
					}
					deviceType = v
					return d.AddJdevice(name, deviceType)
				}
			case "remove_device":
				v, err := param.StringValue()
				if err != nil {
					return err
				}
				return d.RemoveJdevice(v)
			}
		}
	}
	return nil
}

func (d *Driver) Stop(force bool) error {
	return nil
}
