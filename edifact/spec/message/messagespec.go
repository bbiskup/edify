package message

import (
	"bytes"
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
	Parts  []MessageSpecPart
}

func (m *MessageSpec) String() string {
	return fmt.Sprintf(
		"Message %s (%s %s): %d parts",
		m.Id, m.Name, m.Release, m.Count())
}

// Verbose output fo debugging
func (m *MessageSpec) Dump() string {
	count := m.Count()
	var buffer bytes.Buffer

	for i := 0; i < count; i++ {
		buffer.WriteString(m.Parts[i].String() + "\n")
	}
	return buffer.String()
}

// Number of parts
func (m *MessageSpec) Count() int {
	return len(m.Parts)
}

func NewMessageSpec(
	id string, name string,
	version string, release string, contrAgency string,
	revision string, date time.Time, source string,
	parts []MessageSpecPart) *MessageSpec {

	return &MessageSpec{
		Id: id, Name: name,
		Version: version, Release: release, ContrAgency: contrAgency,
		Revision: revision, Date: date, Source: source,
		Parts: parts,
	}
}
