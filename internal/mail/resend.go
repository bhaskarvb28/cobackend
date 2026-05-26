package mail

import (
	"os"

	"github.com/resend/resend-go/v2"
	"fmt"
)

func SendStateAdminInvitationEmail(
	toEmail string,
	inviteLink string,
) error {

	Client := resend.NewClient(
		os.Getenv("RESEND_API_KEY"),
	)
	
	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{toEmail},
		Subject: "You're invited as a State Admin",
		Html: fmt.Sprintf(`
			<h2>State Admin Invitation</h2>

			<p>You have been invited as a state admin.</p>

			<p>
				Click below to complete your account setup:
			</p>

			<a href="%s">
				Complete Setup
			</a>

			<p>
				This link expires in 24 hours.
			</p>
		`, inviteLink),
	}

	_, err := Client.Emails.Send(params)

	if err != nil {
		return err
	}

	return nil
}