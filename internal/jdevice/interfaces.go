package jdevice

import (
	"os"
	"time"

	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/onvif"
)

// Camera stores config and provide REST interface
type Camera interface {
	Enable()           // 开启摄像头
	Disable(wait bool) // 停止摄像头
	IsEnabled() bool   // 查看摄像头是否正在工作

	CapturePhotoJPG() (*os.File, error) // 取实时截图文件，很少用
	GetCapturePath() string             // 取实时截图文件地址
	GetImagePaths() []string            // 取定时存图地址列表
	GetVideoPaths() []string            // 取定时存视频地址列表
	GetStreamAddr() string              // 取推流地址

	// Json config
	MergeConfig(configPatch []byte) error // 整合新配置并重启摄像头
	GetConfigure() []byte                 // 获取当前摄像头配置

	// For channel manage
	AddChannel() error                    // 当前摄像头增加channel
	RemoveChannel(channelId string) error // 当前摄像头移除channel
	ListChannels() []string               // 列出当前摄像头所有channel列表
}

type Onvif interface {
	MergeConfig(configPatch []byte) error // 整合新配置

	ContinuousMove(time time.Duration, move onvif.Move) error // 移动摄像头
	Stop() error                                              // 停止摄像头

	SetHomePosition() error // 设置原点
	Reset() error           // 回到原点

	GetPresets() string            // 获取预置位信息
	SetPreset(number int64) error  // 设置预置位
	GotoPreset(number int64) error // 移动到预置位

	SyncTime() error // 与机器时间同步
}
