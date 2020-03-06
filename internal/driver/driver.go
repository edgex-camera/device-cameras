package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

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
	d.lc.Info(fmt.Sprint("Read command of device: ", deviceName))
	now := time.Now().UnixNano()
	for _, req := range reqs {
		switch req.DeviceResourceName {
		case "capture_path":
			{
				capturePath := d.JDevices[deviceName].Camera.GetCapturePath()
				capturePathJson, _ := json.Marshal(map[string]interface{}{"capture_path": capturePath})
				cv := dsModels.NewStringValue(reqs[0].DeviceResourceName, now, string(capturePathJson))
				res = append(res, cv)
			}
		case "stream_addr":
			{
				streamAddr := d.JDevices[deviceName].Camera.GetStreamAddr()
				streamAddrJson, _ := json.Marshal(map[string]interface{}{"stream_addr": streamAddr})
				cv := dsModels.NewStringValue(reqs[0].DeviceResourceName, now, string(streamAddrJson))
				res = append(res, cv)
			}
		case "video_paths":
			{
				videoPaths := d.JDevices[deviceName].Camera.GetVideoPaths()
				videoPathsJson, _ := json.Marshal(map[string]interface{}{"video_paths": videoPaths})
				cv := dsModels.NewStringValue(reqs[0].DeviceResourceName, now, string(videoPathsJson))
				res = append(res, cv)
			}
		case "image_paths":
			{
				imagePaths := d.JDevices[deviceName].Camera.GetImagePaths()
				imagePathsJson, _ := json.Marshal(map[string]interface{}{"image_paths": imagePaths})
				cv := dsModels.NewStringValue(reqs[0].DeviceResourceName, now, string(imagePathsJson))
				res = append(res, cv)
			}
		case "config":
			{
				cv := dsModels.NewStringValue(reqs[0].DeviceResourceName, now, string(d.JDevices[deviceName].Camera.GetConfigure()))
				res = append(res, cv)
			}
		case "channels":
			{
				channels := d.JDevices[deviceName].Camera.ListChannels()
				channelsJson, _ := json.Marshal(map[string]interface{}{"channels": channels})
				cv := dsModels.NewStringValue(reqs[0].DeviceResourceName, now, string(channelsJson))
				res = append(res, cv)
			}
		case "presets":
			{
				presets := d.JDevices[deviceName].Control.GetPresets()
				cv := dsModels.NewStringValue(reqs[0].DeviceResourceName, now, presets)
				res = append(res, cv)
			}
		}
	}
	return res, err
}

func (d *Driver) HandleWriteCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest, params []*dsModels.CommandValue) error {
	d.lc.Info(fmt.Sprintf("HandleWriteCommands %s", params[0].DeviceResourceName))
	start := time.Now()

	if deviceName == CAMERA_FACTORY {
		// 设备增删
		for _, param := range params {
			switch param.DeviceResourceName {
			case "add_device", "device_type":
				var name, deviceType string
				for _, param := range params {
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
					}
				}
				return d.AddJdevice(name, deviceType)
			case "remove_device":
				v, err := param.StringValue()
				if err != nil {
					return err
				}
				return d.RemoveJdevice(v)
			}
		}
	} else {
		// 摄像头操作
		for _, param := range params {
			switch param.DeviceResourceName {
			case "config":
				{
					v, err := param.StringValue()
					if err != nil {
						return err
					}
					return d.JDevices[deviceName].Camera.PutConfig([]byte(v))
				}
			default:
				if d.JDevices[deviceName].Control == nil {
					return errors.New("Current device does not support control protocols")
				}
				moveHandled := false

				var err error
				switch param.DeviceResourceName {
				case "pan", "tilt", "zoom", "timeout":
					if !moveHandled {
						moveHandled = true
						err = d.HandleMoveCommand(deviceName, params)
					}
				case "stop":
					err = d.JDevices[deviceName].Control.Stop()
				case "set_home_position":
					err = d.JDevices[deviceName].Control.SetHomePosition()
				case "reset_position":
					err = d.JDevices[deviceName].Control.Reset()
				case "set_preset":
					{
						v, err := param.StringValue()
						vint, err := strconv.ParseInt(v, 10, 64)
						if err != nil {
							return err
						}
						err = d.JDevices[deviceName].Control.SetPreset(vint)
					}
				case "goto_preset":
					{
						v, err := param.StringValue()
						vint, err := strconv.ParseInt(v, 10, 64)
						if err != nil {
							return err
						}
						err = d.JDevices[deviceName].Control.GotoPreset(vint)
					}
				}
				if err != nil {
					return err
				}
			}
		}
	}
	elapsed := time.Since(start)
	d.lc.Info(fmt.Sprintf("HandleWriteCommands took %s", elapsed))
	return nil
}

func (d *Driver) Stop(force bool) error {
	return nil
}
