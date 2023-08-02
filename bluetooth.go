package main

import (
	"encoding/binary"
	"log"
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
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
		log.Printf("[BLUETOOTH] Can't use device: %s", err)
		return
	}

	var device *bluetooth.Device

	log.Println("[BLUETOOTH] Starting scan...")

	ch := make(chan bluetooth.ScanResult, 1)

	err := d.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if strings.EqualFold(result.Address.String(), config.DeviceID) {
			log.Printf("[BLUETOOTH] Found device: %s (%d, %s)", result.Address.String(), result.RSSI, result.LocalName())

			if err := adapter.StopScan(); err != nil {
				log.Printf("[BLUETOOTH] Can't use device: %s", err)
				return
			}

			ch <- result
		}
	})
	if err != nil {
		log.Printf("[BLUETOOTH] Can't use device: %s", err)
		return
	}

	select {
	case result := <-ch:
		device, err = d.Connect(result.Address, bluetooth.ConnectionParams{})
		if err != nil {
			log.Printf("[BLUETOOTH] [ %s ] connection failed: %s", result.Address.String(), err)
			return
		}

		log.Printf("[BLUETOOTH] [ %s ] is connected ...", result.Address.String())

		defer device.Disconnect()
	}

	log.Println("[BLUETOOTH] Discovering profile...")

	// Start the exploration.
	explore(device)

	log.Printf("[BLUETOOTH] Discovery done, waiting %d seconds before disconnecting.\n", (config.Sub.AsTimeDuration() / 1e9))
	time.Sleep(config.Sub.AsTimeDuration())
}

func explore(p *bluetooth.Device) {
	// First we subscribe to the metrics
	services, err := p.DiscoverServices(nil)
	if err != nil {
		log.Printf("[BLUETOOTH] Error exploring: %s", err)
		return
	}

	for _, s := range services {
		chars, err := s.DiscoverCharacteristics(nil)
		if err != nil {
			log.Printf("[BLUETOOTH] Error exploring: %s", err)
			continue
		}

		for _, c := range chars {
			if !c.UUID().Is16Bit() {
				continue
			}

			s := strings.Trim(strings.Split(c.UUID().String(), "-")[0], "0")
			log.Printf("[BLUETOOTH] Discovering service: %s", s)

			switch s {
			case "8a21", "8a22", "8a82":
				err := c.EnableNotifications(func(buf []byte) {
					log.Printf("data: %#v\n", buf)
					go decodeData(buf)
				})
				if err != nil {
					log.Printf("[BLUETOOTH] Error exploring: %s", err)
					continue
				}
			case "8a81":
				log.Printf("[BLUETOOTH] Sending the time... ")

				thetime := time.Now().Unix()

				binarytime := generateTime(thetime)
				if _, err := c.WriteWithoutResponse(binarytime); err != nil {
					log.Printf("[BLUETOOTH] Error while writing command: %+v\n", err)
					continue
				}
			}

			log.Printf("[BLUETOOTH] done.\n")
		}
	}
}

func generateTime(therealtime int64) []byte {
	bs := make([]byte, 4)
	thetime := uint32(therealtime) - 1262304000
	binary.LittleEndian.PutUint32(bs, thetime)
	bs = append([]byte{2}, bs...)

	return bs
}
