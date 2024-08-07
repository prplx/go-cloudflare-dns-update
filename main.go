package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"encoding/json"

	"github.com/joho/godotenv"
	"github.com/rdegges/go-ipify"
	"golang.org/x/sync/errgroup"
)

func main() {
	var ip string

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	recordIds := strings.Split(os.Getenv("CF_DNS_RECORD_ID"), ",")

	for {
		time.Sleep(time.Second * 1)

		currentIP, err := getCurrentIpWithIpify()
		if err != nil {
			log.Printf("Failed to get current IP with ipify, trying to get from freecurrentip: %s\n", err)

			currentIP, err = getCurrentIpWithFreeIPAPI()
			if err != nil {
				log.Printf("Failed to get current IP with freecurrentip: %s\n", err)

				continue
			}
		}

		if currentIP == ip {
			continue
		}

		eg := errgroup.Group{}
		for _, recordId := range recordIds {
			func(recId string) {
				eg.Go(func() error {
					return updateDnsRecord(recId, currentIP)
				})
			}(recordId)
		}

		if err := eg.Wait(); err != nil {
			log.Printf("Failed to update DNS record: %s\n", err)
		} else {
			fmt.Printf("DNS record has been updated to: %s\n", currentIP)
			ip = currentIP
		}

	}
}

func getCurrentIpWithIpify() (string, error) {
	ip, err := ipify.GetIp()
	if err != nil {
		return "", err
	}
	return ip, nil
}

func getCurrentIpWithFreeIPAPI() (string, error) {
	type Response struct {
		IP string `json:"ipAddress"`
	}
	resp, err := http.Get("https://freeipapi.com/api/json")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.IP, nil
}

func updateDnsRecord(recordId, ip string) error {
	apiURL := os.Getenv("CF_API_URL")
	apiToken := os.Getenv("CF_API_TOKEN")
	zoneId := os.Getenv("CF_ZONE_ID")

	resp := struct {
		Success bool `json:"success"`
	}{}
	jsonBody := []byte(fmt.Sprintf(`{"content": "%s"}`, ip))
	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/zones/%s/dns_records/%s", apiURL, zoneId, recordId), bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code received %d", res.StatusCode)
	}

	if res.Body != nil {
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &resp)
		if err != nil {
			return err
		}

		if !resp.Success {
			return fmt.Errorf("result not successful")
		}
	}

	return nil
}
