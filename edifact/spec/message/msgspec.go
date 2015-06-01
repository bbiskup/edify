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
type MsgSpec struct {
	Id   string
	Name string

	Version     string
	Release     string
	ContrAgency string
	Revision    string
	Date        time.Time

	Source        string
	TopLevelGroup *MsgSpecSegGrpPart
}

type MsgSpecs []*MsgSpec

func (m *MsgSpec) String() string {
	var partsStr = m.PartsStr()
	if len(partsStr) > 0 {
		partsStr = " - " + partsStr
	}
	return fmt.Sprintf(
		"Message %s (%s %s): %d parts%s",
		m.Id, m.Name, m.Release, m.Count(), partsStr)
}

func (m *MsgSpec) TopLevelParts() []MsgSpecPart {
	return m.TopLevelGroup.Children()
}

func (m *MsgSpec) TopLevelPart(index int) MsgSpecPart {
	return m.TopLevelGroup.Children()[index]
}

// Verbose output fo debugging
func (m *MsgSpec) Dump() string {
	count := m.Count()
	var buffer bytes.Buffer

	for i := 0; i < count; i++ {
		buffer.WriteString(m.TopLevelParts()[i].String() + "\n")
	}
	return buffer.String()
}

func (m *MsgSpec) PartsStr() string {
	result := []string{}
	for _, part := range m.TopLevelParts() {
		result = append(result, part.Id())
	}
	return strings.Join(result, ", ")
}

// Number of parts
func (m *MsgSpec) Count() int {
	return len(m.TopLevelParts())
}

// from sort.Interface
func (m MsgSpecs) Len() int {
	return len(m)
}

// from sort.Interface
func (m MsgSpecs) Less(i, j int) bool {
	return m[i].Id < m[j].Id
}

// from sort.Interface
func (m MsgSpecs) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func NewMsgSpec(
	id string, name string,
	version string, release string, contrAgency string,
	revision string, date time.Time, source string,
	parts []MsgSpecPart) *MsgSpec {

	return &MsgSpec{
		Id: id, Name: name,
		Version: version, Release: release, ContrAgency: contrAgency,
		Revision: revision, Date: date, Source: source,
		TopLevelGroup: NewMsgSpecSegGrpPart(
			TopLevelSegGroupName, parts, 1, true, nil),
	}
}
