package yandex

import (
	"fmt"
	"os"
	"testing"
)

var yx *Yandex

func TestMain(m *testing.M) {
	var err error
	yx, err = New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "initializing new err=%v", err)
		os.Exit(-1)
	}
	os.Exit(m.Run())
}
