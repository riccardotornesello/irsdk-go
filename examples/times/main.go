package main

import (
	"fmt"

	"github.com/riccardotornesello/irsdk-go"
)

func main() {
	sdk := irsdk.Init(nil)
	defer sdk.Close()

	lapTimes, err := sdk.GetVar("CarIdxLastLapTime")
	if err != nil {
		fmt.Println(err)
		return
	}

	session := sdk.GetSession()

	for i := range session.DriverInfo.Drivers {
		fmt.Println(session.DriverInfo.Drivers[i].CarNumber, irsdk.FloatToTimeStr(lapTimes.Value.([]float32)[i]))
	}
}
