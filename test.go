package main

import (
	"./skytraq"
	"fmt"
	"github.com/serg-2/libs-go/seriallib"
)


func main() {

	// Workaround to save glonass state
	var gloState = 6

	var needToDecodeByGPSDecode bool

	reader := seriallib.GetPortReader("/dev/ttyUSB0", 19200)

	//Get first not Full message
	seriallib.GetSkyTraqMessage(reader, false)

	for {
		//Get Good message
		//reply  := getMessage(reader, true)
		reply, crcOk := seriallib.GetSkyTraqMessage(reader, false)

		//fmt.Printf("%v\n", reply)
		// WRITE TO FILE
		// skytraqlib.WriteBytesToFile("/tmp/rtcm.log", reply)

		if crcOk {
			needToDecodeByGPSDecode, gloState = skytraqlib.DecodeSkyTraq(reply, gloState)

			if needToDecodeByGPSDecode {
				fmt.Printf("%v\n", string(skytraqlib.DecodeByGPSDecode(reply)))
			}

		}
	}
}

