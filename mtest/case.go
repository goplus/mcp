/*
 * Copyright (c) 2025 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mtest

import (
	"context"
	"testing"

	"github.com/qiniu/x/test"
)

// -----------------------------------------------------------------------------

type CaseT = test.CaseT

type CaseApp struct {
	test.Case
	*App
	ctx context.Context
}

// Gopt_CaseApp_TestMain is required by Go+ compiler as the entry of a YAP test case.
func Gopt_CaseApp_TestMain(c interface{ initCaseApp(*App, CaseT) }, t *testing.T) {
	app := new(App).initApp()
	c.initCaseApp(app, test.NewT(t))
	c.(interface{ Main() }).Main()
}

func (p *CaseApp) initCaseApp(app *App, t CaseT) {
	p.App = app
	p.CaseT = t
	p.ctx = context.TODO()
}

// -----------------------------------------------------------------------------
