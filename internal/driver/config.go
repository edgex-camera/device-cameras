package driver

import (
	"encoding/json"
	"fmt"

	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice/normalcam"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/camera"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/camera/cmder"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/onvif"
)

func (d *Driver) OnvifConfigChange(oldConf map[string]string, newConf map[string]string) {
	fmt.Println("<<<<<<<<<<<<<<<<<<<", oldConf)
	fmt.Println("<<<<<<<<<<<<<<<<<<<", newConf)
	allDevices, ok := newConf[ALL_DEVICES_KEY]
	// 无任何设备
	if !ok {
		return
	}

	allDevicesMap := make(map[string]bool)
	json.Unmarshal([]byte(allDevices), &allDevicesMap)

	for deviceName := range allDevicesMap {
		_, ok := d.JDevices[deviceName]
		if ok {
			// 目前运行设备配置与旧配置相同，不对设备做操作
			if newConf[deviceName] != oldConf[deviceName] {
				d.addOrModDeviceByConfig(deviceName, newConf[deviceName])
			}
		} else {
			// 目前未运行该设备，根据配置加入
			d.addOrModDeviceByConfig(deviceName, newConf[deviceName])
		}
	}
}

func (d *Driver) addOrModDeviceByConfig(deviceName, conf string) {
	fmt.Println("<<<<<<<<<<<<<<<<<< Conf: ", conf)
	confMap := make(map[string]string)
	json.Unmarshal([]byte(conf), &confMap)

	basicStr, ok := confMap["basic"]
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
	if basicConf.Onvif {
		onvifStr, ok := confMap["onvif"]
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
		jDevice.Onvif = deviceOnvif
	}

	// 创建camera实例
	cameraStr, ok := confMap["camera"]
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
		cameraMap := make(map[string]string)
		json.Unmarshal([]byte(cameraStr), &cameraMap)
		for id := range cameraMap {
			channelId = id
			json.Unmarshal([]byte(cameraMap[id]), &cc)
		}
		deviceCamera, err := normalcam.NewCamera(deviceName, channelId, d.lc, deviceCmder, cc)
		if err != nil {
			d.lc.Info(fmt.Sprint("Failed to create camera: ", err))
			return
		}
		jDevice.Camera = deviceCamera
		// 运行camera
		deviceCamera.MergeConfig([]byte(cameraMap[channelId]))
		// 将jDevice加入到driver
		d.JDevices[deviceName] = jDevice
	}
}
