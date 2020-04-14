package skytraqlib

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
)

func DecodeByGPSDecode (reply []byte) []byte {
	// DECODE WITH GPSDECODE
	var out bytes.Buffer
	cmd := exec.Command("gpsdecode")
	cmd.Stdin = bytes.NewBuffer(reply)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("Error during running GPS DECODE")
	}

	return out.Bytes()
}

func WriteBytesToFile (filename string, reply []byte) {

	err := ioutil.WriteFile(filename, reply, 0644)
	if err != nil {
		panic(err)
	}

}