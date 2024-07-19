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

const (
	irsdkVarTypeChar int = iota
	irsdkVarTypeBool
	irsdkVarTypeInt
	irsdkVarTypeBitField
	irsdkVarTypeFloat
	irsdkVarTypeDouble
)

var irsdkVarTypeBytes = map[int]int{
	irsdkVarTypeChar:     1,
	irsdkVarTypeBool:     1,
	irsdkVarTypeInt:      4,
	irsdkVarTypeBitField: 4,
	irsdkVarTypeFloat:    4,
	irsdkVarTypeDouble:   8,
}

type telemetryVariable struct {
	VarType     int
	offset      int // offset fron start of buffer row
	Count       int // number of entrys (array) so length in bytes would be irsdk_VarTypeBytes[type] * count
	countAsTime bool
	Name        string
	Desc        string
	Unit        string
	Value       interface{}
	rawBytes    []byte
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
			Byte4ToInt(rbuf[0:4]),
			Byte4ToInt(rbuf[4:8]),
		}
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
			Byte4ToInt(rbuf[0:4]),
			Byte4ToInt(rbuf[4:8]),
			Byte4ToInt(rbuf[8:12]),
			int(rbuf[12]) > 0,
			BytesToString(rbuf[16:48]),
			BytesToString(rbuf[48:112]),
			BytesToString(rbuf[112:144]),
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
		bufferSize := irsdkVarTypeBytes[v.VarType] * v.Count
		rbuf := make([]byte, bufferSize)
		_, err := reader.ReadAt(rbuf, int64(vb.bufOffset+v.offset))
		if err != nil {
			log.Fatal(err)
		}

		// Convert raw bytes to value based on type and count (if > 1 it's an array)
		switch v.VarType {
		case irsdkVarTypeChar:
			if v.Count > 1 {
				v.Value = make([]string, v.Count)
				for i := 0; i < v.Count; i++ {
					v.Value.([]string)[i] = BytesToString(rbuf[i*1 : i*1+1])
				}
			} else {
				v.Value = BytesToString(rbuf)
			}
		case irsdkVarTypeBool:
			if v.Count > 1 {
				v.Value = make([]bool, v.Count)
				for i := 0; i < v.Count; i++ {
					v.Value.([]bool)[i] = rbuf[i] > 0
				}
			} else {
				v.Value = rbuf[0] > 0
			}
		case irsdkVarTypeInt:
			if v.Count > 1 {
				v.Value = make([]int, v.Count)
				for i := 0; i < v.Count; i++ {
					v.Value.([]int)[i] = Byte4ToInt(rbuf[i*4 : i*4+4])
				}
			} else {
				v.Value = Byte4ToInt(rbuf)
			}
		case irsdkVarTypeBitField:
			if v.Count > 1 {
				v.Value = make([]string, v.Count)
				for i := 0; i < v.Count; i++ {
					v.Value.([]string)[i] = Byte4toBitField(rbuf[i*4 : i*4+4])
				}
			} else {
				v.Value = Byte4toBitField(rbuf)
			}
		case irsdkVarTypeFloat:
			if v.Count > 1 {
				v.Value = make([]float32, v.Count)
				for i := 0; i < v.Count; i++ {
					v.Value.([]float32)[i] = Byte4ToFloat(rbuf[i*4 : i*4+4])
				}
			} else {
				v.Value = Byte4ToFloat(rbuf)
			}
		case irsdkVarTypeDouble:
			if v.Count > 1 {
				v.Value = make([]float64, v.Count)
				for i := 0; i < v.Count; i++ {
					v.Value.([]float64)[i] = Byte8ToFloat(rbuf[i*8 : i*8+8])
				}
			} else {
				v.Value = Byte8ToFloat(rbuf)
			}
		default:
			log.Fatal(fmt.Sprintf("Unknown variable type: %d", v.VarType))
		}

		v.rawBytes = rbuf
		telemetry[varName] = v
	}

	return time.Now().Unix()
}
