package o3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

// NewGroupMemberLeftMessages returns a slice of GroupMemberLeftMessages ready to be encrypted
func NewGroupMemberLeftMessages(sc *SessionContext, group Group) []GroupMemberLeftMessage {
	gml := make([]GroupMemberLeftMessage, len(group.Members))

	for i := 0; i < len(group.Members); i++ {
		gml[i] = GroupMemberLeftMessage{
			groupMessageHeader{
				creatorID: group.CreatorID,
				groupID:   group.GroupID},
			messageHeader{
				sender:    sc.ID.ID,
				recipient: group.Members[i],
				id:        NewMsgID(),
				time:      time.Now(),
				pubNick:   sc.ID.Nick}}

	}

	return gml

}

//Serialize returns a fully serialized byte slice of a GroupMemberLeftMessage
func (gml GroupMemberLeftMessage) MarshalBinary() ([]byte, error) {
	padding, err := genPadding()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, GROUPMEMBERLEFTMESSAGE)
	binary.Write(&buf, binary.LittleEndian, gml.groupMessageHeader)
	binary.Write(&buf, binary.LittleEndian, padding)

	return buf.Bytes(), nil
}

//GroupMemberLeftMessage represents a group leaving message
type GroupMemberLeftMessage struct {
	groupMessageHeader
	messageHeader
}

// NewDeliveryReceiptMessage returns a TextMessage ready to be encrypted
func NewDeliveryReceiptMessage(sc *SessionContext, recipient string, msgID uint64, msgStatus MsgStatus) (DeliveryReceiptMessage, error) {
	recipientID := NewIDString(recipient)

	dm := DeliveryReceiptMessage{
		messageHeader{
			sender:    sc.ID.ID,
			recipient: recipientID,
			id:        NewMsgID(),
			time:      time.Now(),
			pubNick:   sc.ID.Nick,
		},
		deliveryReceiptMessageBody{
			msgID:  msgID,
			status: msgStatus},
	}
	return dm, nil
}

type deliveryReceiptMessageBody struct {
	status MsgStatus
	msgID  uint64
}

// DeliveryReceiptMessage represents a delivery receipt as sent e2e encrypted to other threema users when a message has been received
type DeliveryReceiptMessage struct {
	messageHeader
	deliveryReceiptMessageBody
}

// GetPrintableContent returns a printable represantion of a DeliveryReceiptMessage.
func (dm DeliveryReceiptMessage) GetPrintableContent() string {
	return fmt.Sprintf("Delivered: %x", dm.msgID)
}

//Serialize returns a fully serialized byte slice of a SeliveryReceiptMessage
func (dm DeliveryReceiptMessage) MarshalBinary() ([]byte, error) {
	padding, err := genPadding()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, DELIVERYRECEIPT)
	binary.Write(&buf, binary.LittleEndian, dm.status)
	binary.Write(&buf, binary.LittleEndian, dm.msgID)
	binary.Write(&buf, binary.LittleEndian, padding)

	return buf.Bytes(), nil
}

// Status returns the messages status
func (dm DeliveryReceiptMessage) Status() MsgStatus {
	return dm.status
}

// MsgID returns the message id
func (dm DeliveryReceiptMessage) MsgID() uint64 {
	return dm.msgID
}

// GROUP MANAGEMENT MESSAGES
////////////////////////////////////////////////////////////////
// TODO: Implement message interface
type groupManageMessageHeader struct {
	groupID [8]byte
}

func (gmh groupManageMessageHeader) GroupID() [8]byte {
	return gmh.groupID
}

// NewGroupManageSetMembersMessages returns a slice of GroupManageSetMembersMessages ready to be encrypted
func NewGroupManageSetMembersMessages(sc *SessionContext, group Group) []GroupManageSetMembersMessage {
	gms := make([]GroupManageSetMembersMessage, len(group.Members))

	for i := 0; i < len(group.Members); i++ {
		gms[i] = GroupManageSetMembersMessage{
			groupManageMessageHeader{
				groupID: group.GroupID},
			messageHeader{
				sender:    sc.ID.ID,
				recipient: group.Members[i],
				id:        NewMsgID(),
				time:      time.Now(),
				pubNick:   sc.ID.Nick},
			groupManageSetMembersMessageBody{
				groupMembers: group.Members}}

	}

	return gms

}

type groupManageSetMembersMessageBody struct {
	groupMembers []IDString
}

