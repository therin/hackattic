package help_me_unpack

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"

	"github.com/therin/hackattic/pkg/tools"
)

type Solution struct {
	Int               int32   `json:"int"`
	Uint              uint32  `json:"uint"`
	Short             int16   `json:"short"`
	Float             float64 `json:"float"`
	Double            float64 `json:"double"`
	Big_endian_double float64 `json:"big_endian_double"`
}

func Unpack() {

	bytissimo := tools.GetProblem("https://hackattic.com/challenges/help_me_unpack/problem?access_token=" + tools.AccessToken)

	// convert to JSON string interface
	var anyJson map[string]interface{}
	json.Unmarshal(bytissimo, &anyJson)

	fmt.Println(anyJson["bytes"].(string))

	stringReceived := anyJson["bytes"].(string)

	fmt.Println(tools.Base64Decode(stringReceived))

	bytesToExtract := tools.Base64Decode(stringReceived)

	fmt.Println(bytesToExtract)

	// extract int
	integer := int32(binary.LittleEndian.Uint32(bytesToExtract[0:4]))
	fmt.Println(integer)

	// extract uint
	unsignedInteger := binary.LittleEndian.Uint32(bytesToExtract[4:8])
	fmt.Println(unsignedInteger)

	// extract short
	short := int16(binary.LittleEndian.Uint32(bytesToExtract[8:12]))
	fmt.Println(short)

	// extract float
	unsignedIntegerTemp := binary.LittleEndian.Uint32(bytesToExtract[12:16])
	float := float64(math.Float32frombits(unsignedIntegerTemp))
	fmt.Printf("%0.15f", float)

	// extract double
	unsigned64IntegerTemp := binary.LittleEndian.Uint64(bytesToExtract[16:24])
	double := math.Float64frombits(unsigned64IntegerTemp)
	fmt.Println(double)

	// extract big-endian double
	unsigned64IntegerTempBig := binary.BigEndian.Uint64(bytesToExtract[24:32])
	doubleBig := math.Float64frombits(unsigned64IntegerTempBig)
	fmt.Println(doubleBig)

	// build solution struct
	solution := Solution{integer, unsignedInteger, short, float, double, doubleBig}
	fmt.Println(solution)

	solutionJson, err := json.Marshal(&solution)
	if err != nil {
		fmt.Println(err)
	}
	tools.SubmitSolution(solutionJson, "https://hackattic.com//challenges/help_me_unpack/solve?access_token="+tools.AccessToken)

}
