package reading_qr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"

	"github.com/kdar/goquirc"
	"github.com/therin/hackattic/pkg/tools"
)

type Solution struct {
	Code string `json:"code"`
}

func ReadingQR() {

	bytissimo := tools.GetProblem("https://hackattic.com/challenges/reading_qr/problem?access_token=" + tools.AccessToken)

	// convert to JSON string interface
	var anyJson map[string]interface{}
	json.Unmarshal(bytissimo, &anyJson)

	imageURL := anyJson["image_url"].(string)

	fmt.Println(imageURL)

	// download image
	workingDir := "../../pkg/reading_qr/work"

	err := tools.DownloadFile(workingDir+"/image.jpg", imageURL)
	if err != nil {
		fmt.Println(err)
	}

	// load as image
	imgdata, err := ioutil.ReadFile(workingDir + "/image.jpg")
	if err != nil {
		fmt.Println(err)
	}

	m, _, err := image.Decode(bytes.NewReader(imgdata))
	if err != nil {
		fmt.Println(err)
	}

	// decode
	d := goquirc.New()
	defer d.Destroy()
	datas, err := d.Decode(m)
	if err != nil {
		fmt.Println(err)
	}
	// load data from image
	for _, data := range datas {
		fmt.Printf("%s\n", data.Payload[:data.PayloadLen])
		answer := data.Payload[:data.PayloadLen]

		// move response to solution struct and submit
		solution := Solution{string(answer)}
		fmt.Printf("%+v\n", solution)

		solutionJson, err := json.Marshal(&solution)
		if err != nil {
			fmt.Println(err)
		}
		tools.SubmitSolution(solutionJson, "https://hackattic.com/challenges/reading_qr/solve?access_token="+tools.AccessToken)

	}

}