// GroupManageSetImageMessage represents the message sent e2e-encrypted by a group's creator to all members to set the group image
type GroupManageSetImageMessage struct {
	groupManageMessageHeader
	messageHeader
	groupImageMessageBody
}

// NewGroupManageSetImageMessages returns a slice of GroupManageSetImageMessages ready to be encrypted
func NewGroupManageSetImageMessages(sc *SessionContext, group Group, filename string) []GroupManageSetImageMessage {
	gms := make([]GroupManageSetImageMessage, len(group.Members))

	for i := 0; i < len(group.Members); i++ {
		gms[i] = GroupManageSetImageMessage{
			groupManageMessageHeader{
				groupID: group.GroupID},
			messageHeader{}, //TODO:
			groupImageMessageBody{},
		}

		err := gms[i].SetImageData(filename)
		if err != nil {
			//TODO: pretty sure this isn't a good idea
			return nil
		}
	}

	return gms
}

// GetImageData returns the decrypted Image
func (im GroupManageSetImageMessage) GetImageData(sc SessionContext) ([]byte, error) {
	return downloadAndDecryptSym(im.BlobID, im.Key)
}

// SetImageData encrypts the given image symmetrically and adds it to the message
func (im *GroupManageSetImageMessage) SetImageData(filename string) error {
	return im.groupImageMessageBody.setImageData(filename)
}

//Serialize returns a fully serialized byte slice of an ImageMessage
func (im GroupManageSetImageMessage) MarshalBinary() ([]byte, error) {
	padding, err := genPadding()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, GROUPSETIMAGEMESSAGE)
	binary.Write(&buf, binary.LittleEndian, im.GroupID())
	binary.Write(&buf, binary.LittleEndian, im.BlobID)
	binary.Write(&buf, binary.LittleEndian, im.Size)
	binary.Write(&buf, binary.LittleEndian, im.Key)
	binary.Write(&buf, binary.LittleEndian, padding)

	return buf.Bytes(), nil
}

// GroupManageSetMembersMessage represents the message sent e2e encrypted by a group's creator to all members
type GroupManageSetMembersMessage struct {
	groupManageMessageHeader
	messageHeader
	groupManageSetMembersMessageBody
}

//Members returns a byte slice of IDString of all members contained in the message
func (gmm GroupManageSetMembersMessage) Members() []IDString {
	return gmm.groupMembers
}

//Serialize returns a fully serialized byte slice of a GroupManageSetMembersMessage
func (gmm GroupManageSetMembersMessage) MarshalBinary() ([]byte, error) {
	padding, err := genPadding()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, GROUPSETMEMEBERSMESSAGE)
	binary.Write(&buf, binary.LittleEndian, gmm.GroupID())

	for _, member := range gmm.Members() {
		binary.Write(&buf, binary.LittleEndian, member)
	}

	binary.Write(&buf, binary.LittleEndian, padding)

	return buf.Bytes(), nil
}

// NewGroupManageSetNameMessages returns a slice of GroupMenageSetNameMessages ready to be encrypted
func NewGroupManageSetNameMessages(sc *SessionContext, group Group) []GroupManageSetNameMessage {
	gms := make([]GroupManageSetNameMessage, len(group.Members))

	for i := 0; i < len(group.Members); i++ {
		gms[i] = GroupManageSetNameMessage{
			groupManageMessageHeader{
				groupID: group.GroupID},
			messageHeader{
				sender:    sc.ID.ID,
				recipient: group.Members[i],
				id:        NewMsgID(),
				time:      time.Now(),
				pubNick:   sc.ID.Nick},
			groupManageSetNameMessageBody{
				groupName: group.Name}}

	}

	return gms

}

type groupManageSetNameMessageBody struct {
	groupName string
}

func (gmm groupManageSetNameMessageBody) Name() string {
	return gmm.groupName
}

//Serialize returns a fully serialized byte slice of a GroupManageSetNameMessage
func (gmm GroupManageSetNameMessage) MarshalBinary() ([]byte, error) {
	padding, err := genPadding()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, GROUPSETNAMEMESSAGE)
	binary.Write(&buf, binary.LittleEndian, gmm.GroupID())
	binary.Write(&buf, binary.LittleEndian, gmm.Name())
	binary.Write(&buf, binary.LittleEndian, padding)

	return buf.Bytes(), nil
}

//GroupManageSetNameMessage represents a group management messate to set the group name
type GroupManageSetNameMessage struct {
	groupManageMessageHeader
	messageHeader
	groupManageSetNameMessageBody
}
