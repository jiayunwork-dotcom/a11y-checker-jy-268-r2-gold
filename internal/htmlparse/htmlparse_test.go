package htmlparse

import (
	"io"
	"strings"
	"testing"
)

func TestParse_ExtractsTagsAndAttrs(t *testing.T) {
	doc := ParseString(`<html lang="en"><img src="a.png" alt="logo"></html>`)
	if len(doc.Tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(doc.Tags))
	}
	img := doc.Tags[1]
	if img.Name != "img" {
		t.Fatalf("second tag should be img, got %q", img.Name)
	}
	if img.AttrMap["alt"] != "logo" {
		t.Fatalf("alt attribute wrong: %q", img.AttrMap["alt"])
	}
	if doc.Tags[0].AttrMap["lang"] != "en" {
		t.Fatalf("html lang wrong: %q", doc.Tags[0].AttrMap["lang"])
	}
}

func TestParse_MissingAttrReturnsEmpty(t *testing.T) {
	doc := ParseString(`<img src="a.png">`)
	img := doc.Tags[0]
	if _, ok := img.AttrMap["alt"]; ok {
		t.Fatal("alt should be absent")
	}
	if v := img.AttrMap["alt"]; v != "" {
		t.Fatalf("missing attr should map to empty string, got %q", v)
	}
}

func TestInnerText_OfLink(t *testing.T) {
	doc := ParseString(`<a href="/x">click me</a>`)
	if doc.Tags[0].Name != "a" {
		t.Fatalf("expected a tag")
	}
	if got := doc.InnerText(0); got != "click me" {
		t.Fatalf("inner text = %q", got)
	}
}

func TestInnerText_Nested(t *testing.T) {
	doc := ParseString(`<a href="/x">outer <b>bold</b> end</a>`)
	if got := doc.InnerText(0); !strings.Contains(got, "outer") || !strings.Contains(got, "end") {
		t.Fatalf("inner text should include nested text, got %q", got)
	}
}

type failReader struct{}

func (failReader) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestParse_ReadError(t *testing.T) {
	_, err := Parse(failReader{})
	if err == nil {
		t.Fatal("expected error when the reader fails")
	}
}

func TestParse_FirstTagAttrsIndependent(t *testing.T) {
	doc := ParseString(`<img alt="one"><img alt="two">`)
	if len(doc.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(doc.Tags))
	}
	a0 := doc.Tags[0].Attrs
	if len(a0) == 0 || a0[0].Value != "one" {
		t.Fatalf("first img alt should stay one, attrs=%+v", a0)
	}
	if doc.Tags[1].AttrMap["alt"] != "two" {
		t.Fatalf("second img alt=%q", doc.Tags[1].AttrMap["alt"])
	}
}
