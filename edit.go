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

func isEditable(value interface{}) bool {
	t := value.(string)
	return strings.HasSuffix(strings.ToLower(t),".json") || strings.HasSuffix(strings.ToLower(t),".md")
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func findArchetype(root, file string) string {
	base := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	dir := strings.Split(filepath.Dir(file), "/")

	for _ , target := range dir {
		targetFileName := path.Join(root, "archetypes", target + ".json")
		if _, err := os.Stat(targetFileName); os.IsNotExist(err) {
			// path/to/whatever does not exist
		} else {
			return target
		}
	}


	targetFileName := path.Join(root, "archetypes", base + ".json")
	if _, err := os.Stat(targetFileName); os.IsNotExist(err) {
		return ""	// path/to/whatever does not exist
	} else {
		return base
	}
}

type editHandler struct {
	root string
	tmpl string
}

func editServer(root, template string) http.Handler {
	return &editHandler{root, template}
}

var defaultSchema = `{
          type: "object",
	  //format: "grid",
          title: "Blog Post",
          properties: {
	    frontmatter: {
		type: "object",
		properties: {
            		title: {
              			type: "string"
            		}
		}
	    },
            body: {
              type: "string",
              format: "html",
              //format: "markdown",
              options: {
                wysiwyg: true
              }
            }
          }
        }
`
func (u *editHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	dir := strings.TrimPrefix(r.URL.Path, cfgRoot + "/edit/")
	filename := r.FormValue("filename")
	if filename != "" {
		dir = path.Join(dir, filename)
	}

	schema := defaultSchema
	archetype := findArchetype(u.root, dir)
	if archetype != "" {
		log.Info("Found archetype: %s", archetype)
		schemaFileName := path.Join(u.root, "archetypes", archetype + ".json")
		targetSchema, err := ioutil.ReadFile(schemaFileName)
		if err != nil {
			log.Error(err)
		} else {
			schema = string(targetSchema)
		}

	} else {
		log.Info("Archetype not found")
	}

	var body []byte

	targetfilename := path.Join(u.root, dir)
	extension := filepath.Ext(targetfilename)
	if stringInSlice(extension, []string{".jpg",".gif"}) {
		http.ServeFile(w, r, targetfilename)
		return
	}

	if r.Method == "POST" {
		var y map[string]interface{}
		body, _ = ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &y)
		//log.Info(extension)
		//log.Info(body)

		switch {
			case extension == ".json":
				var z interface{}
				err2 := json.Unmarshal(body, &z)
				if err2 != nil {
					log.Error(err2)
				}
				log.Info("%v", z)
				body, _ = json.MarshalIndent(z, "", "   ")
				err:=ioutil.WriteFile(targetfilename, body, 0644)
				if err != nil {
					log.Warning(err)
				}
			case extension == ".md":
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
				newbody := y["body"]
				body, _ = json.MarshalIndent(z, "", "   ")
				result := fmt.Sprintf("%s\n\n%s", body, newbody)
				err=ioutil.WriteFile(targetfilename, []byte(result), 0644)
				if err != nil {
					log.Warning(err)
				}
		}

	} else {
		if _, err := os.Stat(targetfilename); os.IsNotExist(err) {
			// path/to/whatever does not exist
		} else {

			switch {
			case extension == ".json":
				body, err =ioutil.ReadFile(targetfilename)
				if err != nil {
					log.Error(err)
				}
			case extension == ".md":
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
				log.Warning(string(p.Content()))
				log.Warning(string(p.FrontMatter()))
				metadata, err := p.Metadata()
				if err != nil {
					log.Warning(err)
				}
				j := make(map[string]interface{})
				j["frontmatter"]=metadata
				j["body"]=string(p.Content())
				body, _ = json.Marshal(j)
			default:
				http.ServeFile(w, r, targetfilename)
				return
			}
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
	//url := r.Header.Get("Referer")
	url :=  path.Join(cfgRoot + "/ui" , path.Dir(dir))

	t.Execute(w, struct {
		FileName	string
		StartValue      string
		Referer		string
		Schema		string
		RootUrl         string
	}{
		dir,
		string(body),
		url,
		schema,
		cfgRoot,
	})
}
