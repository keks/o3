package o3

//TypingNotificationMessage represents a typing notifiaction message
type TypingNotificationMessage struct {
	messageHeader
	typingNotificationBody
}

type typingNotificationBody struct {
	OnOff byte
}

//Serialize returns a fully serialized byte slice of a TypingNotificationMessage
func (tn TypingNotificationMessage) MarshalBinary() ([]byte, error) {
	return []byte{tn.OnOff}, nil
}
