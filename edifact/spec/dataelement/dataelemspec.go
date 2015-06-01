package dataelement

type DataElemSpec interface {
	Id() string
	Name() string
	String() string
}
