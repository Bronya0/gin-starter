package request

import "errors"

// Pagination Pagination.
type Pagination struct {
	PageSize int `form:"page_size" json:"page_size"`
	PageNum  int `form:"page_num" json:"page_num"`
}

// Verify the value of pageNum and pageSize.
func (p *Pagination) Verify() error {
	if p.PageNum < 0 {
		return errors.New("pageNum err")
	} else if p.PageNum == 0 {
		p.PageNum = 1
	}
	if p.PageSize < 0 {
		return errors.New("pageSize err")
	} else if p.PageSize == 0 {
		p.PageSize = 10
	}
	return nil
}
