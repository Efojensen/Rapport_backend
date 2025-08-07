package mail

import (
	"fmt"
	"strings"

	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
)

// SendSOSEmail sends an SOS email with optional location data
func SendSOSEmail(c *fiber.Ctx, mailService *models.EmailService, sosUser models.User, sosReport *models.SOSReport) error {
	user := sosUser.UserDetails()
	fIndex := strings.IndexByte(user.Name, ' ')
	var alertMsg string
	var locationInfo string

	// Get location information if available
	if sosReport != nil {
		location, err := sosReport.GetLatLongAddress()
		if err == nil && location != nil {
			locationInfo = fmt.Sprintf("<br><strong>📍Location Details</strong><br>    Address: %s<br>    City: %s<br>    Suburb: %s<br>    Country: %s<br>    Coordinates: %.6f, %.6f",
				location.Results[0].Properties.Address, location.Results[0].Properties.City, location.Results[0].Properties.Suburb, location.Results[0].Properties.Country,
				sosReport.GeoLocation.Latitude, sosReport.GeoLocation.Longitude)
		} else {
			// Fallback to just coordinates if reverse geocoding fails
			locationInfo = fmt.Sprintf("<br><strong>📍Location Coordinates:** %.6f, %.6f</strong>",
				sosReport.GeoLocation.Latitude, sosReport.GeoLocation.Longitude)
		}
	}

	// Get sender's email from the user data
	var senderEmail string
	switch v := sosUser.(type) {
	case *models.Student:
		senderEmail = v.Email
	case *models.TeachAsst:
		senderEmail = v.Email
	case *models.Lecturer:
		senderEmail = v.Email
	case *models.Other:
		senderEmail = v.Email
	default:
		senderEmail = "unknown@rapport.edu" // fallback
	}

	switch user.Role {
	case "student":
		collegeInfo := "Unknown College"
		hostelInfo := "Unknown Hostel"
		if user.College != nil {
			collegeInfo = *user.College
		}
		if user.Hostel != nil {
			hostelInfo = *user.Hostel
		}
		alertMsg = fmt.Sprintf(
			"Student %s of the college of %s residing at %s is requesting immediate assistance. Contact %s at %s%s",
			user.Name, collegeInfo, hostelInfo, user.Name[:fIndex], user.PhoneNumber, locationInfo,
		)
	case "TA":
		collegeInfo := "Unknown College"
		if user.College != nil {
			collegeInfo = *user.College
		}
		alertMsg = fmt.Sprintf(
			"Teaching Assistant %s from the college of %s is requesting immediate assistance. Contact %s at %s.%s",
			user.Name, collegeInfo, user.Name[:fIndex], user.PhoneNumber, locationInfo,
		)
	case "lecturer":
		collegeInfo := "Unknown College"
		if user.College != nil {
			collegeInfo = *user.College
		}
		alertMsg = fmt.Sprintf(
			"Lecturer %s from the college of %s is requesting immediate assistance. Contact %s at %s.%s",
			user.Name, collegeInfo, user.Name[:fIndex], user.PhoneNumber, locationInfo,
		)
	default:
		alertMsg = fmt.Sprintf(
			"A member of staff by name %s is requesting immediate assistance. Contact %s at %s.%s",
			user.Name, user.Name[:fIndex], user.PhoneNumber, locationInfo,
		)
	}

	// Always send SOS emails to rapportSafety@gmail.com
	err := mailService.SendEmail(
		senderEmail,
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
					.location-box { background: #e2f3ff; border: 1px solid #b3d9ff; border-radius: 6px; padding: 16px; margin: 16px 0; }
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

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"msg": "SOS email sent successfully with location data",
	})
}

// SendGenericEmail sends a generic SOS email without location data (kept for backward compatibility)
func SendGenericEmail(c *fiber.Ctx, mailService *models.EmailService, sosUser models.User) error {
	user := sosUser.UserDetails()
	fIndex := strings.IndexByte(user.Name, ' ')
	var alertMsg string

	// Get sender's email from the user data
	var senderEmail string
	switch v := sosUser.(type) {
	case *models.Student:
		senderEmail = v.Email
	case *models.TeachAsst:
		senderEmail = v.Email
	case *models.Lecturer:
		senderEmail = v.Email
	case *models.Other:
		senderEmail = v.Email
	default:
		senderEmail = "unknown@rapport.edu" // fallback
	}

	switch user.Role {
	case "student":
		collegeInfo := "Unknown College"
		hostelInfo := "Unknown Hostel"
		if user.College != nil {
			collegeInfo = *user.College
		}
		if user.Hostel != nil {
			hostelInfo = *user.Hostel
		}
		alertMsg = fmt.Sprintf(
			"Student %s of the college of %s residing at %s is requesting immediate assistance. Contact %s at %s",
			user.Name, collegeInfo, hostelInfo, user.Name[:fIndex], user.PhoneNumber,
		)
	case "TA":
		collegeInfo := "Unknown College"
		if user.College != nil {
			collegeInfo = *user.College
		}
		alertMsg = fmt.Sprintf(
			"Teaching Assistant %s from the college of %s is requesting immediate assistance. Contact %s at %s.",
			user.Name, collegeInfo, user.Name[:fIndex], user.PhoneNumber,
		)
	case "lecturer":
		collegeInfo := "Unknown College"
		if user.College != nil {
			collegeInfo = *user.College
		}
		alertMsg = fmt.Sprintf(
			"Lecturer %s from the college of %s is requesting immediate assistance. Contact %s at %s.",
			user.Name, collegeInfo, user.Name[:fIndex], user.PhoneNumber,
		)
	default:
		alertMsg = fmt.Sprintf(
			"A member of staff by name %s is requesting immediate assistance. Contact %s at %s.",
			user.Name, user.Name[:fIndex], user.PhoneNumber,
		)
	}

	// Always send SOS emails to rapportSafety@gmail.com
	err := mailService.SendEmail(
		senderEmail,
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
						<div style="display: flex; justify-content: center;">
							<a href="tel:1234567890" class="button">Call for Help</a>
						</div>
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

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"msg": "Email sent successfully",
	})
}
