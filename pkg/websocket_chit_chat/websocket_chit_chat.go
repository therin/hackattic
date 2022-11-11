package websocket_chit_chat

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/therin/hackattic/pkg/tools"
)

type Problem struct {
	Token string `json:"token"`
}

type Solution struct {
	Secret string `json:"secret"`
}

func FindBestDuration(number int64) int64 {
	intervals := [5]int64{700, 1500, 2000, 2500, 3000}
	floatNumber := float64(number)

	minimumDistance := math.Inf(1)

	var answer int64

	for _, interval := range intervals {
		floatInterval := float64(interval)

		curDistance := math.Abs(floatInterval - floatNumber)
		if curDistance < minimumDistance {
			minimumDistance = curDistance
			answer = interval
		}

	}

	fmt.Printf("Found matching duration: %d \n", answer)
	return answer

}

func TimeCounter(stopChan chan struct{}, timeValueChan chan time.Duration) {

	fmt.Println("Starting counter")
	startTime := time.Now()

	for {
		select {
		case <-stopChan:
			endTime := time.Now()
			timeElapsed := endTime.Sub(startTime)
			timeValueChan <- timeElapsed
			return
		}
	}

}

func WebSocketChat() {

	done := make(chan struct{})
	stopChan := make(chan struct{})
	timeValueChan := make(chan time.Duration)
	sendMessageChan := make(chan int64)

	bytissimo := tools.GetProblem("https://hackattic.com/challenges/websocket_chit_chat/problem?access_token=" + tools.AccessToken)

	problemJson := Problem{}

	json.Unmarshal(bytissimo, &problemJson)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: "hackattic.com", Path: "/_/ws/" + problemJson.Token}

	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	// Capture time when connection is estableshed
	fmt.Println("Starting timer")
	go TimeCounter(stopChan, timeValueChan)

	if err != nil {
		log.Fatal("Dialing failed: ", err)
	}

	defer c.Close()

	go func() {
		defer close(done)

		for {
			_, message, err := c.ReadMessage()

			if err != nil {
				log.Fatal("Error receiving message: ", err)
			}

			log.Printf("Parsing: %s", message)
			m := string(message)

			if strings.Contains(m, "hello!") {
				log.Printf("Found hello in: %s", m)

			} else if strings.Contains(m, "ping!") {
				log.Printf("Found ping in: %s, stopping counter", m)
				close(stopChan)
				timeTaken := <-timeValueChan
				close(timeValueChan)
				log.Printf("Recorded time: %d", int64(timeTaken/time.Millisecond))

				// start counter again
				fmt.Println("Restarting counter")
				stopChan = make(chan struct{})
				timeValueChan = make(chan time.Duration)
				go TimeCounter(stopChan, timeValueChan)

				bestDuration := FindBestDuration(int64(timeTaken / time.Millisecond))
				sendMessageChan <- bestDuration

			} else if strings.Contains(m, "congratulations!") {

				log.Printf("Found secret in: %s", m)

				var rgx = regexp.MustCompile(`\"(.*?)\"`)
				answer := rgx.FindStringSubmatch(m)[1]

				log.Printf("Preparing answer: %s", answer)
				solution := Solution{answer}
				solutionJson, err := json.Marshal(solution)
				if err != nil {
					log.Println("Cannot marshal solution struct to JSON", err)
				}
				tools.SubmitSolution([]byte(solutionJson), "https://hackattic.com/challenges/websocket_chit_chat/solve?access_token="+tools.AccessToken)
				close(done)

			} else {
				log.Println("No action required")
			}

		}

	}()

	for {
		select {
		case <-done:
			return
		case calculatedDuration := <-sendMessageChan:
			fmt.Printf("Sending: %s \n", strconv.Itoa(int(calculatedDuration)))
			err := c.WriteMessage(1, []byte(strconv.Itoa(int(calculatedDuration))))
			if err != nil {
				log.Println("Cannot send duration message:", err)
			}
		case <-interrupt:
			log.Println("Interrupt signal received")

			// close the connection cleanly by sending a message and waiting for a server to close it
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Closed writing channel:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Millisecond):
			}
			return
		}

	}
}
