// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package pages

type eventRow struct {
	Id            string `json:"id"`
	DeviceName    string `json:"deviceName"`
	ProfileName   string `json:"profileName"`
	Created       int64  `json:"created"`
	Origin        int64  `json:"origin"`
	ReadingsCount int64  `json:"readingsCount"`
	Tags          string `json:"tags,omitempty"`

	Json string `json:"json"`
}

type readingRow struct {
	Id           string `json:"id"`
	Created      int64  `json:"created"`
	Origin       int64  `json:"origin"`
	DeviceName   string `json:"deviceName"`
	ResourceName string `json:"resourceName"`
	ProfileName  string `json:"profileName"`
	ValueType    string `json:"valueType"`
	BinaryValue  string `json:"binaryValue"`
	MediaType    string `json:"mediaType"`
	Value        string `json:"value"`

	Json string `json:"json"`
}
