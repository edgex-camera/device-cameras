package cmder

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/camera"
)

type cmder struct {
	template              processCmdTemplate
	videoLengthMultiplier int
}

func NewCmder(processMethod string) camera.CameraCmder {
	switch {
	case processMethod == "gst-launch-1.0":
		return &cmder{
			template:              gstreamer,
			videoLengthMultiplier: 1000000000,
		}
	case processMethod == "ffmpeg":
		return &cmder{
			template:              ffmpeg,
			videoLengthMultiplier: 1,
		}
	case isGstAvail():
		return &cmder{
			template:              gstreamer,
			videoLengthMultiplier: 1000000000,
		}
	case isFFmpegAvail():
		return &cmder{
			template:              ffmpeg,
			videoLengthMultiplier: 1,
		}
	default:
		panic("no supported video processor")
	}
}

func isGstAvail() bool {
	_, err := exec.LookPath("gst-launch-1.0")
	return err == nil
}

func isFFmpegAvail() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

func (c *cmder) GetCmdProducers(cc camera.CameraConfig) []func() *exec.Cmd {
	isGst := c.template.processor == "gst-launch-1.0"

	var template cmdTemplate
	var inputAddr string
	isWebcam := false
	switch {
	case strings.HasPrefix(cc.InputAddr, "rtsp://"):
		template = c.template.rtsp
		inputAddr = buildInputAddr(cc.InputAddr, cc.Auth)
	case strings.HasPrefix(cc.InputAddr, "/"):
		template = c.template.webcam
		inputAddr = cc.InputAddr
		isWebcam = true
	default:
		log.Printf("input not supported: %s", cc.InputAddr)
		return []func() *exec.Cmd{}
	}

	outputCapture := cc.CaptureConfig.Enabled && cc.CaptureConfig.Path != ""
	outputStream := cc.StreamConfig.Enabled && cc.StreamConfig.Address != ""
	outputVideo := cc.VideoConfig.Enabled && cc.VideoConfig.Path != ""
	if !outputCapture && !outputStream && !outputVideo {
		return []func() *exec.Cmd{}
	}
	cmdStr := ""
	if isWebcam && !isGst {
		cmdStr += fmt.Sprintf(template.base, cc.QualityConfig.ImageWidth, cc.QualityConfig.ImageHeight, inputAddr)
	} else {
		cmdStr += fmt.Sprintf(template.base, inputAddr)
	}
	if isWebcam && isGst {
		cmdStr += fmt.Sprintf(template.quality, cc.QualityConfig.ImageWidth, cc.QualityConfig.ImageHeight, cc.QualityConfig.FrameRate)
	}
	if outputCapture && isGst {
		cmdStr += fmt.Sprintf(template.capture, cc.CaptureConfig.Rate, cc.CaptureConfig.Path)
	}
	if outputStream || outputVideo {
		if isWebcam {
			cmdStr += fmt.Sprintf(template.h264, cc.QualityConfig.VideoWidth, cc.QualityConfig.VideoHeight)
		}
	}
	if outputStream {
		cmdStr += fmt.Sprintf(template.stream, cc.StreamConfig.Address)
	}
	if outputVideo {
		cmdStr += fmt.Sprintf(template.video, cc.VideoConfig.Length*c.videoLengthMultiplier, cc.VideoConfig.Path)
	}
	if outputCapture && !isGst {
		cmdStr += fmt.Sprintf(template.capture, cc.CaptureConfig.Rate, cc.CaptureConfig.Path)
	}

	cmdList := strings.Fields(cmdStr)
	func1 := func() *exec.Cmd { return exec.Command(c.template.processor, cmdList...) }
	return []func() *exec.Cmd{func1}
}

func buildInputAddr(inputAddr string, auth camera.Auth) string {
	prefix := "rtsp://"
	if auth.Username == "" {
		return inputAddr
	}
	suffix := strings.TrimPrefix(inputAddr, prefix)
	fullAddr := fmt.Sprintf("%s%s:%s@%s", prefix, auth.Username, auth.Password, suffix)
	return fullAddr
}
