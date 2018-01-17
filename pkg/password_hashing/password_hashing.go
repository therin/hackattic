package password_hashing

import (
	"../../pkg/tools"
	// "encoding/base64"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	// "os"
	// "os/exec"
	// "strings"
)

type Solution struct {
	Sha256 string `json:"sha256"`
	Hmac   string `json:"hmac"`
	Pbkdf2 string `json:"pbkdf2"`
	Scrypt string `json:"scrypt"`
}

func Password_hashing() {

	bytissimo := tools.GetProblem("https://hackattic.com/challenges/password_hashing/problem?access_token=" + tools.AccessToken)

	// convert to JSON string interface
	var anyJson map[string]interface{}
	json.Unmarshal(bytissimo, &anyJson)

	password := anyJson["password"].(string)
	salt := anyJson["salt"].(string)

	fmt.Println(password, salt)
	sha256Hex := shaify(password)

	fmt.Println(sha256Hex)

}

func shaify(password string) string {
	fmt.Println("hashing", fmt.Sprintf(password))
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString((hasher.Sum(nil)))
}
