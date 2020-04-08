package main

import (
	"encoding/binary"
	"math"
	"fmt"
)

func Decode_223(payload []byte) {
	fmt.Println("==================================")
	fmt.Println("Navigation state (0xDF 223):")
	if len(payload) != 80 {
		fmt.Println("RECEIVED BAD LENGTH FOR TYPE 223.")
		return
	}
	fmt.Printf("IOD. Issue of DATA %v\n", payload[0])

	switch navState:= payload[1]; navState {
	case byte(0):
		fmt.Println("Navigation state: NO_FIX")
	case byte(1):
		fmt.Println("Navigation state: FIX_PREDICTION")
	case byte(2):
		fmt.Println("Navigation state: FIX_2D")
	case byte(3):
		fmt.Println("Navigation state: FIX_3D")
	case byte(4):
		fmt.Println("Navigation state: FIX_DIFFERENTIAL")
	default:
		fmt.Println("Navigation state: UNKNOWN")
	}

	var WN int16
	WN |= int16(payload[3])
	WN |= int16(payload[2]) << 8
	fmt.Printf("GPS Week Number: %v\n", WN)

	TOW := math.Float64frombits(binary.BigEndian.Uint64(payload[4:12]))
	fmt.Printf("GPS Time of Week: %v\n", TOW)

	ECEF_X := math.Float64frombits(binary.BigEndian.Uint64(payload[12:20]))
	ECEF_Y := math.Float64frombits(binary.BigEndian.Uint64(payload[20:28]))
	ECEF_Z := math.Float64frombits(binary.BigEndian.Uint64(payload[28:36]))
	// fmt.Printf("ECEF Coordinates: X:%v Y:%v Z:%v\n", ECEF_X, ECEF_Y, ECEF_Z)
	r := math.Sqrt(ECEF_X*ECEF_X + ECEF_Y*ECEF_Y + ECEF_Z*ECEF_Z)
	koLatitude := math.Atan(math.Sqrt(ECEF_X*ECEF_X + ECEF_Y*ECEF_Y) / ECEF_Z )
	Longitude := math.Atan(ECEF_Y/ECEF_X)
	fmt.Printf("Radius: %v meter. Latitude: %v Longitude: %v\n", r, 90 - koLatitude*180/math.Pi, Longitude*180/math.Pi)

	ECEF_VEL_X := math.Float32frombits(binary.BigEndian.Uint32(payload[36:40]))
	ECEF_VEL_Y := math.Float32frombits(binary.BigEndian.Uint32(payload[40:44]))
	ECEF_VEL_Z := math.Float32frombits(binary.BigEndian.Uint32(payload[44:48]))
	fmt.Printf("ECEF SPEED: X:%v Y:%v Z:%v\n", ECEF_VEL_X, ECEF_VEL_Y, ECEF_VEL_Z)

	/*
		The third largest error which can be caused by the receiver clock, is its oscillator.
		Both a receiver’s measurement of phase differences and its generation of replica codes depend on
		the reliability of this internal frequency standard.
	*/
	dT := math.Float64frombits(binary.BigEndian.Uint64(payload[48:56]))
	fmt.Printf("Clock BIAS of RECEIVER: %v meter\n", dT)

	C_DRIFT := math.Float32frombits(binary.BigEndian.Uint32(payload[56:60]))
	fmt.Printf("Clock DRIFT of RECEIVER: %v m/s\n", C_DRIFT)

	/*
			HDOP (Horizontal Dilution of Precision) — снижение точности в горизонтальной плоскости
		    VDOP (Vertical) — снижение точности в вертикальной плоскости
		    PDOP (Position) — снижение точности по местоположению
		    TDOP (Time) — снижение точности по времени
		    GDOP (Geometric) — суммарное геометрическое снижение точности по местоположению и времени
	*/
	GDOP := math.Float32frombits(binary.BigEndian.Uint32(payload[60:64]))
	PDOP := math.Float32frombits(binary.BigEndian.Uint32(payload[64:68]))
	HDOP := math.Float32frombits(binary.BigEndian.Uint32(payload[68:72]))
	VDOP := math.Float32frombits(binary.BigEndian.Uint32(payload[72:76]))
	TDOP := math.Float32frombits(binary.BigEndian.Uint32(payload[76:80]))
	fmt.Printf("GDOP: %v PDOP: %v HDOP: %v VDOP: %v TDOP: %v\n", GDOP, PDOP,HDOP,VDOP,TDOP)
	fmt.Println("==================================")
}

func Decode_224(payload []byte) {
	/*
		This is the information about the GPS subframe data bits currently collected in the receiver. The data bits are
		composed from the 24 higher bits of each of the navigation words and the parity bits are not included in the output.
		Only when all 10 navigation words have been verified by parity checking, the data bits in the subframe are output.
		Before being sent out to the host, the data bits are also polarity-adjusted. The 8 preamble bits of a subframe, for
		example, can be obtained from the first byte of the 3-byte field of navigation word 1. This message is sent from the
		receiver to host.
	*/
	fmt.Println("==================================")
	fmt.Println("GPS Subframe Data (0xE0 224):")
	if len(payload) != 32 {
		fmt.Println("RECEIVED BAD LENGTH FOR TYPE 224.")
		return
	}

	// SVN - space vehicle numbers
	// PRN - pseudo-random noise

	// SVID
	fmt.Printf("GPS Satellite PRN (tSV) %v\n", payload[0])
	// SFID
	fmt.Printf("Sub-frame ID (1-5)(frame) %v\n", payload[1])

	for index:=2;index < 32; index += 3 {
		fmt.Printf("24 parity-checked and polarity-adjusted bits of subframe word %v : %v\n", (index-2)/3 +1, payload[index:index+3])
	}

	fmt.Println("==================================")
}


