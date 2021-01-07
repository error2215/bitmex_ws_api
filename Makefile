format:
	go get golang.org/x/tools/cmd/goimports
	find . -name \*.go -not -path ./vendor -exec goimports -w {} \;

build_add:
	cd cmd/api && go build -o ../bin/api

build_all: create_build_dir build_add replace_env_file

create_build_dir:
	mkdir -p ../bin/

replace_env_file:
	cp ./.env ./cmd/bin/.env
test:
	go test ./...
