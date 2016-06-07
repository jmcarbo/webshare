run: build
	docker run -ti -v $$PWD/test:/static -p 8888:8888 webshare /go/bin/webshare --title "JMCA upload" /static 

build:
	docker build -t webshare .
