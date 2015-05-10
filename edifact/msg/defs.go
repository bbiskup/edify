package msg

// http://www.unece.org/trade/untdid/texts/d423.htm
// Chapter 2.2 "Syntax separators, terminator and release character"
const (
	UNH               = "UNH"
	UNT               = "UNT"
	SegTerm           = '\''
	SegTagDataElemSep = '+'

	CompDataElemSep    = ':'
	CompDataElemSepStr = string(CompDataElemSep)

	ReleaseChar = '?'
)
