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
//    2024-02-26: V1.0.0: Created.
//

package hashsignature

import (
	"filesigner/slicehelper"
)

// ******** Private constants ********

// beginFence is a constant byte slice that sits to the left of the hash value.
// These bytes have been generated by the BCryptGenRandom Windows API xored with values from the RDRAND instruction.
var beginFence = []byte{
	0x44, 0x97, 0x72, 0xda, 0xb6, 0xa9, 0x2b, 0x43,
	0xc5, 0x06, 0xc4, 0x92, 0x06, 0x37, 0x58, 0xe4,
}

// endFence is a constant byte slice that sits to the right of the hash value.
// These bytes have been generated by the BCryptGenRandom Windows API xored with values from the RDRAND instruction.
var endFence = []byte{
	0xb8, 0x16, 0x17, 0x05, 0x8d, 0x38, 0xc4, 0x50,
	0x2b, 0x01, 0x2f, 0xf9, 0x49, 0x9e, 0x2d, 0xdc,
}

// ******** Private functions ********

// paddedHash returns the supplied hashValue sitting in the middle between two padding values.
func paddedHash(hashValue []byte) []byte {
	return slicehelper.Concat(beginFence, hashValue, endFence)
}
