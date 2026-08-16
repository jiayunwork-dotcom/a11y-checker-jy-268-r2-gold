// Package htmlparse 提供轻量级的 HTML 标签提取能力。
// Go 标准库没有 DOM 解析器，这里使用最小化扫描提取标签名、属性与元素内文本，
// 足以支撑静态无障碍规则检查（不解析完整 DOM 树）。
package htmlparse

import (
	"io"
	"regexp"
	"strings"
)

// Attr 是一个 HTML 属性。
type Attr struct {
	Name     string
	Value    string
	HasValue bool
}

// Tag 是文档中的一个标签（开始 / 结束 / 自闭合）。
type Tag struct {
	Name    string
	IsClose bool
	IsSelf  bool
	Attrs   []Attr
	AttrMap map[string]string // 仅包含带值的属性，供规则按名查找
	Start   int               // 在源中的起始下标
	End     int               // 结束 '>' 之后的下标
}

// Document 是解析后的文档。
type Document struct {
	Tags []Tag
	Src  string
}

var reTag = regexp.MustCompile(`<(/?)([a-zA-Z][a-zA-Z0-9]*)([^>]*?)(/?)>`)
var reAttr = regexp.MustCompile(`([a-zA-Z_:][-a-zA-Z0-9_:.]*)(?:\s*=\s*"([^"]*)"|\s*=\s*'([^']*)'|\s*=\s*(\S+))?`)

// Parse 从 reader 读取 HTML 并提取标签序列。
func Parse(r io.Reader) (*Document, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ParseString(string(data)), nil
}

// ParseString 解析 HTML 字符串。
func ParseString(s string) *Document {
	doc := &Document{Src: s}
	locs := reTag.FindAllStringSubmatchIndex(s, -1)
	for _, loc := range locs {
		fullStart, fullEnd := loc[0], loc[1]
		// 可选捕获组匹配「空串」时，Go 会报告 start==end（而非 -1），
		// 因此必须用 end>start 判断该组是否真正捕获到 "/"。
		slash := loc[2] != -1 && loc[3] > loc[2] // 组1命中非空 => 结束标签
		nameStart, nameEnd := loc[4], loc[5]
		attrStart, attrEnd := loc[6], loc[7]
		self := loc[8] != -1 && loc[9] > loc[8] // 组4命中非空 => 自闭合
		name := strings.ToLower(s[nameStart:nameEnd])
		t := Tag{
			Name:    name,
			IsClose: slash,
			IsSelf:  self,
			Start:   fullStart,
			End:     fullEnd,
		}
		if !slash && attrStart != -1 {
			t.Attrs, t.AttrMap = parseAttrs(s[attrStart:attrEnd])
		} else {
			t.AttrMap = map[string]string{}
		}
		doc.Tags = append(doc.Tags, t)
	}
	return doc
}

func parseAttrs(s string) ([]Attr, map[string]string) {
	var attrs []Attr
	m := map[string]string{}
	for _, am := range reAttr.FindAllStringSubmatch(s, -1) {
		name := strings.ToLower(am[1])
		val := am[2]
		if val == "" {
			val = am[3]
		}
		if val == "" {
			val = am[4]
		}
		hasVal := am[2] != "" || am[3] != "" || am[4] != ""
		attrs = append(attrs, Attr{Name: name, Value: val, HasValue: hasVal})
		if hasVal {
			m[name] = val
		}
	}
	return attrs, m
}

// InnerText 返回第 idx 个标签到其匹配结束标签之间的纯文本（已去首尾空白）。
// 对自闭合或结束标签返回空字符串。
func (d *Document) InnerText(idx int) string {
	if idx < 0 || idx >= len(d.Tags) {
		return ""
	}
	open := d.Tags[idx]
	if open.IsClose || open.IsSelf {
		return ""
	}
	depth := 0
	for j := idx + 1; j < len(d.Tags); j++ {
		t := d.Tags[j]
		if t.Name != open.Name {
			continue
		}
		if t.IsClose {
			if depth == 0 {
				return strings.TrimSpace(d.Src[open.End:t.Start])
			}
			depth--
		} else {
			depth++
		}
	}
	return strings.TrimSpace(d.Src[open.End:])
}
