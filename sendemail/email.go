package email

import (
    "fmt"
    "github.com/sendgrid/sendgrid-go"
)

func SendEmail(emailRecipient string) {
    sg := sendgrid.NewSendGridClient("kernkw", "password")
    message := sendgrid.NewMail()
    message.AddTo(emailRecipient)
    message.SetSubject("SendGrid Testing")
    message.SetHTML("<html><body>link</body></html>")
    message.SetText("link")
    message.SetFrom("help@sendgrid.com")
    if r := sg.Send(message); r == nil {
        fmt.Println("Email sent!")
    } else {
        fmt.Println(r)
    }
}