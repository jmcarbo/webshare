run: build
	docker run -ti -v $$PWD/test:/static -p 8888:8888 webshare /go/bin/webshare --title "JMCA upload" /static 

build:
	docker build -t webshare .

justrun: 
	docker run -ti -v $$PWD/test:/static -p 8888:8888 webshare /go/bin/webshare --title "JMCA upload" /static 

dev: 
	docker run -ti --rm -v $$PWD:/go/src/github.com/jmcarbo/webshare jmcarbo/golangdev /bin/bash


compile:
	#go get github.com/jteeuwen/go-bindata/...
	#go get github.com/elazarl/go-bindata-assetfs/...
	go-bindata-assetfs static/...
	#go get ./...
	go build

