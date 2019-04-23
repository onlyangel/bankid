package main

import (
	"log"
	"os"
	"time"

	"github.com/jfernstad/bankid"
)

func main() {

	caTestPath := "../CA/test.crt"
	rpCrtPath := "../rp/bankid_rp_test.crt" // NOTE: Replace with your RP (Relaying Partner) certificate
	rpKeyPath := "../rp/bankid_rp_test.key" // NOTE: Replace with your RP key

	// bankid.TestBaseURL or bankid.ProductionBaseURL
	env, err := bankid.NewEnvironment(bankid.TestBaseURL, caTestPath, rpCrtPath, rpKeyPath)
	if err != nil {
		log.Printf(" !! Could not create TestEnvironment: %s", err.Error())
		os.Exit(1)
	}

	// Remove non-digits from personal number
	personalNumber := "198001010109" // NOTE: Replace with a real personal number
	ipAddr := "127.0.0.1"            // IP of your mobile phone with BankID app on it

	rsp, err := bankid.Auth(env, personalNumber, ipAddr)
	if err != nil {
		log.Printf(" !! Could not connect to server: %s\n", err.Error())
		os.Exit(1)
	}

	// Auth started!
	// Pull out your BankID app and sign the Auth request
	collectResponse := &bankid.CollectResponse{}
	done := false
	for !done {
		collectResponse, err = bankid.Collect(env, rsp.OrderRef)
		if err != nil {
			log.Printf(" !! Could not collect: %s\n", err.Error())
			os.Exit(1)
		}

		switch collectResponse.Status {
		case bankid.OrderPending:
			switch collectResponse.HintCode {
			case bankid.PendOutstandingTransaction:
				{
					log.Println(" >> Auth has begun")
				}
			case bankid.PendStarted:
				{
					log.Println(" >> App started")
				}
			case bankid.PendUserSign:
				{
					log.Println(" >> User signing")
				}
			}
		case bankid.OrderFailed:
			{
				done = true
				switch collectResponse.HintCode {
				case bankid.FailCancelled:
					{
						log.Println(" !! Auth got cancelled")
						break
					}
				case bankid.FailUserCancel:
					{
						log.Println(" !! User cancelled auth")
						break
					}
				case bankid.FailExpiredTransaction:
					{
						log.Println(" !! Auth took too long")
						break
					}
				}
			}
		case bankid.OrderComplete:
			{
				done = true
				log.Println(" >> 😎 Auth Complete ")
				log.Printf("%s signed in!\n", collectResponse.CompletionData.User.Name)
				break
			}
		}
		// Don't spam the service plz
		time.Sleep(2 * time.Second)
	}

	if collectResponse.Status != bankid.OrderComplete {
		_, err = bankid.Cancel(env, rsp.OrderRef)
		if err != nil {
			log.Printf(" !! Could not cancel request: %s\n", err.Error())
		}
		log.Printf(" >> Auth cancelled\n")
	}
}
