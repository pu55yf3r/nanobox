// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package auth

type (
	auth struct{}
	Auth interface {
		Authenticate() (string, string)
		Reauthenticate() (string, string)
	}
)

var (
	Default Auth = auth{}
)

func (auth) Authenticate() (string, string) {
	return Authenticate()
}

func (auth) Reauthenticate() (string, string) {
	return Reauthenticate()
}