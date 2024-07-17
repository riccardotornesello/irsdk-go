package winevents

import (
	"log"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const HWND_BROADCAST = 0xffff

var (
	moduser32                 = windows.NewLazySystemDLL("user32.dll")
	procSendNotifyMessage     = moduser32.NewProc("SendNotifyMessageW")
	procRegisterWindowMessage = moduser32.NewProc("RegisterWindowMessageW")
)

var eventHandle windows.Handle

func MAKELONG(a, b uint16) uint32 {
	return uint32(a) | (uint32(b) << 16)
}

func OpenEvent(eventName string) error {
	evtName, err := windows.UTF16PtrFromString(eventName)
	if err != nil {
		return err
	}

	eventHandle, err = windows.OpenEvent(windows.SYNCHRONIZE, false, evtName)
	if err != nil {
		return err
	}
	return nil
}

func WaitForSingleObject(timeout time.Duration) bool {
	t0 := time.Now().UnixNano()
	timeoutInt := uint32(timeout / time.Millisecond)
	r, err := windows.WaitForSingleObject(eventHandle, timeoutInt)
	if err != nil {
		return false
	}
	if r != windows.WAIT_OBJECT_0 {
		remainingTimeout := timeoutInt - uint32((time.Now().UnixNano()-t0)/1000000)
		if remainingTimeout > 0 {
			time.Sleep(time.Duration(remainingTimeout) * time.Millisecond)
		}
		return false
	}
	return true
}

func BroadcastMsg(msgName string, msg int, p1 int, p2 interface{}, p3 int) bool {
	var p2Int int
	switch v := p2.(type) {
	case int, int8, int16, int32, int64:
		p2Int = v.(int)
	case float32, float64:
		p2Int = (int)(v.(float64) * 65536.0)
	default:
		log.Fatal("Second param must be an int or a float")
	}
	msgNameUTF16, err := windows.UTF16PtrFromString(msgName)
	if err != nil {
		log.Fatal(err)
	}

	msgID, _, _ := procRegisterWindowMessage.Call(uintptr(unsafe.Pointer(msgNameUTF16)))
	if msgID == 0 {
		return false
	}

	ret, _, _ := procSendNotifyMessage.Call(
		uintptr(HWND_BROADCAST),
		msgID,
		uintptr(MAKELONG(uint16(msg), uint16(p1))),
		uintptr(MAKELONG(uint16(p2Int), uint16(p3))),
	)
	return ret != 0
}
