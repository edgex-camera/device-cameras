package driver

import (
	dsModels "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice"
)

type Driver struct {
	lc      logger.LoggingClient
	asyncCh chan<- *dsModels.AsyncValues
	Devices map[string]jdevice.JDevice
}

func (d *Driver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues) error {
	return nil
}

func (d *Driver) HandleReadCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {
	return nil, err
}

func (d *Driver) HandleWriteCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest, params []*dsModels.CommandValue) error {
	return nil
}

func (d *Driver) Stop(force bool) error {
	return nil
}
