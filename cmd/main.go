package main

import (
	"github.com/womat/debug"
	"os"
	"os/signal"
	"syscall"

	"pv/global"
	"pv/pkg/pv"
)

type pm struct {
	data *pv.Measurements
}

func main() {
	debug.SetDebug(global.Config.Debug.File, global.Config.Debug.Flag)

	global.Measurements = pv.New()
	global.Measurements.SetMeterURL(global.Config.MeterURL)

	runtime := &pm{
		data: global.Measurements,
	}

	go runtime.calcRuntime(global.Config.DataCollectionInterval)

	// capture exit signals to ensure resources are released on exit.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	// wait for am os.Interrupt signal (CTRL C)
	sig := <-quit
	debug.InfoLog.Printf("Got %s signal. Aborting...\n", sig)
	os.Exit(1)
}
