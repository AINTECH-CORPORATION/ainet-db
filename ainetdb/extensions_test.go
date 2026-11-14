package ainetdb

import "testing"

func TestDefaultExtensions(t *testing.T) {
	exts := DefaultExtensions()
	if len(exts) != 4 {
		t.Fatalf("expected 4 default extensions, got %d", len(exts))
	}
	for _, ext := range exts {
		if ext.ID == "" || ext.Upstream == "" || !ext.Enabled {
			t.Fatalf("invalid extension descriptor: %+v", ext)
		}
	}
}
