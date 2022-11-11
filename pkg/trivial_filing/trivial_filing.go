package trivialfiling

import (
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
var port string = "9876"

func TrivialFiling() {
	go tftpServer()

	bytissimo := tools.GetProblem("https://hackattic.com/challenges/trivial_filing/problem?access_token=" + tools.AccessToken)
	fmt.Println(string(bytissimo))

	var files map[string]map[string]string
	json.Unmarshal(bytissimo, &files)

	for file, content := range files["files"] {
		fmt.Printf("file[%s] content[%s]\n", file, content)
		writeFile(file, strings.NewReader(content))
	}

	// post tftp url to solution endpoint
	solution := Solution{externalAddress, port}
	solutionJson, err := json.Marshal(&solution)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Submitting solution endpoint: %s", solutionJson)
	tools.SubmitSolution(solutionJson, "https://hackattic.com/challenges/trivial_filing/solve?access_token="+tools.AccessToken)

	fmt.Println("Sleeping for 10 seconds to let the attic do the download")
	time.Sleep(10 * time.Second)

}

func tftpServer() {
	s := tftp.NewServer(readHandler, writeHandler)
	s.SetTimeout(20 * time.Second)                                           // optional
	err := s.ListenAndServe(strings.Join([]string{localAddress, port}, ":")) // blocks until s.Shutdown() is called
	if err != nil {
		fmt.Fprintf(os.Stdout, "server: %v\n", err)
		os.Exit(1)
	}
}

// readHandler is called when client starts file download from server
func readHandler(filename string, rf io.ReaderFrom) error {
	fmt.Printf("Got the request for file: %s \n", filename)
	raddr := rf.(tftp.OutgoingTransfer).RemoteAddr()
	log.Println("RRQ from", raddr.String())
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("Reading from the file: %s \n", filename)
	n, err := rf.ReadFrom(file)
	fmt.Print(n)
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

	c, err := tftp.NewClient(strings.Join([]string{localAddress, port}, ":"))
	c.SetTimeout(20 * time.Second) // optional
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	c.SetTimeout(5 * time.Second) // optional
	rf, err := c.Send(name, "octet")
	n, err := rf.ReadFrom(content)
	fmt.Printf("%d bytes sent\n", n)

}
