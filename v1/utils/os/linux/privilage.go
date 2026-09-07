package linux

import (
	"fmt"
	"os"
	"os/exec"
)

// IsAdmin checks if the effective user ID (EUID) is root (0).
func IsAdmin() bool {
	return os.Geteuid() == 0
}

// CheckPrivilege verifies if the process is running as root or has effective CAP_SYS_ADMIN capabilities.
// func CheckPrivilege() bool {
// 	if IsAdmin() {
// 		return true
// 	}

// 	// Check if the binary has active CAP_SYS_ADMIN capabilities via libcap/unix
// 	header := unix.CapUserHeader{
// 		Version: unix.LINUX_CAPABILITY_VERSION_3,
// 		Pid:     int32(os.Getpid()),
// 	}

// 	// LINUX_CAPABILITY_VERSION_3 requires 2 elements (64 bits total)
// 	var data [2]unix.CapUserData

// 	if err := unix.Capget(&header, &data[0]); err == nil {
// 		// CAP_SYS_ADMIN is capability bit 21 (located in the first 32-bit word, index 0)
// 		const capSysAdminBit = 21
// 		if (data[0].Effective & (1 << capSysAdminBit)) != 0 {
// 			return true
// 		}
// 	}
// 	return false
// }

// RunAsAdmin executes a target command with elevated privileges using sudo or pkexec.
func RunAsAdmin(command string, args []string) error {
	var helper string
	var helperArgs []string

	// Determine available elevation agent (sudo for CLI, pkexec for Desktop/GUI environments)
	if _, err := exec.LookPath("sudo"); err == nil {
		helper = "sudo"
		helperArgs = append([]string{command}, args...)
	} else if _, err := exec.LookPath("pkexec"); err == nil {
		helper = "pkexec"
		helperArgs = append([]string{command}, args...)
	} else {
		return fmt.Errorf("no privilege elevation helper (sudo/pkexec) found in PATH")
	}

	cmd := exec.Command(helper, helperArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// RequestAdminPrivileges re-executes the current binary using sudo/pkexec.
func RequestAdminPrivileges() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to determine executable path: %w", err)
	}

	return RunAsAdmin(exe, os.Args[1:])
}

// ElevatePrivilege attempts elevation via sudo/pkexec and exits current unprivileged process.
func ElevatePrivilege() error {
	if err := RequestAdminPrivileges(); err != nil {
		return err
	}
	// Exit unprivileged parent process upon spawning the elevated process
	os.Exit(0)
	return nil
}

// CheckAndElevate checks if the process has privileges, requesting elevation if absent.
// func CheckAndElevate() error {
// 	if !CheckPrivilege() {
// 		if err := ElevatePrivilege(); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
