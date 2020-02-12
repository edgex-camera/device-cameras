package process

import (
	"fmt"
	"os/exec"
	"testing"
	"time"

	"gitlab.jiangxingai.com/applications/edgex/edgex-utils/logger"
)

var lc = logger.NewPrintClient()

func getSleepCmd(t float32) func() *exec.Cmd {
	return func() *exec.Cmd {
		return exec.Command("sleep", fmt.Sprint(t))
	}
}

func TestNoRestart(t *testing.T) {
	p := NewProcess(lc, getSleepCmd(0.5), RestartPolicyNo)
	p.Start()
	if !p.inStart {
		t.Error("process not running")
	}

	time.Sleep(time.Duration(100 * time.Millisecond))
	err := p.Stop(5*time.Second, true)
	if err != nil {
		t.Error(err)
	}

	if p.inStart {
		t.Error("process not stopped")
	}
}

func TestRestart(t *testing.T) {
	p := NewProcess(lc, getSleepCmd(0.5), RestartPolicyAlways)
	p.Start()
	if !p.inStart {
		t.Error("process not running")
	}

	time.Sleep(time.Duration(1000 * time.Millisecond))
	if !p.inStart {
		t.Error("process not restarted")
	}

	err := p.Stop(5*time.Second, true)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Duration(200 * time.Millisecond))
	if p.inStart {
		t.Error("process not stopped")
	}
}

func getTopCmd() *exec.Cmd {
	return exec.Command("bash", "./test_scripts/ignore_signal.sh")
}

func TestKill(t *testing.T) {
	p := NewProcess(lc, getTopCmd, RestartPolicyAlways)
	p.Start()
	if !p.inStart {
		t.Error("process not running")
	}

	time.Sleep(time.Duration(200 * time.Millisecond))
	err := p.Stop(2*time.Second, false)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Duration(200 * time.Millisecond))
	if !p.inStart {
		t.Error("process should not respond to interrupt")
	}

	time.Sleep(time.Duration(2000 * time.Millisecond))

	if p.inStart {
		t.Error("process not force killed")
	}
}
