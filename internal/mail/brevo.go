package mail

import (
	"context"
	"fmt"
	"os"

	brevo "github.com/getbrevo/brevo-go/lib"
)

func SendStateAdminInvitationEmailBrevo(
	toEmail string,
	inviteLink string,
) error {

	cfg := brevo.NewConfiguration()

	cfg.AddDefaultHeader(
		"api-key",
		os.Getenv("BREVO_API_KEY"),
	)

	client := brevo.NewAPIClient(cfg)

	htmlContent := fmt.Sprintf(`
		<h2>State Admin Invitation</h2>

		<p>
			You have been invited as a state admin.
		</p>

		<p>
			Click below to complete your account setup:
		</p>

		<p>
			<a href="%s">
				Complete Setup
			</a>
		</p>

		<p>
			This link expires in 24 hours.
		</p>
	`, inviteLink)

	sender := brevo.SendSmtpEmailSender{
		Name:  "CO Backend",
		Email: os.Getenv("BREVO_SENDER_EMAIL"),
	}

	to := []brevo.SendSmtpEmailTo{
		{
			Email: toEmail,
		},
	}

	email := brevo.SendSmtpEmail{
		Sender:      &sender,
		To:          to,
		Subject:     "You're invited as a State Admin",
		HtmlContent: htmlContent,
	}

	_, httpResponse, err := client.TransactionalEmailsApi.
		SendTransacEmail(
			context.Background(),
			email,
		)

	if err != nil {

		if httpResponse != nil {
			fmt.Println("Brevo status:", httpResponse.Status)
		}

		fmt.Println("Brevo error:", err)

		return err
	}

	return nil
}

func SendDistrictAdminInvitationEmailBrevo(
	toEmail string,
	inviteLink string,
) error {

	cfg := brevo.NewConfiguration()

	cfg.AddDefaultHeader(
		"api-key",
		os.Getenv("BREVO_API_KEY"),
	)

	client := brevo.NewAPIClient(cfg)

	htmlContent := fmt.Sprintf(`
		<h2>District Admin Invitation</h2>

		<p>
			You have been invited as a district admin.
		</p>

		<p>
			Click below to complete your account setup:
		</p>

		<p>
			<a href="%s">
				Complete Setup
			</a>
		</p>

		<p>
			This link expires in 24 hours.
		</p>
	`, inviteLink)

	sender := brevo.SendSmtpEmailSender{
		Name:  "CO Backend",
		Email: os.Getenv("BREVO_SENDER_EMAIL"),
	}

	to := []brevo.SendSmtpEmailTo{
		{
			Email: toEmail,
		},
	}

	email := brevo.SendSmtpEmail{
		Sender:      &sender,
		To:          to,
		Subject:     "You're invited as a District Admin",
		HtmlContent: htmlContent,
	}

	_, httpResponse, err := client.TransactionalEmailsApi.
		SendTransacEmail(
			context.Background(),
			email,
		)

	if err != nil {

		if httpResponse != nil {
			fmt.Println("Brevo status:", httpResponse.Status)
		}

		fmt.Println("Brevo error:", err)

		return err
	}

	return nil
}