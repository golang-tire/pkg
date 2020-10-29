package config

// Int type interface
type Int interface {
	Int() int
	Int64() int64
}

// String type interface
type String interface {
	String() string
}

// Float type interface
type Float interface {
	Float32() float32
	Float64() float64
}

// Bool type interface
type Bool interface {
	Bool() bool
}

type intHolder struct {
	value *int64
}

type stringHolder struct {
	value *string
}

type floatHolder struct {
	value *float64
}

type boolHolder struct {
	value *bool
}

// String will return value of string variable
func (sh stringHolder) String() string {
	return *sh.value
}

// Int will return value of int variable
func (ih intHolder) Int() int {
	return int(*ih.value)
}

// Int64 will return value of int64 variable
func (ih intHolder) Int64() int64 {
	return *ih.value
}

// Float32 will return value of float32 variable
func (fh floatHolder) Float32() float32 {
	return float32(*fh.value)
}

// Float64 will return value of float64 variable
func (fh floatHolder) Float64() float64 {
	return *fh.value
}

// Bool will return value of bool variable
func (bh boolHolder) Bool() bool {
	return *bh.value
}
