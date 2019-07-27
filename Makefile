DATA_FILES := $(shell find conf | sed 's/ /\\ /g')
LESS_FILES := $(wildcard public/less/anony.less public/less/_*.less)
GENERATED  := pkg/bindata/bindata.go public/css/anony.css

OS := $(shell uname)

TAGS = ""
BUILD_FLAGS = "-v"

.PHONY: build bindata less clean
.IGNORE: public/css/anony.css

all: build

build: $(GENERATED)
	go install $(BUILD_FLAGS) -tags '$(TAGS)'
	cp '$(GOPATH)/bin/anonymoe' .

bindata: pkg/bindata/bindata.go

pkg/bindata/bindata.go: $(DATA_FILES)
	go-bindata -o=$@ -pkg=bindata conf/...

less: public/css/anony.css

public/css/anony.css: $(LESS_FILES)
	@type lessc >/dev/null 2>&1 && lessc $< >$@ || echo "lessc command not found, skipped."

clean:
	go clean -i ./...


