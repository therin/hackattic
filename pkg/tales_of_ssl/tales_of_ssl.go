package tales_of_ssl

import (
	"../tools"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"time"
)

type Solution struct {
	Certificate string `json:"certificate"`
}

type Problem struct {
	Private_Key   string `json:"private_key'`
	Required_Data struct {
		Domain        string `json:"domain'`
		Serial_Number string `json:"serial_number'`
		Country       string `json:"country'`
	}
}

func getKeyPairs(private string) (privKey *rsa.PrivateKey, pubKey crypto.PublicKey) {
	decoded, _ := base64.StdEncoding.DecodeString(private)
	privKey, _ = x509.ParsePKCS1PrivateKey([]byte(decoded))
	pubKey = privKey.Public()
	fmt.Println("privKey", privKey)
	fmt.Println("pubKey", pubKey)
	return
}

func TalesOfSSL() {

	bytissimo := tools.GetProblem("https://hackattic.com/challenges/tales_of_ssl/problem?access_token=" + tools.AccessToken)

	// Decode the json object
	p := &Problem{}
	err := json.Unmarshal(bytissimo, &p)
	if err != nil {
		panic(err)
	}

	serialNumber, _ := new(big.Int).SetString(p.Required_Data.Serial_Number, 0)

	// request country code from user:
	country := p.Required_Data.Country
	fmt.Printf("Enter country code for '%s': ", p.Required_Data.Country)
	fmt.Scanf("%s\n", &country)

	priv, pub := getKeyPairs(p.Private_Key)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:    []string{country},
			CommonName: p.Required_Data.Domain,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(110, 0, 0),

		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		SignatureAlgorithm:    x509.SHA256WithRSA,
		IsCA:                  false,
	}

	template.DNSNames = append(template.DNSNames, p.Required_Data.Domain)

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, pub, priv)
	if err != nil {
		fmt.Printf("Failed to create certificate: %s", err)
	}

	solution := Solution{tools.Base64Encode([]byte(derBytes))}

	solutionJson, err := json.Marshal(&solution)
	if err != nil {
		panic(err)
	}

	tools.SubmitSolution(solutionJson, "https://hackattic.com//challenges/tales_of_ssl/solve?access_token="+tools.AccessToken)

}
