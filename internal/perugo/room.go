package perugo

type Room struct {
	players []Player
	Id      string `json:"id"`
	Owner   Player `json:"-"`
}
