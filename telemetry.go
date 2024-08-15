package irsdk

import (
	"fmt"
	"log"
	"time"
)

type VarType = int

const (
	VarTypeChar VarType = iota
	VarTypeBool
	VarTypeInt
	VarTypeBitField
	VarTypeFloat
	VarTypeDouble
	VarTypeCount //index, don't use
)

var VarTypeBytes = [VarTypeCount]int{
	1, // char
	1, // bool
	4, // int
	4, // bitField
	4, // float
	8, // double
}

type varHeader struct {
	Type        VarType
	Offset      int
	Count       int
	CountAsTime bool
	Pad         [3]byte
	Name        string
	Desc        string
	Unit        string
}

type TelemetryVar struct {
	Header   varHeader
	RawValue []byte
}

const varHeaderSize = 16 + MaxString + MaxDesc + MaxString

func (v TelemetryVar) Value() interface{} {
	if v.Header.Count > 1 {
		return v.Array()
	} else {
		return v.Single()
	}
}

func (v TelemetryVar) Single() interface{} {
	switch v.Header.Type {
	case VarTypeChar:
		return v.RawValue[0]
	case VarTypeBool:
		return v.RawValue[0] > 0
	case VarTypeInt:
		return Byte4ToInt(v.RawValue)
	case VarTypeBitField:
		return Byte4toBitField(v.RawValue)
	case VarTypeFloat:
		return Byte4ToFloat(v.RawValue)
	case VarTypeDouble:
		return Byte8ToFloat(v.RawValue)
	}
	return nil
}

func (v TelemetryVar) Array() interface{} {
	switch v.Header.Type {
	case VarTypeChar:
		return v.RawValue
	case VarTypeBool:
		arr := make([]bool, v.Header.Count)
		for i := 0; i < v.Header.Count; i++ {
			arr[i] = v.RawValue[i] > 0
		}
		return arr
	case VarTypeInt:
		arr := make([]int, v.Header.Count)
		for i := 0; i < v.Header.Count; i++ {
			arr[i] = Byte4ToInt(v.RawValue[i*4 : (i+1)*4])
		}
		return arr
	case VarTypeBitField:
		arr := make([]uint32, v.Header.Count)
		for i := 0; i < v.Header.Count; i++ {
			arr[i] = Byte4toBitField(v.RawValue[i*4 : (i+1)*4])
		}
		return arr
	case VarTypeFloat:
		arr := make([]float32, v.Header.Count)
		for i := 0; i < v.Header.Count; i++ {
			arr[i] = Byte4ToFloat(v.RawValue[i*4 : (i+1)*4])
		}
		return arr
	case VarTypeDouble:
		arr := make([]float64, v.Header.Count)
		for i := 0; i < v.Header.Count; i++ {
			arr[i] = Byte8ToFloat(v.RawValue[i*8 : (i+1)*8])
		}
		return arr
	}
	return nil
}

func (v TelemetryVar) String() string {
	return fmt.Sprintf("%v", v.Value())
}

func (v TelemetryVar) Time() time.Duration {
	switch v.Header.Type {
	case VarTypeFloat:
		return FloatToTime(v.Single().(float32))
	case VarTypeDouble:
		return DoubleToTime(v.Single().(float64))
	default:
		return 0
	}
}

func (v TelemetryVar) TimeStr() string {
	return TimeToStr(v.Time())
}

func readVariableHeaders(r reader, h *header) map[string]varHeader {
	vars := make(map[string]varHeader, h.NumVars)

	for i := 0; i < h.NumVars; i++ {
		rbuf := make([]byte, varHeaderSize)

		_, err := r.ReadAt(rbuf, int64(h.VarHeaderOffset+i*varHeaderSize))
		if err != nil {
			log.Fatal(err)
		}

		v := varHeader{
			Byte4ToInt(rbuf[0:4]),
			Byte4ToInt(rbuf[4:8]),
			Byte4ToInt(rbuf[8:12]),
			int(rbuf[12]) > 0,
			[3]byte{rbuf[12], rbuf[13], rbuf[14]},
			BytesToString(rbuf[16:48]),
			BytesToString(rbuf[48:112]),
			BytesToString(rbuf[112:144]),
		}
		vars[v.Name] = v
	}
	return vars
}

// Return which variable buffer has the latest tick count.
// It might be the same as the last one read.
func findLatestBuffer(h *header) *varBuf {
	lastTickIndex := 0

	for i := 1; i < h.NumBuf; i++ {
		if h.VarBuf[i].TickCount > h.VarBuf[lastTickIndex].TickCount {
			lastTickIndex = i
		}
	}

	return &h.VarBuf[lastTickIndex]
}

// This function updates LastTickCount and Telemetry fields of the IRSDK struct.
func updateTelemetryVariables(sdk *IRSDK) bool {
	vb := findLatestBuffer(sdk.Header)

	// If the tick count is the same as the last one read, return false.
	// If it's lower than the last one read, it means the data has been reset, maybe because the sim has been restarted.
	if vb.TickCount == sdk.LastTickCount {
		return false
	}

	headers := readVariableHeaders(sdk.Reader, sdk.Header)
	vars := make(map[string]TelemetryVar, len(headers))

	for varName, v := range headers {
		bufferSize := VarTypeBytes[v.Type] * v.Count

		rbuf := make([]byte, bufferSize)

		_, err := sdk.Reader.ReadAt(rbuf, int64(vb.BufOffset+v.Offset))
		if err != nil {
			log.Fatal(err)
		}

		vars[varName] = TelemetryVar{
			v,
			rbuf,
		}
	}

	sdk.LastTickCount = vb.TickCount
	sdk.Telemetry = vars
	sdk.LastDataTime = time.Now().Unix()

	return true
}
