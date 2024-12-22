package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"golang.org/x/net/context"
)

const (
	CharPersonShort  = "8a82"                                 // person data
	CharWeightShort  = "8a21"                                 // weight data
	CharBodyShort    = "8a22"                                 // body data
	CharCommandShort = "8a81"                                 // command register
	CharPerson       = "00008a82-0000-1000-8000-00805f9b34fb" // person data
	CharWeight       = "00008a21-0000-1000-8000-00805f9b34fb" // weight data
	CharBody         = "00008a22-0000-1000-8000-00805f9b34fb" // body data
	CharCommand      = "00008a81-0000-1000-8000-00805f9b34fb" // command register
)

func showError(err error) {
	if err == nil {
		return
	}

	bluetoothLogger().Errorf("error: %s", err)
}

func bluetoothLogger() log.FieldLogger {
	return log.WithField("component", "bluetooth")
}

// StartBluetooth runs the bluetooth cycle forever, scanning for some time and processing results
func StartBluetooth() {
	d, err := dev.NewDevice(config.Device)
	if err != nil {
		bluetoothLogger().Fatalf("can't use device: %s", err)
	}

	ble.SetDefaultDevice(d)

	filter := func(a ble.Advertisement) bool {
		return strings.EqualFold(a.Addr().String(), config.DeviceID)
	}

	for {
		if err := scan(filter); err != nil {
			if !errors.Is(err, context.DeadlineExceeded) {
				bluetoothLogger().Fatalf("%s", err)
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func scan(filter ble.AdvFilter) error {
	bluetoothLogger().Infoln("starting scan...")

	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), config.ScanDuration.AsTimeDuration()))

	// Calculate the new timestamp by adding n seconds
	newTimestamp := time.Now().Add(config.ScanDuration.AsTimeDuration())

	bluetoothLogger().Infof("connecting... I should be back by %s", newTimestamp)

	ch := make(chan bool, 1)
	go checkBluetoothHangs(ch)

	cln, err := ble.Connect(ctx, filter)
	if err != nil {
		ch <- false
		return fmt.Errorf("scan timeout: %w", err)
	}

	defer disconnect(cln)

	ch <- true

	// Normally, the connection is disconnected by us after our exploration.
	// However, it can be asynchronously disconnected by the remote peripheral.
	// So we wait(detect) the disconnection in the go routine.
	go func(cln ble.Client) {
		select {
		case <-cln.Disconnected():
			bluetoothLogger().Infof("[ %s ] is disconnected ", cln.Addr())
		case <-time.After(config.Sub.AsTimeDuration()):
			bluetoothLogger().Infof("[ %s ] timed out", cln.Addr())
		}
	}(cln)

	bluetoothLogger().Infof("[ %s ] is connected ...", cln.Addr())
	bluetoothLogger().Infoln("discovering profile...")

	p, err := cln.DiscoverProfile(true)
	if err != nil {
		return fmt.Errorf("can't discover profile: %w", err)
	}

	bluetoothLogger().Infof("address: '%s'; name: '%s'", cln.Addr(), cln.Name())

	// Start the exploration.
	showError(explore(cln, p))

	bluetoothLogger().Infof("discovery done, waiting %.0f seconds before disconnecting.", config.Sub.AsTimeDuration().Seconds())
	time.Sleep(config.Sub.AsTimeDuration())

	go debounce()

	return nil
}

func checkBluetoothHangs(ch chan bool) {
	// Add a few extra seconds head start to prevent race conditions
	time.Sleep(10 * time.Second)

	select {
	case <-ch:
		bluetoothLogger().Infof("scan finished in time")
		return
	case <-time.After(config.ScanDuration.AsTimeDuration()):
		bluetoothLogger().Fatal("scan did not finish in time - hardware hang?")
	}
}

func disconnect(cln ble.Client) {
	// Disconnect the connection. (On OS X, this might take a while.)
	bluetoothLogger().Infof("disconnecting [ %s ]...", cln.Addr())

	showError(cln.ClearSubscriptions())
	showError(cln.CancelConnection())

	bluetoothLogger().Infof("disconnected!")
}

func explore(cln ble.Client, p *ble.Profile) error {
	// First we subscribe to the metrics
	for _, s := range p.Services {
		for _, c := range s.Characteristics {
			switch c.UUID.String() {
			case CharWeightShort, CharBodyShort, CharPersonShort:
				h := func(req []byte) { go decodeData(req) }

				if err := cln.Subscribe(c, true, h); err != nil {
					return fmt.Errorf("subscribe failed: %w", err)
				}
			}
		}
	}

	// Then we send the current time (which triggers data transmission)
	for _, s := range p.Services {
		for _, c := range s.Characteristics {
			if c.UUID.String() == CharCommandShort {
				bluetoothLogger().Infof("sending the time... ")

				binarytime := generateTime()

				if err := cln.WriteCharacteristic(c, binarytime, true); err != nil {
					return fmt.Errorf("error while writing command: %w", err)
				}

				bluetoothLogger().Infof("done.")
			}
		}
	}

	return nil
}

func generateTime() []byte {
	therealtime := time.Now().Unix()

	bs := make([]byte, 4)
	thetime := uint32(therealtime) - 1262304000
	binary.LittleEndian.PutUint32(bs, thetime)
	bs = append([]byte{2}, bs...)

	return bs
}
