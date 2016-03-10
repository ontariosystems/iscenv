PKGDIR=pkg
CACHEDIR=cache
compile=GO15VENDOREXPERIMENT=1 GOOS=$(1) GOARCH=$(2) go build -ldflags "-X github.com/ontariosystems/iscenv/internal/iscenv.Version=$(VERSION)" -o=$(PKGDIR)/$(1)/$(2)/bin/iscenv
compile_plugin=GO15VENDOREXPERIMENT=1 GOOS=$(1) GOARCH=$(2) go build -o=$(PKGDIR)/$(1)/$(2)/bin/iscenv-$(3)-$(4) plugins/$(3)/$(4)/*.go

.PHONY: all clean prep version build buid-iscenv build-versions-local build-license-key build-hgcache buid-homedir build-service-bindings buid-shm build-statler-key

all: clean build

clean:
	rm -rf $(PKGDIR)


version:
	$(eval VERSION := $(shell git describe --tags --dirty))

# While we are only building one platform right now, we are considering making iscenv work on windows, so why not just use our fairly standard cross-compile methodology
build: version build-iscenv build-versions-local build-license-key build-hgcache build-homedir build-service-bindings build-shm build-statler-key
	echo PRODUCT_VERSION=$(VERSION) > pkg/versions.properties

build-iscenv:
	$(call compile,linux,amd64)

build-versions-local:
	$(call compile_plugin,linux,amd64,versions,local)

build-license-key:
	$(call compile_plugin,linux,amd64,start,license-key)

build-hgcache:

build-homedir:
	$(call compile_plugin,linux,amd64,start,homedir)

build-service-bindings:
	$(call compile_plugin,linux,amd64,start,service-bindings)

build-shm:
	$(call compile_plugin,linux,amd64,start,shm)

build-statler-key:

# This is a temporary target until we sort out a good single Travis-like build system
build	@curl -Ss -o /usr/local/share/ca-certificates/os_root_ca.crt http://statler.ontsys.com/v2/OS%20Certificate%20Bundle/1.0/os_root_ca.crt
	@update-ca-certificates
	@echo "[trusted]\nusers = $(shell stat -c "%u" .)\n" > /etc/mercurial/hgrc.d/trust.rc
