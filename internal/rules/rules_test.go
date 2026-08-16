package rules

import (
	"testing"

	"a11y-checker/internal/htmlparse"
)

func TestCheckImgAlt(t *testing.T) {
	doc := htmlparse.ParseString(`<img src="a.png">`)
	if issues := CheckAll(doc); len(issues) != 1 || issues[0].Rule != "img-alt" {
		t.Fatalf("expected 1 img-alt issue, got %+v", issues)
	}
	doc2 := htmlparse.ParseString(`<img src="a.png" alt="logo">`)
	if issues := checkImgAlt(doc2); len(issues) != 0 {
		t.Fatalf("img with alt should have no issue, got %+v", issues)
	}
}

func TestCheckHtmlLang(t *testing.T) {
	doc := htmlparse.ParseString(`<html><body></body></html>`)
	if issues := checkHtmlLang(doc); len(issues) != 1 {
		t.Fatalf("expected 1 html-lang issue, got %d", len(issues))
	}
	doc2 := htmlparse.ParseString(`<html lang="en"></html>`)
	if issues := checkHtmlLang(doc2); len(issues) != 0 {
		t.Fatalf("html with lang should have no issue, got %+v", issues)
	}
}

func TestCheckInputLabel_Missing(t *testing.T) {
	doc := htmlparse.ParseString(`<input type="text">`)
	issues := checkInputLabel(doc)
	if len(issues) != 1 {
		t.Fatalf("expected 1 input-label issue, got %+v", issues)
	}
}

func TestCheckHeadingOrder_StartsAtH2(t *testing.T) {
	doc := htmlparse.ParseString(`<h2>Section</h2>`)
	issues := checkHeadingOrder(doc)
	if len(issues) == 0 {
		t.Fatal("expected heading-order issue when page starts at h2")
	}
}

func TestCheckInputLabel_Hidden(t *testing.T) {
	doc := htmlparse.ParseString(`<input type="hidden" name="csrf">`)
	if issues := checkInputLabel(doc); len(issues) != 0 {
		t.Fatalf("hidden input should have no issue, got %+v", issues)
	}
}

func TestCheckInputLabel_WithLabel(t *testing.T) {
	doc := htmlparse.ParseString(`<label for="n">Name</label><input type="text" id="n">`)
	if issues := checkInputLabel(doc); len(issues) != 0 {
		t.Fatalf("input with matching label should have no issue, got %+v", issues)
	}
}

func TestCheckHeadingOrder_Skip(t *testing.T) {
	doc := htmlparse.ParseString(`<h1>Main</h1><h3>Skip</h3>`)
	issues := checkHeadingOrder(doc)
	found := false
	for _, i := range issues {
		if i.Rule == "heading-order" && i.Message != "" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a heading-order issue for h1->h3 skip, got %+v", issues)
	}
}

func TestCheckLinkText_Empty(t *testing.T) {
	doc := htmlparse.ParseString(`<a href="/x"></a>`)
	if issues := checkLinkText(doc); len(issues) != 1 {
		t.Fatalf("expected 1 link-text issue, got %+v", issues)
	}
	doc2 := htmlparse.ParseString(`<a href="/x">home</a>`)
	if issues := checkLinkText(doc2); len(issues) != 0 {
		t.Fatalf("link with text should have no issue, got %+v", issues)
	}
}
