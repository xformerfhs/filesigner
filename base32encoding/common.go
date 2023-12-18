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
var enc = base32.NewEncoding("23456789CDGHJKNPTVXZcdghjknptvxz").WithPadding(base32.NoPadding)
var encKey = base32.NewEncoding("B9C8D7F6G5H4J3K2L1M0NPQRSTUVWXYZ").WithPadding(base32.NoPadding)
