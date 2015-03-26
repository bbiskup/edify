package specparser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var reprRE = regexp.MustCompile(`^(a|n|an)(\.\.)?([0-9]+)`)

type ReprType int

const (
	Num = iota
	Alpha
	AlphaNum
)

type Repr struct {
	Typ   ReprType
	Range bool
	Max   uint32
}

func (r *Repr) String() string {
	var typeStr string
	switch r.Typ {
	case Num:
		typeStr = "n"
	case Alpha:
		typeStr = "a"
	case AlphaNum:
		typeStr = "an"
	}

	var rangeStr string
	if r.Range {
		rangeStr = ".."
	} else {
		rangeStr = ""
	}

	return fmt.Sprintf("%s%s%d", typeStr, rangeStr, r.Max)
}

func NewRepr(typ ReprType, range_ bool, max_ uint32) *Repr {
	return &Repr{
		Typ:   typ,
		Range: range_,
		Max:   max_,
	}
}

func ParseRepr(elemStr string) (*Repr, error) {
	reprMatch := reprRE.FindStringSubmatch(elemStr)
	if reprMatch == nil {
		return nil, errors.New(
			fmt.Sprintf("Missing repr section in line '%s'",
				elemStr))
	}

	typStr := reprMatch[1]
	var typ ReprType
	if typStr == "a" {
		typ = Alpha
	} else if typStr == "an" {
		typ = AlphaNum
	} else if typStr == "n" {
		typ = Num
	} else {
		return nil, errors.New(fmt.Sprintf("Unknown repr type '%s'", typStr))
	}

	var range_ bool = false
	if reprMatch[2] == ".." {
		range_ = true
	}

	max_, err := strconv.Atoi(reprMatch[3])
	if err != nil {
		return nil, err
	}

	return NewRepr(typ, range_, uint32(max_)), nil
}
