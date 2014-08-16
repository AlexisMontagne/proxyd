PREFIX=/usr/local
DESTDIR=
GOFLAGS=
BINDIR=${PREFIX}/bin

BINARIES = proxyd_endpoint proxyd_balancer
BLDDIR = build

all: $(BINARIES)

$(BLDDIR)/%:
	@mkdir -p $(dir $@)
	go build ${GOFLAGS} -o $(abspath $@) ./$*

$(BINARIES): %: $(BLDDIR)/%

clean:
	rm -fr $(BLDDIR)

.PHONY: clean all

