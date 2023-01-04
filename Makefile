
# Build locally, will make ./tfm available to run
# Adds in some build flags to identify the binary in the future
build-local:
	go build -v \
	-ldflags="\
	-X 'github.com/hashicorp-services/tfm/version.Version=x.x.x' \
	-X 'github.com/hashicorp-services/tfm/version.Prerelease=alpha' \
	-X 'github.com/hashicorp-services/tfm/version.Build=local' \
	-X 'github.com/hashicorp-services/tfm/version.BuiltBy=$(shell whoami)' \
	-X 'github.com/hashicorp-services/tfm/version.Date=$(shell date)'"

# Updated go packages (will touch go.mod and go.sum)
update:
	go get -update
	go mod tidy

format:
	go fmt