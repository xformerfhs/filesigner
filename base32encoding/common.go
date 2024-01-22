package base32encoding

import (
	"encoding/base32"
	"sync"
)

// ******** Private constants ********

const keySeparator = '-'
const keyGroupSize = 4

// ******** Private variables ********

var emx sync.Mutex

// enc is a base32 encoder that uses the word-safe alphabet
var enc = base32.NewEncoding("23456789CFGHJMPQRVWXcfghjmpqrvwx").WithPadding(base32.NoPadding)

// encKey is a base32 encoder which uses a custom alphabet to encode keys
var encKey = base32.NewEncoding("B9C8D7F6G5H4J3K2L1M0NPQRSTUVWXYZ").WithPadding(base32.NoPadding)
