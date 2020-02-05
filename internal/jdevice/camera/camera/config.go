package camera

import (
	"encoding/json"
	"path/filepath"

	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/jdevice/utils"
	jsonpatch "gopkg.in/evanphx/json-patch.v4"
)

type StreamConfig struct {
	Enabled bool   `json:"enabled"`
	Address string `json:"addr"`
}

type CaptureConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
	Rate    int    `json:"rate"`
}

type VideoConfig struct {
	Enabled    bool   `json:"enabled"`
	Path       string `json:"path"`
	Length     int    `json:"length"`
	KeepRecord int    `json:"keep"`
}

// 画质选项
type QualityConfig struct {
	ImageHeight int `json:"image_height"`
	ImageWidth  int `json:"image_width"`
	VideoHeight int `json:"video_height"`
	VideoWidth  int `json:"video_width"`
	FrameRate   int `json:"frame_rate"`
}

// 用户认证(主要用于rtsp摄像头)
type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

var defaultConf = CameraConfig{
	Enabled:       true,
	StreamConfig:  StreamConfig{Enabled: true},
	CaptureConfig: CaptureConfig{Enabled: true},
	VideoConfig: VideoConfig{
		Enabled:    true,
		Length:     600,
		KeepRecord: 3,
	},
}

var DefaultConf []byte

func init() {
	DefaultConf, _ = json.Marshal(defaultConf)
}

func (c *camera) MergeConfig(configPatch []byte) error {
	if c.IsEnabled() {
		c.Disable(true)
	}

	old, err := json.Marshal(c.CameraConfig)
	if err != nil {
		return err
	}
	new, err := jsonpatch.MergePatch(old, configPatch)
	if err != nil {
		return err
	}
	err = json.Unmarshal(new, &c.CameraConfig)
	if err != nil {
		return err
	}

	err = utils.MakeDirsIfNotExist(filepath.Dir(c.CameraConfig.CaptureConfig.Path))
	if err != nil {
		return err
	}

	err = utils.MakeDirsIfNotExist(filepath.Dir(c.CameraConfig.VideoConfig.Path))
	if err != nil {
		return err
	}

	if c.CameraConfig.Enabled {
		c.Enable()
	}
	return nil
}

func (c *camera) GetConfigure() []byte {
	config, _ := json.Marshal(c.CameraConfig)
	return config
}
