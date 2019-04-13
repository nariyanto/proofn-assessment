package models

import "time"

type Message struct {
	ID             int64     `json:"id"`
	SenderId       int64     `json:"senderId"`
	ReceiverId     int64     `json:"receiverId"`
	Content        string    `json:"content"`
	CreatedDate    time.Time `json:"orderDate"`
	ReceiverStatus int8      `json:"receiverStatus"`
	SenderStatus   int8      `json:"senderStatus"`
}

type MessagesResp struct {
	Messages []Message `json:"message"`
}
