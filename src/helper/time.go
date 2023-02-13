package helper

/*
Copyright The Rotate Authors.

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

import "time"

// mustParseTime é uma função auxiliar para converter uma string no formato "YYYY-MM-DD HH:MM:SS" para time.Time
func MustParseTime(str string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		panic(err)
	}
	return t
}
