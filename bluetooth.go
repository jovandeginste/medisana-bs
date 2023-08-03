package main

import (
	"encoding/binary"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"tinygo.org/x/bluetooth"
)

const (
	CharPerson  = "00008a82-0000-1000-8000-00805f9b34fb" // person data
	CharWeight  = "00008a21-0000-1000-8000-00805f9b34fb" // weight data
	CharBody    = "00008a22-0000-1000-8000-00805f9b34fb" // body data
	CharCommand = "00008a81-0000-1000-8000-00805f9b34fb" // command register
)

// StartBluetooth runs the bluetooth cycle forever, scanning for some time and processing results
func StartBluetooth() {
	for {
		loop()

		time.Sleep(5 * time.Second)
	}
}

func loop() {
	d := bluetooth.DefaultAdapter

	if err := d.Enable(); err != nil {
		log.Errorf("[BLUETOOTH] Can't use device: %s", err)
		return
	}

	var device *bluetooth.Device

	log.Infoln("[BLUETOOTH] Starting scan...")

	var r bluetooth.ScanResult

	err := d.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if !strings.EqualFold(result.Address.String(), config.DeviceID) {
			return
		}

		log.Infof("[BLUETOOTH] Found device: %s (%d, %s)", result.Address.String(), result.RSSI, result.LocalName())

		if err := adapter.StopScan(); err != nil {
			log.Errorf("[BLUETOOTH] Can't cancel the scan: %s", err)
			return
		}

		r = result
	})
	if err != nil {
		log.Errorf("[BLUETOOTH] Can't use device: %s", err)
		return
	}

	log.Infof("[BLUETOOTH] Connecting to: %s", r.Address.String())

	device, err = d.Connect(r.Address, bluetooth.ConnectionParams{})
	if err != nil {
		log.Errorf("[BLUETOOTH] [ %s ] connection failed: %s", r.Address.String(), err)
		return
	}

	defer func() {
		if err := device.Disconnect(); err != nil {
			log.Errorf("Error disconnecting: %s", err)
		}
	}()

	log.Infof("[BLUETOOTH] [ %s ] is connected ...", r.Address.String())

	explore(device)

	log.Infoln("[BLUETOOTH] Discovery, disconnecting...")
}

func explore(p *bluetooth.Device) {
	log.Infoln("[BLUETOOTH] Discovering profile...")

	// First we subscribe to the metrics
	services, err := p.DiscoverServices(nil)
	if err != nil {
		log.Errorf("[BLUETOOTH] Error exploring: %s", err)
		return
	}

	for _, s := range services {
		chars, err := s.DiscoverCharacteristics(nil)
		if err != nil {
			log.Errorf("[BLUETOOTH] Error exploring: %s", err)
			continue
		}

		for _, c := range chars {
			s := c.UUID().String()
			log.Tracef("[BLUETOOTH] Discovering service: %s", s)

			switch s {
			case CharPerson, CharWeight, CharBody:
				log.Tracef("[BLUETOOTH] Receiving data... ")

				err := c.EnableNotifications(func(buf []byte) {
					go decodeData(buf)
				})
				if err != nil {
					log.Errorf("[BLUETOOTH] Error exploring: %s", err)
					continue
				}
			case CharCommand:
				log.Infof("[BLUETOOTH] Sending the time... ")

				binarytime := generateTime()
				if _, err := c.WriteWithoutResponse(binarytime); err != nil {
					log.Errorf("[BLUETOOTH] Error while writing command: %+v\n", err)
					continue
				}
			}

			log.Tracef("[BLUETOOTH] done.")
		}
	}

	time.Sleep(config.Sub.AsTimeDuration())
}

func generateTime() []byte {
	therealtime := time.Now().Unix()
	bs := make([]byte, 4)
	thetime := uint32(therealtime) - 1262304000
	binary.LittleEndian.PutUint32(bs, thetime)
	bs = append([]byte{2}, bs...)

	return bs
}
