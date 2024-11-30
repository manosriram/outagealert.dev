package integration

type (
	NotificationType   string
	SendGridTemplateId string
)

const (
	MONITOR_DOWN        NotificationType = "monitor_down"
	MONITOR_UP          NotificationType = "monitor_up"
	VERIFY_EMAIL        NotificationType = "verify_email"
	FORGOT_PASSWORD_OTP NotificationType = "forgot_password_otp"
)

const (
	MONITOR_DOWN_TEMPLATE_ID    SendGridTemplateId = "d-cf3e6ff9cbd54df696985ac7ea08475e"
	MONITOR_UP_TEMPLATE_ID      SendGridTemplateId = "d-5b63ef78deee4a37ae62776a3949e0d9"
	VERIFY_EMAIL_TEMPLATE_ID    SendGridTemplateId = "d-c50ac0a5dccb454fbbb6eac650b5e680"
	FORGOT_PASSWORD_TEMPLATE_ID SendGridTemplateId = "d-038cf4d4bd6a492ca28d19f6d8fe3b24"
)

var NotificationTypeVsTemplateId map[NotificationType]SendGridTemplateId = map[NotificationType]SendGridTemplateId{
	MONITOR_UP:          MONITOR_UP_TEMPLATE_ID,
	MONITOR_DOWN:        MONITOR_DOWN_TEMPLATE_ID,
	VERIFY_EMAIL:        VERIFY_EMAIL_TEMPLATE_ID,
	FORGOT_PASSWORD_OTP: FORGOT_PASSWORD_TEMPLATE_ID,
}

var NotificationVsShouldMarkNotificationSent map[NotificationType]bool = map[NotificationType]bool{
	MONITOR_UP:          false,
	MONITOR_DOWN:        true,
	VERIFY_EMAIL:        true,
	FORGOT_PASSWORD_OTP: true,
}
