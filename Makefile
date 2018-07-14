.PNONY:  dist clean

GIT_HASH=$(shell git rev-parse HEAD)
HOSTOS=$(shell go env GOHOSTOS)
LDFLAGS=-X main.Commit=$(GIT_HASH)
TAR=gnutar

build/snap.linux: main.go snap.go
	GOOS="linux" go build -v -o $@ -ldflags "$(LDFLAGS)" 

build/snap.darwin: main.go snap.go
	GOOS="darwin" go build -v -o $@ -ldflags "$(LDFLAGS)" 

build/snap.exe: main.go snap.go
	GOOS=windows go build -o $@ -v -ldflags "$(LDFLAGS)"

dist: dist/snap-linux-distribution.tar.gz dist/snap-windows-distribution.zip

build/snap.man: build/snap.$(HOSTOS)
	./$< --help-man > build/snap.man

dist/snap-linux-distribution.tar.gz: build/snap.linux README.md LICENSE build/snap.man
	mkdir -p dist
	$(TAR) --xform="s/build\///" --xform="s/snap.linux/snap/" --show-transformed -cvzf dist/snap-linux-distribution.tar.gz $^

dist/snap-windows-distribution.zip: build/snap.exe README.md LICENSE
	mkdir -p dist
	zip -9 dist/snap-windows-distribution.zip $^

clean:
	$(RM) -r build
	$(RM) -r dist