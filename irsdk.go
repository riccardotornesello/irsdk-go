package irsdk

import (
	"io"
	"log"

	"github.com/hidez8891/shm"
)

type reader interface {
	io.Reader
	io.ReaderAt
	io.Closer
}

type IRSDK struct {
	Reader reader

	LastTickCount int
	LastDataTime  int64

	Header    *header
	Telemetry map[string]TelemetryVar
	Session   *Session
}

func Init(r reader) *IRSDK {
	if r == nil {
		var err error
		r, err = shm.Open(MemMapFile, MemMapSize)
		if err != nil {
			log.Fatal(err)
		}
	}

	header := readHeader(r)

	sdk := IRSDK{
		Reader:        r,
		LastTickCount: 0,
		Header:        header,
		Telemetry:     make(map[string]TelemetryVar),
		Session:       nil,
	}

	sdk.Update(true)

	return &sdk
}

func (sdk *IRSDK) IsConnected() bool {
	return sdk.Header.Status&stConnected > 0
}

func (sdk *IRSDK) Update(withSession bool) bool {
	// Update the header to get the last data about the variable buffers.
	header := readHeader(sdk.Reader)
	sdk.Header = header

	// Update the session data.
	if withSession {
		updateSessionData(sdk)
	}

	// If the tick count is the same as the last one read, return false.
	// Otherwise update the data.
	return updateTelemetryVariables(sdk)
}

func (sdk *IRSDK) Close() {
	sdk.Reader.Close()
}

func (sdk *IRSDK) GetVar(name string) (interface{}, bool) {
	v, ok := sdk.Telemetry[name]
	if !ok {
		return nil, false
	}
	return v.Value(), true
}
