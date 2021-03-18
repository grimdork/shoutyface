package main

import (
	"fmt"
	"net"
	"net/smtp"
	"time"
)

// RunMailQueue until quit signal is received.
func (srv *Shoutyface) RunMailQueue(q chan interface{}) {
	srv.Add(1)
	srv.L("Mailqueue started.")
	for {
		select {
		case msg := <-srv.mailqueue:
			emails := srv.GetSubscribers(msg.Channel)
			if len(emails) > 0 {
				srv.L("Message to channel %s of severity %s, titled '%s'", msg.Channel, msg.Severity, msg.Title)
				subject := fmt.Sprintf("[%s] %s", msg.Severity, msg.Title)
				err := srv.sendmail(emails, subject, msg.Body)
				if err != nil {
					srv.E("Error sending mail: %s", err.Error())
				}
			} else {
				srv.E("No recipients for channel %s, or no such channel.", msg.Channel)
			}

		case <-srv.mailquit:
			srv.L("Mailqueue stopped.")
			srv.Done()
			return

		default:
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func (srv *Shoutyface) sendmail(to []string, subject, body string) error {
	msg := fmt.Sprintf("Subject: %s\r\nFrom: %s\r\n\r\n%s\r\n", subject, srv.mailfrom, body)
	auth := smtp.PlainAuth("", srv.mailuser, srv.mailpass, srv.mailhost)
	return smtp.SendMail(net.JoinHostPort(srv.mailhost, srv.mailport), auth, srv.mailfrom, to, []byte(msg))
}
