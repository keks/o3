package o3

import (
	"bytes"
	"encoding"
	"encoding/binary"
	mrand "math/rand"
	"time"
)

// MsgType determines the type of message that is sent or received. Users usually
// won't use this directly and rather use message generator functions.
type MsgType uint8

// MsgType mock enum
const (
	TEXTMESSAGE             MsgType = 0x1  //indicates a text message
	IMAGEMESSAGE            MsgType = 0x2  //indicates a image message
	AUDIOMESSAGE            MsgType = 0x14 //indicates a audio message
	POLLMESSAGE             MsgType = 0x15 //indicates a poll message
	LOCATIONMESSAGE         MsgType = 0x16 //indicates a location message
	FILEMESSAGE             MsgType = 0x17 //indicates a file message
	GROUPTEXTMESSAGE        MsgType = 0x41 //indicates a group text message
	GROUPIMAGEMESSAGE       MsgType = 0x43 //indicates a group image message
	GROUPSETMEMEBERSMESSAGE MsgType = 0x4A //indicates a set group member message
	GROUPSETNAMEMESSAGE     MsgType = 0x4B //indicates a set group name message
	GROUPMEMBERLEFTMESSAGE  MsgType = 0x4C //indicates a group member left message
	GROUPSETIMAGEMESSAGE    MsgType = 0x50 //indicates a group set image message
	DELIVERYRECEIPT         MsgType = 0x80 //indicates a delivery receipt sent by the threema servers
	TYPINGNOTIFICATION      MsgType = 0x90 //indicates a typing notifiaction message
	//GROUPSETIMAGEMESSAGE msgType = 76
)

// MsgStatus represents the single-byte status field of DeliveryReceiptMessage
type MsgStatus uint8

//MsgStatus mock enum
const (
	MSGDELIVERED   MsgStatus = 0x1 //indicates message was received by peer
	MSGREAD        MsgStatus = 0x2 //indicates message was read by peer
	MSGAPPROVED    MsgStatus = 0x3 //indicates message was approved (thumb up) by peer
	MSGDISAPPROVED MsgStatus = 0x4 //indicates message was disapproved (thumb down) by peer
)

//TODO: figure these out
type msgFlags struct {
	PushMessage                    bool
	NoQueuing                      bool
	NoAckExpected                  bool
	MessageHasAlreadyBeenDelivered bool
	GroupMessage                   bool
}

func (flags msgFlags) MarshalBinary() ([]byte, error) {
	var flagsByte byte

	if flags.PushMessage {
		flagsByte |= (1 << 0)
	}
	if flags.NoQueuing {
		flagsByte |= (1 << 1)
	}
	if flags.NoAckExpected {
		flagsByte |= (1 << 2)
	}
	if flags.MessageHasAlreadyBeenDelivered {
		flagsByte |= (1 << 3)
	}
	if flags.GroupMessage {
		flagsByte |= (1 << 4)
	}

	return []byte{flagsByte}, nil
}

// NewMsgID returns a randomly generated message ID (not cryptographically secure!)
// TODO: Why mrand?
func NewMsgID() uint64 {
	mrand.Seed(int64(time.Now().Nanosecond()))
	msgID := uint64(mrand.Int63())
	return msgID
}

// NewGrpID returns a randomly generated group ID (not cryptographically secure!)
// TODO: Why mrand?
func NewGrpID() [8]byte {
	mrand.Seed(int64(time.Now().Nanosecond()))

	grpIDbuf := make([]byte, 8)
	mrand.Read(grpIDbuf)

	var grpID [8]byte
	copy(grpID[:], grpIDbuf)

	return grpID
}

// Message representing the various kinds of e2e ecrypted messages threema supports
type Message interface {
	encoding.BinaryMarshaler

	//Sender returns the message's sender ID
	Sender() IDString

	header() messageHeader
}

type messageHeader struct {
	sender    IDString
	recipient IDString
	id        uint64
	time      time.Time
	pubNick   PubNick
}

func (mh messageHeader) Sender() IDString {
	return mh.sender
}

func (mh messageHeader) Recipient() IDString {
	return mh.recipient
}

func (mh messageHeader) ID() uint64 {
	return mh.id
}

func (mh messageHeader) Time() time.Time {
	return mh.time
}

func (mh messageHeader) PubNick() PubNick {
	return mh.pubNick
}

//TODO: WAT?
func (mh messageHeader) header() messageHeader {
	return mh
}

//--------8<--------8<--------8<--------8<--------8<--------8<--------8<--------8<--------8<----

//--------8<--------8<--------8<--------8<--------8<--------8<--------8<--------8<--------8<----
//--------8<--------8<--------8<--------8<--------8<--------8<--------8<--------8<--------8<----

//--------8<--------8<--------8<--------8<--------8<--------8<--------8<--------8<--------8<----

//--------8<--------8<--------8<--------8<--------8<--------8<--------8<--------8<--------8<----

type groupMessageHeader struct {
	creatorID IDString
	groupID   [8]byte
}

func (gh groupMessageHeader) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, gh.creatorID)
	binary.Write(&buf, binary.LittleEndian, gh.groupID)

	return buf.Bytes(), nil
}

type groupImageMessageBody struct {
	BlobID   [16]byte
	ServerID byte
	Size     uint32
	Key      [32]byte
}

// GroupCreator returns the ID of the groups admin/creator as string
func (gmh groupMessageHeader) GroupCreator() IDString {
	return gmh.creatorID
}

// GroupID returns the ID of the group the message belongs to
func (gmh groupMessageHeader) GroupID() [8]byte {
	return gmh.groupID
}
