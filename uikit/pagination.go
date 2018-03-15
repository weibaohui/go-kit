package uikit

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Pagination struct {
	PageSize  int         `json:"page_size,omitempty" `
	PageIndex int         `json:"page_index,omitempty"`
	Total     int         `json:"total,omitempty"`
	Sidx      string      `json:"sidx,omitempty"`
	Order     string      `json:"order,omitempty"`
	Keywords  string      `json:"keywords,omitempty"`
	List      interface{} `json:"list,omitempty"`
}

func GetPage(c *gin.Context) (*Pagination, error) {
	//check  pagination exists,use url query  method

	var pager Pagination
	err := c.MustBindWith(&pager, binding.Query)

	if pager.PageSize == 0 {
		pager.PageSize = 15
	}
	if pager.PageIndex == 0 {
		pager.PageIndex = 1
	}
	return &pager, err
}

//NewPagination default pagination
func NewPagination() *Pagination {
	page := &Pagination{
		PageSize:  15,
		PageIndex: 1,
	}
	return page
}

func (p *Pagination) SetTotal(total int) *Pagination {
	p.Total = total
	return p
}

func (p *Pagination) SetList(list interface{}) *Pagination {
	p.List = list
	return p
}
