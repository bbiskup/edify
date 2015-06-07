package msg

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

const msg1 = `
UNH+1+ORDERS:D:96A:UN'
BGM+220+B10001'
DTM+4:20060620:102'
NAD+BY+++Bestellername+Strasse+Stadt++23436+xx'
LIN+1++Produkt Schrauben:SA'
QTY+1:1000'
UNS+S'
CNT+2:1'
UNT+9+1'
    `

func TestRawMsg1(t *testing.T) {
	p := NewParser()
	rawMessage, err := p.ParseRawMsg(msg1)
	require.NotNil(t, rawMessage)
	require.Nil(t, err)
	assert.Equal(t, 9, len(rawMessage.Segs))
	assert.Equal(t, "UNH", rawMessage.Segs[0].Id())
	assert.Equal(t, "UNT", rawMessage.Segs[8].Id())
	assert.Equal(t, "1", rawMessage.Segs[0].Elems[0].Values[0])
	assert.Equal(t, "ORDERS", rawMessage.Segs[0].Elems[1].Values[0])
	assert.Equal(t, "D", rawMessage.Segs[0].Elems[1].Values[1])
	assert.Equal(t, "9", rawMessage.Segs[8].Elems[0].Values[0])
}

func TestParseINVOIC(t *testing.T) {
	msgStr, err := ioutil.ReadFile("../../testdata/messages/INVOIC_1.txt")
	require.Nil(t, err)

	p := NewParser()
	rawMessage, err := p.ParseRawMsg(string(msgStr))
	require.Nil(t, err)
	require.NotNil(t, rawMessage)

	assert.Equal(t, "UNH", rawMessage.Segs[0].Id())
	assert.Equal(t, "UNT", rawMessage.Segs[len(rawMessage.Segs)-1].Id())

	expectedMultilineStr := "If 0% VAT is charged and your VAT ID number is displayed above, this is either an exempt or a reverse charge transaction."
	assert.Equal(t, expectedMultilineStr, rawMessage.Segs[4].Elems[3].Values[0])
}

var elemSpecs = []struct {
	elemStr  string
	expected []string
}{
	{"a", []string{"a"}},
	{"a:b", []string{"a", "b"}},
}

func TestParseElem(t *testing.T) {
	for _, spec := range elemSpecs {
		p := NewParser()
		dataElem := p.ParseElem(spec.elemStr)
		assert.Equal(t, spec.expected, dataElem.Values)
	}
}

const msgWithDataElemRepetition = `
UNH+1+XYZ:D:96A:UN'
ABC+220*221*222+B10001'
UNT+9+1'`

func TestMessageWithRepetitionSeparator(t *testing.T) {
	p := NewParser()
	rawMessage, err := p.ParseRawMsg(msgWithDataElemRepetition)
	assert.NotNil(t, err)
	assert.Nil(t, rawMessage)
}
