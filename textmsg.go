package o3

import (
	"bytes"
	"encoding/binary"
	"time"
)

//TextMessage represents a text message as sent e2e encrypted to other threema users
type TextMessage struct {
	messageHeader
	textMessageBody
}

type textMessageBody struct {
	text string
}

// NewTextMessage returns a TextMessage ready to be encrypted
func NewTextMessage(sc *SessionContext, recipient string, text string) (TextMessage, error) {
	recipientID := NewIDString(recipient)

	tm := TextMessage{
		messageHeader{
			sender:    sc.ID.ID,
			recipient: recipientID,
			id:        NewMsgID(),
			time:      time.Now(),
			pubNick:   sc.ID.Nick,
		},
		textMessageBody{text: text},
	}
	return tm, nil
}

// Text returns the message text
func (tm TextMessage) Text() string {
	return tm.text
}

// String returns the message text as string
func (tm TextMessage) String() string {
	return tm.Text()
}

//Serialize returns a fully serialized byte slice of a TextMessage
func (tm TextMessage) MarshalBinary() ([]byte, error) {
	padding, err := genPadding()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, TEXTMESSAGE)
	binary.Write(&buf, binary.LittleEndian, tm.text)
	binary.Write(&buf, binary.LittleEndian, padding)

	return buf.Bytes(), nil
}

// NewGroupTextMessages returns a slice of GroupMemberTextMessages ready to be encrypted
func NewGroupTextMessages(sc *SessionContext, group Group, text string) ([]GroupTextMessage, error) {
	gtm := make([]GroupTextMessage, len(group.Members))
	var tm TextMessage
	var err error

	for i, member := range group.Members {
		tm, err = NewTextMessage(sc, member.String(), text)
		if err != nil {
			return []GroupTextMessage{}, err
		}

		gtm[i] = GroupTextMessage{
			groupMessageHeader{
				creatorID: group.CreatorID,
				groupID:   group.GroupID},
			tm}
	}

	return gtm, nil

}

//GroupTextMessage represents a group text message as sent e2e encrypted to other threema users
type GroupTextMessage struct {
	groupMessageHeader
	TextMessage
}

// Serialize : returns byte representation of serialized group text message
func (gtm GroupTextMessage) MarshalBinary() ([]byte, error) {
	padding, err := genPadding()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, GROUPTEXTMESSAGE)
	binary.Write(&buf, binary.LittleEndian, gtm.groupMessageHeader)
	binary.Write(&buf, binary.LittleEndian, gtm.text)
	binary.Write(&buf, binary.LittleEndian, padding)

	return buf.Bytes(), nil
}
