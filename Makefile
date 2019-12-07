#.DEFAULT_GOAL := all
name := "goms-bin"
build_destination := "."
goms_file_name := $(build_destination)/$(name)
docker_image := "thobe/goms:latest"

build_time := $(shell date '+%Y-%m-%d_%H-%M-%S')
rev  := $(shell git rev-parse --short HEAD)
flag := $(shell git diff-index --quiet HEAD -- || echo "_dirty";)
tag := $(shell git describe --tags 2> /dev/null)
branch := $(shell git branch | grep \* | cut -d ' ' -f2)
revision := $(rev)$(flag)
build_info := $(build_time)_$(revision)

all: tools test build finish

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"


test: sep ## Runs all unittests and generates a coverage report.
	@echo "--> Run the unit-tests"
	@go test ./ ./logging ./api -covermode=count -coverprofile=coverage.out

build: sep ## Builds the goms binary.
	@echo "--> Build the $(name) in $(build_destination)"
	@go build -v -ldflags "-X main.version=$(tag) -X main.buildTime=$(build_time) -X main.revision=$(revision) -X main.branch=$(branch)" -o $(goms_file_name) .


run: sep build ## Builds + runs goms
	@echo "--> Run $(goms_file_name)"
	$(goms_file_name)

docker.build: sep ## Builds the goms docker image.
	@echo "--> Build docker image $(docker_image)"
	@docker build -t thobe/goms -f ci/Dockerfile .

docker.run: sep ## Runs the goms docker image.
	@echo "--> Run docker image $(docker_image)"
	@docker run --rm --name=goms -p 11000:11000 $(docker_image)

docker.push: sep ## Pushes the goms docker image to docker-hub
	@echo "--> Tag image to thobe/goms:$(tag)"
	@docker tag thobe/goms:latest thobe/goms:$(tag)
	@echo "--> Push image thobe/goms:latest"
	@docker push thobe/goms:latest 
	@echo "--> Push image thobe/goms:$(tag)"
	@docker push thobe/goms:$(tag)

infra.up: ## Starts up the infra components
	make -C infra up

infra.down: ## Stops up the infra components
	make -C infra down

ui.open:
	@xdg-open http://localhost:15672 # rabbit-ui
	@xdg-open http://localhost:3000 # grafana
	@xdg-open http://localhost:9090 # prometheus
	@xdg-open http://localhost:8180 # keycloak

sep:
	@echo "----------------------------------------------------------------------------------"

finish:
	@echo "=================================================================================="