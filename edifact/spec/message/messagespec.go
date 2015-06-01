package message

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

const TopLevelSegGroupName = "_Group_0"

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

	Source        string
	TopLevelGroup *MessageSpecSegmentGroupPart
}

type MessageSpecs []*MessageSpec

func (m *MessageSpec) String() string {
	var partsStr = m.PartsStr()
	if len(partsStr) > 0 {
		partsStr = " - " + partsStr
	}
	return fmt.Sprintf(
		"Message %s (%s %s): %d parts%s",
		m.Id, m.Name, m.Release, m.Count(), partsStr)
}

func (m *MessageSpec) TopLevelParts() []MessageSpecPart {
	return m.TopLevelGroup.Children()
}

func (m *MessageSpec) TopLevelPart(index int) MessageSpecPart {
	return m.TopLevelGroup.Children()[index]
}

// Verbose output fo debugging
func (m *MessageSpec) Dump() string {
	count := m.Count()
	var buffer bytes.Buffer

	for i := 0; i < count; i++ {
		buffer.WriteString(m.TopLevelParts()[i].String() + "\n")
	}
	return buffer.String()
}

func (m *MessageSpec) PartsStr() string {
	result := []string{}
	for _, part := range m.TopLevelParts() {
		result = append(result, part.Id())
	}
	return strings.Join(result, ", ")
}

// Number of parts
func (m *MessageSpec) Count() int {
	return len(m.TopLevelParts())
}

// from sort.Interface
func (m MessageSpecs) Len() int {
	return len(m)
}

// from sort.Interface
func (m MessageSpecs) Less(i, j int) bool {
	return m[i].Id < m[j].Id
}

// from sort.Interface
func (m MessageSpecs) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
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
		TopLevelGroup: NewMessageSpecSegmentGroupPart(
			TopLevelSegGroupName, parts, 1, true, nil),
	}
}
