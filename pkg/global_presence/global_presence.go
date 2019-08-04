package global_presence

import (
	"../../pkg/tools"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Problem struct {
	PresenceToken string `json:"presence_token"`
}

func LoadProxies(proxyFile string) []string {
	/*
	 Expect a file with a list of proxies:
	 127.0.0.0:80
	 127.0.0.0:80
	*/

	proxyList := make([]string, 0)

	fmt.Printf("Loading proxies from %s \n", proxyFile)

	file, err := os.Open(proxyFile)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		proxyList = append(proxyList, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return proxyList
}

func CallEndpoint(requestUrl string, proxy string) string {

	fmt.Printf("Calling %s with proxy %s \n", requestUrl, proxy)

	url, err := url.Parse(requestUrl)
	if err != nil {
		log.Println(err)
	}

	var client *http.Client

	if proxy != "" {
		proxyURL, err := url.Parse("http://" + proxy)
		if err != nil {
			log.Println(err)
		}

		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}

		client = &http.Client{
			Transport: transport,
		}

	} else {
		client = &http.Client{}

	}

	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Println(err)
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return ""
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return ""
	}

	log.Println(string(data))

	defer response.Body.Close()

	return string(data)
}

func GlobalPresence() {

	bytissimo := tools.GetProblem("https://hackattic.com/challenges/a_global_presence/problem?access_token=" + tools.AccessToken)

	problemJson := Problem{}
	json.Unmarshal(bytissimo, &problemJson)

	fmt.Printf(problemJson.PresenceToken)

	requestUrl := "https://hackattic.com/_/presence/" + problemJson.PresenceToken
	proxyList := LoadProxies("../../pkg/global_presence/proxy.txt")

	// Call endpoint via all loaded proxies
	for _, proxy := range proxyList {
		go CallEndpoint(requestUrl, proxy)
		time.Sleep(2 * time.Second)

	}

	// Keep pinging endpoint until we reach 7 countries
	for {
		response := CallEndpoint(requestUrl, "")
		if len(response) != 0 {
			fmt.Println(response)
			responseLength := len(strings.Split(response, ","))
			fmt.Println(responseLength)
			if responseLength >= 7 {
				break
			}

		}

	}
	// Submit empty json when all countries are pinged
	tools.SubmitSolution([]byte("{}"), "https://hackattic.com/challenges/a_global_presence/solve?access_token="+tools.AccessToken)

}
