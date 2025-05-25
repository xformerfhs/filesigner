//
// SPDX-FileCopyrightText: Copyright 2025 Frank Schwab
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
//    2025-05-25: V1.0.0: Created.
//

package main

// commonMsgBase is the base number for all messages in common.
// Reserved numbers are 10-19.
const commonMsgBase = 10

// signCmdMsgBase is the base number for all messages in sign_command.
// Reserved numbers 20-39.
const signCmdMsgBase = 20

// verifyCmdMsgBase is the base number for all messages in verify_command.
// Reserved numbers are 40-59.
const verifyCmdMsgBase = 40

// mainMsgBase is the base number for all messages in main.
// Reserved numbers are 60-69.
const mainMsgBase = 60

// errorMsgBase is the base number for all messages in error.
// Reserved numbers are 70-79.
const errorMsgBase = 70

// handlerMsgBase is the base number for all messages in handlers.
// Reserved numbers are 80-89.
const handlerMsgBase = 80
