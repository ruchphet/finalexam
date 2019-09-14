package entity

//ErrorMessage for response message
type ErrorMessage struct {
	Message string `json:"message"`
}

func (eMessage ErrorMessage) Error() string {
	return eMessage.Message
}
