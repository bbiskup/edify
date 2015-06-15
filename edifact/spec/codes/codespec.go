package codes

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
)

const (
	maxDescrDisplayLen = 10
)

// EDIFACT Code as defined e.g. in UNCL.14B
type CodeSpec struct {
	Id    string
	Name  string
	Descr string
}

func NewCodeSpec(id string, name string, descr string) *CodeSpec {
	return &CodeSpec{id, name, descr}
}

func (s *CodeSpec) String() string {
	descrStr := util.Ellipsis(s.Descr, maxDescrDisplayLen)
	return fmt.Sprintf("%s %s %s", s.Id, s.Name, descrStr)
}
