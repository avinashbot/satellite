package download

import (
	"testing"
	"time"
)

func TestModifiedSince(t *testing.T) {
	filename, err := Dscovr.ModifiedSince(time.Time{})
	if err != nil {
		t.Error(err)
		return
	}
	if filename == "" {
		t.Error("ModifiedSince should not return empty string in this case.")
	}
}

func TestDownload(t *testing.T) {
	if _, err := Dscovr.Download("epic_1b_20160204202938_00"); err != nil {
		t.Error(t)
	}
}
