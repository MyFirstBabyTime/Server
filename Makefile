.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o first-baby-time app/*.go

.PHONY: image
image:
	docker build . -t first-baby-time:${VERSION}
	docker tag first-baby-time:${VERSION} mspring03/first-baby-time:${VERSION}.RELEASE

.PHONY: upload
upload:
	docker push jinhong0719/first-baby-time:${VERSION}.RELEASE

.PHONY: stack
stack:
	env VERSION=${VERSION} docker stack deploy -c docker-compose.yml DSM_SMS
