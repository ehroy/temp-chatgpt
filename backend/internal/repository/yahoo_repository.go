package repository

import (
	"context"
	"bytes"
	"encoding/base64"
	"crypto/tls"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/textproto"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"emailchatgpt/internal/model"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type YahooConfig struct {
	Host           string
	Port           string
	Username       string
	AppPassword    string
	AllowedFolders []string
	MaxScan        int
}

type YahooRepository struct {
	cfg YahooConfig
}

func NewYahooRepository(cfg YahooConfig) *YahooRepository {
	if cfg.MaxScan <= 0 {
		cfg.MaxScan = 25
	}
	return &YahooRepository{cfg: cfg}
}

func (r *YahooRepository) ListMessages(ctx context.Context, email string) ([]model.EmailMessage, error) {
	if strings.TrimSpace(r.cfg.Username) == "" || strings.TrimSpace(r.cfg.AppPassword) == "" {
		return nil, fmt.Errorf("credential imap belum diatur")
	}

	client, err := r.dial()
	if err != nil {
		return nil, err
	}
	defer client.Logout()

	if err := client.Login(r.cfg.Username, r.cfg.AppPassword); err != nil {
		return nil, err
	}

	var results []model.EmailMessage
	for _, folder := range r.allowedFolders() {
		messages, err := r.listFolder(ctx, client, folder, email)
		if err != nil {
			continue
		}
		results = append(results, messages...)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].ReceivedAt.After(results[j].ReceivedAt)
	})
	return results, nil
}

func (r *YahooRepository) dial() (*client.Client, error) {
	address := net.JoinHostPort(r.cfg.Host, r.cfg.Port)
	conn, err := tls.Dial("tcp", address, &tls.Config{ServerName: r.cfg.Host})
	if err != nil {
		return nil, err
	}
	cli, err := client.New(conn)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (r *YahooRepository) allowedFolders() []string {
	filtered := make([]string, 0, len(r.cfg.AllowedFolders))
	for _, folder := range r.cfg.AllowedFolders {
		if trimmed := strings.TrimSpace(folder); trimmed != "" {
			filtered = append(filtered, trimmed)
		}
	}
	return filtered
}

func (r *YahooRepository) listFolder(ctx context.Context, c *client.Client, folder, targetEmail string) ([]model.EmailMessage, error) {
	_, err := c.Select(folder, false)
	if err != nil {
		return nil, err
	}

	status, err := c.Status(folder, []imap.StatusItem{imap.StatusMessages})
	if err != nil {
		return nil, err
	}
	if status.Messages == 0 {
		return nil, nil
	}

	maxScan := uint32(r.cfg.MaxScan)
	start := uint32(1)
	if status.Messages > maxScan {
		start = status.Messages - maxScan + 1
	}

	set := new(imap.SeqSet)
	set.AddRange(start, status.Messages)

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}
	messages := make(chan *imap.Message, int(maxScan))
	fetchErr := make(chan error, 1)
	go func() {
		fetchErr <- c.Fetch(set, items, messages)
	}()

	results := make([]model.EmailMessage, 0)
	for msg := range messages {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		if msg == nil || msg.Envelope == nil {
			continue
		}
		body := msg.GetBody(section)
		if body == nil {
			continue
		}
		parsed, err := mail.ReadMessage(body)
		if err != nil {
			continue
		}
		messageEmail := strings.TrimSpace(strings.ToLower(targetEmail))
		if !matchesRecipient(parsed.Header, messageEmail) {
			continue
		}

		bodyBytes, _ := io.ReadAll(parsed.Body)
		htmlBody, plainBody := extractEmailBody(parsed.Header.Get("Content-Type"), bodyBytes)
		receivedAt := msg.Envelope.Date
		if receivedAt.IsZero() {
			receivedAt = parseDateHeader(parsed.Header.Get("Date"))
		}

		results = append(results, model.EmailMessage{
			ID:         strconv.FormatUint(uint64(msg.SeqNum), 10),
			Folder:     folder,
			Sender:     headerValue(parsed.Header, "From"),
			Recipient:  headerValue(parsed.Header, "To"),
			Subject:    headerValue(parsed.Header, "Subject"),
			Body:       htmlBody,
			Text:       plainBody,
			ReceivedAt: receivedAt,
		})
	}

	if err := <-fetchErr; err != nil {
		return results, err
	}

	return results, nil
}

