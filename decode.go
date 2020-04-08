package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"os/exec"

	"github.com/tarm/serial"
)


func getMessage(reader *bufio.Reader, debug bool) []byte {

	var partMessage byte

	var checkMessage, fullMessage []byte

	begin1, _ := hex.DecodeString("0D")
	begin2, _ := hex.DecodeString("0A")

	begin3, _ := hex.DecodeString("A0")
	begin4, _ := hex.DecodeString("A1")

	frameString,_ := hex.DecodeString("0D0AA0A1")

	// wait for start


		for len(checkMessage) < 4 {
			partMessage, _ = reader.ReadByte()
			checkMessage = append(checkMessage, partMessage)
			fullMessage = append(fullMessage, partMessage)
		}

		for !(checkMessage[0] == frameString[0] && checkMessage[1] == frameString[1] && checkMessage[2] == frameString[2] && checkMessage[3] == frameString[3]) {
			if debug{
				fmt.Printf("%v\n", checkMessage)
			}
			checkMessage = checkMessage[1:]
			partMessage, _ = reader.ReadByte()
			checkMessage = append(checkMessage, partMessage)
			fullMessage = append(fullMessage, partMessage)
		}
		//fmt.Println("Packet Border found")

		fullMessage = fullMessage[:len(fullMessage)-4]
		fullMessage = append(begin4, fullMessage...)
		fullMessage = append(begin3, fullMessage...)
		fullMessage = append(fullMessage, begin1[0])
		fullMessage = append(fullMessage, begin2[0])

	return fullMessage
}


func calculateCRC(payload []byte) byte {
	var crc byte
	crc = 0
	for i := 0; i < len(payload); i++ {
		crc ^= payload[i]
	}
	//fmt.Printf("Received Payload: %v\n", payload)

	return crc
}

func decode (message []byte) bool{
	var length int16
	length |= int16(message[3])
	length |= int16(message[2]) << 8

	//fmt.Printf("Message length: %v\n", length)

	payload := message[4:len(message)-3]

	// Check CRC
	crc:= calculateCRC(payload)
	if crc != message[len(message)-3]{
		fmt.Println("CRC NOT OK!!!")
		return true
	}


	switch messageID:= payload[0]; messageID{
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
		// can be decoded by gpsdecode
		//Decode_223(payload[1:])
		//return false
	case byte(224):
		// can be decoded by gpsdecode
		//Decode_224(payload[1:])
		//return false
	case byte(225):
		// cannot be decoded by gpsdecode
		Decode_225(payload[1:])
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
		return false
	}


	return true

}


func main() {

	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 19200}

    s, err := serial.OpenPort(c)
    if err != nil {
		log.Printf("Unable to handle error.\n")
		log.Fatal(err)
	}


    reader := bufio.NewReader(s)

	//Get first bad message
    getMessage(reader, false)

    for {
		//Get Good message
		//reply  := getMessage(reader, true)
		reply := getMessage(reader, false)

		//fmt.Printf("%v\n", reply)

		/* WRITE TO FILE
		err = ioutil.WriteFile("/tmp/rtcm.log", reply, 0644)
		if err != nil {
			panic(err)
		}
		*/

		decoded := decode(reply)

		if !decoded {
			// DECODE WITH GPSDECODE
			var out bytes.Buffer
			cmd := exec.Command("gpsdecode")
			cmd.Stdin = bytes.NewBuffer(reply)
			cmd.Stdout = &out
			err = cmd.Run()
			if err != nil {
				fmt.Println("Error")
			}
			fmt.Printf("%v\n", string(out.Bytes()))
		}
	}

}

