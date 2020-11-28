package sms

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// SmsAero implements smsaero.ru client
type SmsAero struct {
	client    http.Client
	login     string
	apiKey    string
	signature string
}

// NewSmsAero returns new SmsAero instance
func NewSmsAero(login, apiKey, signature string) *SmsAero {
	return &SmsAero{
		client:    http.Client{Timeout: time.Second * 5},
		login:     login,
		apiKey:    apiKey,
		signature: signature,
	}
}

// Send SMS to specified number
func (sa *SmsAero) Send(number, text string) error {
	theURL, _ := url.Parse("https://gate.smsaero.ru/v2/sms/send")
	q := url.Values{}
	q.Add("number", number)
	q.Add("text", text)
	q.Add("sign", sa.signature)
	theURL.RawQuery = q.Encode()

	req, err := http.NewRequest("POST", theURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create smsaero request: %w", err)
	}

	req.SetBasicAuth(sa.login, sa.apiKey)

	resp, err := sa.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform smsaero request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("smsaero returned unexpected code: %d", resp.StatusCode)
	}

	return nil
}
