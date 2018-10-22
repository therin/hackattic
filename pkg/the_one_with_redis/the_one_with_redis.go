package the_one_with_redis

import (
	"../../pkg/tools"
	"encoding/json"
	"fmt"
	"github.com/matthewjhe/rdb"
	"io/ioutil"
	"regexp"
	"time"
)

var dbCounter int = 0
var typeToBeFound string
var expiryTimestamp int
var emojiValue string

type Problem struct {
	Rdb          string `json:"rdb"`
	Requirements struct {
		CheckTypeOf string `json:"check_type_of"`
	}
}

type MyFilter struct {
	Problem
}

func (f MyFilter) Key(k rdb.Key) bool {

	if k.Expiry != -1 {
		fmt.Println("key:", k)
		fmt.Println("key expiry:", k.Expiry)
		expiryTimestamp = k.Expiry
	}

	return false
}

func (f MyFilter) Set(v *rdb.Set) {
	if v.Key.Key == f.Problem.Requirements.CheckTypeOf {
		typeToBeFound = "set"

	}
}

func (f MyFilter) Type(t rdb.Type) bool {
	return false
}

func (f MyFilter) Database(db rdb.DB) bool {
	dbCounter += 1
	return false
}

func (f MyFilter) List(v *rdb.List) {
	if v.Key.Key == f.Problem.Requirements.CheckTypeOf {
		typeToBeFound = "list"

	}
}

func (f MyFilter) Hash(v *rdb.Hash) {
	if v.Key.Key == f.Problem.Requirements.CheckTypeOf {
		typeToBeFound = "hash"

	}
	fmt.Println("hash:", v.Key.Key, v.Values)
}

func (f MyFilter) String(v *rdb.String) {
	fmt.Println("string:", v.Key.Key, v.Value)

	// check for emojineess
	var emojiRx = regexp.MustCompile(`[\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}]`)
	var s = emojiRx.FindString(v.Key.Key)
	if len(s) != 0 {
		fmt.Println("found emoji!:", s)
		emojiValue = v.Value
	}

	if v.Key.Key == f.Problem.Requirements.CheckTypeOf {
		typeToBeFound = "string"

	}
}

func (f MyFilter) SortedSet(v *rdb.SortedSet) { fmt.Println("sortedset:", v.Key.Key, v.Values) }

func TheOneWithRedis() {

	bytissimo := tools.GetProblem("https://hackattic.com/challenges/the_redis_one/problem?access_token=" + tools.AccessToken)

	// Convert to JSON string interface
	problemJson := Problem{}

	json.Unmarshal(bytissimo, &problemJson)

	// Fix corrupt RDB dump header
	fixedProblem := "REDIS" + string(tools.Base64Decode(problemJson.Rdb))[5:]

	// Write the body to file to enable the use of RDB library (lame)
	err := ioutil.WriteFile("./dump.rdb", []byte(fixedProblem), 0777)
	if err != nil {
		panic(err)
	}

	const file = "./dump.rdb"
	reader, err := rdb.NewBufferReader(file, 0)
	if err != nil {
		panic(err)
	}

	if err := rdb.Parse(reader, rdb.WithFilter(MyFilter{problemJson})); err != nil {
		panic(err)
	}

	fmt.Println("Located databases:", dbCounter)
	fmt.Println("Located type to be found:", typeToBeFound)
	fmt.Println("Located expiry timestamp:", time.Unix(0, int64(expiryTimestamp)*int64(time.Millisecond)))
	fmt.Println("Located emoji value:", emojiValue)

	// Dynamically build solution JSON:
	dynamicSolution := make(map[string]interface{})

	dynamicSolution[problemJson.Requirements.CheckTypeOf] = typeToBeFound
	dynamicSolution["db_count"] = dbCounter
	dynamicSolution["emoji_key_value"] = emojiValue
	dynamicSolution["expiry_millis"] = expiryTimestamp

	// Marshal to JSON
	jsonSolution, err := json.Marshal(&dynamicSolution)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", dynamicSolution)
	fmt.Printf("%+v\n", string(jsonSolution))

	tools.SubmitSolution([]byte(jsonSolution), "https://hackattic.com/challenges/the_redis_one/solve?access_token="+tools.AccessToken)

}
