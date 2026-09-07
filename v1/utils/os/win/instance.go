package win

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32          = syscall.NewLazyDLL("kernel32.dll")
	procCreateMutexW  = kernel32.NewProc("CreateMutexW")
	singleInstanceHdl windows.Handle
)

// EnsureSingleInstance checks if another instance is running via a Windows Mutex.
// Runs existingInstanceHandler and returns if a duplicate instance is detected.
func EnsureSingleInstance(appUniqueName string, newInstanceHandler func(), existingInstanceHandler func()) error {
	mutexNamePtr, err := syscall.UTF16PtrFromString("Local\\" + appUniqueName)
	if err != nil {
		if existingInstanceHandler != nil {
			existingInstanceHandler()
		}
		return fmt.Errorf("invalid mutex name: %w", err)
	}

	// Create Mutex: BOOL bInitialOwner = TRUE
	ret, _, err := procCreateMutexW.Call(
		0,
		1,
		uintptr(unsafePointer(mutexNamePtr)),
	)

	handle := windows.Handle(ret)
	if handle == 0 {
		if existingInstanceHandler != nil {
			existingInstanceHandler()
		}
		return fmt.Errorf("failed to create system mutex: %w", err)
	}

	// Check if Windows reports the mutex already existed
	if err != nil && err.(syscall.Errno) == windows.ERROR_ALREADY_EXISTS {
		windows.CloseHandle(handle)
		if existingInstanceHandler != nil {
			existingInstanceHandler() // Triggers os.Exit(0) for secondary instance
		}
		return nil
	}

	// Primary instance active: store handle in memory
	singleInstanceHdl = handle
	if newInstanceHandler != nil {
		newInstanceHandler()
	}

	return nil
}

// ReleaseSingleInstance closes the single-instance mutex handle upon exit.
func ReleaseSingleInstance() {
	if singleInstanceHdl != 0 {
		windows.CloseHandle(singleInstanceHdl)
		singleInstanceHdl = 0
	}
}

func unsafePointer(ptr *uint16) unsafe.Pointer {
	return unsafe.Pointer(ptr)
}
