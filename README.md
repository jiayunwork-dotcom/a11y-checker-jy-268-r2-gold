# a11y-checker

网页无障碍（Accessibility）静态合规检测命令行工具（纯标准库、离线、本地 HTML 文件）。

## 能力

针对本地 HTML 文件做静态规则检查（不抓 URL、不依赖无头浏览器）：

- `img-alt`：img 元素缺少 alt 属性
- `html-lang`：根 html 元素缺少 lang 属性
- `input-label`：input 缺少可访问名称（label[for] 或 aria-label）
- `heading-order`：标题层级跳跃 / 页面未从 h1 开始
- `link-text`：链接缺少可读文本

报告支持终端文本与 JSON 两种格式，按严重等级（critical/serious/moderate/minor）分类。

## 用法

```bash
# 文本报告
a11y-checker check example/sample.html

# JSON 报告
a11y-checker check example/sample.html --format json
```

无参数或未知子命令会打印用法并以非零退出码结束（受控退出，不崩溃）。
