package tpl

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type tplFile struct {
	abspath string
	path    string
	parent  *string
	content []byte
}

var re = regexp.MustCompile(`(?m)\{\{(\s*)extends\s+["'](.+)["'](\s*)\}\}`)

func parseDir(dir string, exts ...string) (map[string]*tplFile, error) {
	e := make(map[string]struct{}, len(exts))
	for i := range exts {
		e[exts[i]] = struct{}{}
	}

	files := make(map[string]*tplFile)

	err := filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if _, ok := e[path.Ext(file)]; !ok {
			return nil
		}

		f, err := parseFile(dir, file)
		if err != nil {
			return err
		}

		if f != nil {
			files[f.path] = f
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func parseFile(dir string, file string) (*tplFile, error) {
	body, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	f := &tplFile{
		abspath: file,
		path:    strings.TrimPrefix(filepath.ToSlash(file), dir),
	}

	r := bufio.NewReader(bytes.NewReader(body))

	line, _, err := r.ReadLine()
	if err != nil && err != io.EOF {
		return nil, err
	}

	m := re.FindSubmatch(line)
	if m != nil {
		s := string(m[2])
		f.parent = &s
		f.content = body[len(line):]

		return f, nil
	}

	f.content = body

	return f, nil
}
