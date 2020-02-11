build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o auth-get

image: build
	docker build -t clayz95/auth-get .

push: image
	docker push clayz95/auth-get:latest

image-tar: image
	docker save -o auth-get.tar clayz95/auth-get