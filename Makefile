bin=medisana-bs


.PHONY: build test clean all

all: arm6 arm7 linux32 linux64

pi-service:
	rsync -vaiz dist/medisana-bs.service root@scale-pi:/etc/systemd/system/medisana-bs.service
	ssh root@scale-pi systemctl daemon-reload

pi: pi-service arm6
	rsync -vaiz build/medisana-bs.arm6 root@scale-pi:/opt/medisana-bs/
	ssh root@scale-pi systemctl restart medisana-bs
	sleep 5
	ssh root@scale-pi systemctl status medisana-bs

build:
	$(BUILDOPTS) go build -mod vendor -o build/$(bin).$(EXT)

arm6:
	@$(MAKE) build BUILDOPTS="GOOS=linux GOARCH=arm GOARM=6" EXT=$(@)

arm7:
	@$(MAKE) build BUILDOPTS="GOOS=linux GOARCH=arm GOARM=7" EXT=$(@)

linux64:
	@$(MAKE) build BUILDOPTS="GOOS=linux GOARCH=amd64" EXT=$(@)

linux32:
	@$(MAKE) build BUILDOPTS="GOOS=linux GOARCH=386" EXT=$(@)

clean:
	rm -vf build/*
