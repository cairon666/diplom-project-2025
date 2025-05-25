package models

type TelegramAuthData struct {
	ID        int64
	FirstName string
	LastName  string
	Username  string
	PhotoURL  string
	AuthDate  int
	Hash      string
}
