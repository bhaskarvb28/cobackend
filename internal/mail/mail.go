package mail

import (
	"context"
	"fmt"
	"os"

	brevo "github.com/getbrevo/brevo-go/lib"
)

func SendBrevoEmail(
	toEmail string,
	subject string,
	htmlContent string,
) error {

	//------------------------------------------------
	// Configuration
	//------------------------------------------------

	cfg := brevo.NewConfiguration()

	cfg.AddDefaultHeader(
		"api-key",
		os.Getenv("BREVO_API_KEY"),
	)

	client := brevo.NewAPIClient(
		cfg,
	)

	//------------------------------------------------
	// Sender
	//------------------------------------------------

	sender := brevo.SendSmtpEmailSender{
		Name:  "CO Backend",
		Email: os.Getenv(
			"BREVO_SENDER_EMAIL",
		),
	}

	//------------------------------------------------
	// Recipient
	//------------------------------------------------

	to := []brevo.SendSmtpEmailTo{
		{
			Email: toEmail,
		},
	}

	//------------------------------------------------
	// Email Payload
	//------------------------------------------------

	email := brevo.SendSmtpEmail{
		Sender:      &sender,
		To:          to,
		Subject:     subject,
		HtmlContent: htmlContent,
	}

	//------------------------------------------------
	// Send Email
	//------------------------------------------------

	_, httpResponse, err := client.
		TransactionalEmailsApi.
		SendTransacEmail(
			context.Background(),
			email,
		)

	if err != nil {

		if httpResponse != nil {

			fmt.Println(
				"Brevo status:",
				httpResponse.Status,
			)
		}

		fmt.Println(
			"Brevo error:",
			err,
		)

		return err
	}

	return nil
}

func SendInvitationEmail(
	toEmail string,
	inviteLink string,
	role string,
) error {

	subject := fmt.Sprintf(
		"You're invited as %s",
		role,
	)

	htmlContent := fmt.Sprintf(`
		<h2>Invitation</h2>

		<p>
			You have been invited as
			<strong>%s</strong>.
		</p>

		<p>
			Click below to complete
			your account setup:
		</p>

		<p>
			<a href="%s">
				Complete Setup
			</a>
		</p>

		<p>
			This link expires in
			24 hours.
		</p>
	`,
		role,
		inviteLink,
	)

	return SendBrevoEmail(
		toEmail,
		subject,
		htmlContent,
	)
}