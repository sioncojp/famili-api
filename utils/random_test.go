package utils

import (
	"testing"
)

func TestMakeRandomString(t *testing.T) {
	t.Parallel()
	random1 := MakeRandomString(10)
	random2 := MakeRandomString(10)

	if random1 == random2 {
		t.Error("not random string")
	}
}
