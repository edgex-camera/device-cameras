package camera

import (
	"fmt"
	"os"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice/utils/process"
)

type camera struct {
	lc              logger.LoggingClient
	cmder           CameraCmder
	CameraConfig    CameraConfig
	processes       []process.Process
	videoMaintainer *videoMaintainer
	enabled         bool
}

func NewCamera(lc logger.LoggingClient, cmder CameraCmder, cc CameraConfig) Camera {
	return &camera{
		lc:           lc,
		cmder:        cmder,
		CameraConfig: cc,
		processes:    nil,
		enabled:      false,
	}
}

func (c *camera) Enable() {
	if c.enabled {
		c.lc.Error("camera already enabled")
		return
	}
	c.lc.Error("camera enabled")

	c.enabled = true

	producers := c.cmder.GetCmdProducers(c.CameraConfig)
	if c.processes == nil {
		c.processes = []process.Process{}
		for _, producer := range producers {
			process := process.NewProcess(c.lc, producer, process.RestartPolicyAlways, "", "")
			c.processes = append(c.processes, process)
			process.Start()
		}
	}

	if c.CameraConfig.VideoConfig.Enabled {
		videoMaintainer, err := newVideoMaintainer(c.lc, c.CameraConfig.VideoConfig.Path, c.CameraConfig.VideoConfig.KeepRecord)
		if err != nil {
			c.lc.Error(fmt.Sprintf("video maintainer init failed: %v", err.Error()))
		}
		c.videoMaintainer = videoMaintainer
	}
}

func (c *camera) Disable(wait bool) {
	if !c.enabled {
		c.lc.Error("camera already disabled")
		return
	}
	c.lc.Error("camera disabled")
	c.enabled = false

	if c.processes != nil {
		// TODO: ffmpeg video record cannot be ended gracefully: Failure occurred when ending segment './test/2019-08-07-13-32-26.mp4'
		for _, process := range c.processes {
			err := process.Stop(10*time.Second, wait)
			if err != nil {
				c.lc.Error(err.Error())
			}
		}
		c.processes = nil
	}

	if c.videoMaintainer != nil {
		c.videoMaintainer.stop()
		c.videoMaintainer = nil
	}
}

func (c *camera) IsEnabled() bool {
	return c.enabled
}

func (c *camera) CapturePhotoJPG() (file *os.File, err error) {
	if !c.CameraConfig.CaptureConfig.Enabled {
		return nil, fmt.Errorf("video capture not enabled")
	}

	for i := 0; i < 5; i++ {
		file, err = os.Open(c.CameraConfig.CaptureConfig.Path)
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
		c.lc.Error(fmt.Sprintf("open file failed (%s), retry after 100 milliseconds", err.Error()))
	}
	return file, err
}

func (c *camera) GetCapturePath() string {
	return c.CameraConfig.CaptureConfig.Path
}

func (c *camera) GetVideoPaths() []string {
	return c.videoMaintainer.getFileList()
}

func (c *camera) GetStreamAddr() string {
	return c.CameraConfig.StreamConfig.Address
}
