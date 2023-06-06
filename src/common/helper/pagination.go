package helper

type Pagination struct {
	Limit      int         `json:"limit,omitempty"`
	Page       int         `json:"page,omitempty"`
	TotalRows  int         `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Result     interface{} `json:"result"`
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetLimit() int {
	if p.Limit < 1 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}
