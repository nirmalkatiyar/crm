package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"

	"github.com/nirmal/crm/models"
)


func SendInteractionNotificationWithEmail(interaction models.Interaction, emailTo, meetingStartTime string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_MAIL")
	password := os.Getenv("SMTP_PASSWORD")

	//	example@example.com		EXAMPLE_PASSWORD	smtp.example.com
	auth := smtp.PlainAuth(
		"",
		from,
		password,
		host,
	)

	subject := fmt.Sprintf("Meeting Notification: %s", *interaction.Title)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
        }
        .container {
            padding: 20px;
            border: 1px solid #ddd;
            border-radius: 8px;
            background-color: #f9f9f9;
        }
        .header {
            font-size: 18px;
            font-weight: bold;
            color: #333;
            margin-bottom: 10px;
        }
        .content {
            font-size: 16px;
            color: #555;
        }
        .footer {
            font-size: 14px;
            color: #888;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">Meeting Notification</div>
        <div class="content">
            <p>Dear User,</p>
            <p>You have a scheduled meeting with the following details:</p>
            <p><strong>Interaction ID:</strong> %s</p>
            <p><strong>Title:</strong> %s</p>
            <p><strong>Description:</strong> %s</p>
            <p><strong>Start Time:</strong> %s</p>
            <p>Please ensure you are prepared for the meeting.</p>
        </div>
        <div class="footer">
            <p>Thank you,</p>
            <p>Support Team</p>
        </div>
    </div>
</body>
</html>
`, interaction.CustomerID.Hex(), *interaction.Title, *interaction.Description, meetingStartTime)

	// subject := "Ticket Created: " + *interaction.Title
	// body := fmt.Sprintf("Dear User,\n\nYour Interaction with ID %s has been created.\n\nDetails:\nDescription: %s\n\nThank you,\nSupport Team", interaction.CustomerID, *interaction.Title, *interaction.Description)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	client, err := smtp.Dial(host + ":" + port)
	if err != nil {
		fmt.Println("error while connecting to smtp...")
		return err
	}

	client.StartTLS(tlsConfig)

	if err := client.Auth(auth); err != nil {
		fmt.Println("error while authenticating to SMTP...")
		return err
	}

	if err := client.Mail(from); err != nil {
		fmt.Println("error while initiating mail transaction...")
		return err
	}

	if err := client.Rcpt(emailTo); err != nil {
		fmt.Println("error while issuing RCPT(CMD) to SMTP...")
		return err
	}

	w, err := client.Data()
	if err != nil {
		fmt.Println("error while issuing DATA(CMD) to SMTP...")
		return err
	}

	_, err = w.Write([]byte(
		fmt.Sprintf("MIME-Version: %v\r\n", "1.0") +
			fmt.Sprintf("Content-type: %v\r\n", "text/html; charset=UTF-8") +
			fmt.Sprintf("From: %v\r\n", from) +
			fmt.Sprintf("To: %v\r\n", emailTo) +
			fmt.Sprintf("Subject: %v\r\n", subject) +
			fmt.Sprintf("%v\r\n", body),
	))

	if err != nil {
		fmt.Println("error while writing data to the header/body")
		return err
	}

	err = w.Close()
	if err != nil {
		fmt.Println("error closing writer")
		return err
	}

	client.Quit()
	return nil
}
