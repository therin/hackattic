package trivialfiling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pin/tftp"
	"github.com/therin/hackattic/pkg/tools"
)

type Challenge struct {
	Files map[string]string `json:"files"`
}

type Solution struct {
	TftpHost string `json:"tftp_host"`
	TftpPort string `json:"tftp_port"`
}

var localAddress string = ""
var externalAddress string = ""
var port string = "69"

func TrivialFiling() {
	// go tftpServer()

	bytissimo := tools.GetProblem("https://hackattic.com/challenges/trivial_filing/problem?access_token=" + tools.AccessToken)
	fmt.Println(string(bytissimo))

	var files map[string]map[string]string
	json.Unmarshal(bytissimo, &files)

	for file, content := range files["files"] {
		fmt.Printf("file[%s] content[%s]\n", file, content)
		writeFile(file, strings.NewReader(content))
	}

	// post tftp url to solution endpoint
	solution := Solution{externalAddress, "69"}
	solutionJson, err := json.Marshal(&solution)
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)
	fmt.Printf("Submitting solution endpoint: %s", solutionJson)
	tools.SubmitSolution(solutionJson, "https://hackattic.com/challenges/trivial_filing/solve?access_token="+tools.AccessToken)

	fmt.Println("Sleeping for 10 seconds to let the attic do the download")
	time.Sleep(10 * time.Second)

}

func tftpServer() {
	s := tftp.NewServer(readHandler, writeHandler)
	s.SetTimeout(20 * time.Second) // optional
	err := s.ListenAndServe(":69") // blocks until s.Shutdown() is called
	if err != nil {
		fmt.Fprintf(os.Stdout, "server: %v\n", err)
		os.Exit(1)
	}
}

// readHandler is called when client starts file download from server
func readHandler(filename string, rf io.ReaderFrom) error {
	fmt.Printf("Rf is: %s \n", rf)
	fmt.Printf("Got the request for file: %s \n", filename)
	raddr := rf.(tftp.OutgoingTransfer).RemoteAddr()

	log.Println("RRQ from", raddr.String())
	rf.(tftp.OutgoingTransfer).SetSize(12)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("Reading from the file: %s \n", filename)
	n, err := rf.ReadFrom(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes sent\n", n)
	return nil
}

// writeHandler is called when client starts file upload to server
func writeHandler(filename string, wt io.WriterTo) error {
	fmt.Printf("Got the write request for file: %s \n", filename)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	n, err := wt.WriteTo(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes received\n", n)
	return nil
}

func writeFile(name string, content io.Reader) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(content)

	err := os.WriteFile(name, buf.Bytes(), 0777)
	if err != nil {
		log.Fatal(err)
	}

}
