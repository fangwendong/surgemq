// Copyright (c) 2014 The SurgeMQ Authors. All rights reserved.
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

package auth

import (
	"fmt"
	"testing"
)

func TestGm3Authenticator(t *testing.T) {
	mp := map[string]*ClientInfo{
		"xxcx": &ClientInfo{
			Token:    "fwd",
			UserName: "v1",
			UserId:   1234,
		},
	}

	f := func(token string) *ClientInfo {
		return mp[token]
	}
	manger, err := NewManager("gm3Auth", f)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	verify, clientInfo := manger.Authenticate("xcx")
	fmt.Println(verify, clientInfo)
}
