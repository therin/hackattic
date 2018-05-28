package jotting_jwts

import (
	"../tools"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Solution struct {
	App_url string `json:"app_url"`
}

type FinalSolution struct {
	Solution string `json:"solution"`
}

func verifyToken(reqToken []byte, mySigningKey []byte) (bool, string) {

	var result string

	fmt.Println("Verifying against: ", string(reqToken), string(mySigningKey))
	token, err := jwt.Parse(string(reqToken), func(t *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err == nil && token.Valid {
		fmt.Println("Valid token")

		if claims, ok := token.Claims.(jwt.MapClaims); ok && claims["append"] != nil {
			result = claims["append"].(string)
			return true, result

		} else {
			return false, "Time to finish"
		}

		return false, result
	} else {
		fmt.Println("Invalid token")
		return false, ""
	}

}

func webServe(mySigningKey []byte) {

	var answer strings.Builder

	var AddFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		reqToken, _ := ioutil.ReadAll(r.Body)

		fmt.Println(string(reqToken))

		valid, stringToAppend := verifyToken(reqToken, mySigningKey)

		if valid == true {
			answer.WriteString(stringToAppend)
			fmt.Println("Current string: ", answer.String())
		} else {
			if stringToAppend == "Time to finish" {
				finalSolution := FinalSolution{answer.String()}
				finalSolutionJson, err := json.Marshal(&finalSolution)
				if err != nil {
					panic(err)
				}
				fmt.Println(finalSolution, finalSolutionJson)
				fmt.Println(string(finalSolutionJson))

				w.Header().Set("Content-Type", "application/json")
				w.Write(finalSolutionJson)
			}

		}

		body, _ := ioutil.ReadAll(r.Body)
		fmt.Println(string(body))
	})

	r := mux.NewRouter()
	r.Handle("/myApp", AddFeedbackHandler).Methods("POST")

	http.ListenAndServe(":3000", handlers.CombinedLoggingHandler(os.Stdout, r))

}

func Jotting_jwts() {

	// get JWT token
	bytissimo := tools.GetProblem("https://hackattic.com/challenges/jotting_jwts/problem?access_token=" + tools.AccessToken)
	var anyJson map[string]interface{}
	json.Unmarshal(bytissimo, &anyJson)

	var mySigningKey = []byte(anyJson["jwt_secret"].(string))

	// run app to serve requests
	go webServe(mySigningKey)

	// post app url to solution endpoint
	solution := Solution{"http://127.0.0.1:3000/myApp"}
	solutionJson, err := json.Marshal(&solution)
	if err != nil {
		panic(err)
	}
	tools.SubmitSolution(solutionJson, "https://hackattic.com/challenges/jotting_jwts/solve?access_token="+tools.AccessToken)

}
