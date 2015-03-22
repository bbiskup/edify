package tests

import (
	edi "edifact_experiments/edifact"
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

	p := edi.NewParser()
	message, err := p.ParseMessage(msg1)
	if err != nil {
		t.Fatalf("Parse error: %s", err)
	}

	if message == nil {
		t.Fatalf("Parser returned nil message")
	}
}
