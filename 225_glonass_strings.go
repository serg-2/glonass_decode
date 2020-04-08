package main

import (
	"math"
	"strconv"
	"fmt"
)

func glo_decode_1 (fullString string) {

	var P1 int8
	var sign string

	// 2 bits - 0 (80-79)
	fullString = fullString[2:]

	// 2 bits - P1 (78-77)
	// Слово Р1 - признак смены оперативной информации.
	// Сообщает величину интервала времени между значениями tb (мин) в данном и предыдущем кадрах.
	switch fullString[0:2] {
	case "00":
		P1 = 0
	case "01":
		P1 = 30
	case "10":
		P1 = 45
	case "11":
		P1 = 60
	}
	fmt.Printf("Величина интервала времени между значениями слова tB %v минут\n", P1)
	fullString = fullString[2:]

	// 12 bits (76-65)
	// tk - время начала кадра внутри текущих суток, исчисляемое в шкале бортового времени НКА.
	// В пяти старших разрядах записывается количество целых часов, прошедших с начала текущих суток
	// в шести средних - число целых минут
	// в младшем - число тридцатисекундных интервалов, прошедших с начала текущей минуты.
	hours, _ := strconv.ParseInt(fullString[:5], 2, 64)
	fullString = fullString[5:]
	mins, _ := strconv.ParseInt(fullString[:6], 2, 64)
	fullString = fullString[6:]
	secs, _ := strconv.ParseInt(fullString[:1], 2, 64)
	fullString = fullString[1:]
	fmt.Printf("Время начала кадра внутри текущих суток: %v часов %v минут %v 30 секундных интервалов.\n",hours,mins,secs)

	// 24 bits (64-41)
	// X Составляющая вектора скорости данного НКА в геодезической системе системе координат ПЗ-90 на момент времени tb.
	speed, _ := strconv.ParseInt(fullString[1:24], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fullString = fullString[24:]
	fmt.Printf("X составляющая ветора скорости в ПЗ-90: %v%v км/c.\n",sign,float64(speed)/(math.Pow(2,20)))

	// 5 bits (40-36)
	// X составляющая ускорения данного НКА в геодезической системе координат ПЗ-90 на момент времени tb, обусловленные действием Луны и Солнца.
	accel, _ := strconv.ParseInt(fullString[1:5], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fullString = fullString[5:]
	fmt.Printf("X составляющая ветора ускорения: %v%v км/c2.\n",sign,float64(accel)/(math.Pow(2,30)))

	// 27 bits (35-9)
	// X координата данного НКА в системе координат ПЗ-90 на момент времени tb.
	coord, _ := strconv.ParseInt(fullString[1:27], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fullString = fullString[27:]
	fmt.Printf("X координата: %v%v км.\n",sign,float64(coord)/(math.Pow(2,11)))

	if len(fullString) != 0 {
		fmt.Println("PANIC!")
		return
	}
}


func glo_decode_2 (fullString string) {

	var sign string

	// 3 bits - bn (80-78)
	// Слово Bn - признак недостоверности кадра n-го НКА.
	// Аппаратурой потребителя анализируется только старший разряд этого слова
	// "1" в котором обозначает факт непригодности данного спутника для проведения сеансов измерений.
	// Второй и третий разряды этого слова аппаратурой потребителя не анализируются

	bn, _ := strconv.ParseInt(fullString[:3], 2, 64)
	if bn == 0 {
		fmt.Printf("Хороший спутник\n")
	} else {
		fmt.Printf("Плохой спутник\n")
	}
	fullString = fullString[3:]

	// 1 bit - P2(77)
	// признак смены.
	// Он представляет собой признак нечетности ("1") или четности ("О")
	// порядкового номера b 30(60) - минутного текущего отрезка времени,
	// середина которого оцифрована числовым значением слова tb.
	p2, _ := strconv.ParseInt(fullString[:1], 2, 64)
	if p2 == 0 {
		fmt.Printf("Чётный порядковый номер минутного текущего отрезка времени\n")
	} else {
		fmt.Printf("Нечётный порядковый номер минутного текущего отрезка времени\n")
	}
	fullString = fullString[1:]

	// 7 bit - tb (76-70)
	// порядковый номер временного интервала внутри текущих суток по шкале системного времени ГЛОНАСС,
	// к середине которого относится передаваемая в кадре оперативная информация.
	// Длительность данного временного интервала и, следовательно, максимальное значение слова tb
	// определяются значением слова P1
	tb, _ := strconv.ParseInt(fullString[:7], 2, 64)
	fmt.Printf("Порядковый номер временного интервала внутри текущих суток %v\n", tb)
	fullString = fullString[7:]

	// 5 bits - 0 (69-65)
	fullString = fullString[5:]

	// 24 bits (64-41)
	// Y Составляющая вектора скорости данного НКА в геодезической системе системе координат ПЗ-90 на момент времени tb.
	speed, _ := strconv.ParseInt(fullString[1:24], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fullString = fullString[24:]
	fmt.Printf("Y составляющая ветора скорости в ПЗ-90: %v%v км/c.\n",sign,float64(speed)/(math.Pow(2,20)))

	// 5 bits (40-36)
	// Y составляющия ускорения данного НКА в геодезической системе координат ПЗ-90 на момент времени tb, обусловленные действием Луны и Солнца.
	accel, _ := strconv.ParseInt(fullString[1:5], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fullString = fullString[5:]
	fmt.Printf("Y составляющая ветора ускорения: %v%v км/c2.\n",sign,float64(accel)/(math.Pow(2,30)))

	// 27 bits (35-9)
	// Y координата данного НКА в системе координат ПЗ-90 на момент времени tb.
	coord, _ := strconv.ParseInt(fullString[1:27], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fullString = fullString[27:]
	fmt.Printf("Y координата: %v%v км.\n",sign,float64(coord)/(math.Pow(2,11)))

	if len(fullString) != 0 {
		fmt.Println("PANIC!")
		return
	}
}

func glo_decode_3 (fullString string) {

	var sign string

	// 1 bit - P3 (80)
	// признак, состояние "1" которого означает, что в данном кадре передается альманах для 5-ти спутников системы,
	// а состояние "0" означает, что в данном кадре передается альманах для 4-х спутников
	p3, _ := strconv.ParseInt(fullString[:1], 2, 64)
	if p3 == 0 {
		fmt.Printf("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@\n")
		fmt.Printf("Это 5 кадр в суперкадре (4 спутника в альманахе)\n")
		fmt.Printf("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@\n")
	} else {
		fmt.Printf("Это с 1 по 4 кадр в суперкадре(5 спутников в альманахе)\n")
	}
	fullString = fullString[1:]

	// 11 bit GammaN(tb) (79-69)
	// относительное отклонение прогнозируемого значения несущей частоты
	// излучаемого навигационного радиосигнала n-го спутника от номинального значения на момент времени tb
	// Формула:
	// Gn(tb) = (прогнозируемое значение несущей частоты излучаемого навигационного радиосигнала n-го спутника с учетом гравитационного и релятивистского эффектов на момент времени tb
	// - номинальное значение несущей частоты навигационного радиосигнала n-го спутника ) / номинальное значение несущей частоты навигационного радиосигнала n-го спутника.
	gntb, _ := strconv.ParseInt(fullString[1:11], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fullString = fullString[11:]
	fmt.Printf("Относительное отклонение прогнозируемого значения несущей частоты %v%v.\n",sign,float64(gntb)/(math.Pow(2,40)))

	// 1 bits - 0 (68)
	fullString = fullString[1:]

	// 2 bits - P (67-66)
	// признак режима работы НКА по частотно-временной информации.
	switch fullString[0:2] {
	case "00":
		fmt.Println("Ретрансляция параметра tc, ретрансляция параметра tGPS")
	case "01":
		fmt.Println("Ретрансляция параметра tc, размножение параметра tGPS на борту НКА")
	case "10":
		fmt.Println("Размножение параметра tc на борту НКА, ретрансляция параметра tGPS")
	case "11":
		fmt.Println("Размножение параметра tc на борту НКА, размножение параметра tGPS на борту НКА.")
	}
	fullString = fullString[2:]

	// 1 bits - In (65)
	//признак недостоверности кадра n-го НКА; In = 1 означает факт непригодности данного спутника для навигации.
	In, _ := strconv.ParseInt(fullString[:1], 2, 64)
	if In == 0 {
		fmt.Printf("Хороший спутник\n")
	} else {
		fmt.Printf("Плохой спутник\n")
	}
	fullString = fullString[1:]

	// 24 bits (64-41)
	// Z Составляющая вектора скорости данного НКА в геодезической системе системе координат ПЗ-90 на момент времени tb.
	speed, _ := strconv.ParseInt(fullString[1:24], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fullString = fullString[24:]
	fmt.Printf("Z составляющая ветора скорости в ПЗ-90: %v%v км/c.\n",sign,float64(speed)/(math.Pow(2,20)))

	// 5 bits (40-36)
	// Z составляющия ускорения данного НКА в геодезической системе координат ПЗ-90 на момент времени tb, обусловленные действием Луны и Солнца.
	accel, _ := strconv.ParseInt(fullString[1:5], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fullString = fullString[5:]
	fmt.Printf("Z составляющая ветора ускорения: %v%v км/c2.\n",sign,float64(accel)/(math.Pow(2,30)))

	// 27 bits (35-9)
	// Z координата данного НКА в системе координат ПЗ-90 на момент времени tb.
	coord, _ := strconv.ParseInt(fullString[1:27], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fullString = fullString[27:]
	fmt.Printf("Z координата: %v%v км.\n",sign,float64(coord)/(math.Pow(2,11)))

	if len(fullString) != 0 {
		fmt.Println("PANIC!")
		return
	}
}

func glo_decode_4 (fullString string) {

	var sign string
	var TI float32

	// 22 bit - tn(tb) (80-59)
	// сдвиг шкалы времени n спутника t, относительно шкалы времени системы ГЛОНАСС tc,
	// равный смещению по фазе ПСПД излучаемого навигационного радиосигнала n-го спутника
	// относительно системного опорного сигнала на момент времени tb


	tntb, _ := strconv.ParseInt(fullString[1:22], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fullString = fullString[22:]
	fmt.Printf("Сдвиг шкалы времени спутника n, относительно шкалы времени системы ГЛОНАСС: %v%v сек\n",sign,float64(tntb)/(math.Pow(2,30)))

	// 5 bit - dtn(58-54)
	// dtn - смещение излучаемого навигационного радиосигнала поддиапазона L2
	// относительно навигационного радиосигнала поддиапазона L1 для n-го НКА
	// dtn = tf2 – tf1,
	// где tf1, tf2 – аппаратурные задержки в соответствующих поддиапазонах, выраженные в единицах времени
	dtn, _ := strconv.ParseInt(fullString[1:5], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fullString = fullString[5:]
	fmt.Printf("Смещение излучаемого навигационного радиосигнала поддиапазона L2 относительно навигационного радиосигнала поддиапазона L1 для n-го НКА: %v%v сек\n",sign,float64(dtn)/(math.Pow(2,30)))

	// 5 bit - En(53-49)
	// En характеризует "возраст" оперативной информации,
	// то есть интервал времени, прошедший от момента расчета (закладки) оперативной информации
	// до момента времени tb для n-го спутника.
	// Слово Еn формируется на борту НКА.
	En, _ := strconv.ParseInt(fullString[0:5], 2, 64)
	fullString = fullString[5:]
	fmt.Printf("\"Возраст\" оперативной информациии : %v суток\n", En)

	// 14 bits - 0 (48 - 35)
	fullString = fullString[14:]

	// 1 bits - P4 (34)
	// признак, смена состояния "0" или "1" которого означает,
	// что в данном кадре передается обновленная эфемеридная или частотно-временная информация.
	P4, _ := strconv.ParseInt(fullString[:1], 2, 64)
	if P4 == 0 {
		fmt.Printf("Эфемеридная информация\n")} else {fmt.Printf("Частотно-временная информация\n")
	}
	fullString = fullString[1:]

	// 4 bit - Ft(33-30)
	// фактор точности. В виде эквивалентной ошибки характеризует ошибку набора данных,
	// излучаемых в навигационном сообщении в момент времени tb.
	Ft, _ := strconv.ParseInt(fullString[:4], 2, 64)
	switch Ft {
	case 0:
		TI = 0
	case 1:
		TI = 2
	case 2:
		TI = 2.5
	case 3:
		TI = 4
	case 4:
		TI = 5
	case 5:
		TI = 7
	case 6:
		TI = 10
	case 7:
		TI = 12
	case 8:
		TI = 14
	case 9:
		TI = 16
	case 10:
		TI = 32
	case 11:
		TI = 64
	case 12:
		TI = 128
	case 13:
		TI = 256
	case 14:
		TI = 512
	case 15:
		TI = -1
	}
	fmt.Printf("Точность измерений OMEGA: %v метров\n", TI)
	fullString = fullString[4:]

	// 3 bits - 0 (29 - 27)
	fullString = fullString[3:]

	// 11 bits - Nt (26-16)
	// текущая дата, календарный номер суток внутри четырехлетнего интервала,
	// начиная с 1-го января високосного года (1).
	// Алгоритм пересчета от номера суток внутри четырехлетнего интервала к общепринятой форме даты (чч.мм.гг.)
	// предполагает следующие действия.
	// 1). Вычисляется номер текущего года J в четырехлетнем интервале:
	// если 1 ≤ Nt ≤ 366; J = 1;
	// если 367 ≤ Nt ≤ 731; J = 2;
	// если 732 ≤ Nt ≤ 1096; J = 3;
	// если 1097 ≤ Nt ≤ 1461; J = 4.
	// 2). Вычисляется текущий год в общепринятой форме:
	// Y = 1996 + 4*(N4 –1) + (J – 1).
	Nt, _ := strconv.ParseInt(fullString[:11], 2, 64)
	fullString = fullString[11:]
	fmt.Printf("Календарный номер суток внутри четырехлетнего интервала: %v\n", Nt)

	// 5 bits - n (15-11)
	// Номер НКА, излучающего данный навигационный сигнал и соответствующий его рабочей точке
	// внутри орбитальной группировки ГЛОНАСС
	n, _ := strconv.ParseInt(fullString[:5], 2, 64)
	fullString = fullString[5:]
	fmt.Printf("Номер спутника ГЛОНАСС: %v\n", n)

	// 2 bits - M (10-9)
	// модификация НКА, излучающего данный навигационный сигнал.
	// Значение "00" означает НКА ГЛОНАСС, "01" – НКА ГЛОНАСС-М
	switch fullString[0:2] {
	case "00":
		fmt.Println("Тип спутника: НКА ГЛОНАСС")
	case "01":
		fmt.Println("Тип спутника: НКА ГЛОНАСС-М")
	case "10":
		fmt.Println("ERROR!")
	case "11":
		fmt.Println("ERROR!")
	}
	fullString = fullString[2:]

	if len(fullString) != 0 {
		fmt.Println("PANIC!")
		return
	}
}

func glo_decode_5 (fullString string) {

	var sign string

	// 11 bit - Na (80-70)
	// календарный номер суток внутри четырехлетнего периода, начиная с високосного года,
	// к которым относятся поправка tc и данные альманаха системы (альманах орбит и альманах фаз);
	Na, _ := strconv.ParseInt(fullString[:11], 2, 64)
	fullString = fullString[11:]
	fmt.Printf("Календарный номер суток внутри четырехлетнего периода: %v\n", Na)

	// 32 bit - tc (69-38)
	// поправка к шкале времени системы ГЛОНАСС относительно UTC(SU).
	// Поправка tc дана на начало суток с номером Na
	// NB: Предполагается увеличить цену младшего разряда слова tс до 2-31 с
	// ( то есть до 0.46 нс) за счет увеличения в навигационном сообщении спутника ГЛОНАСС-М
	// разрядности tс с 28 до 32 разрядов.
	// Слово будет расположено в 5-ой, 20-ой, 35-ой, 50-ой и 65-ой строках суперкадра с 38-го по 69 разряд.

	tc, _ := strconv.ParseInt(fullString[1:32], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fullString = fullString[32:]
	fmt.Printf("Поправка к шкале времени системы ГЛОНАСС относительно UTC: %v%v сек\n", sign, float64(tc)/(math.Pow(2,31)))

	// 1 bits - 0 (37)
	fullString = fullString[1:]

	// 5 bits - N4 (36-32)
	// N4 - номер четырехлетнего периода, первый год нулевого четырехлетия соответствует 1996 году.
	N4, _ := strconv.ParseInt(fullString[:5], 2, 64)
	fullString = fullString[5:]
	fmt.Printf("Номер 4-х летнего периода с 1996 года: %v\n", N4)

	// 22 bits -  tGPS (31-10)
	// поправка на расхождение системных шкал времени GPS(TGPS) и ГЛОНАСС (ТГЛ)
	// в соответствии со следующим выражением:
	//TGPS – TГЛ = DT + tGPS,
	//где DT - целая часть, а tGPS - дробная часть расхождения шкал времени, выраженного в секундах.
	// Целая часть расхождения DT определяется потребителем из навигационного сообщения системы GPS
	tGPS, _ := strconv.ParseInt(fullString[1:22], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fmt.Printf("Поправка на расхождение системных шкал времени GPS и ГЛОНАСС: %v%v сек\n", sign, float64(tGPS)/(math.Pow(2,30)))
	fullString = fullString[22:]

	// 1 bit - In (9)
	// признак недостоверности кадра n-го НКА; In = 1 означает факт непригодности данного спутника для навигации.
	In, _ := strconv.ParseInt(fullString[:1], 2, 64)
	if In == 0 {
		fmt.Printf("Хороший спутник\n")
	} else {
		fmt.Printf("Плохой спутник\n")
	}
	fullString = fullString[1:]

	if len(fullString) != 0 {
		fmt.Println("PANIC!")
		return
	}
}

func glo_decode_even_almanac (fullString string) {

	var sign string

	// 1 bit - CnA (80)
	// обобщенный признак состояния спутника с номером nA на момент закладки неоперативной информации
	// (альманаха орбит и фаз).
	// Значение признака Сn = 0 указывает на непригодность спутника для использования в сеансах навигационных определений, а значение признака Сn = 1 - на пригодность спутника.
	CnA, _ := strconv.ParseInt(fullString[:1], 2, 64)
	if CnA == 0 {
		fmt.Printf("Плохой спутник\n")
	} else {
		fmt.Printf("Хороший спутник\n")
	}
	fullString = fullString[1:]

	// 2 bit - MnA (78-79)
	// признак модификации n-го НКА (1); "00" – ГЛОНАСС,
	// "01" - ГЛОНАСС-М;
	switch fullString[0:2] {
	case "00":
		fmt.Println("Тип спутника: НКА ГЛОНАСС")
	case "01":
		fmt.Println("Тип спутника: НКА ГЛОНАСС-М")
	case "10":
		fmt.Println("ERROR!")
	case "11":
		fmt.Println("ERROR!")
	}
	fullString = fullString[2:]

	// 5 bit - nA (73-77)
	// условный номер спутника в системе, который соответствует номеру занимаемой спутником рабочей точки
	nA, _ := strconv.ParseInt(fullString[:5], 2, 64)
	fullString = fullString[5:]
	fmt.Printf("Условный номер спутника в системе: %v\n", nA)

	// 10 bit tnA (72-63)
	// грубое значение сдвига шкалы времени спутника с номером nA относительно шкалы времени системы
	// на момент времени tl nA, равное смещению ПСПД излучаемого навигационного радиосигнала
	// относительно номинального положения, выраженному в единицах времени
	tnA, _ := strconv.ParseInt(fullString[1:10], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fmt.Printf("Грубое значение сдвига шкалы времени спутника с номером nA относительно шкалы времени системы на момент времени: %v%v сек\n", sign, float64(tnA)/(math.Pow(2,18)))
	fullString = fullString[10:]

	// 21 bit lnA (62-42)
	// Долгота в системе координат ПЗ-90 первого внутри суток с номером nA восходящего узла орбиты спутника с номером nA
	lnA, _ := strconv.ParseInt(fullString[1:21], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fmt.Printf("Долгота в системе координат ПЗ-90 первого внутри суток с номером nA восходящего узла орбиты спутника с номером nA: %v%v градуса\n", sign, float64(lnA)/(math.Pow(2,20))*180/math.Pi)
	fullString = fullString[21:]

	// 18 bit DinA (41-24)
	// Поправка к среднему значению наклонения орбиты для спутника с номером nA на момент tlnA
	// (среднее значение наклонения орбиты принято равным 63 градуса);
	DinA, _ := strconv.ParseInt(fullString[1:18], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fmt.Printf("Поправка к среднему значению наклонения орбиты для спутника с номером nA на момент tlnA: %v%v градуса\n", sign, float64(DinA)/(math.Pow(2,20))*180/math.Pi)
	fullString = fullString[18:]

	// 15 bit enA (9-23)
	// эксцентриситет орбиты спутника с номером nA на момент времени tlnA
	enA, _ := strconv.ParseInt(fullString[:15], 2, 64)
	fmt.Printf("Эксцентриситет орбиты спутника с номером nA на момент времени tlnA: %v\n", float64(enA)/(math.Pow(2,20)))
	fullString = fullString[15:]

	if len(fullString) != 0 {
		fmt.Println("PANIC!")
		return
	}
}

func glo_decode_odd_almanac (fullString string) {

	var sign string

	// 16 bit - wnA (80-65)
	// аргумент перигея орбиты спутника с номером nA на момент времени tlnA
	wnA, _ := strconv.ParseInt(fullString[1:16], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fmt.Printf("Аргумент перигея орбиты спутника с номером nA на момент времени tlnA: %v%v градуса\n", sign, float64(wnA)/(math.Pow(2,15))*180/math.Pi)
	fullString = fullString[16:]

	// 21 bit - tlnA (64-44)
	// время прохождения первого внутри суток с номером Na восходящего узла орбиты спутника с номером nA
	tlnA, _ := strconv.ParseInt(fullString[0:21], 2, 64)
	fmt.Printf("Время прохождения первого внутри суток с номером Na восходящего узла орбиты спутника с номером nA: %v сек\n", float64(tlnA)/(math.Pow(2,5)))
	fullString = fullString[21:]

	// 22 bit - DTnA (43-22)
	// поправка к среднему значению драконического периода обращения спутника с номером nA
	// на момент времени tlnA
	// (среднее значение драконического периода обращения спутника принято равным 43200 сек)
	DTnA, _ := strconv.ParseInt(fullString[1:22], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fmt.Printf("Поправка к среднему значению драконического периода обращения спутника с номером nA на момент времени tlnA: %v%v сек/виток\n", sign, float64(DTnA)/(math.Pow(2,9)))
	fullString = fullString[22:]

	// 7 bit - DTnA_SPEED (21-15)
	// скорость изменения драконического периода обращения спутника с номером nA
	DTnA_SPEED, _ := strconv.ParseInt(fullString[1:7], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fmt.Printf("Cкорость изменения драконического периода обращения спутника с номером nA: %v%v сек/виток2\n", sign, float64(DTnA_SPEED)/(math.Pow(2,14)))
	fullString = fullString[7:]

	// 5 bit - HnA (14-10)
	// номер несущей частоты навигационного радиосигнала, излучаемого спутником с номером nA
	Hna, _ := strconv.ParseInt(fullString[:5], 2, 64)
	fmt.Printf("Номер несущей частоты навигационного радиосигнала, излучаемого спутником с номером nA: %v\n", Hna)
	fullString = fullString[5:]

	// 1 bit - In (9)
	// признак недостоверности кадра n-го НКА; In = 1 означает факт непригодности данного спутника для навигации.
	In, _ := strconv.ParseInt(fullString[:1], 2, 64)
	if In == 0 {
		fmt.Printf("Хороший спутник\n")
	} else {
		fmt.Printf("Плохой спутник\n")
	}
	fullString = fullString[1:]

	if len(fullString) != 0 {
		fmt.Println("PANIC!")
		return
	}
}

func glo_decode_odd_spec_almanac (fullString string) {

		// 71 bits RESERVED
		fullString = fullString[71:]

		// 1 bit - In (9)
		// признак недостоверности кадра n-го НКА; In = 1 означает факт непригодности данного спутника для навигации.
		In, _ := strconv.ParseInt(fullString[:1], 2, 64)
		if In == 0 {
			fmt.Printf("Хороший спутник\n")
		} else {
			fmt.Printf("Плохой спутник\n")
		}
		fullString = fullString[1:]

		if len(fullString) != 0 {
			fmt.Println("PANIC!")
			return
		}

}

func glo_decode_even_spec_almanac (fullString string) {

	var sign string

	// 11 bits - B1 (80-70)
	// коэффициент для определения D UT1, равный величине расхождения всемирного и координированного времени
	// на начало текущих суток
	B1, _ := strconv.ParseInt(fullString[1:11], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fmt.Printf("Коэффициент для определения D UT1, равный величине расхождения всемирного и координированного времени на начало текущих суток: %v%v сек\n", sign, float64(B1)/(math.Pow(2,10)))
	fullString = fullString[11:]

	// 10 bits - B2 (69-60)
	// коэффициент для определения D UT1, равный величине суточного изменения расхождения D UT1
	B2, _ := strconv.ParseInt(fullString[1:10], 2, 64)
	if string(fullString[0]) == "1" {
		sign = "-"
	} else {
		sign = "+"
	}
	fmt.Printf("Коэффициент для определения D UT1, равный величине суточного изменения расхождения D UT1: %v%v сек/ССС\n", sign, float64(B2)/(math.Pow(2,16)))

	fullString = fullString[10:]

	// 2 bits - KP(59-58)
	// признак ожидаемой секундной коррекции шкалы UTC на величину ± 1 с
	switch fullString[0:2] {
	case "00":
		fmt.Println("В конце текущего квартала коррекции UTC не будет.")
	case "01":
		fmt.Println("В конце текущего квартала будет коррекция на плюс 1 с.")
	case "10":
		fmt.Println("ERROR!")
	case "11":
		fmt.Println("В конце текущего квартала будет коррекция на минус 1 с.")
	}
	fullString = fullString[2:]

	// 49 - RESERVED
	fullString = fullString[49:]

	if len(fullString) != 0 {
		fmt.Println("PANIC!")
		return
	}

}