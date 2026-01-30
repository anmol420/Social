package mailer

const (
	maxRetries             = 3
	UserActivationTemplate = "UserInvitationSocial"
)

type Client interface {
	Send(templateFile, username, email string, data any) error
}
