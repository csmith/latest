package latest

import (
	"strings"
	"testing"

	"github.com/csmith/latest/v3/internal"
)

func TestTransform(t *testing.T) {
	o := internal.ApplyDefaults(&defaultTagOptions, &TagOptions{
		IgnorePreRelease: true,
		TrimPrefixes:     []string{"curl-"},
		Transform: func(s string) string {
			return strings.ReplaceAll(s, "_", ".")
		},
	})

	tag, err := o.latest([]string{
		"curl-7_88_1",
		"curl-8_12_1",
		"curl-8_13_0",
		"curl-8_14_0-rc1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag != "curl-8_13_0" {
		t.Errorf("got tag %q, want %q", tag, "curl-8_13_0")
	}
}
