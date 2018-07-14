UNIX_PLATFORMS := linux/amd64 linux/arm linux/arm64 darwin/amd64
UNIX_FILES=snap README.md LICENSE
WINDOWS_PLATFORMS := windows/amd64
GIT_HASH=$(shell git rev-parse HEAD)
GIT_TAG=$(shell git describe --tags $(git rev-list --tags --max-count=1) 2&>/dev/null || echo  "0.0.0")
HOSTOS=$(shell go env GOHOSTOS)
HOSTARCH=$(shell go env GOHOSTARCH)
GOARCH=$(shell go env GOARCH)
LDFLAGS=-X main.Commit=$(GIT_HASH) -X main.Version=$(GIT_TAG)
TAR = $(shell which tar)
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	TAR=$(shell which gnutar)
endif

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

.PHONY: all clean test $(UNIX_PLATFORMS)

all: test $(UNIX_PLATFORMS) $(WINDOWS_PLATFORMS)

$(UNIX_PLATFORMS): | main.go snap.go
	GOOS="$(os)" GOARCH=$(arch) go build -v -o build/$(os)/$(arch)/snap -ldflags "$(LDFLAGS)"
	cp LICENSE README.md build/$(os)/$(arch)
	mkdir -p dist
	$(TAR) -cvzf dist/snap-$(os)-$(arch)-$(GIT_TAG).tar.gz -C build/$(os)/$(arch) .


$(WINDOWS_PLATFORMS): | main.go snap.go
	GOOS="$(os)" GOARCH=$(arch) go build -v -o build/$(os)/$(arch)/snap.exe -ldflags "$(LDFLAGS)"
	mkdir -p dist
	cp LICENSE README.md build/$(os)/$(arch)
	zip -D dist/snap-$(os)-$(arch)-$(GIT_TAG).zip build/$(os)/$(arch)/*

clean:
	$(RM) -r build
	$(RM) -r dist