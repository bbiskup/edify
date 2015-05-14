package dataelement

type DataElementSpec interface {
	Id() string
	Name() string
	String() string
}
