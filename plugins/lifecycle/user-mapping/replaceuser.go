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
	"os/exec"
	"os/user"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type userInfo struct {
	user  string
	group string
	uid   int
	gid   int
}

func replaceUser(ui userInfo) error {
	l := plog.WithFields(log.Fields{
		"user":  ui.user,
		"group": ui.group,
	})

	u, err := user.Lookup(ui.user)
	if err != nil {
		return err
	}
	oldUID, err := strconv.Atoi(u.Uid)
	if err != nil {
		return err
	}

	g, err := user.LookupGroup(ui.group)
	if err != nil {
		return err
	}
	oldGID, err := strconv.Atoi(g.Gid)
	if err != nil {
		return err
	}

	l = l.WithFields(log.Fields{
		"oldUID": oldUID,
		"oldGID": oldGID,
		"newUID": strconv.Itoa(ui.uid),
		"newGID": strconv.Itoa(ui.gid),
	})

	if ui.uid == 0 {
		l.Warn("Refusing to switch ISC manager UID to 0 (root)")
		return nil
	}

	if ui.gid == 0 {
		l.Warn("Refusing to switch ISC manager GID to 0 (root)")
		return nil
	}

	if out, err := exec.Command("usermod", "-o", "-u", strconv.Itoa(ui.uid), ui.user).CombinedOutput(); err != nil {
		l.WithField("output", out).WithError(err).Error("Failed to execute usermod")
		return err
	}

	if out, err := exec.Command("groupmod", "-o", "-g", strconv.Itoa(ui.gid), ui.group).CombinedOutput(); err != nil {
		l.WithField("output", out).WithError(err).Error("Failed to execute groupmod")
		return err
	}
	l.Info("Replaced user and group ids")

	l.Info("Searching file system for files owned by old IDs and changing ownership")
	if err := swapOwnersOnDevice("/", oldUID, oldGID, ui.uid, ui.gid); err != nil {
		return nil
	}

	l.Info("Replaced file system ownership")
	return nil
}
