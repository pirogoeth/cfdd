.PHONY: build get-deps
.DEFAULT_GOAL := install

GITHUB_ORG=pirogoeth
GITHUB_REPO=cfdd

UPSTREAM=$(GITHUB_ORG)/$(GITHUB_REPO)

VERSION=0.1.2
SHA=$(shell git rev-parse HEAD | cut -b1-9)

LDFLAGS="-X main.Version=$(VERSION) -X main.BuildHash=$(SHA)"

clean:
	rm -rf release/

get-deps:
	go get -u github.com/Masterminds/glide
	glide install

build:
	go build -v -ldflags $(LDFLAGS) ./cmd/cfdd

install:
	go install -v -ldflags $(LDFLAGS) ./cmd/cfdd

build-release: build-release-cfdd

build-release-cfdd: release/cfdd-linux-amd64 release/cfdd-darwin-amd64

release/cfdd-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags $(LDFLAGS) \
		 -o release/cfdd-darwin-amd64 cmd/cfdd/main.go

release/cfdd-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) \
		 -o release/cfdd-linux-amd64 cmd/cfdd/main.go

release: clean build-release
	@if [ "$(VERSION)" = "" ]; then \
		echo " # 'VERSION' variable not set! To preform a release do the following"; \
		echo "  git tag v1.0.0"; \
		echo "  git push --tags"; \
		echo "  make release VERSION=v1.0.0"; \
		echo ""; \
		exit 1; \
	fi
	@if ! which github-release 2>&1 >> /dev/null; then \
		echo " # github-release not found in path; install and create a github token with 'repo' access"; \
		echo " # See (https://help.github.com/articles/creating-an-access-token-for-command-line-use)"; \
		echo " go get github.com/aktau/github-release"; \
		echo " export GITHUB_TOKEN=<your-token>";\
		echo ""; \
		exit 1; \
	fi
	@github-release release \
		--user $(GITHUB_ORG) \
		--repo $(GITHUB_REPO) \
		--tag $(VERSION)
	@github-release upload \
		--user $(GITHUB_ORG) \
		--repo $(GITHUB_REPO) \
		--tag $(VERSION) \
		--name "cfdd-linux-amd64" \
		--file release/cfdd-linux-amd64
	@github-release upload \
		--user $(GITHUB_ORG) \
		--repo $(GITHUB_REPO) \
		--tag $(VERSION) \
		--name "cfdd-darwin-amd64" \
		--file release/cfdd-darwin-amd64
