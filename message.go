package main

import "github.com/francoispqt/gojay"

// Message is an assembled email based on a template.
type Message struct {
	// Subject is the Subject: part of a standard e-mail.
	Channel string
	// Title is the subject.
	Title string
	// Body is the message.
	Body string
	// Severity of the message.
	Severity string
}

// UnmarshalJSONObject decodes this message from JSON via gojay.
func (msg *Message) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "channel":
		return dec.String(&msg.Channel)
	case "title":
		return dec.String(&msg.Title)
	case "body":
		return dec.String(&msg.Body)
	case "severity":
		return dec.String(&msg.Severity)
	}
	return nil
}

// NKeys is required to unmarshal.
func (msg *Message) NKeys() int {
	return 4
}
