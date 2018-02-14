package collision_course

import (
	"../../pkg/tools"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// https://github.com/brimstone/fastcoll

type Solution struct {
	Files []string `json:"files"`
}

func Collision_course() {

	bytissimo := tools.GetProblem("https://hackattic.com/challenges/collision_course/problem?access_token=" + tools.AccessToken)

	// convert to JSON string interface
	var anyJson map[string]interface{}
	json.Unmarshal(bytissimo, &anyJson)

	stringReceived := anyJson["include"].(string)

	fmt.Println("base string for collision:", string(stringReceived))

	// Write received string to file
	file, err := os.Create("input")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	if _, err := file.Write([]byte(stringReceived)); err != nil {
		panic(err)
	}
	if err := file.Sync(); err != nil {
		panic(err)
	}
	// Run docker container with fastcoll tool to generate solution files
	cmd := exec.Command("sh", "-c", "/usr/local/bin/docker run --rm -i -u $UID:$GID -v $PWD:/work -w /work brimstone/fastcoll --prefixfile input -o msg1.bin msg2.bin")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Output:", outb.String(), "Error:", errb.String())

	// Load solution files
	file1, err := ioutil.ReadFile("msg1.bin")
	if err != nil {
		fmt.Print(err)
	}

	file2, err := ioutil.ReadFile("msg2.bin")
	if err != nil {
		fmt.Print(err)
	}

	// Move response to solution struct and submit
	solution := Solution{[]string{tools.Base64Encode(file1), tools.Base64Encode(file2)}}
	fmt.Printf("%+v\n", solution)

	solutionJson, err := json.Marshal(&solution)
	if err != nil {
		fmt.Println(err)
	}
	tools.SubmitSolution(solutionJson, "https://hackattic.com/challenges/collision_course/solve?access_token="+tools.AccessToken)

}
