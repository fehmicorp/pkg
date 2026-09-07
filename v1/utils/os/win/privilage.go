package win

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

// IsAdmin checks if the user account belongs to the local Administrators group.
func IsAdmin() bool {
	var sid *windows.SID

	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid,
	)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}

	return member
}

// CheckPrivilege checks if the current process token is actively UAC elevated.
func CheckPrivilege() bool {
	token := windows.Token(0)
	return token.IsElevated()
}

// RunAsAdmin executes a target executable with elevated privileges via UAC prompt ("runas").
func RunAsAdmin(command string, args []string) error {
	verbPtr, err := syscall.UTF16PtrFromString("runas")
	if err != nil {
		return fmt.Errorf("failed to encode verb: %w", err)
	}

	exePtr, err := syscall.UTF16PtrFromString(command)
	if err != nil {
		return fmt.Errorf("failed to encode executable path: %w", err)
	}

	argPtr, err := syscall.UTF16PtrFromString(strings.Join(args, " "))
	if err != nil {
		return fmt.Errorf("failed to encode arguments: %w", err)
	}

	cwd, _ := os.Getwd()
	cwdPtr, err := syscall.UTF16PtrFromString(cwd)
	if err != nil {
		return fmt.Errorf("failed to encode working directory: %w", err)
	}

	var showCmd int32 = windows.SW_NORMAL

	err = windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		return fmt.Errorf("failed to execute process as admin: %w", err)
	}

	return nil
}

// RequestAdminPrivileges triggers a UAC prompt to relaunch the current executable as Admin.
func RequestAdminPrivileges() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to obtain current executable path: %w", err)
	}

	// Resolve symlinks to prevent temp directory issues
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	args := os.Args[1:]
	return RunAsAdmin(exe, args)
}

// ElevatePrivilege requests admin privileges and exits the current non-elevated process.
func ElevatePrivilege() error {
	if err := RequestAdminPrivileges(); err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

// CheckAndElevate checks for UAC elevation; requests elevation and exits if missing.
func CheckAndElevate() error {
	admin := IsAdmin()
	if !admin {
		if err := ElevatePrivilege(); err != nil {
			return err
		}
	}
	return nil
}
