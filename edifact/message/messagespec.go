package message

import (
	"fmt"
	"time"
)

// A message specification
// (e.g. edmd/BALANC_D.14B)
type MessageSpec struct {
	Id   string
	Name string

	Version     string
	Release     string
	ContrAgency string
	Revision    string
	Date        time.Time

	Source string
}

func (m *MessageSpec) String() string {
	return fmt.Sprintf("Message %s (%s %s)", m.Id, m.Name, m.Release)
}

func NewMessageSpec(
	Id string, Name string,
	Version string, Release string, ContrAgency string,
	Revision string, Date time.Time, Source string) *MessageSpec {

	return &MessageSpec{
		Id: Id, Name: Name,
		Version: Version, Release: Release, ContrAgency: ContrAgency,
		Revision: Revision, Date: Date, Source: Source,
	}
}