func Decode_225(payload []byte) {
	/*
		This is the information about the string data bits currently collected by the receiver. This message is composed of
		GLONASS satellite slot number, string number and bit 80 to bit 09 in relative bi-binary code of the string. The output
		data bits (bit 80 to bit 09) of each string were already checked as correct by the Hamming code data verification
		algorithm before output by the receiver. The 8 Hamming code check bits (bit 08 to bit 01) are not included in the
		message. The data bits (bit 80 to bit 09) have been polarity-adjusted. This message is sent from the receiver to host.
	*/
	/*
		Это информация о строчных данных битах собранная ресивером.
		Это сообщение состоит из номера спутника, номера строки и битов с 80 по 9 в относительном бидвоичном виде.
		Исходящие биты (с 80 по 09) каждой строки уже были проверены и скорректированы кодом Хемминга перед выходом из ресивера.
		8 бит кода Хемминга (с 8 до 1) не включены в сообщение.
		Биты данных (с 80 по 9) были скорректированы по полярности.
	*/
	fmt.Println("==================================")
	fmt.Println("GPS Subframe Data (0xE1 225):")
	if len(payload) != 11 {
		fmt.Println("RECEIVED BAD LENGTH FOR TYPE 225.")
		return
	}

	// SVID+64
	fmt.Printf("GLONASS satellite slot number: %v\n", payload[0]-64)
	// Суперкадр имеет длительность 2,5 мин и состоит из 5 кадров длительностью 30 с
	// 1-4 кадр - идентичны
	// 5 кадр - особый.

	// Каждый кадр состоит из 15 строк длительностью 2 с.
	// В пределах каждого суперкадра передается полный альманах для всех 24 НКА системы ГЛОНАСС.
	// 1 - 4 строка - Оперативная информация для Навигационного Космического Аппарата(НКА)
	// 5 -относится к неоперативной информации и повторяется в каждом кадре суперкадра.
	// 6 - 15 строка - Альманах (неоперативная информация) (2 строки на спутник)

	// Оперативная информация:
	// оцифровку меток времени НКА;
	// сдвиг шкалы времени НКА относительно шкалы времени системы ГЛОНАСС;
	// относительное отличие несущей частоты излучаемого навигационного радиосигнала от номинального значения;
	// эфемериды НКА и другие параметры.

	// Альманах:
	// данные о состоянии всех НКА системы (альманах состояния);
	// сдвиг шкалы времени каждого НКА относительно шкалы времени системы ГЛОНАСС (альманах фаз);
	// параметры орбит всех НКА системы (альманах орбит);
	// сдвиг шкалы времени системы ГЛОНАСС относительно UTC(SU) и другие параметры.
	// 1 кадр  - 1 – 5 НКА
	// 2 кадр  - 6 – 10 НКА
	// 3 кадр  - 11 – 15 НКА
	// 4 кадр  - 16 – 20 НКА
	// 5 кадр  - 21 – 24 НКА
	// 5 кадр, так как 4 спутника, последние 2 строки - резерв.

	// 85 bit - 0 (for radio transmission)
	// 84 - 81 bit - String number

	var completeString string

	fmt.Printf("String number of navigation message:  %v\n", payload[1])

	for i:=2; i < len(payload); i++{
		completeString += fmt.Sprintf("%08b", payload[i])
	}

	//j:=2
	//for i:=80; i>8; i -=8 {
	//	fmt.Printf("Data bit number %2d-%2d (relative bi-binary):  %08b\n", i, i-7, payload[j])
	//	j++
	//}

	//fmt.Printf("Complete String: %s\n", completeString)
	switch payload[1]{
	case 1:
		glo_decode_1(completeString)
	case 2:
		glo_decode_2(completeString)
	case 3:
		glo_decode_3(completeString)
	case 4:
		glo_decode_4(completeString)
	case 5:
		glo_decode_5(completeString)
	case 6:
		glo_decode_even_almanac(completeString)
	case 7:
		glo_decode_odd_almanac(completeString)
	case 8:
		glo_decode_even_almanac(completeString)
	case 9:
		glo_decode_odd_almanac(completeString)
	case 10:
		glo_decode_even_almanac(completeString)
	case 11:
		glo_decode_odd_almanac(completeString)
	case 12:
		glo_decode_even_almanac(completeString)
	case 13:
		glo_decode_odd_almanac(completeString)
	case 14:
		// TODO: DIFFERENT FUNCTIONS
		//glo_decode_even_almanac(completeString)
		glo_decode_even_spec_almanac(completeString)
	case 15:
		// TODO: DIFFERENT FUNCTIONS
		//glo_decode_odd_almanac(completeString)
		glo_decode_odd_spec_almanac(completeString)
	}

	fmt.Println("==================================")
}
