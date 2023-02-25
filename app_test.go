package caddy_token_auth

import (
	"testing"
)

func TestFooBar(t *testing.T) {
	t.Run("/ping", func(t *testing.T) {
		got := 1
		want := 1
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
