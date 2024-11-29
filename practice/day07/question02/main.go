package main

import "fmt"

/**
	Define an interface named Notifier with two methods:
    SendEmail(): A method that takes a string argument representing the email content and returns a string indicating the email has been sent.
    SendSMS(): A method that takes a string argument representing the SMS content and returns a string indicating the SMS has been sent.
    Implement this interface in a struct called NotificationService.
**/

type Notifier interface {
	SendEmail(string) string
	SendSMS(string) string
}

type NotificationService struct {
	conn string
}

func (service *NotificationService) SendEmail(content string) (status string) {
	if content == "" {
		return "No Mail content"
	} else {
		return "Mail Sent"
	}
}

func (service *NotificationService) SendSMS(content string) (status string) {
	if content == "" {
		return "No SMS content"
	} else {
		return "SMS Sent"
	}
}

func main() {
	service := NotificationService{conn: "SMS Service"}
	mailStatus := service.SendEmail("This is email title")
	fmt.Println(mailStatus)
	smsStatus := service.SendSMS("This is SMS text")
	fmt.Println(smsStatus)

	mailStatus = service.SendEmail("")
	fmt.Println(mailStatus)
	smsStatus = service.SendSMS("")
	fmt.Println(smsStatus)
}
