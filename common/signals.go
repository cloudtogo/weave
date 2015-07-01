package common

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// A subsystem/server/... that can be stopped or queried about the status with a signal
type SignalReceiver interface {
	Stop() error
	Status() string
}

func SignalHandlerLoop(ss ...SignalReceiver) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGUSR1)
	buf := make([]byte, 1<<20)
	for {
		switch <-sigs {
		case syscall.SIGINT:
			Log.Infof("=== received SIGINT ===\n*** exiting\n")
			for _, subsystem := range ss {
				subsystem.Stop()
			}
			return
		case syscall.SIGQUIT:
			stacklen := runtime.Stack(buf, true)
			Log.Infof("=== received SIGQUIT ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stacklen])
		case syscall.SIGUSR1:
			for _, subsystem := range ss {
				Log.Infof("=== received SIGUSR1 ===\n*** status...\n%s\n*** end\n", subsystem.Status())
			}
		}
	}
}
