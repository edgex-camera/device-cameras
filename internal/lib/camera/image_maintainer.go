package camera

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

type imageMaintainer struct {
	lc        logger.LoggingClient
	seconds   int
	number    int
	path      string
	ext       string
	prefix    string
	dir       string
	imageList []string
	enabled   bool
}

func newImageMaintainer(lc logger.LoggingClient, path string, seconds, number int) (*imageMaintainer, error) {
	var imageList []string
	im := &imageMaintainer{
		lc:        lc,
		seconds:   seconds,
		number:    number,
		path:      path,
		ext:       filepath.Ext(path),
		prefix:    strings.TrimSuffix(path, filepath.Ext(path)),
		dir:       filepath.Dir(path),
		imageList: imageList,
	}
	err := im.init()
	if err != nil {
		return nil, err
	}
	return im, nil
}

func (im *imageMaintainer) init() error {
	files, err := ioutil.ReadDir(im.dir)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})
	fmt.Println(files)

	for _, file := range files {
		path := im.dir + "/" + file.Name()
		if strings.HasPrefix(path, im.prefix) && path != im.path {
			im.imageList = append(im.imageList, path)
			if len(im.imageList) > im.number {
				os.Remove(im.imageList[0])
				im.imageList = im.imageList[1:]
			}
		} else {
			os.Remove(path)
		}
	}
	im.enabled = true
	return nil
}

func (im *imageMaintainer) start() error {
	for {
		err := im.loop()
		if err != nil {
			im.lc.Error(fmt.Sprint("storing images error: ", err))
		}
		time.Sleep(time.Duration(im.seconds) * time.Second)
	}
	return nil
}

func (im *imageMaintainer) loop() error {
	if !im.enabled {
		return nil
	}
	from, err := os.Open(im.path)
	if err != nil {
		return err
	}
	defer from.Close()

	toPath := fmt.Sprintf("%v%d%v", strings.TrimSuffix(im.path, im.ext), time.Now().Unix(), im.ext)
	to, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()
	io.Copy(to, from)

	im.imageList = append(im.imageList, toPath)
	if len(im.imageList) > im.number {
		os.Remove(im.imageList[0])
		im.imageList = im.imageList[1:]
	}
	return nil
}

func (im *imageMaintainer) stop() {
	im.enabled = false
}

func (im *imageMaintainer) getFileList() []string {
	return im.imageList
}
