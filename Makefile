bin=medisana-bs

pi: arm6
	rsync -vaiz build/medisana-bs.arm6 root@scale-pi:
	ssh -t root@scale-pi /root/medisana-bs.arm6

build:
	$(BUILDOPTS) go build -o build/$(bin).$(EXT)

arm6:
	@$(MAKE) build BUILDOPTS="GOOS=linux GOARCH=arm GOARM=6" EXT=$(@)

arm7:
	@$(MAKE) build BUILDOPTS="GOOS=linux GOARCH=arm GOARM=7" EXT=$(@)

linux64:
	@$(MAKE) build BUILDOPTS="GOOS=linux GOARCH=amd64" EXT=$(@)

linux32:
	@$(MAKE) build BUILDOPTS="GOOS=linux GOARCH=386" EXT=$(@)

all: arm6 arm7 linux32 linux64

.PHONY: build test
