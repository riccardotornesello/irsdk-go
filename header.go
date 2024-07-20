package irsdk

import (
	"log"
)

type varBuf struct {
	TickCount int
	BufOffset int
	Pad       [2]int
}

type header struct {
	Version           int
	Status            int
	TickRate          int
	SessionInfoUpdate int
	SessionInfoLen    int
	SessionInfoOffset int
	NumVars           int
	VarHeaderOffset   int
	NumBuf            int
	BufLen            int
	Pad1              [2]int
	VarBuf            [MaxBufs]varBuf
}

const varBufSize = 4 * 4
const headerSize = 12*4 + MaxBufs*varBufSize

func readHeader(r reader) *header {
	rbuf := make([]byte, headerSize)
	_, err := r.ReadAt(rbuf, 0)
	if err != nil {
		log.Fatal(err)
	}

	h := header{
		Byte4ToInt(rbuf[0:4]),
		Byte4ToInt(rbuf[4:8]),
		Byte4ToInt(rbuf[8:12]),
		Byte4ToInt(rbuf[12:16]),
		Byte4ToInt(rbuf[16:20]),
		Byte4ToInt(rbuf[20:24]),
		Byte4ToInt(rbuf[24:28]),
		Byte4ToInt(rbuf[28:32]),
		Byte4ToInt(rbuf[32:36]),
		Byte4ToInt(rbuf[36:40]),
		[2]int{Byte4ToInt(rbuf[40:44]), Byte4ToInt(rbuf[44:48])},
		[MaxBufs]varBuf{},
	}

	for i := range h.VarBuf {
		h.VarBuf[i].TickCount = Byte4ToInt(rbuf[48+i*varBufSize : 52+i*varBufSize])
		h.VarBuf[i].BufOffset = Byte4ToInt(rbuf[52+i*varBufSize : 56+i*varBufSize])
	}

	return &h
}
