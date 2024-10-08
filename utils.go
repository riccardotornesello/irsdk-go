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

func Byte4toBitField(in []byte) uint32 {
	return binary.LittleEndian.Uint32(in)
}

func BytesToString(in []byte) string {
	return strings.TrimRight(string(in), "\x00")
}

func Byte8ToTimeStr(in []byte) string {
	return DoubleToTimeStr(Byte8ToFloat(in))
}

func FloatToTime(f float32) time.Duration {
	return time.Duration(f * float32(time.Second))
}

func DoubleToTime(f float64) time.Duration {
	return time.Duration(f * float64(time.Second))
}

func TimeToStr(t time.Duration) string {
	return fmt.Sprintf("%d:%02d.%03d", int(t.Minutes()), int(t.Seconds())%60, int(t.Milliseconds())%1000)
}

func FloatToTimeStr(f float32) string {
	return TimeToStr(FloatToTime(f))
}

func DoubleToTimeStr(f float64) string {
	return TimeToStr(DoubleToTime(f))
}
