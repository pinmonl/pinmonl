package pubsub

// Message stores information of broadcasting.
type Message struct {
	UserID string      `json:"-"`
	Topic  string      `json:"topic"`
	Data   interface{} `json:"data"`
}

// NewMessage creates pubsub Message.
func NewMessage(userID, topic string, data interface{}) *Message {
	return &Message{
		UserID: userID,
		Topic:  topic,
		Data:   data,
	}
}
