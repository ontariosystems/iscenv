PKGDIR=pkg
CACHEDIR=cache
compile=GO15VENDOREXPERIMENT=1 GOOS=$(1) GOARCH=$(2) go build -ldflags "-X main.Version=$(VERSION)" -o=$(PKGDIR)/$(1)/$(2)/bin/iscenv_$(1)_$(2)$(3) .

.PHONY: all clean prep version build

all: clean build

clean:
	rm -rf $(PKGDIR)


version:
	$(eval VERSION := $(shell git describe --tags --dirty))

# While we are only building one platform right now, we are considering making iscenv work on windows, so why not just use our fairly standard cross-compile methodology
build: version
	mkdir -p $(PKGDIR)
	$(call compile,linux,amd64,)
	echo PRODUCT_VERSION=$(VERSION) > pkg/versions.properties

# This is a temporary target until we sort out a good single Travis-like build system
build	@curl -Ss -o /usr/local/share/ca-certificates/os_root_ca.crt http://statler.ontsys.com/v2/OS%20Certificate%20Bundle/1.0/os_root_ca.crt
	@update-ca-certificates
	@echo "[trusted]\nusers = $(shell stat -c "%u" .)\n" > /etc/mercurial/hgrc.d/trust.rc
