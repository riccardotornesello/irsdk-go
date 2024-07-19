package irsdk

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"
)

func Byte4ToInt(in []byte) int {
	return int(binary.LittleEndian.Uint32(in))
}

func Byte4ToFloat(in []byte) float32 {
	bits := binary.LittleEndian.Uint32(in)
	return math.Float32frombits(bits)
}

func Byte8ToFloat(in []byte) float64 {
	bits := binary.LittleEndian.Uint64(in)
	return math.Float64frombits(bits)
}

func Byte4toBitField(in []byte) string {
	return fmt.Sprintf("0x%x", int(binary.LittleEndian.Uint32(in)))
}

func BytesToString(in []byte) string {
	return strings.TrimRight(string(in), "\x00")
}

func FloatToTime(f float32) time.Duration {
	return time.Duration(f * float32(time.Second))
}

func TimeToStr(t time.Duration) string {
	return fmt.Sprintf("%d:%02d.%03d", int(t.Minutes()), int(t.Seconds())%60, int(t.Milliseconds())%1000)
}

func FloatToTimeStr(f float32) string {
	return TimeToStr(FloatToTime(f))
}
