package intents

import "fmt"

// MockSender is a mock for Sender interface
type MockSender struct{}

// Send just outputs message to stdout
func (*MockSender) Send(text, phoneNumber string) error {
	if text == "error" {
		return fmt.Errorf("test error, phone number %s", phoneNumber)
	}

	fmt.Printf("Sending to %s: %s", phoneNumber, text)

	return nil
}
