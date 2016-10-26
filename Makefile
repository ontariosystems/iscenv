PKGDIR=pkg
CACHEDIR=cache
compile=GO15VENDOREXPERIMENT=1 GOOS=$(1) GOARCH=$(2) go build -ldflags "-X github.com/ontariosystems/iscenv/iscenv.Version=$(VERSION)" -o=$(PKGDIR)/iscenv
compile_plugin=GO15VENDOREXPERIMENT=1 GOOS=$(1) GOARCH=$(2) go build -o=$(PKGDIR)/iscenv-$(3)-$(4) plugins/$(3)/$(4)/*.go

.PHONY: all clean version build build-external-test-plugin publish

all: clean build

clean:
	rm -rf $(PKGDIR)


version:
	$(eval VERSION := $(shell git describe --tags --dirty))

# While we are only building one platform right now, we are considering making iscenv work on windows, so why not just use our fairly standard cross-compile methodology
build: version
	$(call compile,linux,amd64)
	echo PRODUCT_VERSION=$(VERSION) > pkg/versions.properties

build-external-test-plugin:
	$(call compile_plugin,linux,amd64,lifecycle,external-test-plugin)
