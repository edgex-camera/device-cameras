package onvif

import (
	"testing"
	"time"

	"gitlab.jiangxingai.com/applications/edgex/edgex-utils/logger"
)

var lc = logger.NewPrintClient()
var address = "10.54.128.132:8899"

// var rtspAddr = "rtsp://10.54.128.132:554/mpeg4"

func setup() (Onvif, error) {
	onvif, err := NewOnvif(lc, address)
	return onvif, err
}

func TestContinuousMove(t *testing.T) {
	camera, err := setup()
	if err != nil {
		t.Error(err)
	}

	err = camera.ContinuousMove(1000*time.Millisecond, Move{Vector2D{X: 1, Y: 1}, 0})
	if err != nil {
		t.Error(err)
	}

	time.Sleep(500 * time.Millisecond)

	err = camera.ContinuousMove(1000*time.Millisecond, Move{Vector2D{X: -1, Y: -1}, 0})
	if err != nil {
		t.Error(err)
	}

	time.Sleep(100 * time.Second)

	camera.Stop()
}

func TestSetHomePosition(t *testing.T) {
	camera, err := setup()
	if err != nil {
		t.Error(err)
	}

	camera.SetHomePosition()
}

func TestReset(t *testing.T) {
	camera, err := setup()
	if err != nil {
		t.Error(err)
	}

	camera.Reset()
}
