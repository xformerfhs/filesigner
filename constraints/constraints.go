// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package constraints defines a set of useful constraints to be used
// with type parameters.

// This is a copy of golang.org/x/exp/constraints to get rid of the
// golang.org/x/exp dependency which contains a lot of weird stuff.

package constraints

// Signed is a constraint that permits any signed integer type.
// If future releases of Go add new predeclared signed integer types,
// this constraint will be modified to include them.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint that permits any unsigned integer type.
// If future releases of Go add new predeclared unsigned integer types,
// this constraint will be modified to include them.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a constraint that permits any integer type.
// If future releases of Go add new predeclared integer types,
// this constraint will be modified to include them.
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint that permits any floating-point type.
// If future releases of Go add new predeclared floating-point types,
// this constraint will be modified to include them.
type Float interface {
	~float32 | ~float64
}

// OrderedNumber is a constraint that permits any integer or floating-point type.
// These types allow an ordering of numbers.
// This is a type that is not present in the original constraints package.
// If future releases of Go add new predeclared ordered number types,
// this constraint will be modified to include them.
type OrderedNumber interface {
	Integer | Float
}

// Complex is a constraint that permits any complex numeric type.
// If future releases of Go add new predeclared complex numeric types,
// this constraint will be modified to include them.
type Complex interface {
	~complex64 | ~complex128
}

// Number is a constraint that permits any integer, floating-point or complex type.
// This is a type that is not present in the original constraints package.
// If future releases of Go add new predeclared number types,
// this constraint will be modified to include them.
type Number interface {
	OrderedNumber | Complex
}

// Ordered is a constraint that permits any ordered type: any type
// that supports the operators < <= >= >.
// If future releases of Go add new ordered types,
// this constraint will be modified to include them.
type Ordered interface {
	OrderedNumber | ~string
}
