package query

import (
	"fmt"
	"strconv"
)

type ItemKind int

const (
	MessageKind ItemKind = iota
	SegKind
	SegGrpKind
	CompositeDataElemKind
	SimpleDataElemKind
)

// means wildcard (all items)
const noIndex = -1

// Single part of a query expression
type QueryPart struct {
	ItemKind ItemKind
	Id       string
	Index    int
}

func (q *QueryPart) itemKindStr() string {
	switch q.ItemKind {
	case MessageKind:
		return "msg"
	case SegKind:
		return "seg"
	case SegGrpKind:
		return "grp"
	case CompositeDataElemKind:
		return "cmp"
	case SimpleDataElemKind:
		return "smp"
	default:
		panic(fmt.Sprintf("Unknown item kind: '%d'", q.ItemKind))
	}
}

func (q *QueryPart) String() string {
	var indexPart string
	if q.Index != noIndex {
		indexPart = strconv.Itoa(q.Index)
	} else {
		indexPart = "*"
	}
	return fmt.Sprintf("QueryPart %s %s %s", q.itemKindStr(), q.Id, indexPart)
}

func NewQueryPart(kind ItemKind, id string, index int) *QueryPart {
	return &QueryPart{kind, id, index}
}
