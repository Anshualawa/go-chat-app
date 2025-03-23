package golang_social_chat

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type Chat struct {
	ID        int    `json:"id"`
	User1ID   int    `json:"user1_id"`
	User2ID   int    `json:"user2_id"`
	CreatedAt string `json:"created_at"`
}

type Message struct {
	ID        int    `json:"id"`
	ChatID    int    `json:"chat_id"`
	SenderID  int    `json:"sender_id"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

type RedisMessage struct {
	SenderID  int    `json:"sender_id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type WebSocketMessage struct {
	Type      string `json:"type"` // "message", "typing", "online"
	ChatID    int    `json:"chat_id"`
	SenderID  int    `json:"sender_id"`
	Message   string `json:"message,omitempty"`
	Timestamp string `json:"timestamp"`
}
