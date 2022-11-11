package password_hashing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/therin/hackattic/pkg/tools"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
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

	// Calculate SHA256
	// fmt.Println(password, salt)
	sha256Bytes := shaify(password)
	// fmt.Println(hex.EncodeToString(sha256Bytes))

	// Calculate HMAC_SHA256
	fmt.Println("salt:", salt)
	hmacSHA256Bytes := hmacyize(salt, password)
	fmt.Println("hmac:", hmacSHA256Bytes)

	// Calculate PBKDF2
	hash := anyJson["pbkdf2"].(map[string]interface{})["hash"]
	iterations := anyJson["pbkdf2"].(map[string]interface{})["rounds"].(float64)
	fmt.Println("hash to use:", hash)
	fmt.Println("iterations:", iterations)

	dk := pbkdf2.Key([]byte(password), tools.Base64Decode(salt), int(iterations), sha256.Size, sha256.New)
	fmt.Println(dk)
	hexDk := hex.EncodeToString(dk)
	fmt.Println("pbkdf2:", hexDk)

	// Calculate scrypt
	N := anyJson["scrypt"].(map[string]interface{})["N"].(float64)
	p := anyJson["scrypt"].(map[string]interface{})["p"].(float64)
	r := anyJson["scrypt"].(map[string]interface{})["r"].(float64)
	buflen := anyJson["scrypt"].(map[string]interface{})["buflen"].(float64)

	// fmt.Println(N, p, r, buflen)

	scryptDk, err := scrypt.Key([]byte(password), tools.Base64Decode(salt), int(N), int(r), int(p), int(buflen))
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(hex.EncodeToString(scryptDk))

	// Calculate test scrypt
	// _control: example scrypt calculated for password="rosebud", salt="pepper", N=128, p=8, n=4
	// scryptDk, err = scrypt.Key([]byte("rosebud"), []byte("pepper"), 128, 4, 8, 32)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("exampleMine:", hex.EncodeToString(scryptDk))
	// fmt.Println("example:", anyJson["scrypt"].(map[string]interface{})["_control"])

	// Build solution struct and send away
	solution := Solution{hex.EncodeToString(sha256Bytes), hmacSHA256Bytes, hexDk, hex.EncodeToString(scryptDk)}
	fmt.Println(solution)

	solutionJson, err := json.Marshal(&solution)
	if err != nil {
		fmt.Println(err)
	}
	tools.SubmitSolution(solutionJson, "https://hackattic.com//challenges/password_hashing/solve?access_token="+tools.AccessToken)

}

func shaify(password string) []byte {
	fmt.Println("hashing", fmt.Sprintf(password))
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hasher.Sum(nil)
}

func hmacyize(secret string, message string) string {
	secretByte := tools.Base64Decode(secret)
	fmt.Println("decoded salt:", string(secretByte))
	messageByte := []byte(message)

	hash := hmac.New(sha256.New, secretByte)
	hash.Write(messageByte)

	return hex.EncodeToString(hash.Sum(nil))

}
