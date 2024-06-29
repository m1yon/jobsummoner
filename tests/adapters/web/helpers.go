package web

import (
	"strings"
	"testing"
	"time"

	"github.com/go-rod/rod"
)

func findElement(t *testing.T, page *rod.Page, tag string, text string) *rod.Element {
	t.Helper()

	elem, err := page.Timeout(3*time.Second).ElementR(tag, text)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			t.Fatalf("findElement: timed out looking for element with the tag '%v' and the text '%v'", tag, text)
		}

		t.Fatal(err)
	}

	return elem
}

func assertTextExistsInTheDocument(t *testing.T, page *rod.Page, text string) {
	t.Helper()

	_, err := page.Timeout(3*time.Second).ElementR("*", text)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			t.Fatalf("assertTextExistsInTheDocument: could not find the text '%v' in the document", text)
		}

		t.Fatal(err)
	}
}
