package main

import (
	"fmt"

	"github.com/fehmicorp/pkg/v1/utils/os/win"
)

func main() {
	uerr := win.CheckAndElevate()
	if uerr != nil {
		fmt.Printf("Error: %v\n", uerr)
		return
	}
	win.RunDll("tray")
}
