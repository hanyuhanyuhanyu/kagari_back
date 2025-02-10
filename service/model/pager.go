package model

import "strconv"

type Pager struct {
	Limit  uint64
	Offset uint64
}

func PagerFromString(l string, o string) (*Pager, error) {
	p := &Pager{}
	limit, err := strconv.ParseUint(l, 10, 64)
	if err != nil {
		return nil, err
	}
	p.Limit = limit
	offset, err := strconv.ParseUint(o, 10, 64)
	if err != nil {
		return nil, err
	}
	p.Offset = offset
	return p, err

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
