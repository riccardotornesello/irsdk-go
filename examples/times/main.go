package main

import (
	"fmt"

	"github.com/riccardotornesello/irsdk-go"
)

func main() {
	sdk := irsdk.Init(nil)
	defer sdk.Close()

	lapTimes := sdk.Telemetry["CarIdxLastLapTime"]

	session := sdk.Session

	for i := range session.DriverInfo.Drivers {
		fmt.Println(session.DriverInfo.Drivers[i].CarNumber, irsdk.FloatToTimeStr(lapTimes.Array().([]float32)[i]))
	}
}
