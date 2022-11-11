package backup_restore

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/therin/hackattic/pkg/tools"
)

// NOTES:
// createdb -T template0 hackattic
// create user postgres;
// alter user postgres with password 'postgres'
// brew services start postgresql
// pg_ctl -D /usr/local/var/postgres start
// zcat < pg_dump > decompressed.sql

type Solution struct {
	AliveSsns []string `json:"alive_ssns"`
}

func Backup_restore() {

	bytissimo := tools.GetProblem("https://hackattic.com/challenges/backup_restore/problem?access_token=" + tools.AccessToken)

	// convert to JSON string interface
	var anyJson map[string]interface{}
	json.Unmarshal(bytissimo, &anyJson)

	stringReceived := anyJson["dump"].(string)

	bytesToDump, err := base64.StdEncoding.DecodeString(stringReceived)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("pg_dump")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	if _, err := file.Write(bytesToDump); err != nil {
		panic(err)
	}
	if err := file.Sync(); err != nil {
		panic(err)
	}
	// decompress file
	_, err = exec.Command("sh", "-c", "zcat < pg_dump > decompressed.sql").Output()
	if err != nil {
		fmt.Println(err)
	}
	// restore backup to test db
	_, err = exec.Command("sh", "-c", "psql hackattick < decompressed.sql").Output()
	if err != nil {
		fmt.Println(err)
	}

	// run query
	queryResult, queryErr := exec.Command("sh", "-c", "psql -U postgres -d hackattick -c \"select ssn from criminal_records WHERE status = 'alive'\"").Output()
	if queryErr != nil {
		fmt.Println(err)
	}
	queryList := strings.Split(string(queryResult), "\n")[2 : len(strings.Split(string(queryResult), "\n"))-3]
	// queryList1 := queryList[2 : len(queryList)-3]
	solutionSlice := make([]string, 0)

	for i := 0; i < len(queryList); i++ {
		if queryList[i] != "" {
			solutionSlice = append(solutionSlice, strings.TrimSpace(queryList[i]))
		}
	}

	fmt.Println(len(solutionSlice))

	// cleanup
	queryResult, queryErr = exec.Command("sh", "-c", "psql -U postgres -d hackattick -c \"drop table criminal_records\"").Output()
	if queryErr != nil {
		fmt.Println(err)
	}

	// move response to solution struct and submit
	solution := Solution{solutionSlice}
	fmt.Printf("%+v\n", solution)

	solutionJson, err := json.Marshal(&solution)
	if err != nil {
		fmt.Println(err)
	}
	tools.SubmitSolution(solutionJson, "https://hackattic.com/challenges/backup_restore/solve?access_token="+tools.AccessToken)
}
