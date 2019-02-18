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
package db

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPaginatorCreation(t *testing.T) {
	Convey("Given a paginator (perPage: 0, page: 0)", t, func() {
		paginator := NewPaginator(0, 0)
		paginator.SetTotal(0)

		Convey("When checking default values", func() {

			Convey("Then default values should be valid", func() {
				So(paginator.Page, ShouldEqual, 1)
				So(paginator.PerPage, ShouldEqual, DefaultPerPage)
				So(paginator.Total(), ShouldEqual, 0)
				So(paginator.HasNext(), ShouldBeFalse)
				So(paginator.HasPrev(), ShouldBeFalse)
				So(paginator.HasOtherPages(), ShouldBeFalse)
				So(paginator.NumPages(), ShouldEqual, 1)
				So(paginator.Offset(), ShouldEqual, 0)
				So(paginator.NextPage(), ShouldEqual, 1)
				So(paginator.PrevPage(), ShouldEqual, 1)
				So(paginator.CurrentPageCount(), ShouldEqual, 0)
			})
		})
	})
}

func TestPaginator(t *testing.T) {
	Convey("Given a paginator (perPage: 50, page: 1)", t, func() {
		paginator := NewPaginator(1, 50)

		Convey("When there are 20 elements", func() {
			paginator.SetTotal(20)

			Convey("Then there is no pagination enabled", func() {
				So(paginator.Page, ShouldEqual, 1)
				So(paginator.PerPage, ShouldEqual, 50)
				So(paginator.Total(), ShouldEqual, 20)
				So(paginator.NumPages(), ShouldEqual, 1)
				So(paginator.Offset(), ShouldEqual, 0)
				So(paginator.HasNext(), ShouldBeFalse)
				So(paginator.HasPrev(), ShouldBeFalse)
				So(paginator.HasOtherPages(), ShouldBeFalse)
				So(paginator.NextPage(), ShouldEqual, 1)
				So(paginator.PrevPage(), ShouldEqual, 1)
				So(paginator.CurrentPageCount(), ShouldEqual, 20)
			})
		})

		Convey("When there are 52 elements", func() {
			paginator.SetTotal(52)

			Convey("Then pagination should be enabled", func() {
				So(paginator.Page, ShouldEqual, 1)
				So(paginator.PerPage, ShouldEqual, 50)
				So(paginator.Total(), ShouldEqual, 52)
				So(paginator.NumPages(), ShouldEqual, 2)
				So(paginator.Offset(), ShouldEqual, 0)
				So(paginator.HasNext(), ShouldBeTrue)
				So(paginator.HasPrev(), ShouldBeFalse)
				So(paginator.HasOtherPages(), ShouldBeTrue)
				So(paginator.NextPage(), ShouldEqual, 2)
				So(paginator.PrevPage(), ShouldEqual, 1)
				So(paginator.CurrentPageCount(), ShouldEqual, 50)
			})

		})

		Convey("When going to the last page", func() {
			paginator.SetTotal(52)
			paginator.Page = 2

			Convey("Then pagination should be enabled", func() {
				So(paginator.Page, ShouldEqual, 2)
				So(paginator.PerPage, ShouldEqual, 50)
				So(paginator.Total(), ShouldEqual, 52)
				So(paginator.NumPages(), ShouldEqual, 2)
				So(paginator.Offset(), ShouldEqual, 50)
				So(paginator.HasNext(), ShouldBeFalse)
				So(paginator.HasPrev(), ShouldBeTrue)
				So(paginator.HasOtherPages(), ShouldBeTrue)
				So(paginator.NextPage(), ShouldEqual, 2)
				So(paginator.PrevPage(), ShouldEqual, 1)
				So(paginator.CurrentPageCount(), ShouldEqual, 2)
			})

		})
	})
}

func TestAnotherPaginator(t *testing.T) {
	Convey("Given another paginator (perPage: 50, page: 2)", t, func() {
		paginator := NewPaginator(2, 50)

		Convey("When there are 20 elements", func() {
			paginator.SetTotal(20)

			Convey("Then there is no pagination enabled", func() {
				So(paginator.Total(), ShouldEqual, 20)
				So(paginator.HasNext(), ShouldBeFalse)
				So(paginator.HasPrev(), ShouldBeTrue)
				So(paginator.HasOtherPages(), ShouldBeTrue)
				So(paginator.NumPages(), ShouldEqual, 1)
				So(paginator.Offset(), ShouldEqual, 20)
				So(paginator.NextPage(), ShouldEqual, 1)
				So(paginator.PrevPage(), ShouldEqual, 1)
				So(paginator.CurrentPageCount(), ShouldEqual, 0)
			})
		})

		Convey("When there are 52 elements", func() {
			paginator.SetTotal(52)

			Convey("Then pagination should be enabled", func() {
				So(paginator.Total(), ShouldEqual, 52)
				So(paginator.HasNext(), ShouldBeFalse)
				So(paginator.HasPrev(), ShouldBeTrue)
				So(paginator.HasOtherPages(), ShouldBeTrue)
				So(paginator.NumPages(), ShouldEqual, 2)
				So(paginator.Offset(), ShouldEqual, 50)
				So(paginator.NextPage(), ShouldEqual, 2)
				So(paginator.PrevPage(), ShouldEqual, 1)
				So(paginator.CurrentPageCount(), ShouldEqual, 2)
			})

		})
	})
}
