package seed

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (s *Seed) SeedEmailTemplate() {
	bodyAgentApproval := `
<p>Dear {{.AgentName}},</p>

<p>We’re excited to welcome you to The HotelBox! Your registration has been successfully reviewed and approved. You can now access your agent account, explore exclusive hotel deals, and start booking with us.</p>

<p>👉 <a href="{{.LoginLink}}">Login here</a></p>

<p>We’re committed to helping you grow your travel business with access to competitive rates and trusted hotel partners.</p>

<p>If you have any questions, our support team is always ready to assist.</p>

<p>Welcome aboard,<br>The HotelBox Team</p>
`
	bodyAgentRejection := `
<p>Dear {{.AgentName}},</p>

<p>Thank you for applying to become a travel partner with The HotelBox. After reviewing your application, we were unable to approve your registration at this time due to the following reason(s):</p>

<p><strong>{{.RejectionReason}}</strong></p>

<p>We would be delighted to have you join our network, and we encourage you to re-register once the above issue(s) have been resolved.</p>

<p>👉 <a href="{{.ReRegisterLink}}">Re-register here</a></p>

<p>If you need guidance on the process or clarification on the requirements, our support team will be happy to assist you.</p>

<p>We truly value your interest in The HotelBox and look forward to welcoming you soon.</p>

<p>Warm regards,<br>The HotelBox Team</p>
`
	bodyBookingConfirmed := `
<p>Dear {{.AgentName}},</p>

<p>Great news! Your booking request has been <strong>confirmed</strong>.</p>

<h3>Booking Summary:</h3>
<ul>
    <li><strong>Booking Id:</strong> {{.BookingID}}</li>
    <li><strong>Guest Name:</strong> {{.GuestName}}</li>
</ul>

<h3>Itinerary:</h3>
{{range .SubBookings}}
<div style="margin-bottom: 1em;">
    <p><strong>Sub-booking {{.Index}}</strong></p>
    <ul>
        <li><strong>Hotel:</strong> {{.HotelName}}</li>
        <li><strong>Check-in:</strong> {{.CheckIn}}</li>
        <li><strong>Check-out:</strong> {{.CheckOut}}</li>
    </ul>
</div>
{{end}}

<p>You can view full booking details here: <a href="{{.BookingLink}}">{{.BookingLink}}</a></p>

<p>Thank you for booking with <strong>The HotelBox</strong>. We look forward to serving you.</p>

<p>Best regards,<br>The HotelBox Team</p>
`
	bodyBookingRejected := `
<p>Dear {{.AgentName}},</p>

<p>Unfortunately, your booking request <strong>{{.BookingID}}</strong> could not be confirmed due to the following reason(s):</p>

<p><strong>{{.RejectionReason}}</strong></p>

<h3>Booking Summary:</h3>
<ul>
    <li><strong>Booking Id:</strong> {{.BookingID}}</li>
    <li><strong>Guest Name:</strong> {{.GuestName}}</li>
</ul>

<h3>Itinerary Attempted:</h3>
{{range .SubBookings}}
<div style="margin-bottom: 1em;">
    <p><strong>Sub-booking {{.Index}}</strong></p>
    <ul>
        <li><strong>Hotel:</strong> {{.HotelName}}</li>
        <li><strong>Check-in:</strong> {{.CheckIn}}</li>
        <li><strong>Check-out:</strong> {{.CheckOut}}</li>
    </ul>
</div>
{{end}}

<p>We encourage you to try booking again with alternative dates or hotels.</p>

<p>Or if you prefer, you can search for other available hotels here: <a href="{{.HomePageLink}}">{{.HomePageLink}}</a></p>

<p>We value your partnership and are here to help you find the best options for your clients.</p>

<p>Sincerely,<br>The HotelBox Team</p>
`
	bodyHotelBookingRequest := `
<p>Dear Reservation Team,</p>

<p><em>" Warmest Greeting From World Travel Marketing Bali "</em></p>

<p>Please kindly assist us to <strong>BOOK & CONFIRM</strong> reservation with details as below:</p>

<ul>
    <li><strong>NAME:</strong> {{.GuestName}}</li>
    <li><strong>PERIOD:</strong> {{.Period}}</li>
    <li><strong>ROOM:</strong> {{.RoomType}}</li>
    <li><strong>RATE:</strong> {{.Rate}}</li>
    <li><strong>BOOKING CODE:</strong> {{.BookingCode}}</li>
    <li><strong>BENEFIT:</strong> {{.Benefit}}</li>
    <li><strong>REMARK:</strong> {{.Remark}}</li>
    <li><strong>ADDITIONAL:</strong> {{.Additional}}</li>
</ul>

<p>We are looking forward to hearing back from you soon.<br>
Many thanks for your kind attention and assistance.</p>

<p>Best Regards,<br>{{.SystemSignature}}</p>
`
	bodyContactUsGeneral := `
<p>Hello Team,</p>

<p>A new inquiry has been submitted via the Contact Us page.</p>

<p>Details:</p>

<p>
<strong>Name:</strong> {{.UserName}}<br>
<strong>Email:</strong> {{.UserEmail}}<br>
<strong>Subject:</strong> {{.Subject}}<br>
<strong>Message:</strong><br>
{{.UserMessage}}
</p>

<p>Please review and respond accordingly.</p>

<p>Best regards,<br>The HotelBox System</p>
`
	bodyContactUsBooking := `
<p>Dear Support Team,</p>

<p>A help request has been submitted by an agent via the Help Page.</p>

<p><strong>Agent Details:</strong><br>
Name: {{.AgentName}}<br>
Email: {{.AgentEmail}}<br>
Agency: {{.AgencyName}}<br>
Contact Number: {{.AgentPhone}}</p>

<p><strong>Booking Details:</strong><br>
Booking Id: {{.BookingID}}<br>
Guest Name: {{.GuestName}}</p>

<p><strong>Itinerary:</strong><br>
{{range $index, $sb := .SubBookings}}
<strong>Sub-booking {{$index | add1}}</strong><br>
Hotel: {{$sb.Hotel}}<br>
Check-in: {{$sb.CheckIn}}<br>
Check-out: {{$sb.CheckOut}}<br><br>
{{end}}
</p>

<p><strong>Agent Message:</strong><br>
{{.AgentMessage}}</p>

<p>Please review and provide assistance accordingly.</p>

<p>Best regards,<br>
The HotelBox System</p>
`
	bodyForgotPassword := `
<p>Hello {{.FullName}},</p>

<p>We received a request to reset your account password. If you made this request,</p>

<p>
Please click the link below to set a new password:<br>
👉 <a href="{{.ResetLink}}" target="_blank">{{.ResetLink}}</a>
</p>

<p>
For your security, this link will expire in <strong>{{.ExpiresIn}}</strong>.<br>
If you did not request a password reset, please ignore this email — your account will remain secure.
</p>

<p>Best regards,<br>
World Travel Management</p>
`
	bodyAccountActivated := `
<p>Hello {{.FullName}},</p>

<p>Congratulations! Your account has been successfully activated.</p>

<p>
For security reasons, please update your password immediately after your first login.
</p>

<p><strong>Your temporary password is:</strong><br>
{{.TemporaryPassword}}</p>

<p>
We strongly recommend changing this temporary password to a new, secure one to keep your account protected.
</p>

<p>Best regards,<br>
World Travel Management</p>
`
	templates := []model.EmailTemplate{
		{Subject: `🎉 Welcome to The HotelBox – Your Agent Account is Approved!`, Body: bodyAgentApproval, Name: constant.EmailAgentApproved, IsSignatureImage: false},
		{Subject: `Your The HotelBox Registration – Action Required`, Body: bodyAgentRejection, Name: constant.EmailAgentRejected, IsSignatureImage: false},
		{Subject: `Booking Confirmation – {{.BookingID}}`, Body: bodyBookingConfirmed, Name: constant.EmailBookingConfirmed, IsSignatureImage: false},
		{Subject: `Booking Request Update – {{.BookingID}}`, Body: bodyBookingRejected, Name: constant.EmailBookingRejected, IsSignatureImage: false},
		{Subject: `New Booking Request – {{.BookingCode}}`, Body: bodyHotelBookingRequest, Name: constant.EmailHotelBookingRequest, IsSignatureImage: false},
		{Subject: `New Contact Us Submission – The HotelBox - {{.UserName}}`, Body: bodyContactUsGeneral, Name: constant.EmailContactUsGeneral, IsSignatureImage: false},
		{Subject: `Help Request for Booking – {{.BookingID}}`, Body: bodyContactUsBooking, Name: constant.EmailContactUsBooking, IsSignatureImage: false},
		{Subject: `Password Reset Request`, Body: bodyForgotPassword, Name: constant.EmailForgotPassword, IsSignatureImage: false},
		{Subject: `Your Account Has Been Activated – Please Change Your Password Immediately`, Body: bodyAccountActivated, Name: constant.EmailAccountActivated, IsSignatureImage: false},
	}

	for _, tpl := range templates {
		var existing model.EmailTemplate
		err := s.db.
			Where("name = ?", tpl.Name).
			First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.db.Create(&tpl).Error; err != nil {
				log.Printf("Failed to insert email template %s: %s", tpl.Name, err.Error())
			}
		} else if err == nil {
			// Update only if subject/body changed
			if existing.Subject != tpl.Subject || existing.Body != tpl.Body || existing.IsSignatureImage != tpl.IsSignatureImage {
				existing.Subject = tpl.Subject
				existing.Body = tpl.Body
				existing.IsSignatureImage = tpl.IsSignatureImage
				if err := s.db.Save(&existing).Error; err != nil {
					log.Printf("Failed to update email template %s: %s", tpl.Name, err.Error())
				}
			}
		} else {
			log.Printf("Error checking email template %s: %s", tpl.Name, err.Error())
		}
	}

	log.Println("Seeding email templates completed")
}

