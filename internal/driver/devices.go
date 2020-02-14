package driver

import (
	"github.com/edgexfoundry/device-sdk-go/pkg/jxstartup"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
)

func NewDevice(deviceName string) contract.Device {
	deviceServiceName := "camera"
	protocols := make(map[string]contract.ProtocolProperties)
	protocolProperties := make(map[string]string)
	protocolProperties["Address"] = "/api/v1/device/camera"
	protocols["other"] = protocolProperties
	return contract.Device{
		Name:           deviceName,
		AdminState:     "UNLOCKED",
		OperatingState: "ENABLED",
		Protocols:      protocols,
		Service:        contract.DeviceService{Name: deviceServiceName},
		Profile:        contract.DeviceProfile{Name: deviceServiceName},
	}
}

func GetDeviceByName(deviceName string) (contract.Device, error) {
	return jxstartup.Service.GetDeviceByName(deviceName)
}

func GetDevices() []contract.Device {
	return jxstartup.Service.Devices()
}
