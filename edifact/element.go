package edifact

type Element struct {
	Name  string
	Value string
}

func (e *Element) String() string {
	return e.Name + " " + e.Value
}

func NewElement(name string, value string) *Element {
	return &Element{name, value}
}
