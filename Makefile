GOARCH?=amd64
GOOS?=linux
TARBALL?=mister.tar.gz

VERSION=$(shell cat VERSION)

build-tarball:
	$(eval TMP := $(shell mktemp -d /tmp/mister.XXXXXX))

	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(TMP)/mister/misterd cmd/misterd/*.go 
	chmod +x $(TMP)/mister/misterd
	tar -C $(TMP) -zcvf $(TARBALL) mister
	rm -rf $(TMP)


test-unit: 
	echo "TODO"

test-integration:
	echo "TODO"