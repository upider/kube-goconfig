package pkg

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Control interface {
	WaitStop()
	Stop()
}

type SignalController struct {
	SignalChan chan os.Signal `json:"signalChan"`
	ExitTime   time.Duration  `json:"exitTime"`
}

func (ctrl *SignalController) WaitStop() {
	ctrl.SignalChan = make(chan os.Signal, 1)
	signal.Notify(ctrl.SignalChan, os.Interrupt, syscall.SIGTERM)
	<-ctrl.SignalChan
	time.Sleep(ctrl.ExitTime)
}

func (ctrl *SignalController) Stop() {
	ctrl.SignalChan <- syscall.SIGTERM
}

func NewSignalController(exitTime time.Duration) *SignalController {
	return &SignalController{
		ExitTime: exitTime,
	}
}
