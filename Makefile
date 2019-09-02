DATA_FILES := $(shell find conf | sed 's/ /\\ /g')
LESS_FILES := $(wildcard public/less/anony.less public/less/_*.less)
GENERATED  := bindata less uglifycss uglifyjs

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

uglifycss: public/css/anony.min.css

public/css/anony.min.css: public/css/anony.css
	@type uglifycss >/dev/null 2>&1 && uglifycss $< >$@ || echo "uglifycss command not found, skipped."

uglifyjs: public/js/anony.min.js public/serviceworker.min.js

public/js/anony.min.js: public/js/anony.js
	@type terser >/dev/null 2>&1 && terser --compress --mangle -- $< >$@ || echo "terser  command not found, skipped."

public/serviceworker.min.js: public/serviceworker.js
	@type terser >/dev/null 2>&1 && terser --compress --mangle -- $< >$@ || echo "terser  command not found, skipped."

clean:
	go clean -i ./...


