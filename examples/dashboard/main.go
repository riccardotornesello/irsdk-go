package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/riccardotornesello/irsdk-go"
)

var sdk *irsdk.IRSDK
var homeTemplate *template.Template

func main() {
	sdk = irsdk.Init(nil)
	defer sdk.Close()

	h, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal(err)
	}
	homeTemplate = h

	flag.Parse()
	log.SetFlags(0)
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", home)
	log.Printf("Listening on %q", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

type data struct {
	IsConnected            bool
	Weather                string
	RPMLights              rpmLights
	EngineWarnings         interface{}
	TrackAirTemp           interface{}
	TrackSurfaceTemp       interface{}
	DisplayUnits           interface{}
	FuelLevel              interface{}
	Speed                  interface{}
	Gear                   interface{}
	LapLastLapTime         interface{}
	LapBestLapTime         interface{}
	SessionTimeRemain      interface{}
	PlayerCarClassPosition interface{}
	RPM                    interface{}
}

type rpmLights struct {
	First string
	Last  string
	Blink string
	Shift string
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			//log.Println("read:", err)
			break
		}
		online := true
		for {
			sdk.Update(true)

			session := sdk.Session
			weather := session.WeekendInfo.TrackSkies

			// driverIdx, err := sdk.GetSessionData("DriverInfo:DriverCarIdx")
			// checkErr(err)
			// incidentCount, err := sdk.GetSessionData("DriverInfo:Drivers:{" + driverIdx + "}CurDriverIncidentCount")
			// checkErr(err)

			rpmL, err := getRPMData(sdk)
			checkErr(err)

			engineWarnings := sdk.Telemetry["EngineWarnings"]
			airTemp := sdk.Telemetry["AirTemp"]
			trackTemp := sdk.Telemetry["TrackTempCrew"]
			units := sdk.Telemetry["DisplayUnits"]
			fuel := sdk.Telemetry["FuelLevel"]
			speed := sdk.Telemetry["Speed"]
			gear := sdk.Telemetry["Gear"]
			lapLastLapTime := sdk.Telemetry["LapLastLapTime"]
			lapBestLapTime := sdk.Telemetry["LapBestLapTime"]
			sessionTimeRemain := sdk.Telemetry["SessionTimeRemain"]
			playerCarClassPosition := sdk.Telemetry["PlayerCarClassPosition"]
			rpm := sdk.Telemetry["RPM"]

			d := data{
				sdk.IsConnected(),
				weather,
				rpmL,
				engineWarnings.Value(),
				airTemp.Value(),
				trackTemp.Value(),
				units.Value(),
				fuel.Value(),
				speed.Value(),
				gear.Value(),
				lapLastLapTime.Value(),
				lapBestLapTime.Value(),
				sessionTimeRemain.Value(),
				playerCarClassPosition.Value(),
				rpm.Value(),
			}
			message, err = json.Marshal(d)
			if err != nil {
				log.Println("error json: ", err)
				break
			}
			err = c.WriteMessage(mt, message)
			if err != nil {
				//log.Println("error write: ", err)
				break
			}
			if sdk.IsConnected() {
				time.Sleep(50 * time.Millisecond)
				if !online {
					log.Println("iRacing connected!")
				}
				online = true
			} else {
				time.Sleep(5 * time.Second)
				if online {
					log.Println("Waiting for iRacing connection...")
				}
				online = false
			}
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/ws")
}

func checkErr(err error) {
	if err != nil {
		//log.Println(err)
	}
}

func getRPMData(sdl *irsdk.IRSDK) (rpmLights, error) {
	session := sdl.Session

	first := fmt.Sprintf("%d", session.DriverInfo.DriverCarSLFirstRPM)
	last := fmt.Sprintf("%d", session.DriverInfo.DriverCarSLLastRPM)
	blink := fmt.Sprintf("%d", session.DriverInfo.DriverCarSLBlinkRPM)
	shift := fmt.Sprintf("%d", session.DriverInfo.DriverCarSLShiftRPM)

	return rpmLights{first, last, blink, shift}, nil
}
