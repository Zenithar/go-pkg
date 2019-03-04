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

package rethinkdb

import (
	"strings"

	"go.zenithar.org/pkg/db"

	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

// ConvertSortParameters to rethinkdb query string
func ConvertSortParameters(params db.SortParameters) []interface{} {

	var sorts []interface{}
	for k, v := range params {
		switch v {
		case db.Ascending:
			sorts = append(sorts, r.Asc(strings.ToLower(k)))
			break
		case db.Descending:
			sorts = append(sorts, r.Desc(strings.ToLower(k)))
			break
		default:
			sorts = append(sorts, r.Desc(strings.ToLower(k)))
		}
	}

	// Apply default sort
	if len(sorts) == 0 {
		sorts = append(sorts, r.Asc("id"))
	}

	return sorts
}
