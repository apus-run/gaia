package pagination

import (
	"errors"
)

var (
	_                     Pager = (*Pagination)(nil)
	DefaultPaginationSize       = 10
)

// Pagination is the default implementation of Page interface
type Pagination struct {
	// pageNumber 当前页
	pageNumber int `json:"page_number"`
	// pageSize 分页数
	pageSize int `json:"page_size"`
	// total means total page count
	total int `json:"total"`
	// data 数据
	data []interface{} `json:"data"`
	// totalPages 总页数
	totalPages int  `json:"total_pages,omitempty"`
	hasNext    bool `json:"has_next,omitempty"`
}

// PageNumber will return the page number
func (p *Pagination) PageNumber() int {
	return p.pageNumber
}

func (p *Pagination) SetPageNumber(in int) {
	p.pageNumber = in
}

// PageSize will return the page size
func (p *Pagination) PageSize() int {
	return p.pageSize
}

func (p *Pagination) SetPageSize(in int) {
	p.pageSize = in
}

// TotalPages will return the number of total pages
func (p *Pagination) TotalPages() int {
	return p.totalPages
}

// Data will return the data
func (p *Pagination) Data() []interface{} {
	return p.data
}

// DataSize will return the size of data.
// it's len(GetData())
func (p *Pagination) DataSize() int {
	return len(p.Data())
}

// HasNext will return whether has next page
func (p *Pagination) HasNext() bool {
	return p.hasNext
}

// HasData will return whether this page has data.
func (p *Pagination) HasData() bool {
	return p.DataSize() > 0
}

func (p *Pagination) Offset() int {
	return (p.pageNumber - 1) * p.pageSize
}

func (p *Pagination) Valid() error {
	if p.pageNumber == 0 {
		p.pageNumber = 1
	}
	if p.pageSize == 0 {
		p.pageSize = DefaultPaginationSize
	}

	if p.pageNumber < 0 {
		return errors.New("current MUST be larger than 0")
	}

	if p.pageSize < 0 {
		return errors.New("invalid pageSize")
	}
	return nil
}

// New will create an instance
func New(pageNumber int, pageSize int, data []interface{}, total int) *Pagination {
	remain := total % pageSize
	totalPages := total / pageSize
	if remain > 0 {
		totalPages++
	}

	hasNext := total-pageNumber-pageSize > 0

	return &Pagination{
		pageNumber: pageNumber,
		pageSize:   pageSize,
		data:       data,
		total:      total,
		totalPages: totalPages,
		hasNext:    hasNext,
	}
}
