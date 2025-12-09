package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
	"wtm-backend/config"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

// SMTPProvider types
type SMTPProvider string

const (
	ProviderSES     SMTPProvider = "ses"
	ProviderGmail   SMTPProvider = "gmail"
	ProviderMailhog SMTPProvider = "mailhog"
)

// SMTPAccount konfigurasi per provider
type SMTPAccount struct {
	Name        SMTPProvider
	Host        string
	Port        int
	Username    string
	Password    string
	DisableAuth bool

	// TLS settings
	UseImplicitTLS bool // true = port 465, false = port 587/25
	SkipSTARTTLS   bool // true = jangan pakai STARTTLS (untuk relay yang bermasalah)
	TLSConfig      *tls.Config

	// Email settings
	DefaultFrom   string
	SupportEmail  string
	SetReplyTo    bool
	SetReturnPath bool
}

// ScopeRoute defines routing per scope
type ScopeRoute struct {
	Scope    constant.Scope
	From     string
	Provider SMTPAccount
}

// SMTPEmailSender main service
type SMTPEmailSender struct {
	Routes       map[constant.Scope]ScopeRoute
	SupportEmail string

	// Retry & Timeout configuration
	MaxRetries     int
	RetryDelay     time.Duration
	DialTimeout    time.Duration
	SendTimeout    time.Duration
	CommandTimeout time.Duration
	HostIP         string
}

// NewSMTPEmailSender initializes sender
func NewSMTPEmailSender(cfg *config.Config) *SMTPEmailSender {
	accounts := buildAccounts(cfg)

	routes := map[constant.Scope]ScopeRoute{
		constant.ScopeAgent: {
			Scope:    constant.ScopeAgent,
			From:     cfg.EmailFromAgent,
			Provider: accounts[cfg.ProviderAgent],
		},
		constant.ScopeHotel: {
			Scope:    constant.ScopeHotel,
			From:     cfg.EmailFromHotel,
			Provider: accounts[cfg.ProviderHotel],
		},
	}

	return &SMTPEmailSender{
		Routes:         routes,
		SupportEmail:   cfg.EmailContactUs,
		MaxRetries:     3,
		RetryDelay:     cfg.RetryDelay,     // 2s
		DialTimeout:    cfg.DialTimeout,    // 10s
		SendTimeout:    cfg.SendTimeout,    // 30s
		CommandTimeout: cfg.CommandTimeout, // 5s
		HostIP:         cfg.HostIP,
	}
}

// buildAccounts centralizes account creation
func buildAccounts(cfg *config.Config) map[string]SMTPAccount {
	return map[string]SMTPAccount{
		"ses": {
			Name:           ProviderSES,
			Host:           cfg.HostSES,
			Port:           cfg.PortSES,
			Username:       cfg.UsernameSES,
			Password:       cfg.PasswordSES,
			DisableAuth:    cfg.DisableAuthSES,
			UseImplicitTLS: cfg.PortSES == 465,
			TLSConfig:      &tls.Config{ServerName: cfg.HostNameSES},
			DefaultFrom:    cfg.DefaultFromSES,
			SupportEmail:   cfg.SupportEmailSES,
			SetReplyTo:     cfg.ProviderSESReplyTo,
			SetReturnPath:  cfg.ProviderSESReturnPath,
		},
		"gmail": {
			Name:           ProviderGmail,
			Host:           cfg.HostGmail,
			Port:           cfg.PortGmail,
			Username:       cfg.UsernameGmail,
			Password:       cfg.PasswordGmail,
			DisableAuth:    cfg.DisableAuthGmail,
			UseImplicitTLS: cfg.PortGmail == 465,
			SkipSTARTTLS:   cfg.DisableAuthGmail, // Skip STARTTLS untuk Gmail Relay IP-based
			TLSConfig:      &tls.Config{ServerName: cfg.HostGmail},
			DefaultFrom:    cfg.DefaultFromGmail,
			SupportEmail:   cfg.SupportEmailGmail,
			SetReplyTo:     cfg.ProviderGmailReplyTo,
			SetReturnPath:  cfg.ProviderGmailReturnPath,
		},
		"mailhog": {
			Name:           ProviderMailhog,
			Host:           cfg.HostMailhog,
			Port:           cfg.PortMailhog,
			DisableAuth:    true,
			UseImplicitTLS: false,
			DefaultFrom:    cfg.DefaultFromMailhog,
		},
	}
}

