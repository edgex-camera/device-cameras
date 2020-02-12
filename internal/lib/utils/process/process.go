package process

// reference: https://github.com/ochinchina/supervisord/blob/33bc01bd202799917615b06c74b8136ef9e4df80/process/process.go
// TODO: add stdout stderr log

import (
	"fmt"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

type RestartPolicy int

const (
	RestartPolicyAlways RestartPolicy = iota
	RestartPolicyNo
)

type Process interface {
	Start()
	Stop(timeout time.Duration, wait bool) error
}

type process struct {
	lc            logger.LoggingClient
	getCmd        func() *exec.Cmd
	cmd           *exec.Cmd
	restartPolicy RestartPolicy
	inStart       bool
	stopByUser    bool
	lock          sync.Mutex

	stdoutFilename string
	stderrFilename string
}

func NewProcess(lc logger.LoggingClient, getCmd func() *exec.Cmd, restartPolicy RestartPolicy, stdoutFilename string, stderrFilename string) Process {
	return &process{
		lc:             lc,
		getCmd:         getCmd,
		restartPolicy:  restartPolicy,
		stdoutFilename: stdoutFilename,
		stderrFilename: stderrFilename,
	}
}

func (p *process) Start() {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.inStart {
		p.lc.Info("already starting")
		return
	}

	p.inStart = true
	p.stopByUser = false

	go func() {
		for {
			p.cmd = p.getCmd()
			stdout, _ := p.cmd.StdoutPipe()
			stderr, _ := p.cmd.StderrPipe()

			p.lc.Info(fmt.Sprintf("process starting: %v", p.cmd.Args))

			start := time.Now()
			// TODO: sigkill cannot be caught by exec.cmd.Run()
			err := p.cmd.Start()
			if err != nil {
				p.lc.Error(fmt.Sprintf("cannot start process %v with error: %v", p.cmd.Args, err.Error()))
				continue
			}

			err = p.handleOutputs(stdout, stderr)
			if err != nil {
				p.lc.Error(fmt.Sprintf("Failed to handle output: %v", err.Error()))
			}
			err = p.cmd.Wait()
			if err != nil {
				p.lc.Error(fmt.Sprintf("process %v stopped with error: %v", p.cmd.Args[0], err.Error()))
			}

			elapsed := time.Since(start)
			if elapsed < 3*time.Second {
				// prevent restarting use 100% cpu
				time.Sleep(3*time.Second - elapsed)
			}
			if p.stopByUser || p.restartPolicy == RestartPolicyNo {
				break
			}
		}

		p.lock.Lock()
		p.inStart = false
		p.lock.Unlock()
	}()
}

func (p *process) Stop(timeout time.Duration, wait bool) error {
	if !p.inStart || p.cmd == nil || p.cmd.Process == nil {
		return fmt.Errorf("process no started")
	}

	p.lock.Lock()
	p.stopByUser = true
	p.lock.Unlock()

	var stopped = sync.WaitGroup{}
	stopped.Add(1)

	go func() {
		defer stopped.Done()
		err := p.cmd.Process.Signal(syscall.SIGINT)
		if err != nil {
			p.lc.Error(err.Error())
			return
		}

		endTime := time.Now().Add(timeout)
		for endTime.After(time.Now()) {
			if !p.inStart {
				return
			}
			time.Sleep(200 * time.Millisecond)
		}

		// TODO: sigkill cannot be caught by exec.cmd.Run()
		err = p.cmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			p.lc.Error(err.Error())
			return
		}

		p.lock.Lock()
		p.inStart = false
		p.lock.Unlock()
	}()

	if wait {
		stopped.Wait()
	}
	return nil
}
