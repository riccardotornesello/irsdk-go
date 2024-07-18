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
	session       *session
	rawSession    []string
	telemetry     map[string]telemetryVariable
	lastValidData int64
}

func (sdk *IRSDK) WaitForData(timeout time.Duration) bool {
	if !sdk.IsConnected() {
		initIRSDK(sdk)
	}
	if winevents.WaitForSingleObject(timeout) {
		sdk.RefreshTelemetry()
		return true
	}
	return false
}

func (sdk *IRSDK) GetVar(name string) (telemetryVariable, error) {
	if !sessionStatusOK(sdk.header.status) {
		return telemetryVariable{}, fmt.Errorf("Session is not active")
	}

	if v, ok := sdk.telemetry[name]; ok {
		return v, nil
	}

	return telemetryVariable{}, fmt.Errorf("Telemetry variable %q not found", name)
}

func (sdk *IRSDK) GetSession() session {
	return *sdk.session
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

func (sdk *IRSDK) RefreshHeader() {
	sdk.header = readHeader(sdk.reader)
}

func (sdk *IRSDK) RefreshSession() {
	sRaw := readSessionData(sdk.reader, sdk.header)

	newSession := session{}
	err := yaml.Unmarshal([]byte(sRaw), &newSession)
	if err != nil {
		log.Fatal(err)
	}

	sdk.session = &newSession
	sdk.rawSession = strings.Split(sRaw, "\n")
}

func (sdk *IRSDK) RefreshTelemetry() {
	telemetry := readVariableHeaders(sdk.reader, sdk.header)
	lastValidData := readVariableValues(sdk.header, sdk.reader, telemetry)

	sdk.telemetry = telemetry
	sdk.lastValidData = lastValidData
}

func (sdk *IRSDK) Refresh() {
	sdk.RefreshHeader()

	if sessionStatusOK(sdk.header.status) {
		sdk.RefreshSession()
		sdk.RefreshTelemetry()
	}
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
	sdk.RefreshHeader()

	sdk.rawSession = nil
	if sdk.telemetry != nil {
		sdk.telemetry = nil
	}

	if sessionStatusOK(sdk.header.status) {
		sdk.RefreshSession()
		sdk.RefreshTelemetry()
	}
}

func sessionStatusOK(status int) bool {
	return (status & stConnected) > 0
}
