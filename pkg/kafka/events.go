package kafka

type Event struct {
	Entry string `json:"entry"`
	Exit  string `json:"exit"`
	ID    uint64 `json:"id"`
}

type EventDocument struct {
	Event Event `json:"event"`
}

type Response struct {
	Events []*EventDocument
}
