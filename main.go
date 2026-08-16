package main

import (
	"fmt"
	"io"
	"os"

	"a11y-checker/internal/htmlparse"
	"a11y-checker/internal/report"
	"a11y-checker/internal/rules"
)

func usage(w io.Writer) {
	fmt.Fprintln(w, "usage: a11y-checker check <file.html> [--format text|json]")
}

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		os.Exit(1)
	}
}

// run 是可测试的入口：解析参数、读取文件、运行规则并渲染报告。
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) < 1 || args[0] != "check" {
		usage(stderr)
		return fmt.Errorf("missing or unknown command")
	}
	rest := args[1:]
	var file, format string
	format = "text"
	for i := 0; i < len(rest); i++ {
		switch rest[i] {
		case "--format":
			if i+1 < len(rest) {
				format = rest[i+1]
				i++
			}
		default:
			if file == "" {
				file = rest[i]
			}
		}
	}
	if file == "" {
		usage(stderr)
		return fmt.Errorf("check requires a file")
	}
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	doc, err := htmlparse.Parse(f)
	if err != nil {
		return err
	}
	issues := rules.CheckAll(doc)
	rep := &report.Report{Issues: issues}
	switch format {
	case "json":
		return rep.RenderJSON(stdout)
	default:
		rep.RenderText(stdout)
	}
	return nil
}
