package dataelement

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

func IsNumChar(char rune) bool {
	return char >= '0' && char <= '9'
}

// A.1 alphabetic character set: A character set that contains
//     letters and may contain control characters and special
//     characters but not digits (ISO 2382/4)
func IsEDIFACTAlphabetic(char rune) bool {
	// TODO restrictive enough?
	return !IsNumChar(char)
}

func (r *Repr) isPunctuation(char rune) bool {
	return char == '.' || char == ','
}

func (r *Repr) Validate(dataElemStr string) (valid bool, err error) {
	var typ = r.Typ
	var strLen uint32

	for _, c := range dataElemStr {
		strLen++
		if typ == AlphaNum {
			continue
		} else {
			isNum := IsNumChar(c)
			if typ == Num && !isNum && !r.isPunctuation(c) {
				return false, errors.New(fmt.Sprintf("Found non-numeric character '%c'", c))
			} else if typ == Alpha && isNum {
				return false, errors.New(fmt.Sprintf("Found numeric character '%c'", c))
			}
		}
	}
	if r.Range {
		if strLen > r.Max {
			return false, errors.New(fmt.Sprintf("String too long: '%s' (max: %d)",
				dataElemStr, r.Max))
		}
	} else {
		if strLen != r.Max {
			return false, errors.New(fmt.Sprintf("String '%s' should have %d characters)",
				dataElemStr, r.Max))
		}
	}
	return true, nil
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
