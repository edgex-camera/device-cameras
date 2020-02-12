package driver

import (
	"encoding/json"
	"fmt"

	"github.com/edgexfoundry/device-sdk-go/pkg/jxstartup"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice/normalcam"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/camera"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/camera/cmder"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/onvif"
)

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

	jDevice := jdevice.JDevice{Id: id}

	// 普通摄像头、onvif摄像头的摄像头部分
	if deviceType == NORMAL_CAMERA || deviceType == ONVIF_CAMERA {
		var deviceCmder camera.CameraCmder
		if d.ProcessMethod == PROCESS_GST {
			deviceCmder = cmder.NewCmder(PROCESS_GST)
		} else {
			deviceCmder = cmder.NewCmder(PROCESS_FFMPEG)
		}
		cc := camera.CameraConfig{}
		deviceCamera := normalcam.NewCamera(deviceName, d.lc, deviceCmder, cc)
		jDevice.Camera = deviceCamera
	}
	// onvif摄像头onvif部分
	if deviceType == ONVIF_CAMERA {
		config := onvif.OnvifConfig{}
		deviceOnvif, err := jdevice.NewOnvif(deviceName, d.lc, config)
		if err != nil {
			return err
		}
		jDevice.Onvif = deviceOnvif
	}
	if deviceType == DUAL_USB_CAMERA {

	}
	if deviceType == SIMPLE_CAMERA {

	}

	d.JDevices[deviceName] = jDevice
	return setupJdeviceConfig(jDevice)
}

// JDevice基础信息
func setupJdeviceConfig(jDevice jdevice.JDevice) error {
	config := jdevice.JDeviceConfig{}
	config.Name = jDevice.Name
	config.Id = jDevice.Id
	if jDevice.Onvif != nil {
		config.Onvif = true
	}
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	configName := config.Name + "/" + "basic"
	return jxstartup.PutDriverConfig(configName, configBytes)
}
