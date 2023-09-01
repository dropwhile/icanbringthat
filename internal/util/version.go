// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package util

import (
	"fmt"
	"runtime/debug"
)

var Version string

type VersionInfo struct {
	Version   string
	GoVersion string
	Revision  string
	Dirty     bool
}

func GetVersion() (VersionInfo, error) {
	vinfo := VersionInfo{
		Version:   Version,
		Revision:  "unknown",
		GoVersion: "unknown",
		Dirty:     true,
	}

	buildinfo, ok := debug.ReadBuildInfo()
	if !ok {
		return vinfo, fmt.Errorf("could not read buildinfo")
	}

	if Version == "" && buildinfo.Main.Version != "" {
		vinfo.Version = buildinfo.Main.Version
	}

	if buildinfo.GoVersion != "" {
		vinfo.GoVersion = buildinfo.GoVersion
	}

	for _, kv := range buildinfo.Settings {
		if kv.Value == "" {
			continue
		}

		switch kv.Key {
		case "vcs.revision":
			vinfo.Revision = kv.Value
		case "vcs.modified":
			vinfo.Dirty = kv.Value == "true"
		}
	}
	return vinfo, nil
}
