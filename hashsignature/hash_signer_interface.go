//
// SPDX-FileCopyrightText: Copyright 2024 Frank Schwab
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileType: SOURCE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Frank Schwab
//
// Version: 1.0.0
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//

package hashsignature

import (
	"errors"
)

// IsDestroyedErr is returned after the private key of the signer has been destroyed.
var IsDestroyedErr = errors.New(`Signer has been destroyed`)

// HashSigner is the interface that each HashSigner has to use.
type HashSigner interface {
	PublicKey() ([]byte, error)

	SignHash(hashValue []byte) ([]byte, error)

	Destroy()
}
