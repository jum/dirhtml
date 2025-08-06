/*
 * Make simple HTML file for a directory.
 * jum@anubis.han.de
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

var (
	dirListTemplate = template.Must(template.New("dirList").Funcs(
		template.FuncMap{
			"eq": func(a, b string) bool {
				return a == b
			},
			"sz": func(a int64) ByteSize {
				return ByteSize(a)
			},
		}).Parse(`<!doctype html>
<html>
<head>
   <title>{{.Title}}</title>
   <meta name=color-scheme content="light dark">
</head>
<body>
	<table border="1" cellpadding="5">
	<tr>
		<th>Name</th>
		<th>Last Modified</th>
		<th>Size</th>
	</tr>
	<tr>
		<td><a href="..">Parent Directory</a></td>
		<td></td>
		<td>-</td>
	</tr>
	{{range .Files}}
		{{if not (eq .Name "index.html")}}
		<tr>
			<td><a href="{{.Name}}">{{.Name}}</a></td>
			<td>{{.ModTime}}</td>
			<td>{{sz .Size}}</td>
		</tr>
		{{end}}
	{{end}}
	</table>
</body>
</html>
`))
)

const DEBUG = false

func debug(format string, a ...interface{}) {
	if DEBUG {
		fmt.Printf(format, a...)
	}
}

func main() {
	var (
		//err error
		data struct {
			Title string
			Files []os.FileInfo
		}
	)
	flag.Parse()
	for _, d := range flag.Args() {
		debug("d %#v\n", d)
		absDir, err := filepath.Abs(d)
		if err != nil {
			panic(err.Error())
		}
		debug("absDir %#v\n", absDir)
		data.Title = filepath.Base(absDir)
		f, err := os.Open(absDir)
		defer f.Close()
		if err != nil {
			panic(err.Error())
		}
		data.Files, err = f.Readdir(0)
		if err != nil {
			panic(err.Error())
		}
		debug("data %#v\n", data)
		index, err := os.Create(filepath.Join(absDir, "index.html"))
		if err != nil {
			panic(err.Error())
		}
		defer index.Close()
		w := bufio.NewWriter(index)
		defer w.Flush()
		err = dirListTemplate.Execute(w, data)
		if err != nil {
			panic(err.Error())
		}
	}
}

type ByteSize float64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}
