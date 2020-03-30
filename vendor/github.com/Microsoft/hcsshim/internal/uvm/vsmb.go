package uvm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Microsoft/hcsshim/internal/requesttype"
	hcsschema "github.com/Microsoft/hcsshim/internal/schema2"
)

const vsmbSharePrefix = `\\?\VMSMB\VSMB-{dcc079ae-60ba-4d07-847c-3493609c0870}\`

// findVSMBShare finds a share by `hostPath`. If not found returns `ErrNotAttached`.
func (uvm *UtilityVM) findVSMBShare(ctx context.Context, m map[string]*vsmbShare, hostPath string) (*vsmbShare, error) {
	share, ok := m[hostPath]
	if !ok {
		return nil, ErrNotAttached
	}
	return share, nil
}

func (share *vsmbShare) GuestPath() string {
	return vsmbSharePrefix + share.name
}

// AddVSMB adds a VSMB share to a Windows utility VM. Each VSMB share is ref-counted and
// only added if it isn't already. This is used for read-only layers, mapped directories
// to a container, and for mapped pipes.
func (uvm *UtilityVM) AddVSMB(ctx context.Context, hostPath string, guestRequest interface{}, options *hcsschema.VirtualSmbShareOptions) error {
	if uvm.operatingSystem != "windows" {
		return errNotSupported
	}

	uvm.m.Lock()
	defer uvm.m.Unlock()

	// Temporary support to allow single-file mapping. If `hostPath` is a
	// directory, map it without restriction. However, if it is a file, map the
	// directory containing the file, and use `AllowedFileList` to only allow
	// access to that file. If the directory has been mapped before for
	// single-file use, add the new file to the `AllowedFileList` and issue an
	// Update operation.
	st, err := os.Stat(hostPath)
	if err != nil {
		return err
	}
	var file string
	m := uvm.vsmbDirShares
	if !st.IsDir() {
		m = uvm.vsmbFileShares
		file = hostPath
		hostPath = filepath.Dir(hostPath)
		options.RestrictFileAccess = true
		options.SingleFileMapping = true
	}
	hostPath = filepath.Clean(hostPath)
	var requestType = requesttype.Update
	share, err := uvm.findVSMBShare(ctx, m, hostPath)
	if err == ErrNotAttached {
		requestType = requesttype.Add
		uvm.vsmbCounter++
		shareName := "s" + strconv.FormatUint(uvm.vsmbCounter, 16)

		share = &vsmbShare{
			name:         shareName,
			guestRequest: guestRequest,
		}
	}
	newAllowedFiles := share.allowedFiles
	if options.RestrictFileAccess {
		newAllowedFiles = append(newAllowedFiles, file)
	}

	// Update on a VSMB share currently only supports updating the
	// AllowedFileList, and in fact will return an error if RestrictFileAccess
	// isn't set (e.g. if used on an unrestricted share). So we only call Modify
	// if we are either doing an Add, or if RestrictFileAccess is set.
	if requestType == requesttype.Add || options.RestrictFileAccess {
		modification := &hcsschema.ModifySettingRequest{
			RequestType: requestType,
			Settings: hcsschema.VirtualSmbShare{
				Name:         share.name,
				Options:      options,
				Path:         hostPath,
				AllowedFiles: newAllowedFiles,
			},
			ResourcePath: "VirtualMachine/Devices/VirtualSmb/Shares",
		}
		if err := uvm.modify(ctx, modification); err != nil {
			return err
		}
	}

	share.allowedFiles = newAllowedFiles
	share.refCount++
	m[hostPath] = share

	return nil
}

// RemoveVSMB removes a VSMB share from a utility VM. Each VSMB share is ref-counted
// and only actually removed when the ref-count drops to zero.
func (uvm *UtilityVM) RemoveVSMB(ctx context.Context, hostPath string) error {
	if uvm.operatingSystem != "windows" {
		return errNotSupported
	}

	uvm.m.Lock()
	defer uvm.m.Unlock()

	st, err := os.Stat(hostPath)
	if err != nil {
		return err
	}
	m := uvm.vsmbDirShares
	if !st.IsDir() {
		m = uvm.vsmbFileShares
		hostPath = filepath.Dir(hostPath)
	}
	hostPath = filepath.Clean(hostPath)
	share, err := uvm.findVSMBShare(ctx, m, hostPath)
	if err != nil {
		return fmt.Errorf("%s is not present as a VSMB share in %s, cannot remove", hostPath, uvm.id)
	}

	share.refCount--
	if share.refCount > 0 {
		return nil
	}

	modification := &hcsschema.ModifySettingRequest{
		RequestType:  requesttype.Remove,
		Settings:     hcsschema.VirtualSmbShare{Name: share.name},
		ResourcePath: "VirtualMachine/Devices/VirtualSmb/Shares",
	}
	if err := uvm.modify(ctx, modification); err != nil {
		return fmt.Errorf("failed to remove vsmb share %s from %s: %+v: %s", hostPath, uvm.id, modification, err)
	}

	delete(m, hostPath)
	return nil
}

// GetVSMBUvmPath returns the guest path of a VSMB mount.
func (uvm *UtilityVM) GetVSMBUvmPath(ctx context.Context, hostPath string) (string, error) {
	if hostPath == "" {
		return "", fmt.Errorf("no hostPath passed to GetVSMBUvmPath")
	}

	uvm.m.Lock()
	defer uvm.m.Unlock()

	st, err := os.Stat(hostPath)
	if err != nil {
		return "", err
	}
	m := uvm.vsmbDirShares
	f := ""
	if !st.IsDir() {
		m = uvm.vsmbFileShares
		hostPath, f = filepath.Split(hostPath)
	}
	hostPath = filepath.Clean(hostPath)
	share, err := uvm.findVSMBShare(ctx, m, hostPath)
	if err != nil {
		return "", err
	}
	return filepath.Join(share.GuestPath(), f), nil
}