func matchesRecipient(h mail.Header, email string) bool {
	if email == "" {
		return true
	}
	fields := []string{"To", "Cc", "Delivered-To", "X-Original-To"}
	for _, field := range fields {
		if strings.Contains(strings.ToLower(h.Get(field)), email) {
			return true
		}
	}
	return false
}

func headerValue(h mail.Header, key string) string {
	return strings.TrimSpace(h.Get(key))
}

func parseDateHeader(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	parsed, err := mail.ParseDate(value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func extractEmailBody(contentType string, raw []byte) (string, string) {
	html, plain, attachments := extractMIMEPart(contentType, textproto.MIMEHeader{}, raw)
	if html == "" {
		if plain != "" {
			html = "<pre style=\"white-space:pre-wrap;font-family:inherit;margin:0;\">" + htmlEscape(plain) + "</pre>"
		} else {
			html = "<pre style=\"white-space:pre-wrap;font-family:inherit;margin:0;\">" + htmlEscape(string(raw)) + "</pre>"
		}
	}
	return rewriteEmailHTML(wrapEmailHTML(html), attachments), plain
}

func extractMIMEPart(contentType string, headers textproto.MIMEHeader, raw []byte) (html string, plain string, attachments map[string]string) {
	attachments = map[string]string{}
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		mediaType = strings.ToLower(strings.TrimSpace(contentType))
	}

	if strings.HasPrefix(strings.ToLower(mediaType), "multipart/") {
		boundary := params["boundary"]
		if boundary == "" {
			return "", "", attachments
		}
		reader := multipart.NewReader(bytes.NewReader(raw), boundary)
		var altHTML, altPlain string
		for {
			part, err := reader.NextPart()
			if err != nil {
				break
			}
			partType := part.Header.Get("Content-Type")
			partBytes, _ := io.ReadAll(part)
			partHTML, partPlain, partAttachments := extractMIMEPart(partType, part.Header, partBytes)
			for k, v := range partAttachments {
				attachments[k] = v
			}
			if partHTML != "" && altHTML == "" {
				altHTML = partHTML
			}
			if partPlain != "" && altPlain == "" {
				altPlain = partPlain
			}
		}
		if altHTML != "" {
			return altHTML, altPlain, attachments
		}
		return "", altPlain, attachments
	}

	decoded := decodeBody(raw, headers.Get("Content-Transfer-Encoding"))
	switch {
	case strings.HasPrefix(strings.ToLower(mediaType), "text/html"):
		return decoded, "", attachments
	case strings.HasPrefix(strings.ToLower(mediaType), "text/plain"):
		return "", decoded, attachments
	}

	if cid := strings.TrimSpace(headers.Get("Content-ID")); cid != "" {
		cid = strings.Trim(cid, "<>")
		attachments[cid] = dataURI(mediaType, decoded)
	}
	return "", "", attachments
}

func decodeBody(raw []byte, encoding string) string {
	switch strings.ToLower(strings.TrimSpace(encoding)) {
	case "base64":
		decoded, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, bytes.NewReader(raw)))
		if err == nil && len(decoded) > 0 {
			return string(decoded)
		}
	case "quoted-printable":
		decoded, err := io.ReadAll(quotedprintable.NewReader(bytes.NewReader(raw)))
		if err == nil && len(decoded) > 0 {
			return string(decoded)
		}
	}
	return string(raw)
}

func dataURI(mediaType string, decoded string) string {
	if mediaType == "" {
		mediaType = "application/octet-stream"
	}
	return "data:" + mediaType + ";base64," + base64.StdEncoding.EncodeToString([]byte(decoded))
}

func wrapEmailHTML(html string) string {
	lower := strings.ToLower(html)
	if strings.Contains(lower, "<html") {
		return html
	}
	return "<!doctype html><html><head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><style>html,body{margin:0;padding:0;background:#fff;color:#111;font-family:Arial,Helvetica,sans-serif}img{max-width:100%;height:auto}.email-wrapper{padding:16px}</style></head><body><div class=\"email-wrapper\">" + html + "</div></body></html>"
}

func rewriteEmailHTML(html string, attachments map[string]string) string {
	if len(attachments) == 0 {
		return html
	}
	re := regexp.MustCompile(`cid:([^\"' >]+)`)
	return re.ReplaceAllStringFunc(html, func(match string) string {
		cid := strings.TrimPrefix(match, "cid:")
		cid = strings.Trim(cid, "<>")
		if data, ok := attachments[cid]; ok {
			return data
		}
		return match
	})
}

func htmlEscape(s string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	return replacer.Replace(s)
}
