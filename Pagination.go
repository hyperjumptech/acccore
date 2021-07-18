package acccore

// Sort define a sorting information, it specifies the column should be sorted and whether it should be ASCENDING or
// DESCENDING
type Sort struct {
	// Column Name to sort
	Column string

	// Ascending define the ordering, `true` means ASCENDING, `false` means DESCENDING
	Ascending bool
}

// PageRequest define the pagination request when listing a huge dataset.
type PageRequest struct {
	// PageNo is page number to request. The 1st page number is 1 (one)
	PageNo int

	// ItemSize is the maximum item should be contained within a single page
	ItemSize int

	// Sorts define the sorting strategy
	Sorts []Sort
}

// PageResultFor will calculate proper pagination result on a request with total rows in the result set.
func PageResultFor(request PageRequest, count int) PageResult {
	if request.ItemSize == 0 {
		request.ItemSize = 1
	}
	pr := PageResult{}
	pr.TotalEntries = count
	pr.FirstPage = 1
	if count == 0 {
		pr.LastPage = 1
	} else {
		if count%request.ItemSize == 0 {
			pr.LastPage = count / request.ItemSize
		} else {
			pr.LastPage = (count / request.ItemSize) + 1
		}
	}
	if request.PageNo > pr.LastPage {
		pr.Page = pr.LastPage
	} else if request.PageNo < pr.FirstPage {
		pr.Page = pr.FirstPage
	} else {
		pr.Page = request.PageNo
	}
	pr.IsFirst = pr.Page == pr.FirstPage
	pr.IsLast = pr.Page == pr.LastPage
	pr.HaveNext = !pr.IsLast
	pr.HavePrev = !pr.IsFirst
	pr.Offset = (pr.Page - 1) * request.ItemSize
	if pr.IsLast {
		if count == request.ItemSize || count == 0 {
			pr.PageSize = count
		} else {
			if count%request.ItemSize == 0 {
				pr.PageSize = request.ItemSize
			} else {
				pr.PageSize = count % request.ItemSize
			}
		}
		pr.NextPage = pr.Page
	} else {
		pr.PageSize = request.ItemSize
		pr.NextPage = pr.Page + 1
	}
	if pr.IsFirst {
		pr.PreviousPage = pr.Page
	} else {
		pr.PreviousPage = pr.Page - 1
	}
	if pr.TotalEntries == 0 {
		pr.TotalPages = 1
	} else {
		pr.TotalPages = pr.TotalEntries / request.ItemSize
		if pr.TotalEntries%request.ItemSize > 0 {
			pr.TotalPages++
		}
	}
	return pr
}

// PageResult define the pagination result that returned together with the listing.
type PageResult struct {
	// Request define the request that specified the pagination in the first place.
	Request PageRequest

	// TotalEntries total number of rows that is available in the result set.
	// for example if we request a page with maximum 10 entries in a page, but the select is actually
	// resulting 1000 records, that this should show 1000 TotalEntries.
	TotalEntries int

	// TotalPages is total pages calculated by dividing the total entries with page size in the request.
	// for example. If total number of record is 1000,
	// for a 10 PageSize, this should shows 100 total pages.
	// If the TotalEntries is 1001 and 10 PageSize, this should shows 101 pages.
	// If the TotalEntries is 1005 and 10 PageSize, this should shows 101 pages.
	// If the TotalEntries is 0 and 10 PageSize, this should shows 1 pages.
	TotalPages int

	// Page is the current page number. The 1st page number is 1.
	// If TotalEntries is 0, then this should still show 1.
	Page int

	// PageSize shows the current number of items in the this Page
	// If TotalEntries is 0 than this shows 0
	PageSize int

	// NextPage shows the Page + 1 IF the HaveNext is true and IsLast is false
	NextPage int

	// PreviousPage shows the Page + 1 IF the HavePrev is true and IsFirst is false
	PreviousPage int

	// FirstPage shows the number of first page. This always return 1
	FirstPage int

	// LastPage shows the number of last page.
	LastPage int

	// IsFirst an indicator if current Page is at the FirstPage
	IsFirst bool

	// IsLast an indicator if current Page is at the LastPage
	IsLast bool

	// HavePrev is an indicator if current Page is not at the FirstPage
	HavePrev bool

	// HaveNext is an indicator if current Page is not at the LastPage
	HaveNext bool

	// Offset is an offset number of item that shown at the beginning of each page.
	// for example
	// If Request.RequestForItemSize is 10 then the 1st offset is 0
	// If Request.RequestForItemSize is 10 then the 2nd offset is 10
	// If Request.RequestForItemSize is 10 then the 3rd offset is 20
	// If Request.RequestForItemSize is 10 then the 4nd offset is 30
	Offset int
}
