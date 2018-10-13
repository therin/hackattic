package brute_force_zip

import (
	"../../pkg/tools"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Solution struct {
	Secret string `json:"secret"`
}

func BruteForceZip() {

	/*
		Zip contains The Dunwich Horror by H. P. Lovecraft book: http://www.gutenberg.org/ebooks/50133.txt.utf-8 book.
		This allows us to use known plaintext attack with this tool: https://github.com/keyunluo/pkcrack
		Zipped book file included in this package for your convenience
	*/

	workingDir := "../../pkg/brute_force_zip/work"

	bytissimo := tools.GetProblem("https://hackattic.com/challenges/brute_force_zip/problem?access_token=" + tools.AccessToken)

	// Convert to JSON string interface
	var anyJson map[string]interface{}
	json.Unmarshal(bytissimo, &anyJson)

	zipUrl := anyJson["zip_url"].(string)

	fmt.Println("zip url:", zipUrl)

	err := downloadFile(workingDir+"/package.zip", zipUrl)
	if err != nil {
		fmt.Println(err)
	}

	// Run plaintext attack on downloaded file
	runCommand("/usr/local/bin/pkcrack",
		"-a -C package.zip -c dunwich_horror.txt -P unprotected.zip -p dunwich_horror.txt -d decrypted.zip",
		workingDir)

	// Extract decrypted zip file
	runCommand("/usr/bin/unzip",
		"-o decrypted.zip secret.txt",
		workingDir)

	// Load secret.txt to variable
	secret, err := ioutil.ReadFile(workingDir + "/secret.txt")
	if err != nil {
		fmt.Print(err)
	}

	// Move response to solution struct and submit
	solution := Solution{strings.TrimSuffix(string(secret), "\n")}
	fmt.Printf("%+v\n", solution)

	solutionJson, err := json.Marshal(&solution)
	if err != nil {
		fmt.Println(err)
	}
	tools.SubmitSolution(solutionJson, "https://hackattic.com/challenges/brute_force_zip/solve?access_token="+tools.AccessToken)

}

func runCommand(cmd string, args string, dir string) {
	command := exec.Command(cmd, strings.Split(args, " ")...)
	command.Dir = dir

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	command.Stdout = mw
	command.Stderr = mw

	// Execute the command
	if err := command.Run(); err != nil {
		log.Panic(err)
	}

	log.Println(stdBuffer.String())

}

func downloadFile(filepath string, url string) error {

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
