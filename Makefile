GOMOBILE=gomobile
GOBIND=$(GOMOBILE) bind
BUILDDIR=$(shell pwd)/build
ARTIFACT=$(BUILDDIR)/Tun2socks.framework
LDFLAGS='-s -w'
IMPORT_PATH=github.com/eycorsican/go-tun2socks-ios
TUN2SOCKS_PATH=$(GOPATH)/src/github.com/eycorsican/go-tun2socks

BUILD_CMD="cd $(BUILDDIR) && $(GOBIND) -a -ldflags $(LDFLAGS) -target=ios/arm -o $(ARTIFACT) $(IMPORT_PATH)"

all: $(ARTIFACT)

$(ARTIFACT):
	mkdir -p $(BUILDDIR)
	cd $(TUN2SOCKS_PATH) && make copy
	eval $(BUILD_CMD)
	cd $(TUN2SOCKS_PATH) && make clean

clean:
	rm -rf $(BUILDDIR)
	cd $(TUN2SOCKS_PATH) && make clean
