package main

import (
	"encoding/binary"
	"log"
	"math"

	"github.com/jovandeginste/medisana-bs/structs"
)

func decodePerson(data []byte) (person structs.Person) {
	/*
		fixed: byte 0                       [0x84]
		person: byte 2                      [1..8]
		gender: byte 4 (1=male, 2=female)   [1|2]
		age: byte 5                         [year]
		size: byte 6                        [cm]
		activity: byte 8 (0=normal, 3=high) [0|3]
	*/
	person.Valid = (data[0] == 0x84)
	person.Person = decode8(data, 2)
	if data[4] == 1 {
		person.Gender = "male"
	} else {
		person.Gender = "female"
	}
	person.Age = decode8(data, 5)
	person.Size = decode8(data, 6)
	if data[8] == 3 {
		person.Activity = "high"
	} else {
		person.Activity = "normal"
	}
	return
}

func decodeWeight(data []byte) (weight structs.Weight) {
	/*
		fixed: byte: 0                     [0x1d]
		weight: byte: 1 & 2                [kg*100]
		timestamp: byte 5-8                Unix timestamp
		person: byte 13                    [1..8]
	*/
	weight.Valid = (data[0] == 0x1d)
	weight.Weight = float32(decode16(data, 1)) / 100.0
	weight.Timestamp = sanitizeTimestamp(decode32(data, 5))
	weight.Person = decode8(data, 13)
	return
}

func decodeBody(data []byte) (body structs.Body) {
	/*
		fixed: byte 0                      [0x6f]
		timestamp: byte 1-4                Unix timestamp
		person: byte 5                     [1..8]
		kcal: byte 6 & 7                   first nibble = 0xf, [kcal]
		fat: byte 8 & 9                    first nibble = 0xf, [fat*10]
		tbw: byte 10 & 11                  first nibble = 0xf, [tbw*10]
		muscle: byte 12 & 13               first nibble = 0xf, [muscle*10]
		bone: byte 14 & 15                 first nibble = 0xf, [bone*10]
	*/
	body.Valid = (data[0] == 0x6f)
	body.Timestamp = sanitizeTimestamp(decode32(data, 1))
	body.Person = decode8(data, 5)
	body.Kcal = decode16(data, 6)
	body.Fat = smallValue(decode16(data, 8))
	body.Tbw = smallValue(decode16(data, 10))
	body.Muscle = smallValue(decode16(data, 12))
	body.Bone = smallValue(decode16(data, 14))
	return
}

func smallValue(value int) float32 {
	return float32(0x0fff&value) / 10.0
}
func decode8(data []byte, firstByte int) int {
	myUint := data[firstByte]
	return int(myUint)
}
func decode16(data []byte, firstByte int) int {
	myUint := binary.LittleEndian.Uint16(data[firstByte:(firstByte + 2)])
	return int(myUint)
}
func decode32(data []byte, firstByte int) int {
	myUint := binary.LittleEndian.Uint32(data[firstByte:(firstByte + 4)])
	return int(myUint)
}

func sanitizeTimestamp(timestamp int) int {
	retTS := 0
	if timestamp+config.TimeOffset < math.MaxInt32 {
		retTS = timestamp + config.TimeOffset
	} else {
		retTS = timestamp
	}

	if timestamp >= math.MaxInt32 {
		retTS = 0
	}

	return retTS
}

func decodeData(req []byte) {
	result := new(structs.PartialMetric)
	switch req[0] {
	case 0x84:
		person := decodePerson(req)
		result.Person = person
	case 0x1D:
		weight := decodeWeight(req)
		result.Weight = weight
	case 0x6F:
		body := decodeBody(req)
		result.Body = body
	default:
		log.Printf("[DECODE] Unhandled data encountered: [% X]\n", req)
	}
	metricChan <- result
}
