package trivialfiling

import (
	"fmt"

	"github.com/therin/hackattic/pkg/tools"
)

func TrivialFiling() {
	bytissimo := tools.GetProblem("https://hackattic.com/challenges/reading_qr/problem?access_token=" + tools.AccessToken)

	fmt.Println(bytissimo)

}
