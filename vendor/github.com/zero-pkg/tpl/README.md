# tpl

[![build](https://github.com/zero-pkg/tpl/actions/workflows/ci.yml/badge.svg)](https://github.com/zero-pkg/tpl/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/zero-pkg/tpl)](https://goreportcard.com/report/github.com/zero-pkg/tpl)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/zero-pkg/tpl/blob/master/LICENSE)

Provides helpers on top of `html/template` to dynamically parse all templates from the specific directory and provides a unified interface for accessing them. In addition, the package provides the ability to use the Jinja2/Django like `{{ extends }}` tag.

## Install and update

`go get -u github.com/zero-pkg/tpl`

## How to use

```
tmpl := tpl.Must(tpl.New().ParseDir("templates", ".html"))

if err := tmpl.Execute(os.Stdout, "content.html", ""); err != nil {
    panic(err)
}
```

## Explaning extends tag

The `{{ extends }}` tag is the key here. It tells the package that this template “extends” another template.
When the package evaluates this template, it first locates the parent.
The extends tag should be the first tag in the template.

### Nesting extends

```
# parent.html
body: {{ block "content" . }}Hi from parent.{{ end }}

# child.html
{{ extends "parent.html" }}
{{ block "content" . }}Hi from child.{{ end }}

# grandchild.html
{{ extends "child.html" }}
{{ block "content" . }}Hi from grandchild.{{ end }}
```

Rendering `grandchild.html` will give body: Hi from grandchild.

## License

http://www.opensource.org/licenses/mit-license.php
