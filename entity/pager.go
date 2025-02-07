package entity

type Pager struct {
	Limit  uint64
	Offset uint64
}

func (p *Pager) GetLimit() uint64 {
	if p.Limit > 0 {
		return p.Limit
	}
	return 20
}

func (p *Pager) GetOffset() uint64 {
	return p.Offset
}
