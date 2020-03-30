package driver

import (
	"github.com/edgex-camera/device-sdk-go/pkg/camstartup"
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
	return camstartup.Service.GetDeviceByName(deviceName)
}

func GetDevices() []contract.Device {
	return camstartup.Service.Devices()
}
