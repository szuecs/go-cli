default: test.test

all: godep.restore clean build.linux build.osx build.win

clean:
	rm -rf build
	rm -f test/bench.*
	rm -f test/prof.*
	find . -name '*.test' -delete

config:
	@test -d ~/.config/go-cli || mkdir -p ~/.config/go-cli
	@test -e ~/.config/go-cli/config.yaml || cp config.yaml.sample ~/.config/go-cli/config.yaml
	@echo "modify ~/.config/go-cli/config.yaml as you need"

test.all: test.benchmark.new test.test test.vet test.errcheck

test.test:
	GIN_MODE=release go test ./...

test.errcheck:
	errcheck ./...

test.vet:
	go vet -v ./...

test.benchmark.new:
	go list ./... |  xargs go test -run=^$$ -bench=. | tee test/bench.new

test.benchmark.old:
	go list ./... |  xargs go test -run=^$$ -bench=. | tee test/bench.old

test.benchmark.cmp:
	benchcmp test/bench.old test/bench.new

## profiling parameters can be overridden: % make test.profile.cpu N=5 P=/health
# 30s is the default value to profile your web application
N?=30
# / is the default for request path to the web application you want to profile
P?=/
# requires: % go get github.com/DeanThompson/ginpprof
test.profile.cpu: build.local
	scripts/pprof $(N) $(P)

# requires: % go get github.com/laher/gols/...
check.dependencies:
	go-ls -ignore=/vendor/ -exec="depscheck -v" ./...

# requires: % go get github.com/tools/godep
godep.clean:
	rm -rf Godeps

godep.restore:
	git checkout Godeps
	godep restore

godep.recreate: godep.clean
	godep save

prepare:
	@mkdir -p build/linux
	@mkdir -p build/osx
	@mkdir -p build/windows
	@echo created ./build/ directories

# release
build.linux.release: prepare
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -o build/linux/go-cli -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD` -X main.Version=`git describe --tags`" -tags zalandoValidation ./cmd/go-cli

# dev builds
build.service.local: prepare
	godep go build -o build/go-cli -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD` -X main.Version=HEAD" -tags zalandoValidation  ./cmd/go-cli

# OS specific builds
build.linux: prepare
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -o build/linux/go-cli -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD` -X main.Version=HEAD" -tags zalandoValidation ./cmd/go-cli

build.osx: prepare
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 godep go build -o build/osx/go-cli -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD` -X main.Version=HEAD" -tags zalandoValidation ./cmd/go-cli

build.win: prepare
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 godep go build -o build/windows/go-cli -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD` -X main.Version=HEAD" -tags zalandoValidation ./cmd/go-cli

# build and install multi binary project
dev.install:
	godep go install -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD` -X main.Version=HEAD" -tags zalandoValidation ./...
