package msg

// Special segment identifiers
const (
	UNH = "UNH"
	UNT = "UNT"
)

// Service characters
// http://www.unece.org/trade/untdid/texts/d423.htm
// Chapter 2.2 "Syntax separators, terminator and release character"
const (
	SegTerm           = '\''
	SegTagDataElemSep = '+'

	CompDataElemSep    = ':'
	CompDataElemSepStr = string(CompDataElemSep)

	// Definitions
	// A.64    repeating data element: A composite data element or stand-
	// alone data element having a maximum occurrence of greater than
	// one in the segment specification.  [ 64 ]
	//
	// A.65 repetition separator: A service character used to separate
	// adjacent occurrences of a repeating data element.  [ 65 ]
	RepetitionSep = "*"

	ReleaseChar = '?'
)
