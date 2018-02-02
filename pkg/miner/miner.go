package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/bits"
	"net/http"
	"os"
	"time"
)

type TestResponse struct {
	Difficulty int   `json:"difficulty"`
	Block      Block `json:"block"`
}

type Block struct {
	Data  []interface{} `json:"data"`
	Nonce int           `json:"nonce"`
}

type Solution struct {
	Nonce int `json:"nonce"`
}

func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func readGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

func request(url string) TestResponse {

	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := netClient.Get(url)
	if err != nil {
		fmt.Println("Couldn't get problem:", err)
	}

	decoder := json.NewDecoder(resp.Body)
	testResponse := TestResponse{}
	err6 := decoder.Decode(&testResponse)
	if err6 != nil {
		fmt.Printf("Error decoding response:", err6)
	}

	return testResponse
}

func hashIt(str string) []byte {
	h := sha256.New()
	h.Write([]byte(str))
	fmt.Println("working with", str)

	return h.Sum(nil)
}

func hashIt2(block Block) [32]byte {
	fmt.Println("hashing", fmt.Sprintf("", block))
	jsonBuf, err := json.Marshal(&block)
	if err != nil {
		fmt.Println("Couldn't JSON encode block", block, ":", err)
	}
	return sha256.Sum256(jsonBuf)
}

func LeadingZeroBits(buf [32]byte) int {
	zeros := 0
	for _, b := range buf {
		if b == 0 {
			zeros += 8
			continue
		}
		zeros += bits.LeadingZeros8(uint8(b))
		break
	}
	return zeros
}

func calcZeroBits(input [32]byte) int {
	stopFlag := 0
	counter := 0
	for _, item := range input {

		for i := uint(0); i < 8; i++ {
			if (item & (1 << i) >> i) == 0 {
				counter += 1
			} else {
				stopFlag = 1
				break
			}

		}
		if stopFlag == 1 {
			break
		}

	}
	return counter

}

func submitSolution(nonce int) {

	url := "https://hackattic.com//challenges/mini_miner/solve?access_token=" + tools.AccessToken

	var jsonStr = []byte(fmt.Sprintf(`{"nonce":"%d"}`, nonce))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}

func SubmitSolution(nonce int) {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	solution := Solution{nonce}
	solutionJson, err := json.Marshal(&solution)
	if err != nil {
		fmt.Println("Couldn't marshal solution", solution, ":", err)
	}
	url := "https://hackattic.com//challenges/mini_miner/solve?access_token=" + tools.AccessToken
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

	bytes := request("https://hackattic.com/challenges/mini_miner/problem?access_token=" + tools.AccessToken)

	fmt.Printf("Got problem %+v\n", bytes)

	for i := 0; i < 100000; i++ {
		bytes.Block.Nonce = i
		currentHash := hashIt2(bytes.Block)
		currentZeroBits := LeadingZeroBits(currentHash)
		if currentZeroBits == bytes.Difficulty {
			fmt.Printf("working with", currentZeroBits, bytes.Difficulty)
			fmt.Printf("found desired hash: %s and nonce %d ", fmt.Sprintf("%x", currentHash), i)
			SubmitSolution(i)
			break

		}
	}

}
