// MIT License
//
// Copyright (c) 2019 Thibault NORMAND
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package internal

// A FieldType indicates which member of the Field union struct should be used
// and how it should be serialized.
type FieldType uint8

const (
	// UnknownFieldType is the default field type.
	UnknownFieldType FieldType = iota
	// BoolFieldType indicates that the field carries a bool.
	BoolFieldType
	// ByteStringFieldType indicates that the field carries UTF-8 encoded bytes.
	ByteStringFieldType
	// DurationFieldType indicates that the field carries a time.Duration.
	DurationFieldType
	// ErrorFieldType indicates that the field carries an error.
	ErrorFieldType
	// Float64FieldType indicates that the field carries a float64.
	Float64FieldType
	// Float32FieldType indicates that the field carries a float32.
	Float32FieldType
	// Int64FieldType indicates that the field carries an int64.
	Int64FieldType
	// Int32FieldType indicates that the field carries an int32.
	Int32FieldType
	// Int16FieldType indicates that the field carries an int16.
	Int16FieldType
	// Int8FieldType indicates that the field carries an int8.
	Int8FieldType
	// StringFieldType indicates that the field carries a string.
	StringFieldType
	// StringerFieldType indicates that the field carries a fmt.Stringer.
	StringerFieldType
	// TimeFieldType indicates that the field carries a time.Time.
	TimeFieldType
	// Uint64FieldType indicates that the field carries a uint64.
	Uint64FieldType
	// Uint32FieldType indicates that the field carries a uint32.
	Uint32FieldType
	// Uint16FieldType indicates that the field carries a uint16.
	Uint16FieldType
	// Uint8FieldType indicates that the field carries a uint8.
	Uint8FieldType
	// UintptrFieldType indicates that the field carries a uintptr.
	UintptrFieldType
)

// Field declare logger attributes
type Field struct {
	Key       string
	Type      FieldType
	String    string
	Integer   int64
	Interface interface{}
}
