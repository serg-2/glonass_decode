package skytraqlib

import (
	"fmt"
)

func DecodeSkyTraq (message []byte, gloState int) (bool,int) {
	var length int16
	length |= int16(message[3])
	length |= int16(message[2]) << 8

	//fmt.Printf("Message length: %v\n", length)
	// Do we need to check length ?
	// 4 - header. 3 - crc
	// length == len(message) - 4 - 3

	switch messageID:= message[4]; messageID{
	case byte(9):
		fmt.Println("Configure message type received")
	case byte(14):
		fmt.Println("Configure Position update rate received")
	case byte(16):
		fmt.Println("Query Position update rate received")
	case byte(30):
		fmt.Println("Configure Binary Measurement Data Output received")
	case byte(31):
		fmt.Println("Query Binary Measurement Data Output Status")
	case byte(32):
		fmt.Println("Configure Binary RTCM Data Output")
	case byte(33):
		fmt.Println("Query Binary RTCM Data Output Status")
	case byte(34):
		fmt.Println("Configure Base Position")
	case byte(35):
		fmt.Println("Query Base Position")
	case byte(48):
		fmt.Println("Get GPS Ephemeris")
	case byte(65):
		fmt.Println("Set GPS Ephemeris ")
	case byte(91):
		fmt.Println("Get GLONASS ephemeris")
	case byte(92):
		fmt.Println("Set GLONASS ephemeris")
	case byte(106):
		// NEED SUB_ID 4
		fmt.Println("Reset and Re-calculate GLONASS Inter-Frequency Bias (IFB)")
	case byte(128):
		fmt.Println("Software version")
	case byte(129):
		fmt.Println("Software CRC")
	case byte(130):
		fmt.Println("RESERVED")
	case byte(131):
		fmt.Println("Acknoledgment received")
	case byte(132):
		fmt.Println("REJECT received")
	case byte(134):
		fmt.Println("Position update rate")
	case byte(137):
		fmt.Println("Binary Measurement Data Output Status")
	case byte(138):
		fmt.Println("Binary RTCM Data Output Status")
	case byte(139):
		fmt.Println("Base Position")
	case byte(144):
		fmt.Println("GLONASS ephemeris")
	case byte(177):
		fmt.Println("GPS Ephemeris Data")
	case byte(220):
		fmt.Println("Measurement Epoch")
	case byte(221):
		fmt.Println("RAW Raw Measurement")
	case byte(222):
		fmt.Println("SV and channel status")
	case byte(223):
		//fmt.Println("(223) Receiver navigation status")
		// can be decoded by gpsdecode
		//Decode_223(payload[1:])
		//return true
	case byte(224):
		// can be decoded by gpsdecode
		//fmt.Println("(224) GPS Subframe buffer data")
		//Decode_224(payload[1:])
		//return true
	case byte(225):
		fmt.Println("(225) Glonass String buffer data")
		// cannot be decoded by gpsdecode
		gloState = Decode_225(message[5:len(message)-3], gloState)
	case byte(226):
		fmt.Println("Beidou2 D1 Subframe Data")
	case byte(227):
		fmt.Println("Beidou2 D2 Subframe Data")
	case byte(229):
		// INTERESTING
		//fmt.Println("Extended Raw Measurement Data v.1")
		// cannot be decoded by gpsdecode
	default:
		fmt.Printf("Message ID: %v\n", messageID)
		return true, gloState
	}


	return false, gloState

}
