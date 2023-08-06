package main

import (
	"encoding/binary"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"golang.org/x/net/context"
)

func showError(err error) {
	if err == nil {
		return
	}

	log.Errorf("[BLUETOOTH] Error: %s", err)
}

// StartBluetooth runs the bluetooth cycle forever, scanning for some time and processing results
func StartBluetooth() { //nolint:funlen
	d, err := dev.NewDevice(config.Device)
	if err != nil {
		log.Fatalf("[BLUETOOTH] Can't use device: %s", err)
	}

	ble.SetDefaultDevice(d)

	filter := func(a ble.Advertisement) bool {
		return strings.EqualFold(a.Addr().String(), config.DeviceID)
	}

	for {
		log.Infoln("[BLUETOOTH] Starting scan...")

		ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), config.ScanDuration.AsTimeDuration()))

		cln, err := ble.Connect(ctx, filter)
		if err != nil {
			log.Warnf("[BLUETOOTH] Scan timeout: %s", err)
			continue
		}

		// Normally, the connection is disconnected by us after our exploration.
		// However, it can be asynchronously disconnected by the remote peripheral.
		// So we wait(detect) the disconnection in the go routine.
		go func(cln ble.Client) {
			select {
			case <-cln.Disconnected():
				log.Infof("[BLUETOOTH] [ %s ] is disconnected ", cln.Addr())
			case <-time.After(config.Sub.AsTimeDuration()):
				log.Infof("[BLUETOOTH] [ %s ] timed out", cln.Addr())
			}
		}(cln)

		log.Infof("[BLUETOOTH] [ %s ] is connected ...", cln.Addr())
		log.Infoln("[BLUETOOTH] Discovering profile...")

		p, err := cln.DiscoverProfile(true)
		if err != nil {
			log.Errorf("[BLUETOOTH] can't discover profile: %s", err)
			showError(cln.CancelConnection())

			continue
		}

		log.Infof("[BLUETOOTH] Name: '%s'", cln.Name())

		// Start the exploration.
		explore(cln, p)

		log.Infof("[BLUETOOTH] Discovery done, waiting %f seconds before disconnecting.", config.Sub.AsTimeDuration().Seconds())
		time.Sleep(config.Sub.AsTimeDuration())

		// Disconnect the connection. (On OS X, this might take a while.)
		log.Infof("[BLUETOOTH] Disconnecting [ %s ]... (this might take up to few seconds on OS X)", cln.Addr())

		showError(cln.ClearSubscriptions())
		showError(cln.CancelConnection())

		time.Sleep(1 * time.Second)
	}
}

func explore(cln ble.Client, p *ble.Profile) {
	// First we subscribe to the metrics
	for _, s := range p.Services {
		for _, c := range s.Characteristics {
			switch c.UUID.String() {
			case "8a21", "8a22", "8a82":
				h := func(req []byte) { go decodeData(req) }

				if err := cln.Subscribe(c, true, h); err != nil {
					log.Errorf("[BLUETOOTH] subscribe failed: %s", err)
				}
			}
		}
	}

	// Then we send the current time (which triggers data transmission)
	for _, s := range p.Services {
		for _, c := range s.Characteristics {
			if c.UUID.String() == "8a81" {
				log.Infof("[BLUETOOTH] Sending the time... ")

				binarytime := generateTime()

				if err := cln.WriteCharacteristic(c, binarytime, true); err != nil {
					log.Errorf("[BLUETOOTH] Error while writing command: %+v", err)
				}

				log.Infof("[BLUETOOTH] done.")
			}
		}
	}
}

func generateTime() []byte {
	therealtime := time.Now().Unix()

	bs := make([]byte, 4)
	thetime := uint32(therealtime) - 1262304000
	binary.LittleEndian.PutUint32(bs, thetime)
	bs = append([]byte{2}, bs...)

	return bs
}
