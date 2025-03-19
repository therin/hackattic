package dns

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/therin/hackattic/pkg/tools"
)

// ----------------------------------------------------------------------
// Required Environment Variables:
//  export DNS_SERVER_ADDRESS=127.0.0.1
//  export DNS_SERVER_PORT=53
//  export HACKATTIC_ACCESS_TOKEN=xxxxxxx
// ----------------------------------------------------------------------

var records = make(map[string]map[uint16]string)

type dnsHandler struct{}

type Record struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Data string `json:"data"`
}

type DNSChallenge struct {
	Records []Record `json:"records"`
}

type Reply struct {
	DnsIp   string `json:"dns_ip"`
	DnsPort string `json:"dns_port"`
}

// buildRR creates a DNS resource record based on the query name, type, and data.
func buildRR(query string, qtype uint16, data string) (dns.RR, error) {
	switch qtype {
	case dns.TypeA:
		return dns.NewRR(fmt.Sprintf("%s A %s", query, data))
	case dns.TypeAAAA:
		return dns.NewRR(fmt.Sprintf("%s AAAA %s", query, data))
	case dns.TypeTXT:
		return dns.NewRR(fmt.Sprintf("%s TXT \"%s\"", query, data))
	case dns.TypeRP:
		return dns.NewRR(fmt.Sprintf("%s RP %s .", query, data))
	default:
		return nil, fmt.Errorf("unsupported query type %d", qtype)
	}
}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true
	msg.Compress = false

	for _, question := range r.Question {
		normalizedQuery := dns.Fqdn(question.Name)
		log.Printf("Received query for %s of type %d\n", normalizedQuery, question.Qtype)

		// Check for an exact match first.
		if recMap, found := records[normalizedQuery]; found {
			if data, ok := recMap[question.Qtype]; ok {
				rr, err := buildRR(normalizedQuery, question.Qtype, data)
				if err != nil {
					log.Printf("Failed to create DNS record for %s: %v", normalizedQuery, err)
					msg.Rcode = dns.RcodeServerFailure
					continue
				}
				msg.Answer = append(msg.Answer, rr)
				continue
			}
		}

		// If no exact match is found, perform wildcard matching.
		foundWildcard := false
		for key, recMap := range records {
			if strings.HasPrefix(key, "*.") {
				// Remove the "*" so that "*.domain.com." becomes ".domain.com."
				suffix := key[1:]
				if strings.HasSuffix(normalizedQuery, suffix) {
					if data, ok := recMap[question.Qtype]; ok {
						rr, err := buildRR(normalizedQuery, question.Qtype, data)
						if err != nil {
							log.Printf("Failed to create DNS record for %s: %v", normalizedQuery, err)
							msg.Rcode = dns.RcodeServerFailure
							continue
						}
						msg.Answer = append(msg.Answer, rr)
						foundWildcard = true
						break
					}
				}
			}
		}
		if !foundWildcard && len(msg.Answer) == 0 {
			log.Printf("No record found for %s", normalizedQuery)
			msg.Rcode = dns.RcodeNameError
		}
	}

	w.WriteMsg(msg)
}

func Serve() {

	// Load configuration from environment variables.
	dnsServerAddress := os.Getenv("DNS_SERVER_ADDRESS")
	dnsServerPort := os.Getenv("DNS_SERVER_PORT")
	hackatticAccessToken := os.Getenv("HACKATTIC_ACCESS_TOKEN")

	if dnsServerAddress == "" {
		log.Fatal("DNS_SERVER_ADDRESS environment variable not set")
	}
	if dnsServerPort == "" {
		log.Fatal("DNS_SERVER_PORT environment variable not set")
	}
	if hackatticAccessToken == "" {
		log.Fatal("HACKATTIC_ACCESS_TOKEN environment variable not set")
	}

	handler := new(dnsHandler)
	server := &dns.Server{
		Addr:      fmt.Sprintf(":%s", dnsServerPort),
		Net:       "udp",
		Handler:   handler,
		UDPSize:   dns.DefaultMsgSize,
		ReusePort: true,
	}

	fmt.Println("Starting DNS server in the background...")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start DNS server: %s\n", err.Error())
		}
	}()

	// Wait briefly to ensure the server is up before proceeding.
	time.Sleep(100 * time.Millisecond)

	problemURL := "https://hackattic.com/challenges/serving_dns/problem?access_token=" + hackatticAccessToken
	bytissimo := tools.GetProblem(problemURL)
	fmt.Println(string(bytissimo))

	fmt.Println(string(bytissimo))

	// Unmarshal the JSON into the custom Problem struct
	var challenge DNSChallenge
	err := json.Unmarshal(bytissimo, &challenge)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Further manipulation: for example, print out each record
	for _, record := range challenge.Records {
		fmt.Printf("Name: %s, Type: %s, Data: %s\n", record.Name, record.Type, record.Data)
	}

	// Load the fetched records into the global records map
	for _, record := range challenge.Records {
		normalizedName := dns.Fqdn(record.Name)
		log.Printf("Processing record: Name: %s, Type: %s, Data: %s", normalizedName, record.Type, record.Data)
		typeCode, ok := dns.StringToType[record.Type]
		if !ok {
			log.Printf("Unknown record type: %s", record.Type)
			continue
		}
		if records[normalizedName] == nil {
			records[normalizedName] = make(map[uint16]string)
		}
		records[normalizedName][typeCode] = record.Data
	}

	reply := Reply{
		DnsIp:   dnsServerAddress,
		DnsPort: dnsServerPort,
	}

	replyData, err := json.Marshal(reply)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	submitURL := "https://hackattic.com/challenges/serving_dns/solve?access_token=" + hackatticAccessToken
	tools.SubmitSolution(replyData, submitURL)
}
