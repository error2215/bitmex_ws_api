package models

type WsMsg struct {
	Action  string   `json:"action"`
	Symbols []string `json:"symbols"`
}

type BitmexSymbol struct {
	Timestamp string          `json:"timestamp"`
	Symbol    string          `json:"symbol"`
	Price     float64         `json:"price"`
	Clients   []chan<- []byte `json:"-"`
}
