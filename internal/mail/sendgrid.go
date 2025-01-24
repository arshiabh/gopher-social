package mail

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGrid struct {
	client    *sendgrid.Client
	apiKey    string
	fromEmail string
}

func NewSendGrip(apikey string, fromEmail string) *SendGrid {
	client := sendgrid.NewSendClient(apikey)
	return &SendGrid{
		client:    client,
		apiKey:    apikey,
		fromEmail: fromEmail,
	}
}

func (m *SendGrid) Send(username string, userEmail string) error {
	from := mail.NewEmail("GopherSocial", m.fromEmail)
	subject := "Finish User Activations"
	to := mail.NewEmail(username, userEmail)

	content := contentFromTxt()
	message := mail.NewSingleEmail(from, subject, to, string(content), "")
	for i := 0; i < 3; i++ {
		_, err := m.client.Send(message)
		if err != nil {
			log.Printf("failed to send email")
			continue
		}
		log.Printf("email sent successfully")
		return nil
	}
	return fmt.Errorf("failed to send email after 3 times attempt")
}

func contentFromTxt() []byte {
	filePath := "template/content.txt"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()
	// Read the file contents
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return nil
	}

	fileContent := make([]byte, fileInfo.Size())
	_, err = file.Read(fileContent)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	fmt.Println(string(fileContent))
	return fileContent
}
