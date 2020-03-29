BUILDDIR=build/
CHANGELOG=$(BUILDDIR)/CHANGELOG.md

GOPROXY=https://mirrors.aliyun.com/goproxy/
# UPDATE_SUBMODULE:=$(shell git submodule update --init --remote)
GO=CGO_ENABLED=0 GO111MODULE=on go

# office registry
DOCKER_IMAGE_NAME_GST_RK=registry.jiangxingai.com:5000/edgex/device-service/cameras-gst-rkmpp
DOCKER_IMAGE_NAME_FFMPEG=registry.jiangxingai.com:5000/edgex/device-service/cameras-ffmpeg
# harbor测试环境
HARBOR_IMAGE_NAME_GST_RK=harbor.jiangxingai.iotedge/library/edgex-cameras-gst-rkmpp/$(ARCHTAG)/others
HARBOR_IMAGE_NAME_FFMPEG=harbor.jiangxingai.iotedge/library/edgex-cameras-ffmpeg/$(ARCHTAG)/others
# harbor线上环境
HARBOR_IOTEDGE_IMAGE_NAME_GST_RK=harbor.jiangxingai.com/library/edgex-cameras-gst-rkmpp/$(ARCHTAG)/others
HARBOR_IOTEDGE_IMAGE_NAME_FFMPEG=harbor.jiangxingai.com/library/edgex-cameras-ffmpeg/$(ARCHTAG)/others

GOARCH=arm64
ARCHTAG=arm64v8
SRC_PATH=github.com/edgex-camera/device-cameras

VERSION=$(shell git tag -l "v*" --points-at HEAD | tail -n 1 | tail -c +2)
GIT_SHA=$(shell git rev-parse HEAD)
GOFLAGS=-ldflags "-X $(SRC_PATH).Version=$(VERSION)"

MICROSERVICES=device-service
.PHONY: $(MICROSERVICES) frontend

check_version:
ifeq ($(VERSION),)
	$(error No version tag found)
endif
ifeq ($(shell git cat-file -t v$(VERSION)),commit)
	$(error Changelog should be in tag message)
endif

build: check_version $(MICROSERVICES) frontend changelog

build-amd64: GOARCH=amd64 build

build-arm64: GOARCH=arm64 build

device-service: | $(BUILDDIR)
	GOARCH=$(GOARCH) $(GO) build $(GOFLAGS) -o $(BUILDDIR)/bin/$@-$(GOARCH) ./cmd/device-service
	cp -r cmd/device-service/res $(BUILDDIR)/bin/

changelog: check_version | $(BUILDDIR)
	echo "# Changelog\n" > $(CHANGELOG)
	git tag --sort=-taggerdate --format='## %(tag) - %(taggerdate:short)%0a### Author: %(taggername) %(taggeremail)%0a%(contents)%0a' >> $(CHANGELOG)

$(BUILDDIR):
	mkdir -p $(BUILDDIR)

docker-arm64: GOARCH=arm64
docker-arm64: ARCHTAG=arm64v8
docker-arm64: docker

docker: docker-gst-rk docker-ffmpeg

docker-gst-rk: frontend build
	docker build \
				--label "git_sha=$(GIT_SHA)" \
		--build-arg GOARCH=$(GOARCH) \
		-f Dockerfiles/Dockerfile-gst-rkmpp \
		-t $(DOCKER_IMAGE_NAME_GST_RK):$(GIT_SHA) \
		-t $(DOCKER_IMAGE_NAME_GST_RK):$(ARCHTAG)-cpu-$(VERSION) \
		-t $(DOCKER_IMAGE_NAME_GST_RK):$(ARCHTAG)-cpu-latest \
		build

docker-ffmpeg: frontend build
		docker build \
		--label "git_sha=$(GIT_SHA)" \
		--build-arg GOARCH=$(GOARCH) \
		-f Dockerfiles/Dockerfile-ffmpeg \
		-t $(DOCKER_IMAGE_NAME_FFMPEG):$(GIT_SHA) \
		-t $(DOCKER_IMAGE_NAME_FFMPEG):$(ARCHTAG)-cpu-$(VERSION) \
		-t $(DOCKER_IMAGE_NAME_FFMPEG):$(ARCHTAG)-cpu-latest \
		build

deploy-arm64: ARCHTAG=arm64v8
deploy-arm64: deploy

deploy: check_version
	docker push $(DOCKER_IMAGE_NAME_GST_RK):$(ARCHTAG)-cpu-$(VERSION)
	docker push $(DOCKER_IMAGE_NAME_GST_RK):$(ARCHTAG)-cpu-latest
	docker push $(DOCKER_IMAGE_NAME_FFMPEG):$(ARCHTAG)-cpu-$(VERSION)
	docker push $(DOCKER_IMAGE_NAME_FFMPEG):$(ARCHTAG)-cpu-latest

deploy-harbor-arm64: ARCHTAG=arm64v8
deploy-harbor-arm64: deploy-harbor

deploy-harbor:
	docker login harbor.jiangxingai.iotedge --username $(DOCKER_USERNAME) --password $(DOCKER_PASSWORD)
	docker tag $(DOCKER_IMAGE_NAME_GST_RK):$(ARCHTAG)-cpu-$(VERSION) $(HARBOR_IMAGE_NAME_GST_RK):$(VERSION)
	docker push $(HARBOR_IMAGE_NAME_GST_RK):$(VERSION)
	docker tag $(DOCKER_IMAGE_NAME_FFMPEG):$(ARCHTAG)-cpu-$(VERSION) $(HARBOR_IMAGE_NAME_FFMPEG):$(VERSION)
	docker push $(HARBOR_IMAGE_NAME_FFMPEG):$(VERSION)

deploy-harbor-iotedge-arm64: ARCHTAG=arm64v8
deploy-harbor-iotedge-arm64: deploy-harbor-iotedge

deploy-harbor-iotedge:
	docker login harbor.jiangxingai.com --username $(IOTEDGE_DOCKER_USERNAME) --password $(IOTEDGE_DOCKER_PASSWORD)
	docker tag $(DOCKER_IMAGE_NAME_GST_RK):$(ARCHTAG)-cpu-$(VERSION) $(HARBOR_IOTEDGE_IMAGE_NAME_GST_RK):$(VERSION)
	docker push $(HARBOR_IOTEDGE_IMAGE_NAME_GST_RK):$(VERSION)
	docker tag $(DOCKER_IMAGE_NAME_FFMPEG):$(ARCHTAG)-cpu-$(VERSION) $(HARBOR_IOTEDGE_IMAGE_NAME_FFMPEG):$(VERSION)
	docker push $(HARBOR_IOTEDGE_IMAGE_NAME_FFMPEG):$(VERSION)

test:
	$(GO) vet ./...
	gofmt -l .
	$(GO) test -coverprofile=coverage.out ./...

clean:
	rm -rf build/

frontend:
	npm config set registry http://npm.registry.jiangxingai.com:7001/
	cd frontend; \
		npm i; \
		CI=false npm run build
	mkdir -p build
	rm -rf build/frontend
	mv frontend/build build/frontend
