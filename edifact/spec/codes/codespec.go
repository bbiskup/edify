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
	Id          string
	Name        string
	Description string
}

func NewCodeSpec(id string, name string, description string) *CodeSpec {
	return &CodeSpec{id, name, description}
}

func (s *CodeSpec) String() string {
	descriptionStr := util.Ellipsis(s.Description, maxDescrDisplayLen)
	return fmt.Sprintf("%s %s %s", s.Id, s.Name, descriptionStr)
}
