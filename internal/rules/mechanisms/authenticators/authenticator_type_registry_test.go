// Copyright 2022 Dimitrij Drus <dadrus@gmx.de>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package authenticators

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dadrus/heimdall/internal/app"
	"github.com/dadrus/heimdall/internal/validation"
)

func TestCreateAuthenticatorPrototype(t *testing.T) {
	t.Parallel()

	// there are seven authenticators implemented, which should have been registered
	require.Len(t, authenticatorTypeFactories, 6)

	for uc, tc := range map[string]struct {
		typ    string
		assert func(t *testing.T, err error, auth Authenticator)
	}{
		"using known type": {
			typ: AuthenticatorAnonymous,
			assert: func(t *testing.T, err error, auth Authenticator) {
				t.Helper()

				require.NoError(t, err)
				assert.IsType(t, &anonymousAuthenticator{}, auth)
			},
		},
		"using unknown type": {
			typ: "foo",
			assert: func(t *testing.T, err error, _ Authenticator) {
				t.Helper()

				require.Error(t, err)
				require.ErrorIs(t, err, ErrUnsupportedAuthenticatorType)
			},
		},
	} {
		t.Run(uc, func(t *testing.T) {
			// GIVEN
			validator, err := validation.NewValidator()
			require.NoError(t, err)

			appCtx := app.NewContextMock(t)
			appCtx.EXPECT().Validator().Maybe().Return(validator)
			appCtx.EXPECT().Logger().Maybe().Return(log.Logger)

			// WHEN
			auth, err := CreatePrototype(appCtx, "foo", tc.typ, nil)

			// THEN
			tc.assert(t, err, auth)
		})
	}
}
