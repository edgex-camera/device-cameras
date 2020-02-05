package camera

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/fsnotify/fsnotify"
)

type videoMaintainer struct {
	lc         logger.LoggingClient
	keepRecord int // ignored if <=0
	fileList   *list.List
	watcher    *fsnotify.Watcher
	stopped    chan interface{}
	path       string
	ext        string
}

func newVideoMaintainer(lc logger.LoggingClient, path string, keepRecord int) (*videoMaintainer, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	vm := &videoMaintainer{
		lc:         lc,
		keepRecord: keepRecord,
		fileList:   list.New(),
		watcher:    watcher,
		stopped:    make(chan interface{}),
		path:       path,
		ext:        filepath.Ext(path),
	}
	err = vm.init()
	if err != nil {
		return nil, err
	}
	return vm, nil
}

func (vm *videoMaintainer) init() error {
	err := vm.addExistFiles()
	if err != nil {
		return err
	}

	// setup watcher
	err = vm.watcher.Add(filepath.Dir(vm.path))
	if err != nil {
		return err
	}
	go vm.loop()

	return nil
}

func (vm *videoMaintainer) addExistFiles() error {
	files, err := ioutil.ReadDir(filepath.Dir(vm.path))
	if err != nil {
		return err
	}
	// sort by increasing modified time
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	for _, file := range files {
		if filepath.Ext(file.Name()) == vm.ext {
			fullPath := filepath.Join(filepath.Dir(vm.path), file.Name())
			vm.addFile(fullPath)
		}
	}
	return nil
}

func (vm *videoMaintainer) stop() {
	vm.stopped <- struct{}{}
	vm.watcher.Close()
}

func (vm *videoMaintainer) getFileList() []string {
	len := vm.fileList.Len()
	if len < 2 {
		return []string{}
	}

	files := make([]string, 0, len-1)
	// skip last element because it is processing
	for e := vm.fileList.Back().Prev(); e != nil; e = e.Prev() {
		absPath, err := filepath.Abs(e.Value.(string))
		if err != nil {
			vm.lc.Error(err.Error())
		} else {
			files = append(files, absPath)
		}
	}
	return files
}

func (vm *videoMaintainer) loop() {
	for {
		// vm.lc.Debug("videoMaintainer loop")
		select {
		case event, ok := <-vm.watcher.Events:
			if !ok {
				return
			}
			// vm.lc.Debug(fmt.Sprint("event:", event))
			if event.Op&fsnotify.Create == fsnotify.Create && filepath.Ext(event.Name) == vm.ext {
				vm.lc.Debug(fmt.Sprint("created:", event.Name))
				vm.addFile(event.Name)
			}
		case err, ok := <-vm.watcher.Errors:
			if !ok {
				return
			}
			if err != nil {
				vm.lc.Error(err.Error())
			}
		case <-vm.stopped:
			return
		}
	}
}

// push file to back, if len > keepRecord +1, remove front
// +1 because the last video is processing
func (vm *videoMaintainer) addFile(fileName string) {
	vm.fileList.PushBack(fileName)
	if vm.keepRecord > 0 && vm.fileList.Len() > vm.keepRecord+1 {
		front := vm.fileList.Front()
		vm.fileList.Remove(front)
		err := os.Remove(front.Value.(string))
		if err != nil {
			vm.lc.Error(err.Error())
		}
	}
}
