// +build windows

package beater

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	psapi                    = syscall.NewLazyDLL("psapi.dll")
	getWindowText            = user32.NewProc("GetWindowTextW")
	getForegroundWindow      = user32.NewProc("GetForegroundWindow")
	getWindowThreadProcessID = user32.NewProc("GetWindowThreadProcessId")
	getProcessImageFileName  = psapi.NewProc("GetProcessImageFileNameA")
)

type ActispyWin32 struct {
	hwnd syscall.Handle
}

func newActispyWin32() (spy *ActispyWin32, err error) {
	spy = new(ActispyWin32)
	// Get forground window handler
	r0, _, e1 := syscall.Syscall(getForegroundWindow.Addr(), 0, 0, 0, 0)
	if e1 != 0 {
		err = fmt.Errorf("Get foreground window fails with %v", e1)
		return
	}
	spy.hwnd = syscall.Handle(r0)
	return
}

func (s *ActispyWin32) getProcessID() (processID int, err error) {
	getWindowThreadProcessID.Call(
		uintptr(s.hwnd),
		uintptr(unsafe.Pointer(&processID)),
	)
	return
}

func (s *ActispyWin32) getProcessName() (processName string, err error) {
	pid, _ := s.getProcessID()
	handle, e1 := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	defer syscall.CloseHandle(handle)
	if e1 != nil {
		err = fmt.Errorf("OpenProcess fails with %v", e1)
		return
	}
	var nameProc [512]byte
	ret, _, _ := getProcessImageFileName.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&nameProc)),
		uintptr(512),
	)
	if ret == 0 {
		return
	}
	processName = filepath.Base(carrayToString(nameProc))
	return
}

func (s *ActispyWin32) getWindowName() (windowName string, err error) {
	n := make([]uint16, 512)
	p := &n[0]
	r0, _, e1 := syscall.Syscall(getWindowText.Addr(), 3, uintptr(s.hwnd), uintptr(unsafe.Pointer(p)), uintptr(len(n)))
	if r0 == 0 {
		if e1 != 0 {
			err = fmt.Errorf("Get window name fails with %v", e1)
		} else {
			windowName = ""
		}
		return
	}
	windowName = syscall.UTF16ToString(n)
	return
}

func (s *ActispyWin32) getUserName() (userName string, err error) {
	return os.Getenv("USERNAME"), nil
}

func carrayToString(c [512]byte) string {
	end := 0
	for {
		if c[end] == 0 {
			break
		}
		end++
	}
	return string(c[:end])
}
