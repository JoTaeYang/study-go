package prof

import (
	"testing"
	"time"
)

func FuncHI() {
	final := Write()
	defer final()
	time.Sleep(1)
}
func TestProfFunction(t *testing.T) {
	InitProfile()

	FuncHI()

	FuncHI()

	FuncHI()

	FuncHI()

	FuncHI()
	time.Sleep(time.Second * 3)
	Read()
}
