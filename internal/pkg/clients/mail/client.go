package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"strconv"

	"gopkg.in/gomail.v2"

	"boilerplate/internal/model"
	logger_okg "boilerplate/internal/pkg/logger"
)

type Client interface {
	Send(ctx context.Context, to, subject, body string, attachments []*Attachment) error
}

type client struct {
	logger logger_okg.Logger
	config model.ConfigMail
}

func NewClient(logger logger_okg.Logger, config model.ConfigMail) Client {
	return &client{
		logger: logger,
		config: config,
	}
}

func (c *client) Send(ctx context.Context, to, subject, body string, attachments []*Attachment) error {
	c.logger.DebugKV(ctx, "send email", "from", c.config.From, "to", to, "subject", subject)

	smtpPort, err := strconv.Atoi(c.config.SMTPPort)
	if err != nil {
		return fmt.Errorf("convert smtp port: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader(FieldFrom, c.config.From)
	m.SetHeader(FieldTo, to)
	m.SetHeader(FieldSubject, subject)
	m.SetBody(ContentTypeHTML, body)

	for _, attachment := range attachments {
		m.Attach(attachment.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := io.Copy(w, attachment.Content)
			if err != nil {
				return fmt.Errorf("copy attachment content: %w", err)
			}
			return nil
		}))
	}

	d := gomail.NewDialer(c.config.SMTPHost, smtpPort, c.config.Username, c.config.Password)

	if c.config.SSL {
		d.SSL = true
	} else if c.config.TLS {
		d.TLSConfig = &tls.Config{
			InsecureSkipVerify: true, // nolint:gosec
		}
	}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("dial and send mail: %w", err)
	}

	return nil
}