// Send sends email with retry logic
func (s *SMTPEmailSender) Send(
	ctx context.Context,
	scope constant.Scope,
	to string,
	subject string,
	bodyHTML string,
	bodyText string,
) error {
	// Validation
	if to == "" {
		return fmt.Errorf("email recipient is empty")
	}

	route, ok := s.Routes[scope]
	if !ok {
		return fmt.Errorf("no route for scope: %s", scope)
	}

	// Handle support email alias
	if to == constant.SupportEmail {
		to = s.SupportEmail
	}

	acc := route.Provider
	from := s.resolveFrom(route, acc)

	// Build message once
	msg := s.buildMessage(route, to, subject, bodyHTML, bodyText, from)

	logger.Info(ctx, fmt.Sprintf("[email] scope=%s provider=%s from=%s to=%s",
		scope, acc.Name, from, to))

	var lastErr error

	// Retry logic
	for attempt := 0; attempt <= s.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(s.RetryDelay)
			logger.Info(ctx, fmt.Sprintf("[email] retry #%d for provider=%s", attempt, acc.Name))
		}

		err := s.sendEmail(ctx, acc, from, to, msg)
		if err == nil {
			logger.Info(ctx, fmt.Sprintf("[email] ✅ sent via %s to %s", acc.Name, to))
			return nil
		}

		lastErr = err
		logger.Error(ctx, fmt.Sprintf("[email] ❌ attempt=%d provider=%s error=%v",
			attempt+1, acc.Name, err))
	}

	return fmt.Errorf("all attempts failed for %s: %w", acc.Name, lastErr)
}

// sendEmail sends email using SMTP with timeout
func (s *SMTPEmailSender) sendEmail(
	ctx context.Context,
	acc SMTPAccount,
	from string,
	to string,
	msg []byte,
) error {
	// Create context with timeout
	sendCtx, cancel := context.WithTimeout(ctx, s.SendTimeout)
	defer cancel()

	// Channel untuk hasil
	errChan := make(chan error, 1)

	go func() {
		errChan <- s.dialAndSend(ctx, acc, from, to, msg)
	}()

	// Wait dengan timeout
	select {
	case err := <-errChan:
		return err
	case <-sendCtx.Done():
		return fmt.Errorf("send timeout after %v: %w", s.SendTimeout, sendCtx.Err())
	}
}

// dialAndSend performs actual SMTP operations
func (s *SMTPEmailSender) dialAndSend(
	ctx context.Context,
	acc SMTPAccount,
	from string,
	to string,
	msg []byte,
) error {
	addr := fmt.Sprintf("%s:%d", acc.Host, acc.Port)

	// Step 1: Establish connection
	logger.Info(ctx, fmt.Sprintf("[email-debug] Step 1: Dialing %s", addr))
	conn, err := s.dialWithTimeout(ctx, acc, addr)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}
	defer conn.Close()

	// Set deadline yang lebih generous untuk Gmail
	deadline := s.CommandTimeout
	if acc.Name == ProviderGmail {
		deadline = s.CommandTimeout * 3 // 15s untuk Gmail
	}

	if err := conn.SetDeadline(time.Now().Add(deadline)); err != nil {
		return fmt.Errorf("set deadline failed: %w", err)
	}

	// Step 2: Create SMTP client
	logger.Info(ctx, "[email-debug] Step 2: Creating SMTP client")
	client, err := smtp.NewClient(conn, acc.Host)
	if err != nil {
		return fmt.Errorf("smtp client failed: %w", err)
	}
	defer client.Close()

	// Step 3: Send EHLO/HELO
	hostname := s.HostIP
	logger.Info(ctx, fmt.Sprintf("[email-debug] Step 3: Sending EHLO %s", hostname))
	if err := client.Hello(hostname); err != nil {
		return fmt.Errorf("EHLO/HELO failed: %w", err)
	}

	// Step 4: STARTTLS for non-implicit TLS (port 587)
	if !acc.UseImplicitTLS && !acc.SkipSTARTTLS {
		if ok, _ := client.Extension("STARTTLS"); ok {
			logger.Info(ctx, "[email-debug] Step 4: Starting TLS")
			if err := client.StartTLS(acc.TLSConfig); err != nil {
				return fmt.Errorf("starttls failed: %w", err)
			}

			// Reset deadline after STARTTLS
			if err := conn.SetDeadline(time.Now().Add(deadline)); err != nil {
				return fmt.Errorf("set deadline after tls failed: %w", err)
			}
		}
	} else if acc.SkipSTARTTLS {
		logger.Info(ctx, "[email-debug] Step 4: Skipping STARTTLS (relay mode)")
	}

	// Step 5: Authentication
	if !acc.DisableAuth && acc.Username != "" {
		logger.Info(ctx, fmt.Sprintf("[email-debug] Step 5: Authenticating as %s", acc.Username))
		auth := smtp.PlainAuth("", acc.Username, acc.Password, acc.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}

		// Reset deadline after Auth
		if err := conn.SetDeadline(time.Now().Add(deadline)); err != nil {
			return fmt.Errorf("set deadline after auth failed: %w", err)
		}
	}

	// Step 6: Validate From address matches Username for Gmail (HANYA jika pakai Auth)
	if acc.Name == ProviderGmail && !acc.DisableAuth && acc.Username != "" {
		actualFrom := extractEmail(from)
		logger.Info(ctx, fmt.Sprintf("[email-debug] Step 6: Validating From=%s vs Username=%s", actualFrom, acc.Username))
		if actualFrom != acc.Username {
			return fmt.Errorf("From address (%s) must match authenticated user (%s) for Gmail", actualFrom, acc.Username)
		}
	} else if acc.Name == ProviderGmail && acc.DisableAuth {
		logger.Info(ctx, fmt.Sprintf("[email-debug] Step 6: Skipping validation (Gmail Relay with IP auth), From=%s", from))
	}

	// Step 7: Send mail commands with extended deadline
	if err := conn.SetDeadline(time.Now().Add(deadline)); err != nil {
		return fmt.Errorf("set deadline for mail failed: %w", err)
	}

	logger.Info(ctx, fmt.Sprintf("[email-debug] Step 7: MAIL FROM %s", from))
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}

	logger.Info(ctx, fmt.Sprintf("[email-debug] Step 8: RCPT TO %s", to))
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("RCPT TO failed: %w", err)
	}

	// Step 9: Send message body
	logger.Info(ctx, "[email-debug] Step 9: Sending DATA")
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}

	logger.Info(ctx, fmt.Sprintf("[email-debug] Step 10: Writing message (%d bytes)", len(msg)))
	if _, err := wc.Write(msg); err != nil {
		wc.Close()
		return fmt.Errorf("write message failed: %w", err)
	}

	logger.Info(ctx, "[email-debug] Step 11: Closing DATA writer")
	if err := wc.Close(); err != nil {
		return fmt.Errorf("close data writer failed: %w", err)
	}

	logger.Info(ctx, "[email-debug] ✅ All steps completed successfully")
	return nil
}

