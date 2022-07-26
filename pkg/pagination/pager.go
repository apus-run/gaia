package pagination

type Pager interface {

	// PageNumber will return the page number
	PageNumber() int

	// PageSize will return the page size
	PageSize() int

	// TotalPages will return the number of total pages
	TotalPages() int

	// Data will return the data
	Data() []interface{}

	// DataSize will return the size of data.
	// Usually it's len(Data())
	DataSize() int

	// HasNext will return whether has next page
	HasNext() bool

	// HasData will return whether this page has data.
	HasData() bool

	Offset() int
}
