// +build lib
package main

import (
	"bytes"
	// "crypto/sha256"
	// "encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "math/bits"
	"net/http"
	// "os"
	L "../tales_of_ssl/lib"
	"strings"
	"time"
)

// https://golang.org/src/crypto/tls/generate_cert.go

type TestResponse struct {
	PrivateKey   string       `json:"private_key"`
	RequiredData RequiredData `json:"required_data"`
}

type RequiredData struct {
	Domain       string `json:"domain"`
	SerialNumber string `json:"serial_number"`
	Country      string `json:"country"`
}

type Solution struct {
	certificate string `json:"certificate"`
}

// func writeGob(filePath string, object interface{}) error {
// 	file, err := os.Create(filePath)
// 	if err == nil {
// 		encoder := gob.NewEncoder(file)
// 		encoder.Encode(object)
// 	}
// 	file.Close()
// 	return err
// }

// func readGob(filePath string, object interface{}) error {
// 	file, err := os.Open(filePath)
// 	if err == nil {
// 		decoder := gob.NewDecoder(file)
// 		err = decoder.Decode(object)
// 	}
// 	file.Close()
// 	return err
// }

func request(url string) TestResponse {

	// netClient := &http.Client{
	// 	Timeout: time.Second * 10,
	// }

	// resp, err := netClient.Get(url)
	// if err != nil {
	// 	fmt.Println("Couldn't get problem:", err)
	// }

	asdfd := strings.NewReader(testJson)

	// asdfd, _ := ioutil.ReadAll(resp.Body)

	// fmt.Printf(string(asdfd))

	decoder := json.NewDecoder(asdfd)
	testResponse := TestResponse{}
	err6 := decoder.Decode(&testResponse)
	if err6 != nil {
		fmt.Printf("Error decoding response:", err6)
	}

	return testResponse
}

// func hashIt(str string) []byte {
// 	h := sha256.New()
// 	h.Write([]byte(str))
// 	fmt.Println("working with", str)

// 	return h.Sum(nil)
// }

// func hashIt2(block Block) [32]byte {
// 	fmt.Println("hashing", fmt.Sprintf("", block))
// 	jsonBuf, err := json.Marshal(&block)
// 	if err != nil {
// 		fmt.Println("Couldn't JSON encode block", block, ":", err)
// 	}
// 	return sha256.Sum256(jsonBuf)
// }

// func LeadingZeroBits(buf [32]byte) int {
// 	zeros := 0
// 	for _, b := range buf {
// 		if b == 0 {
// 			zeros += 8
// 			continue
// 		}
// 		zeros += bits.LeadingZeros8(uint8(b))
// 		break
// 	}
// 	return zeros
// }

// func calcZeroBits(input [32]byte) int {
// 	stopFlag := 0
// 	counter := 0
// 	for _, item := range input {

// 		for i := uint(0); i < 8; i++ {
// 			if (item & (1 << i) >> i) == 0 {
// 				counter += 1
// 			} else {
// 				stopFlag = 1
// 				break
// 			}

// 		}
// 		if stopFlag == 1 {
// 			break
// 		}

// 	}
// 	return counter

// }

func SubmitSolution(nonce string) {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	solution := Solution{nonce}
	solutionJson, err := json.Marshal(&solution)
	if err != nil {
		fmt.Println("Couldn't marshal solution", solution, ":", err)
	}
	url := "https://hackattic.com/challenges/tales_of_ssl/solve?access_token=" + tools.AccessToken
	resp, err := netClient.Post(url,
		"application/json",
		bytes.NewReader(solutionJson),
	)
	if err != nil {
		fmt.Println("Couldn't do POST:", err)
	}

	fmt.Println("Got response to solution:", resp, resp.Body)

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func main() {

	problem := request("https://hackattic.com/challenges/tales_of_ssl/problem?access_token=" + tools.AccessToken)

	fmt.Printf("Got problem %+v\n", problem)

	L.Generate_cert()

	// for i := 0; i < 100000; i++ {
	// 	bytes.Block.Nonce = i
	// 	currentHash := hashIt2(bytes.Block)
	// 	currentZeroBits := LeadingZeroBits(currentHash)
	// 	if currentZeroBits == bytes.Difficulty {
	// 		fmt.Printf("working with", currentZeroBits, bytes.Difficulty)
	// 		fmt.Printf("found desired hash: %s and nonce %d ", fmt.Sprintf("%x", currentHash), i)
	// 		SubmitSolution(i)
	// 		break

	// 	}
	// }

}
