package onvif

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/yakovlevdmv/goonvif"
	"github.com/yakovlevdmv/goonvif/Device"
	"github.com/yakovlevdmv/goonvif/PTZ"
	"github.com/yakovlevdmv/goonvif/xsd"
	"github.com/yakovlevdmv/goonvif/xsd/onvif"
)

type OnvifDevice interface {
	CallMethod(method interface{}) (*http.Response, error)
}

func (c *OnvifCamera) connect() (err error) {
	if c.profileToken == "" {
		c.profileToken = getToken(c.OnvifConfig)
	}
	defer func() {
		if r := recover(); r != nil {
			c.Lc.Error(fmt.Sprint("Init Onvif camera failed, Recovered in ", r))
			err = fmt.Errorf("Init Onvif camera failed")
		}
	}()

	if c.device != nil { // already connected
		return nil
	}
	device, err := goonvif.NewDevice(c.OnvifConfig.Address)
	if err != nil {
		c.Lc.Error("onvif camera connect error: %v", err)
		return err
	}
	device.Authenticate(c.OnvifConfig.Username, c.OnvifConfig.Password)
	c.device = device
	return nil
}

func (c *OnvifCamera) callMethod(method interface{}) error {
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
	// c.Lc.Info(fmt.Sprintf("onvif callMethod response: %s", buf.String()))
	return nil
}

func (c *OnvifCamera) ContinuousMove(timeout time.Duration, moveSpeed Move) error {
	c.Lc.Info("camera move started")
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

func (c *OnvifCamera) Stop() error {
	c.ensureNoStopTimer()

	c.Lc.Info("camera move stopped")
	req := PTZ.Stop{
		ProfileToken: c.profileToken,
		PanTilt:      true,
		Zoom:         true,
	}

	return c.callMethod(req)
}

func (c *OnvifCamera) SetHomePosition() error {
	c.Lc.Info("camera move reset")
	req := PTZ.SetPreset{
		ProfileToken: c.profileToken,
		PresetToken:  "1",
	}
	return c.callMethod(req)
}

func (c *OnvifCamera) Reset() error {
	c.Lc.Info("camera move reset")

	c.ensureNoStopTimer()
	c.Stop()

	req := PTZ.GotoPreset{
		ProfileToken: c.profileToken,
		PresetToken:  "1",
	}
	return c.callMethod(req)
}

func (c *OnvifCamera) GetPresets() string {
	c.Lc.Info("get presets info")
	return getPresets()
}

func (c *OnvifCamera) SetPreset(number int64) error {
	c.Lc.Info("set preset", number)
	if number == int64(1) {
		return errors.New("cannot set preset 1, it is home position")
	}
	setPreset(c.Name, number)
	req := PTZ.SetPreset{
		ProfileToken: c.profileToken,
		PresetToken:  numberToToken(number),
	}
	return c.callMethod(req)
}

func (c *OnvifCamera) GotoPreset(number int64) error {
	c.Lc.Info("camera move to preset", number)
	c.ensureNoStopTimer()
	req := PTZ.GotoPreset{
		ProfileToken: c.profileToken,
		PresetToken:  numberToToken(number),
	}
	return c.callMethod(req)
}

func (c *OnvifCamera) SyncTime() error {
	c.Lc.Info("time sync")
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
func (c *OnvifCamera) ensureNoStopTimer() {
	c.mutex.Lock()
	if c.stopTimer != nil {
		c.stopTimer.Stop()
		c.stopTimer = nil
	}
	c.mutex.Unlock()
}
