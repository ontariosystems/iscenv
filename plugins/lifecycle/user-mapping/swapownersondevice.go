/*
Copyright 2024 Finvi, Ontario Systems

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package usermappingplugin

import (
	"os"
	"path/filepath"
	"syscall"
)

func swapOwnersOnDevice(root string, oldUID, oldGID, newUID, newGID int) error {
	// Nothing to do, yay
	if oldUID == newUID && oldGID == newGID {
		return nil
	}

	info, err := os.Stat(root)
	if err != nil {
		return err
	}

	stat := info.Sys().(*syscall.Stat_t)
	rootDev := stat.Dev

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// if there was an error walking just abort the whole process
		if err != nil {
			return err
		}

		stat := info.Sys().(*syscall.Stat_t)
		// We're not on the same mount point anymore so skip this directory
		if stat.Dev != rootDev {
			return filepath.SkipDir
		}

		// This is a bit cumbersome, but it's to avoid having to do multiple Chown calls
		uid := int(stat.Uid)
		gid := int(stat.Gid)

		// These aren't the IDs you're looking for
		if uid != oldUID && gid != oldGID {
			return nil
		}

		if uid == oldUID {
			uid = newUID
		}

		if gid == oldGID {
			gid = newGID
		}

		return os.Chown(path, uid, gid)
	})
}
