package camera

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// 获取摄像头支持的分辨率, addr: "/dev/video0"
func GetFrameSizes(addr string) []string {
	cmd_str := fmt.Sprintf("v4l2-ctl --list-framesizes=MJPG -d %v", addr)
	cmd := exec.Command("sh", "-c", cmd_str)
	output, _ := cmd.Output()
	outs := strings.Split(string(output), "\n\tSize: Discrete ")
	return removeDuplicate(outs[1:])
}

// 摄像头是否支持该分辨率，size格式为类似"3264x2448"的字符串
func SupportSize(addr string, size string) bool {
	sizes := GetFrameSizes(addr)
	for i := 0; i < len(sizes); i++ {
		if size == sizes[i] {
			return true
		}
	}
	return false
}

// 摄像头分辨率支持的帧率
func SupportedFps(addr string, width int, height int) []int {
	var res []int

	cmd_str := fmt.Sprintf("v4l2-ctl --list-frameintervals=width=%d,height=%d,pixelformat=MJPG -d %s", width, height, addr)
	cmd := exec.Command("sh", "-c", cmd_str)
	output, _ := cmd.Output()
	outs := strings.Split(string(output), "\n\tInterval: Discrete ")

	for _, i := range outs[1:] {
		strs := strings.Split(i, " ")
		fps_str := strings.TrimPrefix(strs[1], "(")
		fps_str = strings.TrimSuffix(fps_str, ".000")
		fps, _ := strconv.Atoi(fps_str)
		res = append(res, fps)
	}
	return res
}

func removeDuplicate(list []string) []string {
	var res []string
	for _, item := range list {
		item = strings.TrimSuffix(item, "\n")
		if len(res) == 0 {
			res = append(res, item)
		} else {
			for k, v := range res {
				if item == v {
					break
				}
				if k == len(res)-1 {
					res = append(res, item)
				}
			}
		}
	}
	return res
}
