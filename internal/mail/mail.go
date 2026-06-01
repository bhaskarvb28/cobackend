package mail

import (
	"context"
	"fmt"
	"os"
	"strings"

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
		Name: "CO Backend",

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

// ------------------------------------------------
// Invitation Email
// ------------------------------------------------

func SendInvitationEmail(
	toEmail string,
	name string,
	inviteLink string,
	role string,
) error {

	//------------------------------------------------
	// Format Role Label
	//------------------------------------------------

	roleLabel := strings.ReplaceAll(
		role,
		"_",
		" ",
	)

	roleLabel = strings.Title(
		roleLabel,
	)

	//------------------------------------------------
	// Subject
	//------------------------------------------------

	subject := fmt.Sprintf(
		"You've been invited as %s",
		roleLabel,
	)

	//------------------------------------------------
	// HTML Content
	//------------------------------------------------

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>

<head>
	<meta charset="UTF-8" />

	<title>
		Invitation
	</title>
</head>

<body style="
	margin: 0;
	padding: 0;
	background-color: #f4f4f5;
	font-family: Arial, sans-serif;
">

	<table
		width="100%%"
		cellpadding="0"
		cellspacing="0"
		style="
			padding: 40px 20px;
		"
	>
		<tr>
			<td align="center">

				<table
					width="600"
					cellpadding="0"
					cellspacing="0"
					style="
						background-color: white;
						border-radius: 12px;
						padding: 40px;
					"
				>

					<tr>
						<td>

							<h1 style="
								margin: 0;
								font-size: 28px;
								color: #111827;
							">
								You're Invited
							</h1>

							<p style="
								margin-top: 24px;
								font-size: 16px;
								line-height: 1.7;
								color: #374151;
							">
								Hi %s,
							</p>

							<p style="
								font-size: 16px;
								line-height: 1.7;
								color: #374151;
							">
								You have been invited to join the platform as
								<strong>%s</strong>.
							</p>

							<p style="
								font-size: 16px;
								line-height: 1.7;
								color: #374151;
							">
								Click the button below to accept your invitation
								and complete your account setup.
							</p>

							<div style="
								margin-top: 36px;
								margin-bottom: 36px;
								text-align: center;
							">

								<a
									href="%s"
									style="
										display: inline-block;
										background-color: #111827;
										color: white;
										text-decoration: none;
										padding: 14px 28px;
										border-radius: 8px;
										font-size: 16px;
										font-weight: 600;
									"
								>
									Accept Invitation
								</a>
							</div>

							<p style="
								font-size: 14px;
								line-height: 1.7;
								color: #6b7280;
							">
								This invitation link will expire in 24 hours.
							</p>

							<p style="
								margin-top: 32px;
								font-size: 14px;
								color: #9ca3af;
							">
								If you did not expect this invitation,
								you can safely ignore this email.
							</p>

						</td>
					</tr>

				</table>

			</td>
		</tr>
	</table>

</body>
</html>
`,
		name,
		roleLabel,
		inviteLink,
	)

	//------------------------------------------------
	// Send Email
	//------------------------------------------------

	return SendBrevoEmail(
		toEmail,
		subject,
		htmlContent,
	)
}