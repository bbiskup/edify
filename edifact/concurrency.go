package edifact

import (
	"runtime"
)

var NumThreads = runtime.NumCPU()

func init() {
	runtime.GOMAXPROCS(NumThreads)
}
