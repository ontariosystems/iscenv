PKGDIR=pkg
compile=GOOS=$(1) GOARCH=$(2) go build -ldflags "-X github.com/ontariosystems/iscenv/iscenv.Version=$(VERSION)" -o=$(PKGDIR)/iscenv
compile_plugin=GOOS=$(1) GOARCH=$(2) go build -o=$(PKGDIR)/iscenv-$(3)-$(4) plugins/$(3)/$(4)/*.go

.PHONY: all
all: clean build

.PHONY: clean
clean:
	rm -rf $(PKGDIR)

.PHONY: version
version:
	$(eval VERSION := $(shell git describe --tags --dirty))

# While we are only building one platform right now, we are considering making iscenv work on windows, so why not just use our fairly standard cross-compile methodology
.PHONY: build
build: version
	$(call compile,linux,amd64)
	echo PRODUCT_VERSION=$(VERSION) > pkg/versions.properties

.PHONY: build-external-test-plugin
build-external-test-plugin:
	$(call compile_plugin,linux,amd64,lifecycle,external-test-plugin)
