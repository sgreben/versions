VERSION = 1.1.2

APP      := versions
PACKAGES := $(shell go list -f {{.Dir}} ./...)
GOFILES  := $(addsuffix /*.go,$(PACKAGES))
GOFILES  := $(wildcard $(GOFILES))

.PHONY: clean release docker docker-latest README.md

clean:
	rm -rf binaries/
	rm -rf release/

# go get -u github.com/github/hub
release: README.md zip
	git add vendor
	git add Gopkg.toml Gopkg.lock
	git add README.md
	git add Makefile
	git commit -m "Release $(VERSION)" || true
	git push
	hub release create $(VERSION) -m "$(VERSION)" -a release/$(APP)_$(VERSION)_osx_x86_64.tar.gz -a release/$(APP)_$(VERSION)_windows_x86_64.zip -a release/$(APP)_$(VERSION)_linux_x86_64.tar.gz -a release/$(APP)_$(VERSION)_osx_x86_32.tar.gz -a release/$(APP)_$(VERSION)_windows_x86_32.zip -a release/$(APP)_$(VERSION)_linux_x86_32.tar.gz -a release/$(APP)_$(VERSION)_linux_arm64.tar.gz

README.md:
	sed "s/\$${VERSION}/$(VERSION)/g;s/\$${APP}/$(APP)/g;" README.template.md > README.md

docker: binaries/linux_x86_64/$(APP)
	docker build -t quay.io/sergey_grebenshchikov/$(APP):v$(VERSION) .
	docker push quay.io/sergey_grebenshchikov/$(APP):v$(VERSION)

docker-latest: docker
	docker tag quay.io/sergey_grebenshchikov/$(APP):v$(VERSION) quay.io/sergey_grebenshchikov/$(APP):latest
	docker push quay.io/sergey_grebenshchikov/$(APP):latest

zip: release/$(APP)_$(VERSION)_osx_x86_64.tar.gz release/$(APP)_$(VERSION)_windows_x86_64.zip release/$(APP)_$(VERSION)_linux_x86_64.tar.gz release/$(APP)_$(VERSION)_osx_x86_32.tar.gz release/$(APP)_$(VERSION)_windows_x86_32.zip release/$(APP)_$(VERSION)_linux_x86_32.tar.gz release/$(APP)_$(VERSION)_linux_arm64.tar.gz

binaries: binaries/osx_x86_64/$(APP) binaries/windows_x86_64/$(APP).exe binaries/linux_x86_64/$(APP) binaries/osx_x86_32/$(APP) binaries/windows_x86_32/$(APP).exe binaries/linux_x86_32/$(APP)

release/$(APP)_$(VERSION)_osx_x86_64.tar.gz: binaries/osx_x86_64/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_osx_x86_64.tar.gz -C binaries/osx_x86_64 $(APP)

binaries/osx_x86_64/$(APP): $(GOFILES)
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o binaries/osx_x86_64/$(APP) ./cmd/$(APP)

release/$(APP)_$(VERSION)_windows_x86_64.zip: binaries/windows_x86_64/$(APP).exe
	mkdir -p release
	cd ./binaries/windows_x86_64 && zip -r -D ../../release/$(APP)_$(VERSION)_windows_x86_64.zip $(APP).exe

binaries/windows_x86_64/$(APP).exe: $(GOFILES)
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o binaries/windows_x86_64/$(APP).exe ./cmd/$(APP)

release/$(APP)_$(VERSION)_linux_x86_64.tar.gz: binaries/linux_x86_64/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_linux_x86_64.tar.gz -C binaries/linux_x86_64 $(APP)

binaries/linux_x86_64/$(APP): $(GOFILES)
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o binaries/linux_x86_64/$(APP) ./cmd/$(APP)

release/$(APP)_$(VERSION)_osx_x86_32.tar.gz: binaries/osx_x86_32/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_osx_x86_32.tar.gz -C binaries/osx_x86_32 $(APP)

binaries/osx_x86_32/$(APP): $(GOFILES)
	GOOS=darwin GOARCH=386 go build -ldflags "-X main.version=$(VERSION)" -o binaries/osx_x86_32/$(APP) ./cmd/$(APP)

release/$(APP)_$(VERSION)_windows_x86_32.zip: binaries/windows_x86_32/$(APP).exe
	mkdir -p release
	cd ./binaries/windows_x86_32 && zip -r -D ../../release/$(APP)_$(VERSION)_windows_x86_32.zip $(APP).exe

binaries/windows_x86_32/$(APP).exe: $(GOFILES)
	GOOS=windows GOARCH=386 go build -ldflags "-X main.version=$(VERSION)" -o binaries/windows_x86_32/$(APP).exe ./cmd/$(APP)

release/$(APP)_$(VERSION)_linux_x86_32.tar.gz: binaries/linux_x86_32/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_linux_x86_32.tar.gz -C binaries/linux_x86_32 $(APP)

binaries/linux_x86_32/$(APP): $(GOFILES)
	GOOS=linux GOARCH=386 go build -ldflags "-X main.version=$(VERSION)" -o binaries/linux_x86_32/$(APP) ./cmd/$(APP)

release/$(APP)_$(VERSION)_linux_arm64.tar.gz: binaries/linux_arm64/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_linux_arm64.tar.gz -C binaries/linux_arm64 $(APP)

binaries/linux_arm64/$(APP): $(GOFILES)
	GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o binaries/linux_arm64/$(APP) ./cmd/$(APP)
