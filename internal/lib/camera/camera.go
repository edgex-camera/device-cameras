package camera

import (
	"fmt"
	"os"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/utils/process"
)

type RawCamera struct {
	ChannelId       string
	Lc              logger.LoggingClient
	Cmder           CameraCmder
	CameraConfig    CameraConfig
	processes       []process.Process
	videoMaintainer *videoMaintainer
	imageMaintainer *imageMaintainer
	enabled         bool
}

type CameraConfig struct {
	Enabled       bool   `json:"enabled"`
	InputAddr     string `json:"input_addr"`
	Auth          `json:"auth"`
	StreamConfig  `json:"stream"`
	CaptureConfig `json:"capture"`
	VideoConfig   `json:"video"`
	QualityConfig `json:"quality"`
}

func (c *RawCamera) Enable() {
	if c.enabled {
		c.Lc.Error("camera already enabled")
		return
	}
	c.Lc.Error("camera enabled")

	c.enabled = true

	producers := c.Cmder.GetCmdProducers(c.CameraConfig)
	if c.processes == nil {
		c.processes = []process.Process{}
		for _, producer := range producers {
			process := process.NewProcess(c.Lc, producer, process.RestartPolicyAlways, "", "")
			c.processes = append(c.processes, process)
			process.Start()
		}
	}

	if c.CameraConfig.VideoConfig.Enabled {
		videoMaintainer, err := newVideoMaintainer(c.Lc, c.CameraConfig.VideoConfig.Path, c.CameraConfig.VideoConfig.KeepRecord)
		if err != nil {
			c.Lc.Error(fmt.Sprintf("video maintainer init failed: %v", err.Error()))
		}
		c.videoMaintainer = videoMaintainer
	}

	if c.CameraConfig.CaptureConfig.Enabled && c.CameraConfig.CaptureConfig.Storage {
		imageMaintainer, err := newImageMaintainer(c.Lc, c.CameraConfig.CaptureConfig.Path, c.CameraConfig.CaptureConfig.Seconds, c.CameraConfig.CaptureConfig.Number)
		if err != nil {
			c.Lc.Error(fmt.Sprintf("image maintainer init failed: %v", err.Error()))
		}
		c.imageMaintainer = imageMaintainer
		go c.imageMaintainer.start()
	}
}

func (c *RawCamera) Disable(wait bool) {
	if !c.enabled {
		c.Lc.Error("camera already disabled")
		return
	}
	c.Lc.Error("camera disabled")
	c.enabled = false

	if c.processes != nil {
		// TODO: ffmpeg video record cannot be ended gracefully: Failure occurred when ending segment './test/2019-08-07-13-32-26.mp4'
		for _, process := range c.processes {
			err := process.Stop(10*time.Second, wait)
			if err != nil {
				c.Lc.Error(err.Error())
			}
		}
		c.processes = nil
	}

	if c.videoMaintainer != nil {
		c.videoMaintainer.stop()
		c.videoMaintainer = nil
	}
	if c.imageMaintainer != nil {
		c.imageMaintainer.stop()
		c.imageMaintainer = nil
	}
}

func (c *RawCamera) IsEnabled() bool {
	return c.enabled
}

func (c *RawCamera) CapturePhotoJPG() (file *os.File, err error) {
	if !c.CameraConfig.CaptureConfig.Enabled {
		return nil, fmt.Errorf("video capture not enabled")
	}

	for i := 0; i < 5; i++ {
		file, err = os.Open(c.CameraConfig.CaptureConfig.Path)
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
		c.Lc.Error(fmt.Sprintf("open file failed (%s), retry after 100 milliseconds", err.Error()))
	}
	return file, err
}

func (c *RawCamera) GetCapturePath() string {
	return c.CameraConfig.CaptureConfig.Path
}

func (c *RawCamera) GetVideoPaths() []string {
	if c.videoMaintainer == nil {
		return []string{}
	}
	return c.videoMaintainer.getFileList()
}

func (c *RawCamera) GetStreamAddr() string {
	return c.CameraConfig.StreamConfig.Address
}

func (c *RawCamera) GetImagePaths() []string {
	if c.imageMaintainer == nil {
		return []string{}
	}
	return c.imageMaintainer.getFileList()
}
