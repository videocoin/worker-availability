JOB_IMAGE = registry.videocoin.net/worker-availability/job
VERSION ?= $(shell git describe --tags)

.PHONY: build
build:
	go build -o ./build/job ./stats/job
	go build -o ./build/rep ./stats/rep

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: test
test:
	go test ./...

.PHONY: image
image:
	docker build -t ${JOB_IMAGE}:$(VERSION) -f _assets/Dockerfile .
	docker tag ${JOB_IMAGE}:$(VERSION) ${JOB_IMAGE}:latest

.PHONY: push
push:
	docker push ${JOB_IMAGE}:${VERSION}
	docker push ${JOB_IMAGE}:latest

.PHONY: release
release: image push

.PHONY: db
db:
	docker run -p 27017:27017 -ti mongo
