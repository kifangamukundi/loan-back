package templates

import (
	"fmt"
)

func GenerateAgentWelcomeMessage(firstName, lastName, dashboardUrl, supportEmail, supportPhone, companyName string) string {
	return fmt.Sprintf(`
        <p>Dear <strong>%s %s</strong>,</p>
        <p>We are pleased to inform you that you have been successfully added to our system as an agent with <strong>%s</strong>.</p>
        <p>To start managing your tasks, please visit your dashboard by clicking the following link:</p>
        <a href="%s" clicktracking="off">%s</a>
        <p>You can now log in and access all the features available to agents, including managing your assigned tasks and more.</p>
        <p>If you have any questions or need assistance, feel free to reach out to our support team at <strong>%s</strong> or <strong>%s</strong>.</p>
        <p>We are excited to have you on board and look forward to your contributions!</p>
        <p>Best regards,</p>
        <p><strong>%s</strong></p>
    `, firstName, lastName, companyName, dashboardUrl, dashboardUrl, supportEmail, supportPhone, companyName)
}

func GenerateLoanOfficerWelcomeMessage(firstName, lastName, dashboardUrl, supportEmail, supportPhone, companyName string) string {
	return fmt.Sprintf(`
        <p>Dear <strong>%s %s</strong>,</p>
        <p>We are pleased to welcome you to <strong>%s</strong> as a Loan Officer.</p>
        <p>Your role is crucial in helping clients navigate their loan applications, managing financial services, and ensuring smooth loan processing.</p>
        <p>To access your dashboard and start managing loan applications, please click the following link:</p>
        <a href="%s" clicktracking="off">%s</a>
        <p>Through your dashboard, you will be able to review applications, track loan statuses, communicate with clients, and perform other essential tasks.</p>
        <p>If you have any questions or require assistance, our support team is available at <strong>%s</strong> or <strong>%s</strong>.</p>
        <p>We are excited to have you on board and look forward to your contributions in empowering our clients with financial solutions!</p>
        <p>Best regards,</p>
        <p><strong>%s</strong></p>
    `, firstName, lastName, companyName, dashboardUrl, dashboardUrl, supportEmail, supportPhone, companyName)
}

func GenerateMemberAddedByAgentMessage(firstName, lastName, dashboardUrl, supportEmail, supportPhone, companyName string) string {
	return fmt.Sprintf(`
        <p>Dear <strong>%s %s</strong>,</p>
        <p>We are pleased to inform you that you have been added to a new group in <strong>%s</strong> by an agent.</p>
        <p>Your participation in this group will enable you to collaborate with other members, access essential resources, and contribute effectively.</p>
        <p>To get started, please log in to your dashboard using the following link:</p>
        <a href="%s" clicktracking="off">%s</a>
        <p>Through your dashboard, you can engage with other members, manage tasks, and stay updated on important information.</p>
        <p>If you have any questions or need assistance, feel free to reach out to our support team at <strong>%s</strong> or <strong>%s</strong>.</p>
        <p>We look forward to your active participation in the group!</p>
        <p>Best regards,</p>
        <p><strong>%s</strong></p>
    `, firstName, lastName, companyName, dashboardUrl, dashboardUrl, supportEmail, supportPhone, companyName)
}

func GenerateActivationMessage(firstName, lastName, activationUrl, supportEmail, supportPhone, companyName string) string {
	return fmt.Sprintf(`
        <p>Dear <strong>%s %s</strong>,</p>
        <p>Thank you for registering with our website! We're excited to have you as a new member of our community.</p>
        <p>To activate your account, please click the following link:</p>
        <a href="%s" clicktracking="off">%s</a>
        <p>If the link above does not work, please copy and paste the URL below into your browser:</p>
        <a href="%s" clicktracking="off">%s</a>
        <p>Once your account is activated, you'll be able to log in to our website and enjoy all the benefits of membership and our services.</p>
        <p>If you have any questions or concerns, please don't hesitate to reach out to our customer support team at <strong>%s</strong> or <strong>%s</strong>.</p>
        <p>Thank you again for joining us. We look forward to connecting with you soon!</p>
        <p>Best regards,</p>
        <p><strong>%s</strong></p>
    `, firstName, lastName, activationUrl, activationUrl, activationUrl, activationUrl, supportEmail, supportPhone, companyName)
}

// GenerateResetPasswordMessage generates an HTML email message for password reset
func GenerateResetPasswordMessage(firstName, lastName, resetUrl string) string {
	message := fmt.Sprintf(`
		<html>
			<head>
				<title>Password Reset Request</title>
			</head>
			<body>
				<p>Hi %s %s,</p>
				<p>We received a request to reset your password. To complete the process, please click the link below:</p>
				<p><a href="%s">%s</a></p>
				<p>If you did not request a password reset, please ignore this email.</p>
				<p>Best regards,</p>
				<p>Your Team</p>
			</body>
		</html>
	`, firstName, lastName, resetUrl, resetUrl)

	return message
}

// GenerateResetPasswordMessage generates the message body for password reset confirmation
func GenerateChangedPasswordMessage(firstName, lastName, companyName string) string {
	return fmt.Sprintf(`
	<p>Dear <strong>%s %s</strong>,</p>
	<p>Your password has been successfully reset. If you did not initiate this change, please contact support immediately.</p>
	<p>Best regards,</p>
	<p><strong>%s</strong></p>
	`, firstName, lastName, companyName)
}
