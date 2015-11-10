PKGDIR=pkg
CACHEDIR=cache

.PHONY: all clean prep version build

all: clean build

clean:
	rm -rf $(PKGDIR)


version:
	$(eval VERSION := $(shell git describe --tags --dirty))

# While we are only building one platform right now, we are considering making iscenv work on windows, so why not just use our fairly standard cross-compile methodology
build: version
	mkdir -p $(PKGDIR)
	go build -ldflags "-X github.com/ontariosystems/iscenv.VERSION=$(VERSION)" -o $(PKGDIR)/iscenv .

# This is a temporary target until we sort out a good single Travis-like build system
build	@echo "[trusted]\nusers = $(shell stat -c "%u" .)\n" > /etc/mercurial/hgrc.d/trust.rc
