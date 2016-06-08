package main

import (
	"net/http"
	"os/exec"
	"os"
)


var (
	//log    = logging.MustGetLogger("main")
)

type publishHandler struct {
	root string
}

func publishServer(root string) http.Handler {
	return &publishHandler{root}
}

func (u *publishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	curdir, _ := os.Getwd()
	os.Chdir(u.root)
	out, err := exec.Command("hugo").CombinedOutput()
	if err != nil {
		log.Error(err)
	}
	log.Info(string(out))
	os.Chdir(curdir)

	//sort.Sort(byName(files))
	url := r.Header.Get("Referer")

	if url != "" {
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
