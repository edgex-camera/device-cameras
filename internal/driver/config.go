package driver

import (
	"encoding/json"
	"fmt"

	"github.com/edgex-camera/device-cameras/internal/jdevice"
	"github.com/edgex-camera/device-cameras/internal/jdevice/normalcam"
	"github.com/edgex-camera/device-cameras/internal/lib/camera"
	"github.com/edgex-camera/device-cameras/internal/lib/camera/cmder"
	"github.com/edgex-camera/device-cameras/internal/lib/onvif"
)

func (d *Driver) OnConfigChange(oldConf map[string]string, newConf map[string]string) {
	d.lc.Info("Config changed ...")
	allDevices, ok := newConf[ALL_DEVICES_KEY]
	// 无任何设备
	if !ok {
		return
	}

	allDevicesMap := make(map[string]bool)
	json.Unmarshal([]byte(allDevices), &allDevicesMap)

	// 在deviceMap中的处理
	for deviceName := range allDevicesMap {
		_, ok := d.JDevices[deviceName]
		if !ok {
			// 目前未运行该设备，根据配置加入
			d.addDeviceByConfig(deviceName, newConf)
		} else if d.cameraConfigChanged(deviceName, oldConf, newConf) {
			// 目前配置与旧配置不同，根据配置修改
			d.modifyDeviceByConfig(deviceName, newConf)
		}
	}

	// 在jdevices中，但配置中没有的
	for name := range d.JDevices {
		if _, ok := allDevicesMap[name]; !ok {
			d.RemoveJdevice(name)
		}
	}
	d.lc.Info(fmt.Sprint("Current devices: ", d.JDevices))
}

func (d *Driver) cameraConfigChanged(deviceName string, oldConf, newConf map[string]string) bool {
	basicConfName := deviceName + ".camera"
	if oldConf[basicConfName] != newConf[basicConfName] {
		return true
	}
	channelsMap := make(map[string]bool)
	json.Unmarshal([]byte(newConf[basicConfName]), &channelsMap)
	for channelId := range channelsMap {
		channelConfName := deviceName + ".camera." + channelId
		if oldConf[channelConfName] != newConf[channelConfName] {
			return true
		}
	}
	return false
}

func (d *Driver) addDeviceByConfig(deviceName string, conf map[string]string) {
	d.lc.Info(fmt.Sprint("Adding device: ", deviceName))

	basicStr, ok := conf[deviceName]
	if !ok {
		d.lc.Info(fmt.Sprintf("Device with name %v config does not have basic.", deviceName))
		return
	}

	basicConf := jdevice.JDeviceConfig{}
	json.Unmarshal([]byte(basicStr), &basicConf)
	if !basicConf.Enabled {
		d.lc.Info(fmt.Sprintf("Device %v is not enabled.", deviceName))
		return
	}

	// 创建jDevice实例
	jDevice := jdevice.JDevice{Id: basicConf.Id, Name: basicConf.Name, Type: basicConf.Type}

	// 创建onvif实例
	if basicConf.Control == "onvif" {
		onvifStr, ok := conf[deviceName+".onvif.config"]
		if !ok {
			d.lc.Info(fmt.Sprintf("Device with name %v onvif config not exists.", deviceName))
			return
		}
		onvifConf := onvif.OnvifConfig{}
		json.Unmarshal([]byte(onvifStr), &onvifConf)
		deviceOnvif, err := jdevice.NewOnvif(deviceName, d.lc, onvifConf)
		if err != nil {
			d.lc.Info(fmt.Sprint("Create onvif device failed: ", err))
			return
		}
		jDevice.Control = deviceOnvif
	}

	// 创建sip实例
	// TODO
	// if basicConf.Control == "sip" {}

	// 创建camera实例
	cameraStr, ok := conf[deviceName+".camera"]
	if !ok {
		d.lc.Info(fmt.Sprintf("Device with name %v camera config not exists.", deviceName))
		return
	}
	if basicConf.Type == NORMAL_CAMERA || basicConf.Type == ONVIF_CAMERA {
		var deviceCmder camera.CameraCmder
		if d.ProcessMethod == PROCESS_GST {
			deviceCmder = cmder.NewCmder(PROCESS_GST)
		} else {
			deviceCmder = cmder.NewCmder(PROCESS_FFMPEG)
		}
		cc := camera.CameraConfig{}
		var channelId string
		channelsMap := make(map[string]bool)
		json.Unmarshal([]byte(cameraStr), &channelsMap)
		for id := range channelsMap {
			channelId = id
			json.Unmarshal([]byte(conf[deviceName+".camera."+channelId]), &cc)
		}
		deviceCamera, err := normalcam.NewCamera(deviceName, channelId, d.lc, deviceCmder, cc)
		if err != nil {
			d.lc.Info(fmt.Sprint("Failed to create camera: ", err))
			return
		}
		jDevice.Camera = deviceCamera
		// 运行camera
		deviceCamera.MergeConfig(conf)
		// 将jDevice加入到driver
		d.JDevices[deviceName] = jDevice
	}
}

func (d *Driver) modifyDeviceByConfig(deviceName string, conf map[string]string) {
	d.lc.Info(fmt.Sprint("Modifying device: ", deviceName))

	// 更新摄像头配置
	if d.JDevices[deviceName].Camera != nil {
		d.JDevices[deviceName].Camera.MergeConfig(conf)
	}

	// 更新控制协议配置
	if d.JDevices[deviceName].Control != nil {
		d.JDevices[deviceName].Control.MergeConfig(conf)
	}
}