func (s *Seed) SeedingEmailLog() {
	ctx := context.Background()
	statusEmails := []model.StatusEmail{
		{ID: constant.StatusEmailPendingID, Status: constant.StatusEmailPending},
		{ID: constant.StatusEmailSuccessID, Status: constant.StatusEmailSuccess},
		{ID: constant.StatusEmailFailedID, Status: constant.StatusEmailFailed},
	}

	for _, se := range statusEmails {
		var existing model.StatusEmail
		err := s.db.Where("id = ?", se.ID).First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.db.Create(&se).Error; err != nil {
				logger.Info(ctx, fmt.Sprintf("❌ Failed to insert StatusEmail ID %d: %s", se.ID, err.Error()))
			} else {
				logger.Info(ctx, fmt.Sprintf("✅ Inserted StatusEmail ID %d: %s", se.ID, se.Status))
			}
		} else if err == nil {
			if existing.Status != se.Status {
				existing.Status = se.Status
				if err := s.db.Save(&existing).Error; err != nil {
					logger.Info(ctx, fmt.Sprintf("❌ Failed to update StatusEmail ID %d: %s", se.ID, err.Error()))
				} else {
					logger.Info(ctx, fmt.Sprintf("🔄 Updated StatusEmail ID %d to: %s", se.ID, se.Status))
				}
			} else {
				logger.Info(ctx, fmt.Sprintf("⏩ StatusEmail ID %d already up-to-date: %s", se.ID, se.Status))
			}
		} else {
			logger.Info(ctx, fmt.Sprintf("⚠️ Error checking StatusEmail ID %d: %s", se.ID, err.Error()))
		}
	}

	logger.Info(ctx, fmt.Sprintf("📦 Seeding StatusEmail completed with insert/update checks"))
}
