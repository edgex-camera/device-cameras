package driver

import (
	"encoding/json"
	"fmt"

	"github.com/edgexfoundry/device-sdk-go"
	"github.com/edgexfoundry/device-sdk-go/pkg/jxstartup"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice/normalcam"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/camera"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/camera/cmder"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/onvif"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/utils"
)

const ALL_DEVICES_KEY = "all_devices"

// 设备类型
const NORMAL_CAMERA = "normal-camera"     // 普通usb/ip摄像头
const ONVIF_CAMERA = "onvif-camera"       // onvif摄像头(球机)
const DUAL_USB_CAMERA = "dual-usb-camera" // usb双摄
const SIMPLE_CAMERA = "simple-camera"     // 低配摄像头

// 执行程序类型
const PROCESS_GST = "gst-launch-1.0"
const PROCESS_FFMPEG = "ffmpeg"

// Add Jdevice.
func (d *Driver) AddJdevice(deviceName, deviceType string) error {
	device := NewDevice(deviceName)
	id, err := jxstartup.Service.AddDevice(device)
	if err != nil {
		d.lc.Error(fmt.Sprintf("Failed to add device with error: %v", err))
		return err
	}

	jDevice := jdevice.JDevice{Id: id, Type: deviceType, Name: deviceName}

	// 普通摄像头、onvif摄像头的摄像头部分
	if deviceType == NORMAL_CAMERA || deviceType == ONVIF_CAMERA {
		var deviceCmder camera.CameraCmder
		if d.ProcessMethod == PROCESS_GST {
			deviceCmder = cmder.NewCmder(PROCESS_GST)
		} else {
			deviceCmder = cmder.NewCmder(PROCESS_FFMPEG)
		}
		cc := camera.CameraConfig{}
		// 默认生成一个channel的摄像头
		channelId := utils.GenUUID()
		deviceCamera, err := normalcam.NewCamera(deviceName, channelId, d.lc, deviceCmder, cc)
		if err != nil {
			return err
		}
		jDevice.Camera = deviceCamera
	}
	// onvif摄像头onvif部分
	if deviceType == ONVIF_CAMERA {
		config := onvif.OnvifConfig{}
		deviceControl, err := jdevice.NewOnvif(deviceName, d.lc, config)
		if err != nil {
			return err
		}
		jDevice.Control = deviceControl
	}
	if deviceType == DUAL_USB_CAMERA {

	}
	if deviceType == SIMPLE_CAMERA {

	}

	d.JDevices[deviceName] = jDevice
	return setupJdeviceConfig(jDevice, true, deviceType)
}

// Remove Jdevice
func (d *Driver) RemoveJdevice(deviceName string) error {
	if _, ok := d.JDevices[deviceName]; !ok {
		d.lc.Info(fmt.Sprintf("Device to remove is not running "), deviceName)
	} else {
		if d.JDevices[deviceName].Camera != nil {
			d.JDevices[deviceName].Camera.Disable(true)
		}
		delete(d.JDevices, deviceName)
	}
	jDevice := jdevice.JDevice{Name: deviceName}
	err := setupJdeviceConfig(jDevice, false, "")
	if err != nil {
		return err
	}
	return jxstartup.Service.RemoveDeviceByName(deviceName)
}

// JDevice基础信息
func setupJdeviceConfig(jDevice jdevice.JDevice, enabled bool, deviceType string) error {
	// 修改all_devices配置信息
	allDevices := []byte{}
	allDevicesMap := make(map[string]bool)
	if all, ok := device.DriverConfigs()[ALL_DEVICES_KEY]; ok {
		allDevices = []byte(all)
		json.Unmarshal(allDevices, &allDevicesMap)
	}

	allDevicesMap[jDevice.Name] = enabled
	if !enabled {
		delete(allDevicesMap, jDevice.Name)
	}
	allDevices, _ = json.Marshal(allDevicesMap)
	err := jxstartup.PutDriverConfig(ALL_DEVICES_KEY, allDevices)
	if err != nil {
		return err
	}

	// 修改device自身配置
	config := jdevice.JDeviceConfig{}
	config.Enabled = enabled
	config.Name = jDevice.Name
	config.Id = jDevice.Id
	config.Type = deviceType
	if deviceType == ONVIF_CAMERA {
		config.Control = "onvif"
	} else {
		config.Control = "none"
	}
	configName := config.Name
	configBytes, _ := json.Marshal(config)
	return jxstartup.PutDriverConfig(configName, configBytes)
}
