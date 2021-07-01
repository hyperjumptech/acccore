package acccore

import (
	"fmt"
	"testing"
)

type PageTest struct {
	RequestForPageNo   int
	RequestForItemSize int
	TotalResultSet     int
	TotalEntries       int
	TotalPages         int
	Page               int
	PageSize           int
	NextPage           int
	PreviousPage       int
	FirstPage          int
	LastPage           int
	IsFirst            bool
	IsLast             bool
	HavePrev           bool
	HaveNext           bool
	Offset             int
}

var (
	testPageData = []PageTest{
		{
			RequestForPageNo:   0,
			RequestForItemSize: 0,
			TotalResultSet:     0,
			TotalEntries:       0,
			TotalPages:         1,
			Page:               1,
			PageSize:           0,
			NextPage:           1,
			PreviousPage:       1,
			FirstPage:          1,
			LastPage:           1,
			IsFirst:            true,
			IsLast:             true,
			HavePrev:           false,
			HaveNext:           false,
			Offset:             0,
		},
		{
			RequestForPageNo:   0,
			RequestForItemSize: 10,
			TotalResultSet:     0,
			TotalEntries:       0,
			TotalPages:         1,
			Page:               1,
			PageSize:           0,
			NextPage:           1,
			PreviousPage:       1,
			FirstPage:          1,
			LastPage:           1,
			IsFirst:            true,
			IsLast:             true,
			HavePrev:           false,
			HaveNext:           false,
			Offset:             0,
		},
		{
			RequestForPageNo:   0,
			RequestForItemSize: 10,
			TotalResultSet:     10,
			TotalEntries:       10,
			TotalPages:         1,
			Page:               1,
			PageSize:           10,
			NextPage:           1,
			PreviousPage:       1,
			FirstPage:          1,
			LastPage:           1,
			IsFirst:            true,
			IsLast:             true,
			HavePrev:           false,
			HaveNext:           false,
			Offset:             0,
		},
		{
			RequestForPageNo:   0,
			RequestForItemSize: 10,
			TotalResultSet:     15,
			TotalEntries:       15,
			TotalPages:         2,
			Page:               1,
			PageSize:           10,
			NextPage:           2,
			PreviousPage:       1,
			FirstPage:          1,
			LastPage:           2,
			IsFirst:            true,
			IsLast:             false,
			HavePrev:           false,
			HaveNext:           true,
			Offset:             0,
		},
		{
			RequestForPageNo:   1,
			RequestForItemSize: 10,
			TotalResultSet:     15,
			TotalEntries:       15,
			TotalPages:         2,
			Page:               1,
			PageSize:           10,
			NextPage:           2,
			PreviousPage:       1,
			FirstPage:          1,
			LastPage:           2,
			IsFirst:            true,
			IsLast:             false,
			HavePrev:           false,
			HaveNext:           true,
			Offset:             0,
		},
		{
			RequestForPageNo:   2,
			RequestForItemSize: 10,
			TotalResultSet:     15,
			TotalEntries:       15,
			TotalPages:         2,
			Page:               2,
			PageSize:           5,
			NextPage:           2,
			PreviousPage:       1,
			FirstPage:          1,
			LastPage:           2,
			IsFirst:            false,
			IsLast:             true,
			HavePrev:           true,
			HaveNext:           false,
			Offset:             10,
		},
	}
)

func TestPageResultFor(t *testing.T) {
	for i, tst := range testPageData {
		req := PageRequest{
			PageNo:   tst.RequestForPageNo,
			ItemSize: tst.RequestForItemSize,
			Sorts:    nil,
		}
		pr := PageResultFor(req, tst.TotalResultSet)
		if pr.IsLast != tst.IsLast {
			t.Error(fmt.Sprintf("#%d : Expect IsLast %v but %v", i, tst.IsLast, pr.IsLast))
			t.FailNow()
		}
		if pr.IsFirst != tst.IsFirst {
			t.Error(fmt.Sprintf("#%d : Expect IsFirst %v but %v", i, tst.IsFirst, pr.IsFirst))
			t.FailNow()
		}
		if pr.HavePrev != tst.HavePrev {
			t.Error(fmt.Sprintf("#%d : Expect HavePrev %v but %v", i, tst.HavePrev, pr.HavePrev))
			t.FailNow()
		}
		if pr.HaveNext != tst.HaveNext {
			t.Error(fmt.Sprintf("#%d : Expect HaveNext %v but %v", i, tst.HaveNext, pr.HaveNext))
			t.FailNow()
		}
		if pr.TotalEntries != tst.TotalEntries {
			t.Error(fmt.Sprintf("#%d : Expect TotalEntries %v but %v", i, tst.TotalEntries, pr.TotalEntries))
			t.FailNow()
		}
		if pr.PageSize != tst.PageSize {
			t.Error(fmt.Sprintf("#%d : Expect PageSize %v but %v", i, tst.PageSize, pr.PageSize))
			t.FailNow()
		}
		if pr.Page != tst.Page {
			t.Error(fmt.Sprintf("#%d : Expect Page %v but %v", i, tst.Page, pr.Page))
			t.FailNow()
		}
		if pr.Offset != tst.Offset {
			t.Error(fmt.Sprintf("#%d : Expect Offset %v but %v", i, tst.Offset, pr.Offset))
			t.FailNow()
		}
		if pr.LastPage != tst.LastPage {
			t.Error(fmt.Sprintf("#%d : Expect LastPage %v but %v", i, tst.LastPage, pr.LastPage))
			t.FailNow()
		}
		if pr.FirstPage != tst.FirstPage {
			t.Error(fmt.Sprintf("#%d : Expect FirstPage %v but %v", i, tst.FirstPage, pr.FirstPage))
			t.FailNow()
		}
		if pr.NextPage != tst.NextPage {
			t.Error(fmt.Sprintf("#%d : Expect NextPage %v but %v", i, tst.NextPage, pr.NextPage))
			t.FailNow()
		}
		if pr.PreviousPage != tst.PreviousPage {
			t.Error(fmt.Sprintf("#%d : Expect PreviousPage %v but %v", i, tst.PreviousPage, pr.PreviousPage))
			t.FailNow()
		}
		if pr.TotalPages != tst.TotalPages {
			t.Error(fmt.Sprintf("#%d : Expect TotalPages %v but %v", i, tst.TotalPages, pr.TotalPages))
			t.FailNow()
		}
	}
}
