// Package report 提供无障碍检查报告的文本 / JSON 渲染。
package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"a11y-checker/internal/rules"
)

// Report 是检查报告。
type Report struct {
	Issues []rules.Issue
}

// Counts 按严重等级统计问题数量。
func (r *Report) Counts() map[string]int {
	c := map[string]int{}
	for _, i := range r.Issues {
		c[string(i.Severity)]++
	}
	return c
}

// RenderText 以可读文本形式输出报告。
func (r *Report) RenderText(w io.Writer) {
	if len(r.Issues) == 0 {
		fmt.Fprintln(w, "no issues found")
		return
	}
	for _, i := range r.Issues {
		fmt.Fprintf(w, "[%s] %s: %s\n", strings.ToUpper(string(i.Severity)), i.Rule, i.Message)
	}
	c := r.Counts()
	fmt.Fprintf(w, "total: %d (critical=%d serious=%d moderate=%d minor=%d)\n",
		len(r.Issues), c["critical"], c["serious"], c["moderate"], c["minor"])
}

// RenderJSON 以 JSON 形式输出问题列表。
func (r *Report) RenderJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r.Issues)
}
