// Package rules 实现网页无障碍的静态检查规则，返回结构化问题列表。
package rules

import (
	"strconv"
	"strings"

	"a11y-checker/internal/htmlparse"
)

// Severity 表示问题严重等级。
type Severity string

const (
	Critical Severity = "critical"
	Serious  Severity = "serious"
	Moderate Severity = "moderate"
	Minor    Severity = "minor"
)

// Issue 是单条无障碍问题。
type Issue struct {
	Rule     string   `json:"rule"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
}

// CheckAll 运行全部规则，返回发现的问题。
func CheckAll(doc *htmlparse.Document) []Issue {
	var issues []Issue
	issues = append(issues, checkImgAlt(doc)...)
	issues = append(issues, checkHtmlLang(doc)...)
	issues = append(issues, checkInputLabel(doc)...)
	issues = append(issues, checkHeadingOrder(doc)...)
	issues = append(issues, checkLinkText(doc)...)
	return issues
}

func checkImgAlt(doc *htmlparse.Document) []Issue {
	var out []Issue
	for _, t := range doc.Tags {
		if t.IsClose || t.Name != "img" {
			continue
		}
		if _, ok := t.AttrMap["alt"]; !ok {
			out = append(out, Issue{Rule: "img-alt", Severity: Serious, Message: "img 元素缺少 alt 属性"})
		}
	}
	return out
}

func checkHtmlLang(doc *htmlparse.Document) []Issue {
	for _, t := range doc.Tags {
		if t.IsClose || t.Name != "html" {
			continue
		}
		if _, ok := t.AttrMap["lang"]; !ok {
			return []Issue{{Rule: "html-lang", Severity: Serious, Message: "html 元素缺少 lang 属性"}}
		}
		return nil
	}
	return nil
}

func checkInputLabel(doc *htmlparse.Document) []Issue {
	labelFor := map[string]bool{}
	for _, t := range doc.Tags {
		if t.IsClose || t.Name != "label" {
			continue
		}
		if v, ok := t.AttrMap["for"]; ok && v != "" {
			labelFor[v] = true
		}
	}
	var out []Issue
	for _, t := range doc.Tags {
		if t.IsClose || t.Name != "input" {
			continue
		}
		if t.AttrMap["type"] == "hidden" {
			continue
		}
		if _, ok := t.AttrMap["aria-label"]; ok {
			continue
		}
		if _, ok := t.AttrMap["aria-labelledby"]; ok {
			continue
		}
		id := t.AttrMap["id"]
		if id == "" || !labelFor[id] {
			out = append(out, Issue{Rule: "input-label", Severity: Serious, Message: "input 缺少可访问名称（label[for] 或 aria-label）"})
		}
	}
	return out
}

func checkHeadingOrder(doc *htmlparse.Document) []Issue {
	var out []Issue
	prev := 0
	first := true
	for _, t := range doc.Tags {
		if t.IsClose {
			continue
		}
		l, ok := headingLevel(t.Name)
		if !ok {
			continue
		}
		if first {
			if l != 1 {
				out = append(out, Issue{Rule: "heading-order", Severity: Moderate, Message: "页面应从一个 h1 标题开始"})
			}
			first = false
		} else if l > prev+1 {
			out = append(out, Issue{Rule: "heading-order", Severity: Moderate, Message: "标题层级从 h" + strconv.Itoa(prev) + " 跳到 h" + strconv.Itoa(l)})
		}
		prev = l
	}
	return out
}

func headingLevel(name string) (int, bool) {
	if len(name) == 2 && name[0] == 'h' && name[1] >= '1' && name[1] <= '6' {
		return int(name[1] - '0'), true
	}
	return 0, false
}

func checkLinkText(doc *htmlparse.Document) []Issue {
	var out []Issue
	for i, t := range doc.Tags {
		if t.IsClose || t.Name != "a" {
			continue
		}
		if _, ok := t.AttrMap["href"]; !ok {
			continue
		}
		if strings.TrimSpace(doc.InnerText(i)) == "" {
			out = append(out, Issue{Rule: "link-text", Severity: Moderate, Message: "链接缺少可读文本"})
		}
	}
	return out
}
