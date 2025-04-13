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

package rtx

import (
	"context"
	"errors"
)

var (
	// ErrUnknownMethod is returned when the method is not recognized.
	ErrUnknownMethod = errors.New("unknown method")
)

// -----------------------------------------------------------------------------

type M = map[string]any

type RoundTripper interface {
	RoundTrip(ctx context.Context, method string, params M) (resp M, err error)
	OnNotify(notify func(method string, params M))
	Close() error
}

// -----------------------------------------------------------------------------
