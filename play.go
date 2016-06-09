package main

import (
	"net/http"
	"strings"
	"text/template"
	//"path"

	//"github.com/op/go-logging"
)


var (
	//log    = logging.MustGetLogger("main")
)

type playHandler struct {
	root string
	tmpl string
}

func playServer(root, template string) http.Handler {
	return &playHandler{root, template}
}

func (u *playHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dir := strings.TrimPrefix(r.URL.Path, cfgRoot + "/play/")
	//dst := path.Join(u.root, dir)
	html, _ := Asset(u.tmpl)
	funcMap := template.FuncMap{
/*
		"humanizeBytes": humanizeBytes,
		"humanizeTime":  humanizeTime,
		"isZip":  isZip,
*/
	}
	t, err := template.New("").Funcs(funcMap).Parse(string(html))

	if err != nil {
		log.Warning("error %s", err)
	}

	//files, err := ioutil.ReadDir(path.Join(v.root, r.URL.Path))

	//sort.Sort(byName(files))
	url := r.Header.Get("Referer")

	t.Execute(w, struct {
		FileMP3      string
		Referer      string
	}{
		dir,
		url,
	})
}
