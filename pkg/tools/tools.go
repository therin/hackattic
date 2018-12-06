package tools

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var AccessToken string = ""

func GetProblem(url string) []byte {

	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := netClient.Get(url)
	if err != nil {
		fmt.Println("Couldn't get problem:", err)
	}

	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println("Failed to read byte stream")
	}

	// fmt.Println("response Body:", string(bodyBytes))
	defer resp.Body.Close()
	return bodyBytes
}

func SubmitSolution(input []byte, url string) {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := netClient.Post(url,
		"application/json",
		bytes.NewReader(input),
	)
	if err != nil {
		fmt.Println("Couldn't do POST:", err)
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func Base64Encode(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

func Base64Decode(input string) []byte {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		fmt.Println(err)
		return make([]byte, 0)
	}
	return data
}

func DownloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
