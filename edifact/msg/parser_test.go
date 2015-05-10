package msg

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestParser1(t *testing.T) {

	p := NewParser()
	message, err := p.ParseMessage(msg1)
	require.NotNil(t, message)
	require.Nil(t, err)
	assert.Equal(t, 9, len(message.Segments))
	assert.Equal(t, "UNH", message.Segments[0].Name)
	assert.Equal(t, "UNT", message.Segments[8].Name)
	assert.Equal(t, "1", message.Segments[0].Elements[0].Values[0])
	assert.Equal(t, "ORDERS", message.Segments[0].Elements[1].Values[0])
	assert.Equal(t, "D", message.Segments[0].Elements[1].Values[1])
	assert.Equal(t, "9", message.Segments[8].Elements[0].Values[0])
}
