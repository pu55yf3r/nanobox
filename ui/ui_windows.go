// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

// +build windows

package ui

// PPrompt calls prompt, because in windows the lib that hides the typed response
// cant be used
func PPrompt(p string) string {
	return Prompt(p)
}