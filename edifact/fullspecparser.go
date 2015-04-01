package edifact

import (
	"github.com/bbiskup/edify/edifact/codes"
	"log"
	"os"
	"strings"
)

// Parses all relevant parts of EDIFACT spec
type FullSpecParser struct {
	Version string
}

func NewFullSpecParser(version string, dir string) (*FullSpecParser, error) {
	const sepStr = string(os.PathSeparator)
	p := codes.NewCodesSpecParser()

	path := strings.Join([]string{
		dir, "uncl", "UNCL." + version,
	}, string(os.PathSeparator))

	codesSpecs, err := p.ParseSpecFile(path)
	if err != nil {
		return nil, err
	}
	log.Printf("Loaded %d codes specs", len(codesSpecs))
	return &FullSpecParser{version}, nil
}
