package irsdk

import (
	"fmt"
	"log"
	"time"
)

type varBuffer struct {
	tickCount int // used to detect changes in data
	bufOffset int // offset from header
}

type telemetryVariable struct {
	varType     int // irsdk_VarType
	offset      int // offset fron start of buffer row
	count       int // number of entrys (array) so length in bytes would be irsdk_VarTypeBytes[type] * count
	countAsTime bool
	Name        string
	Desc        string
	Unit        string
	Value       interface{}
	rawBytes    []byte
}

func (v telemetryVariable) String() string {
	var ret string
	switch v.varType {
	case 0:
		ret = fmt.Sprintf("%c", v.Value)
	case 1:
		ret = fmt.Sprintf("%v", v.Value)
	case 2:
		ret = fmt.Sprintf("%d", v.Value)
	case 3:
		ret = fmt.Sprintf("%s", v.Value)
	case 4:
		ret = fmt.Sprintf("%f", v.Value)
	case 5:
		ret = fmt.Sprintf("%f", v.Value)
	default:
		ret = fmt.Sprintf("Unknown (%d)", v.varType)
	}
	return ret
}

func findLatestBuffer(r reader, h *header) varBuffer {
	var vb varBuffer
	foundTickCount := 0
	for i := 0; i < h.numBuf; i++ {
		rbuf := make([]byte, 16)
		_, err := r.ReadAt(rbuf, int64(48+i*16))
		if err != nil {
			log.Fatal(err)
		}
		currentVb := varBuffer{
			byte4ToInt(rbuf[0:4]),
			byte4ToInt(rbuf[4:8]),
		}
		//fmt.Printf("BUFF?: %+v\n", currentVb)
		if foundTickCount < currentVb.tickCount {
			foundTickCount = currentVb.tickCount
			vb = currentVb
		}
	}
	return vb
}

func readVariableHeaders(r reader, h *header) map[string]telemetryVariable {
	vars := make(map[string]telemetryVariable, h.numVars)
	for i := 0; i < h.numVars; i++ {
		rbuf := make([]byte, 144)
		_, err := r.ReadAt(rbuf, int64(h.headerOffset+i*144))
		if err != nil {
			log.Fatal(err)
		}
		v := telemetryVariable{
			byte4ToInt(rbuf[0:4]),
			byte4ToInt(rbuf[4:8]),
			byte4ToInt(rbuf[8:12]),
			int(rbuf[12]) > 0,
			bytesToString(rbuf[16:48]),
			bytesToString(rbuf[48:112]),
			bytesToString(rbuf[112:144]),
			nil,
			nil,
		}
		vars[v.Name] = v
	}
	return vars
}

func readVariableValues(header *header, reader reader, telemetry map[string]telemetryVariable) int64 {
	if !sessionStatusOK(header.status) {
		return 0
	}

	// find latest buffer for variables
	vb := findLatestBuffer(reader, header)

	for varName, v := range telemetry {
		var rbuf []byte
		switch v.varType {
		case 0:
			rbuf = make([]byte, 1)
			_, err := reader.ReadAt(rbuf, int64(vb.bufOffset+v.offset))
			if err != nil {
				log.Fatal(err)
			}
			v.Value = string(rbuf[0])
		case 1:
			rbuf = make([]byte, 1)
			_, err := reader.ReadAt(rbuf, int64(vb.bufOffset+v.offset))
			if err != nil {
				log.Fatal(err)
			}
			v.Value = int(rbuf[0]) > 0
		case 2:
			rbuf = make([]byte, 4)
			_, err := reader.ReadAt(rbuf, int64(vb.bufOffset+v.offset))
			if err != nil {
				log.Fatal(err)
			}
			v.Value = byte4ToInt(rbuf)
		case 3:
			rbuf = make([]byte, 4)
			_, err := reader.ReadAt(rbuf, int64(vb.bufOffset+v.offset))
			if err != nil {
				log.Fatal(err)
			}
			v.Value = byte4toBitField(rbuf)
		case 4:
			rbuf = make([]byte, 4)
			_, err := reader.ReadAt(rbuf, int64(vb.bufOffset+v.offset))
			if err != nil {
				log.Fatal(err)
			}
			v.Value = byte4ToFloat(rbuf)
		case 5:
			rbuf = make([]byte, 8)
			_, err := reader.ReadAt(rbuf, int64(vb.bufOffset+v.offset))
			if err != nil {
				log.Fatal(err)
			}
			v.Value = byte8ToFloat(rbuf)
		}
		v.rawBytes = rbuf
		telemetry[varName] = v
	}

	return time.Now().Unix()
}
