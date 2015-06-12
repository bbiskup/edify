package message

type SegGrpStart struct {
	RecordNum    int
	GroupNum     int
	IsMandatory  bool
	MaxCount     int
	NestingLevel int
}