// extractEmail extracts email address from formats like:
// "Display Name <email@domain.com>" -> "email@domain.com"
// "email@domain.com" -> "email@domain.com"
func extractEmail(address string) string {
	if idx := strings.Index(address, "<"); idx >= 0 {
		if endIdx := strings.Index(address[idx:], ">"); endIdx >= 0 {
			return strings.TrimSpace(address[idx+1 : idx+endIdx])
		}
	}
	return strings.TrimSpace(address)
}

// dialWithTimeout establishes connection with proper TLS handling
func (s *SMTPEmailSender) dialWithTimeout(
	ctx context.Context,
	acc SMTPAccount,
	addr string,
) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout: s.DialTimeout,
	}

	if acc.UseImplicitTLS {
		// Port 465: Implicit TLS (direct TLS connection)
		return tls.DialWithDialer(dialer, "tcp", addr, acc.TLSConfig)
	}

	// Port 587/25: Plain connection (will upgrade via STARTTLS)
	dialCtx, cancel := context.WithTimeout(ctx, s.DialTimeout)
	defer cancel()
	return dialer.DialContext(dialCtx, "tcp", addr)
}

// buildMessage constructs RFC 5322 compliant email
func (s *SMTPEmailSender) buildMessage(
	route ScopeRoute,
	to string,
	subject string,
	bodyHTML string,
	bodyText string,
	from string,
) []byte {
	var b strings.Builder

	// Required headers
	b.WriteString(fmt.Sprintf("From: %s\r\n", from))
	b.WriteString(fmt.Sprintf("To: %s\r\n", to))
	b.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))

	// Optional headers (Reply-To, Return-Path)
	acc := route.Provider
	if acc.SetReplyTo {
		replyTo := acc.SupportEmail
		if replyTo == "" {
			replyTo = from
		}
		b.WriteString(fmt.Sprintf("Reply-To: %s\r\n", replyTo))
	}

	if acc.SetReturnPath && acc.SupportEmail != "" {
		b.WriteString(fmt.Sprintf("Return-Path: %s\r\n", acc.SupportEmail))
	}

	b.WriteString("MIME-Version: 1.0\r\n")

	// Multipart setup
	boundary := fmt.Sprintf("----=_Part_%d_%d", time.Now().Unix(), time.Now().UnixNano())
	b.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", boundary))
	b.WriteString("\r\n")

	// Plain text part
	b.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	b.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
	if bodyText != "" {
		b.WriteString(bodyText)
	} else {
		b.WriteString("Please view the HTML version of this email.")
	}
	b.WriteString("\r\n\r\n")

	// HTML part
	if bodyHTML != "" {
		b.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		b.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		b.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
		b.WriteString(bodyHTML)
		b.WriteString("\r\n\r\n")
	}

	// End boundary
	b.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return []byte(b.String())
}

// resolveFrom determines the From address
func (s *SMTPEmailSender) resolveFrom(route ScopeRoute, acc SMTPAccount) string {
	if route.From != "" {
		return route.From
	}
	return acc.DefaultFrom
}
