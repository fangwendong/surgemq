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
	"errors"
	"fmt"
)

var (
	ErrAuthFailure          = errors.New("auth: Authentication failure")
	ErrAuthProviderNotFound = errors.New("auth: Authentication provider not found")

	providers = make(map[string]Authenticator)
)

type ClientInfo struct {
	Token    string
	UserName string
	UserId   string
	Verified bool
}

type Authenticator interface {
	Authenticate(token string) (bool, *ClientInfo)
	SetVerifyFunc(f AuthFunc)
}

func Register(name string, provider Authenticator) {
	if provider == nil {
		panic("auth: Register provide is nil")
	}

	if _, dup := providers[name]; dup {
		panic("auth: Register called twice for provider " + name)
	}

	providers[name] = provider
}

func Unregister(name string) {
	delete(providers, name)
}

type Manager struct {
	p Authenticator
}

func NewManager(providerName string, verifyTokenFunc AuthFunc) (*Manager, error) {
	p, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provider %q", providerName)
	}
	p.SetVerifyFunc(verifyTokenFunc)

	return &Manager{p: p}, nil
}

func (this *Manager) Authenticate(token string) (bool, *ClientInfo) {
	return this.p.Authenticate(token)
}

func (this *Manager) SetVerifyFunc(f AuthFunc) {
	this.p.SetVerifyFunc(f)
}
