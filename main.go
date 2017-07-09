package main

import (
	"encoding/binary"
	"fmt"
	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/examples/lib/dev"
	"golang.org/x/net/context"
	"strconv"
	"strings"
	"time"
)

func main() {
	for i := 1; i < 8; i++ {
		go func(i int) {
			new_weights := ImportCsv("csv/" + strconv.Itoa(i) + ".csv")
			updateData(i, new_weights)
		}(i)
	}

	d, err := dev.NewDevice(device)
	if err != nil {
		fmt.Printf("Can't use new device: %s", err)
	}
	ble.SetDefaultDevice(d)

	filter := func(a ble.Advertisement) bool {
		return strings.ToUpper(a.Address().String()) == strings.ToUpper(deviceID)
	}

	for {
		ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), scanDuration))
		cln, err := ble.Connect(ctx, filter)
		if err != nil {
			fmt.Printf("Timeout: %s\n", err)
		} else {
			// Make sure we had the chance to print out the message.
			done := make(chan struct{})
			// Normally, the connection is disconnected by us after our exploration.
			// However, it can be asynchronously disconnected by the remote peripheral.
			// So we wait(detect) the disconnection in the go routine.
			go func() {
				<-cln.Disconnected()
				fmt.Printf("[ %s ] is disconnected \n", cln.Address())
				close(done)
			}()

			fmt.Printf("Discovering profile...\n")
			p, err := cln.DiscoverProfile(true)
			if err != nil {
				fmt.Printf("can't discover profile: %s", err)
			}

			// Start the exploration.
			explore(cln, p)

			time.Sleep(sub)

			// Disconnect the connection. (On OS X, this might take a while.)
			fmt.Printf("Disconnecting [ %s ]... (this might take up to few seconds on OS X)\n", cln.Address())
			cln.CancelConnection()

			<-done
		}
	}
}

func updateData(person int, new_weights BodyMetrics) {
	csvFile := "csv/" + strconv.Itoa(person) + ".csv"
	cur_weights := ImportCsv(csvFile)

	cur_weights = mergeSort(cur_weights, new_weights)

	ExportCsv(csvFile+".new", cur_weights)
}

func generateTime(therealtime int64) []byte {
	bs := make([]byte, 4)
	thetime := uint32(therealtime) - 1262304000
	binary.LittleEndian.PutUint32(bs, thetime)
	bs = append([]byte{2}, bs...)

	return bs
}
