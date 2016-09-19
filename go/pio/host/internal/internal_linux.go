// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package internal

// CPUInfo returns parsed data from /proc/cpuinfo.
func CPUInfo() map[string]string {
	return makeCPUInfo()
}

// CPUInfo returns parsed data from /etc/os-release.
func OSRelease() map[string]string {
	return makeOSRelease()
}
