package segment

import (
	"fmt"
	"testing"
)

func TestParseSpec(t *testing.T) {
	_ = NewSegmentSpecParser(nil, nil)
}

const dataElemSpecStr = "010    C817 ADDRESS USAGE                              C    1"

func TestParseDataElemSpec(t *testing.T) {
	p := NewSegmentSpecParser(nil, nil)

	pos, id, dataElementKind, count, isMandatory, err := p.parseDataElemSpec(dataElemSpecStr)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error parsing data elem spec: %s", err))
	}

	if pos != 10 {
		t.Errorf(fmt.Sprintf("Incorrect position: %d", pos))
	}
	if id != "C817" {
		t.Errorf(fmt.Sprintf("Incorrect ID: %s", id))
	}

	if count != 1 {
		t.Errorf(fmt.Sprintf("Incorrect count: %d", count))
	}

	if isMandatory {
		t.Errorf("Should be condiional")
	}

	if dataElementKind != Composite {
		t.Errorf("Should be composite")
	}
}
