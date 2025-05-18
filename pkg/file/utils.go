/*
Copyright 2023 cuisongliu@qq.com.

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

package file

import (
	"os"
	"strings"
)

func GetFileNameFromSubStr(path, substr string) (names []string, err error) {
	_, err = os.Stat(path)
	if err != nil {
		return
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), substr) {
			names = append(names, file.Name())
		}
	}
	return names, err
}

func GetFileNameFromVersion(s []string, t string) (o string, err error) {
	for _, v := range s {
		if strings.Contains(v, t) {
			return v, nil
		}
	}
	return
}
