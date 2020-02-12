package process

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/lib/utils"
)

func (p *process) handleOutputs(stdout io.ReadCloser, stderr io.ReadCloser) error {

	if p.stdoutFilename == "" {
		p.stdoutFilename = fmt.Sprintf("/data/logs/camera/%s.stdout", filepath.Base(p.cmd.Path))
	}

	if p.stderrFilename == "" {
		p.stderrFilename = fmt.Sprintf("/data/logs/camera/%s.stderr", filepath.Base(p.cmd.Path))
	}

	handleOutput := func(pipe io.ReadCloser, filename string) error {
		if err := utils.MakeDirsIfNotExist(filepath.Dir(filename)); err != nil {
			return err
		}
		p.lc.Info(fmt.Sprintf("log writing to %s", filename))
		go pipeToFile(p.lc, pipe, filename)
		return nil
	}

	if err := handleOutput(stdout, p.stdoutFilename); err != nil {
		return err
	}
	if err := handleOutput(stderr, p.stderrFilename); err != nil {
		return err
	}
	return nil
}

func pipeToFile(lc logger.LoggingClient, pipe io.ReadCloser, filename string) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		lc.Error(err.Error())
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, pipe); err != nil {
		lc.Error(err.Error())
	}
}
