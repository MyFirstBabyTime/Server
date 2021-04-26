.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o first-baby-time app/*.go

.PHONY: image
image:
	docker build . -t first-baby-time:${VERSION}
	docker tag first-baby-time:${VERSION} mspring03/first-baby-time:${VERSION}.RELEASE

.PHONY: upload
upload:
	docker push mspring03/first-baby-time:${VERSION}.RELEASE

.PHONY: deploy
deploy:
	curl -X POST -H "User-Agent: linux bla bla" -H "Content-Type: application/json" \
	-d " \
	{ \
	 \"cloud_management_key\":\"dbaudcjf3116\", \
	 \"image\":\"mspring03/first-baby-time:1.0.0.RELEASE\"\
	} \
	" \
	http://54.180.121.144:80/redeploy
