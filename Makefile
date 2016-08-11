PKGDIR=pkg
CACHEDIR=cache
compile=GO15VENDOREXPERIMENT=1 GOOS=$(1) GOARCH=$(2) go build -ldflags "-X github.com/ontariosystems/iscenv/iscenv.Version=$(VERSION)" -o=$(PKGDIR)/$(1)/$(2)/bin/iscenv
compile_plugin=GO15VENDOREXPERIMENT=1 GOOS=$(1) GOARCH=$(2) go build -o=$(PKGDIR)/$(1)/$(2)/bin/iscenv-$(3)-$(4) plugins/$(3)/$(4)/*.go

.PHONY: all clean version build build-external-test-plugin

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
	$(call compile_plugin,linux,amd64,start,external-test-plugin)

# This is a temporary target until we sort out a good single Travis-like build system
build	@curl -Ss -o /usr/local/share/ca-certificates/os_root_ca.crt http://statler.ontsys.com/v2/OS%20Certificate%20Bundle/1.0/os_root_ca.crt
	@update-ca-certificates
	# On the build server, this will add the "jenkins" user to the mercurial trust, since the container is running as root and the mapped repo is owned by the "jenkins" user. 
	@echo "[trusted]\nusers = $(shell stat -c "%u" .)\n" > /etc/mercurial/hgrc.d/trust.rc
