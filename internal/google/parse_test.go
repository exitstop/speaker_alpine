package google_test

import (
	"fmt"
	"testing"

	"github.com/exitstop/speaker_alpine/internal/google"
	"github.com/exitstop/speaker_alpine/internal/google/data_test"
)

// ParseGoogle test
// go test ./internal/google/parse_test.go -run TestParseGoogle -v
func TestParseGoogle(t *testing.T) {
	for i, it := range data_test.Text {
		parseText, err := google.ParseGoogle7("hello world. Hello world? Hello world?", it)
		if err != nil {
			t.Errorf("[%d]: %s", i, err.Error())
		}
		fmt.Printf("[%d][PARSE OK]: %s\n", i, parseText)
	}
}
