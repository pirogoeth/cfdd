.PHONY: build get-deps
.DEFAULT_GOAL := install

GITHUB_ORG=pirogoeth
GITHUB_REPO=cfdd

UPSTREAM=$(GITHUB_ORG)/$(GITHUB_REPO)

SHA=$(shell git rev-parse HEAD | cut -b1-9)

LDFLAGS="-X main.Version=$(VERSION) -X main.BuildHash=$(SHA)"

clean:
	rm -rf release/

get-deps:
	go get -u github.com/itchio/gothub
	dep ensure -v

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
	@if ! which gothub 2>&1 >> /dev/null; then \
		echo " # gothub not found in path; install and create a github token with 'repo' access"; \
		echo " # See (https://help.github.com/articles/creating-an-access-token-for-command-line-use)"; \
		echo " make get-deps -OR- go get github.com/itchio/gothub"; \
		echo " export GITHUB_TOKEN=<your-token>";\
		echo ""; \
		exit 1; \
	fi
	@gothub release \
		--user $(GITHUB_ORG) \
		--repo $(GITHUB_REPO) \
		--tag $(VERSION)
	@gothub upload \
		--user $(GITHUB_ORG) \
		--repo $(GITHUB_REPO) \
		--tag $(VERSION) \
		--name "cfdd-linux-amd64" \
		--file release/cfdd-linux-amd64
	@gothub upload \
		--user $(GITHUB_ORG) \
		--repo $(GITHUB_REPO) \
		--tag $(VERSION) \
		--name "cfdd-darwin-amd64" \
		--file release/cfdd-darwin-amd64
	cd release && shasum -a 256 ./* | tee SHA256SUMS
	@gothub upload \
		--user $(GITHUB_ORG) \
		--repo $(GITHUB_REPO) \
		--tag $(VERSION) \
		--name "SHA256SUMS" \
		--file release/SHA256SUMS
	gothub info \
		--user $(GITHUB_ORG) \
		--repo $(GITHUB_REPO) \
		--tag $(VERSION)
