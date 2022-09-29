compile=GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 go build -tags netgo -ldflags "-X github.com/ontariosystems/iscenv/iscenv.Version=$(version)" -o=pkg/iscenv
compile_plugin=GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 go build -tags netgo -o=pkg/iscenv-$(3)-$(4) plugins/$(3)/$(4)/*.go

.PHONY: all
all: clean build

.PHONY: clean
clean:
	rm -rf pkg

.PHONY: version
version:
	$(eval version := $(shell git describe --tags --dirty | sed 's/^v//'))

pkg:
	mkdir -p $@

# While we are only building one platform right now, we are considering making iscenv work on windows, so why not just use our fairly standard cross-compile methodology
.PHONY: build
build: version | pkg
	$(call compile,linux,amd64)
	echo PRODUCT_VERSION=$(version) > pkg/versions.properties

.PHONY: build-external-test-plugin
build-external-test-plugin: | pkg
	$(call compile_plugin,linux,amd64,lifecycle,external-test-plugin)
