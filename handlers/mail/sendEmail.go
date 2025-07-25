package mail

import (
	"fmt"
	"strings"

	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
)

func SendGenericEmail(c *fiber.Ctx, mailService *models.EmailService, sosUser models.User) {
	user := sosUser.UserDetails()
	fIndex := strings.IndexByte(user.Name, ' ')
	var alertMsg string
	switch user.Role {
	case "student":
		alertMsg = fmt.Sprintf(
			"Student %s of the college of %s residing at %s is requesting immediate assistance from %s. Contact %s at %s",
			user.Name, *user.College, *user.Hostel, "loc", user.Name[:fIndex], user.PhoneNumber,
		)
	case "TA":
		alertMsg = fmt.Sprintf(
			"Teaching Assistant %s from the college of %s, is requesting immediate assistance from %s. Contact %s at %s.",
			user.Name, *user.College, "loc", user.Name[:fIndex], user.PhoneNumber,
		)
	case "lecturer":
		alertMsg = fmt.Sprintf(
			"Lecturer %s from the college of %s, is requesting immediate assistance from %s. Contact %s at %s.",
			user.Name, *user.College, "loc", user.Name[:fIndex], user.PhoneNumber,
		)
	default:
		alertMsg = fmt.Sprintf(
			"A member of staff by name %s is requesting immediate assistance from %s. Contact %s at %s.",
			user.Name, "loc", user.Name[:fIndex], user.PhoneNumber,
		)
	}

	mailService.SendEmail(
		"someEmailAddress",
		"SOS Distress call",
		fmt.Sprintf(
			`
			<!DOCTYPE html>
			<html>
			<head>
				<meta charset="UTF-8">
				<title>SOS Distress Call</title>
				<style>
					body { font-family: Arial, sans-serif; background: #f8d7da; color: #721c24; margin: 0; padding: 0; }
					.container { max-width: 600px; margin: 40px auto; background: #fff; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); padding: 32px; }
					.header { background: #f5c6cb; padding: 16px; border-radius: 8px 8px 0 0; text-align: center; }
					.header h1 { margin: 0; color: #721c24; }
					.content { margin-top: 24px; }
					.button {
						display: inline-block;
						padding: 12px 24px;
						background: #c82333;
						color: #fff;
						text-decoration: none;
						border-radius: 4px;
						font-weight: bold;
						margin-top: 24px;
					}
					.footer { margin-top: 32px; font-size: 0.9em; color: #856404; text-align: center; }
				</style>
			</head>
			<body>
				<div class="container">
					<div class="header">
						<h1>🚨 SOS Distress Alert 🚨</h1>
					</div>
					<div class="content">
						<p>Alert! Alert!</p>
						<p><strong>This is an urgent distress notification.</strong></p>
						<p>
							%s
						</p>
						<a href="tel:1234567890" class="button">Call for Help</a>
					</div>
					<div class="footer">
						This message was sent automatically by the Rapport SOS system.<br>
						If you believe this was sent in error, please disregard.
					</div>
				</div>
			</body>
			</html>
			`, alertMsg,
		),
	)
}
