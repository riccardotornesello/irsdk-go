package irsdk

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/hidez8891/shm"
	"github.com/riccardotornesello/irsdk-go/lib/winevents"
)

// IRSDK is the main SDK object clients must use
type IRSDK struct {
	reader        reader
	header        *header
	session       Session
	rawSession    []string
	telemetry     *TelemetryVars
	lastValidData int64
}

func (sdk *IRSDK) WaitForData(timeout time.Duration) bool {
	if !sdk.IsConnected() {
		initIRSDK(sdk)
	}
	if winevents.WaitForSingleObject(timeout) {
		return readVariableValues(sdk)
	}
	return false
}

func (sdk *IRSDK) GetVar(name string) (variable, error) {
	if !sessionStatusOK(sdk.header.status) {
		return variable{}, fmt.Errorf("Session is not active")
	}
	sdk.telemetry.mux.Lock()
	if v, ok := sdk.telemetry.vars[name]; ok {
		sdk.telemetry.mux.Unlock()
		return v, nil
	}
	sdk.telemetry.mux.Unlock()
	return variable{}, fmt.Errorf("Telemetry variable %q not found", name)
}

func (sdk *IRSDK) GetSession() Session {
	return sdk.session
}

func (sdk *IRSDK) GetLastVersion() int {
	if !sessionStatusOK(sdk.header.status) {
		return -1
	}
	sdk.telemetry.mux.Lock()
	last := sdk.telemetry.lastVersion
	sdk.telemetry.mux.Unlock()
	return last
}

func (sdk *IRSDK) GetSessionData(path string) (string, error) {
	if !sessionStatusOK(sdk.header.status) {
		return "", fmt.Errorf("Session not connected")
	}
	return getSessionDataPath(sdk.rawSession, path)
}

func (sdk *IRSDK) IsConnected() bool {
	if sdk.header != nil {
		if sessionStatusOK(sdk.header.status) && (sdk.lastValidData+connTimeout > time.Now().Unix()) {
			return true
		}
	}

	return false
}

// ExportTo exports current memory data to a file
func (sdk *IRSDK) ExportIbtTo(fileName string) error {
	rbuf := make([]byte, fileMapSize)
	_, err := sdk.reader.ReadAt(rbuf, 0)
	if err != nil {
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(rbuf)
	if err != nil {
		return err
	}

	return nil
}

// ExportTo exports current session yaml data to a file
func (sdk *IRSDK) ExportSessionTo(fileName string) error {
	y := strings.Join(sdk.rawSession, "\n")

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(y))
	if err != nil {
		return err
	}

	return nil
}

func (sdk *IRSDK) BroadcastMsg(msg Msg) {
	if msg.P2 == nil {
		msg.P2 = 0
	}
	winevents.BroadcastMsg(broadcastMsgName, msg.Cmd, msg.P1, msg.P2, msg.P3)
}

// Close clean up sdk resources
func (sdk *IRSDK) Close() {
	sdk.reader.Close()
}

// Init creates a SDK instance to operate with
func Init(r reader) IRSDK {
	if r == nil {
		var err error
		r, err = shm.Open(fileMapName, fileMapSize)
		if err != nil {
			log.Fatal(err)
		}
	}

	sdk := IRSDK{reader: r, lastValidData: 0}
	winevents.OpenEvent(dataValidEventName)
	initIRSDK(&sdk)
	return sdk
}

func initIRSDK(sdk *IRSDK) {
	h := readHeader(sdk.reader)
	sdk.header = &h
	sdk.rawSession = nil
	if sdk.telemetry != nil {
		sdk.telemetry.vars = nil
	}
	if sessionStatusOK(h.status) {
		sRaw := readSessionData(sdk.reader, &h)
		err := yaml.Unmarshal([]byte(sRaw), &sdk.session)
		if err != nil {
			log.Fatal(err)
		}
		sdk.rawSession = strings.Split(sRaw, "\n")
		sdk.telemetry = readVariableHeaders(sdk.reader, &h)
		readVariableValues(sdk)
	}
}

func sessionStatusOK(status int) bool {
	return (status & stConnected) > 0
}
