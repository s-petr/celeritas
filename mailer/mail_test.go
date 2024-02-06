package mailer

import (
	"testing"
)

func getTestMsg() Message {
	return Message{
		From:        "test@test.com",
		FromName:    "Test Sender",
		To:          "recipient@test.com",
		Subject:     "Mailer Test",
		Template:    "test",
		Attachments: []string{"./testdata/croc.jpg"},
		Data:        nil,
	}
}

func TestMail_SendUsingSMTPFunc(t *testing.T) {
	msg := getTestMsg()

	if err := mailer.SendSMTPMessage(msg); err != nil {
		t.Error(err)
	}
}

func TestMail_SendUsingSTMPChan(t *testing.T) {
	msg := getTestMsg()

	mailer.Jobs <- msg
	res := <-mailer.Results
	if res.Error != nil {
		t.Error("failed to send over channel")
	}

	msg.To = "invalid-email"

	mailer.Jobs <- msg
	res = <-mailer.Results
	if res.Error == nil {
		t.Error("attempted to send to an invalid email, expected error, received none")
	}
}

func TestMail_SendUsingAPI(t *testing.T) {
	msg := getTestMsg()

	mailer.API = "unknown"
	mailer.APIKey = "test1234"
	mailer.APIURL = "https://faketest.com"

	if err := mailer.SendUsingAPI(msg, "unknown"); err == nil {
		t.Error("attempted to send email using invalid API details, expected error, received none")
	}

	mailer.API = ""
	mailer.APIKey = ""
	mailer.APIURL = ""
}

func TestMail_BuildHTMLMessage(t *testing.T) {
	msg := getTestMsg()

	if _, err := mailer.buildHTMLMessage(msg); err != nil {
		t.Error("failed to build HTML message:", err)
	}
}

func TestMail_BuildPlainTextMessage(t *testing.T) {
	msg := getTestMsg()

	if _, err := mailer.buildPlainTextMessage(msg); err != nil {
		t.Error("failed to build plain text message:", err)
	}
}

func TestMail_Send(t *testing.T) {
	msg := getTestMsg()

	mailer.API = "unknown"
	mailer.APIKey = "test1234"
	mailer.APIURL = "https://faketest.com"

	if err := mailer.Send(msg); err == nil {
		t.Error("attempted to send email using invalid API details, expected error, received none")
	}

	mailer.API = ""
	mailer.APIKey = ""
	mailer.APIURL = ""

	if err := mailer.Send(msg); err != nil {
		t.Error("failed to send message:", err)
	}

}

func TestMail_ChooseAPI(t *testing.T) {
	msg := getTestMsg()

	mailer.API = "unknown"

	if err := mailer.ChooseAPI(msg); err == nil {
		t.Error("attempted to send email using invalid API details, expected error, received none")
	}

}
