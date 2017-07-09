package main

import (
	"encoding/binary"
	"encoding/hex"
	"log"
)

const MaxInt = 2147483647
const time_offset = 1262304000

func test() {
	/*
		From CSV:
		1499291857,80.2,19.1,45.0,4.8,59.8,1788,25.3

		Results here:
		{Valid:true Person:2 Gender:male Age:33 Size:178 Activity:normal}
		{Valid:true BodyMetric:80.2 Timestamp:1499291857 Person:2}
		{Valid:true Timestamp:1499291857 Person:2 Kcal:1788 Fat:19.1 Tbw:59.8 Muscle:45 Bone:4.8}
	*/

	person_data, _ := hex.DecodeString("845302800121B2E0000000000000000000000000")
	weight_data, _ := hex.DecodeString("1d541f00fed125200e00000000020900000000")
	body_data, _ := hex.DecodeString("6fd125200e02fc06bff056f2c2f130f0000000")

	person := decodePerson(person_data)
	weight := decodeWeight(weight_data)
	body := decodeBody(body_data)

	log.Printf("%+v\n", person)
	log.Printf("%+v\n", weight)
	log.Printf("%+v\n", body)
}

func decodePerson(data []byte) (person Person) {
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

func decodeWeight(data []byte) (weight Weight) {
	/*
		fixed: byte: 0                     [0x1d]
		weight: byte: 1 & 2                [kg*100]
		timestamp: byte 5-8                Unix timestamp
		person: byte 13                    [1..8]
	*/
	weight.Valid = (data[0] == 0x1d)
	weight.BodyMetric = float32(decode16(data, 1)) / 100.0
	weight.Timestamp = sanitize_timestamp(decode32(data, 5))
	weight.Person = decode8(data, 13)
	return
}

func decodeBody(data []byte) (body Body) {
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
	body.Timestamp = sanitize_timestamp(decode32(data, 1))
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
	my_uint := data[firstByte]
	return int(my_uint)
}
func decode16(data []byte, firstByte int) int {
	my_uint := binary.LittleEndian.Uint16(data[firstByte:(firstByte + 2)])
	return int(my_uint)
}
func decode32(data []byte, firstByte int) int {
	my_uint := binary.LittleEndian.Uint32(data[firstByte:(firstByte + 4)])
	return int(my_uint)
}

func sanitize_timestamp(timestamp int) int {
	retTS := 0
	if timestamp+time_offset < MaxInt {
		retTS = timestamp + time_offset
	} else {
		retTS = timestamp
	}

	if timestamp >= MaxInt {
		retTS = 0
	}

	return retTS
}
