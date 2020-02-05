package onvif

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/yakovlevdmv/goonvif"
	"github.com/yakovlevdmv/goonvif/Device"
	"github.com/yakovlevdmv/goonvif/PTZ"
	"github.com/yakovlevdmv/goonvif/xsd"
	"github.com/yakovlevdmv/goonvif/xsd/onvif"
)

type onvifDevice interface {
	CallMethod(method interface{}) (*http.Response, error)
}

type onvifCamera struct {
	device       onvifDevice
	lc           logger.LoggingClient
	OnvifConfig  OnvifConfig
	stopTimer    *time.Timer
	mutex        *sync.Mutex
	profileToken onvif.ReferenceToken
}

func NewOnvif(lc logger.LoggingClient, config OnvifConfig) (cam Onvif, err error) {
	return &onvifCamera{
		lc:           lc,
		device:       nil,
		OnvifConfig:  config,
		mutex:        &sync.Mutex{},
		profileToken: "",
	}, nil
}

func (c *onvifCamera) connect() (err error) {
	if c.profileToken == "" {
		c.profileToken = getToken(c.OnvifConfig)
	}
	defer func() {
		if r := recover(); r != nil {
			c.lc.Error(fmt.Sprint("Init Onvif camera failed, Recovered in ", r))
			err = fmt.Errorf("Init Onvif camera failed")
		}
	}()

	if c.device != nil { // already connected
		return nil
	}
	device, err := goonvif.NewDevice(c.OnvifConfig.Address)
	if err != nil {
		c.lc.Error("onvif camera connect error: %v", err)
		return err
	}
	device.Authenticate(c.OnvifConfig.Username, c.OnvifConfig.Password)
	c.device = device
	return nil
}

func (c *onvifCamera) callMethod(method interface{}) error {
	err := c.connect()
	if err != nil {
		return err
	}
	_, err = c.device.CallMethod(method)
	if err != nil {
		return err
	}

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(response.Body)
	// c.lc.Info(fmt.Sprintf("onvif callMethod response: %s", buf.String()))
	return nil
}

func (c *onvifCamera) ContinuousMove(timeout time.Duration, moveSpeed Move) error {
	c.lc.Info("camera move started")
	req := PTZ.ContinuousMove{
		ProfileToken: c.profileToken,
		Velocity: onvif.PTZSpeed{
			PanTilt: onvif.Vector2D{
				X:     moveSpeed.PanTiltSpeed.X,
				Y:     moveSpeed.PanTiltSpeed.Y,
				Space: xsd.AnyURI("http://www.onvif.org/ver10/tptz/PanTiltSpaces/GenericSpeedSpace"),
			},
			Zoom: onvif.Vector1D{
				X:     moveSpeed.Zoom,
				Space: xsd.AnyURI("http://www.onvif.org/ver10/tptz/ZoomSpaces/ZoomGenericSpeedSpace"),
			},
		},
		// timeout not working
	}

	c.ensureNoStopTimer()

	err := c.callMethod(req)
	c.mutex.Lock()
	if c.stopTimer == nil {
		c.stopTimer = time.AfterFunc(timeout, func() { _ = c.Stop() })
	}
	c.mutex.Unlock()
	return err
}

func (c *onvifCamera) Stop() error {
	c.ensureNoStopTimer()

	c.lc.Info("camera move stopped")
	req := PTZ.Stop{
		ProfileToken: c.profileToken,
		PanTilt:      true,
		Zoom:         true,
	}

	return c.callMethod(req)
}

func (c *onvifCamera) SetHomePosition() error {
	c.lc.Info("camera move reset")
	req := PTZ.SetPreset{
		ProfileToken: c.profileToken,
		PresetToken:  "1",
	}
	return c.callMethod(req)
}

func (c *onvifCamera) Reset() error {
	c.lc.Info("camera move reset")

	c.ensureNoStopTimer()
	c.Stop()

	req := PTZ.GotoPreset{
		ProfileToken: c.profileToken,
		PresetToken:  "1",
	}
	return c.callMethod(req)
}

func (c *onvifCamera) GetPresets() string {
	c.lc.Info("get presets info")
	return getPresets()
}

func (c *onvifCamera) SetPreset(number int64) error {
	c.lc.Info("set preset", number)
	if number == int64(1) {
		return errors.New("cannot set preset 1, it is home position")
	}
	setPreset(number)
	req := PTZ.SetPreset{
		ProfileToken: c.profileToken,
		PresetToken:  numberToToken(number),
	}
	return c.callMethod(req)
}

func (c *onvifCamera) GotoPreset(number int64) error {
	c.lc.Info("camera move to preset", number)
	c.ensureNoStopTimer()
	req := PTZ.GotoPreset{
		ProfileToken: c.profileToken,
		PresetToken:  numberToToken(number),
	}
	return c.callMethod(req)
}

func (c *onvifCamera) SyncTime() error {
	c.lc.Info("time sync")
	now := time.Now().UTC()
	req := Device.SetSystemDateAndTime{
		TimeZone: onvif.TimeZone{
			TZ: "CST-8",
		},
		UTCDateTime: onvif.DateTime{
			Time: onvif.Time{
				Hour:   xsd.Int(now.Hour()),
				Minute: xsd.Int(now.Minute()),
				Second: xsd.Int(now.Second()),
			},
			Date: onvif.Date{
				Year:  xsd.Int(now.Year()),
				Month: xsd.Int(now.Month()),
				Day:   xsd.Int(now.Day()),
			},
		},
	}
	return c.callMethod(req)
}

// helpers
func (c *onvifCamera) ensureNoStopTimer() {
	c.mutex.Lock()
	if c.stopTimer != nil {
		c.stopTimer.Stop()
		c.stopTimer = nil
	}
	c.mutex.Unlock()
}
