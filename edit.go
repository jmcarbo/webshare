package main

import (
	"net/http"
	"strings"
	"text/template"
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"
	"os"
	"fmt"
	"github.com/spf13/hugo/parser"

	//"github.com/op/go-logging"
)


var (
	//log    = logging.MustGetLogger("main")
)

type editHandler struct {
	root string
	tmpl string
}

func editServer(root, template string) http.Handler {
	return &editHandler{root, template}
}

func (u *editHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	dir := strings.TrimPrefix(r.URL.Path, "/edit/")
	filename := r.FormValue("filename")
	if filename != "" {
		dir = path.Join(dir, filename)
	}

	var body []byte

	targetfilename := path.Join(u.root, dir)
	if r.Method == "POST" {
		var y map[string]interface{}
		body, _ = ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &y)
		extension := filepath.Ext(targetfilename)
		log.Info(extension)
		var z interface{}
		if val, ok := y["frontmatter"]; ok {
			z = val
		} else {
			z = y
		}
		body, err := json.Marshal(z)
		if err != nil {
			log.Warning(err)
		}

		switch {
			case extension == ".json":
				err:=ioutil.WriteFile(targetfilename, body, 0644)
				if err != nil {
					log.Warning(err)
				}
			case extension == ".md":
				newbody := y["body"]
				body, _ = json.Marshal(z)
				result := fmt.Sprintf("%s\n\n%s", body, newbody)
				err:=ioutil.WriteFile(targetfilename, []byte(result), 0644)
				if err != nil {
					log.Warning(err)
				}
		}
	} else {
		if _, err := os.Stat(targetfilename); os.IsNotExist(err) {
			// path/to/whatever does not exist
		} else {
			f, err := os.Open(targetfilename)
			if err != nil {
				log.Warning(err)
			}
			defer f.Close()
			//body, err=ioutil.ReadFile(targetfilename)
			p, err := parser.ReadFrom(f)
			if err != nil {
				log.Warning(err)
			}
			metadata, err := p.Metadata()
			if err != nil {
				log.Warning(err)
			}
			j := make(map[string]interface{})
			j["frontmatter"]=metadata
			j["body"]=p.Content()
			body, _ = json.Marshal(j)
		}
	}

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
		StartValue      string
		Referer      string
	}{
		string(body),
		url,
	})
}
