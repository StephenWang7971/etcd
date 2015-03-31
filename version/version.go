// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package version

import (
	"os"
	"path"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/coreos/etcd/pkg/types"
)

var (
	Version = "2.0.8"
)

// WalVersion is an enum for versions of etcd logs.
type DataDirVersion string

const (
	DataDirUnknown  DataDirVersion = "Unknown WAL"
	DataDir0_4      DataDirVersion = "0.4.x"
	DataDir2_0      DataDirVersion = "2.0.0"
	DataDir2_0Proxy DataDirVersion = "2.0 proxy"
	DataDir2_0_1    DataDirVersion = "2.0.1"
)

func DetectDataDir(dirpath string) (DataDirVersion, error) {
	names, err := fileutil.ReadDir(dirpath)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		// Error reading the directory
		return DataDirUnknown, err
	}
	nameSet := types.NewUnsafeSet(names...)
	if nameSet.Contains("member") {
		ver, err := DetectDataDir(path.Join(dirpath, "member"))
		if ver == DataDir2_0 {
			return DataDir2_0_1, nil
		} else if ver == DataDir0_4 {
			// How in the blazes did it get there?
			return DataDirUnknown, nil
		}
		return ver, err
	}
	if nameSet.ContainsAll([]string{"snap", "wal"}) {
		// .../wal cannot be empty to exist.
		walnames, err := fileutil.ReadDir(path.Join(dirpath, "wal"))
		if err == nil && len(walnames) > 0 {
			return DataDir2_0, nil
		}
	}
	if nameSet.ContainsAll([]string{"proxy"}) {
		return DataDir2_0Proxy, nil
	}
	if nameSet.ContainsAll([]string{"snapshot", "conf", "log"}) {
		return DataDir0_4, nil
	}
	if nameSet.ContainsAll([]string{"standby_info"}) {
		return DataDir0_4, nil
	}

	return DataDirUnknown, nil
}