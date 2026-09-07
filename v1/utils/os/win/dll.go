package win

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

var dllDir = "dll"
var programDir = "C:\\Program Files\\Fehmi"

func RunDll(pkg string) {
	dll := GetDll(true, dllDir, pkg)
	defer dll.Release()
	runProc, err := dll.FindProc("Run")
	if err != nil {
		fmt.Printf("Failed to find 'Run' procedure in DLL: %v\n", err)
		return
	}
	go func() {
		runProc.Call()
	}()
	select {}
}

func GetDll(current bool, dllDir string, pkg string) *syscall.DLL {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	var dllPath string
	if current {
		dllPath = filepath.Join(currentDir, dllDir, fmt.Sprintf("%s.dll", pkg))
	} else {
		dllPath = filepath.Join(programDir, dllDir, fmt.Sprintf("%s.dll", pkg))
	}
	dll, err := syscall.LoadDLL(dllPath)
	if err != nil {
		fmt.Printf("Failed to load DLL from %s: %v\n", dllPath, err)
		return nil
	}
	return dll
}
