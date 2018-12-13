//
// Package o3 functions to prepare and send packets. All preparation required to transmit a
// packet takes place in the packet's respective dispatcher function. Functions
// from packetserializer are used to convert from struct to byte buffer form that
// can then be transmitted on the wire. Errors from packetserializer bubble up here
// in the form of panics that have to be passed on to communicationhandler for
// conversion to go errors.
//
package o3

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/box"
)

func dispatcherPanicHandler(context string, i interface{}) error {
	if _, ok := i.(string); ok {
		return fmt.Errorf("%s: error occurred dispatching %s", context, i)
	}
	return fmt.Errorf("%s: unknown dispatch error occurred: %#v", context, i)
}

func writeHelper(wr io.Writer, buf *bytes.Buffer) error {
	i, err := wr.Write(buf.Bytes())
	if i != buf.Len() {
		return errors.New("not enough bytes were transmitted")
	}

	return err
}

func (sc *SessionContext) dispatchClientHello(wr io.Writer) error {
	defer func() {
		if r := recover(); r != nil {
			panic(dispatcherPanicHandler("client hello", r))
		}
	}()
	var ch clientHelloPacket

	ch.ClientSPK = sc.clientSPK
	ch.NoncePrefix = sc.clientNonce.prefix()

	return binary.Write(wr, binary.LittleEndian, ch)
}

//not necessary on the client side
func (sc *SessionContext) dispatchServerHello(wr io.Writer) {}

func (sc *SessionContext) dispatchAuthMsg(wr io.Writer) error {
	var app authPacketPayload
	var ap authPacket

	app.Username = sc.ID.ID
	//TODO System Data: app.SysData = ..
	app.ServerNoncePrefix = sc.serverNonce.prefix()
	app.RandomNonce = newNonce()

	//create payload ciphertext
	ct := box.Seal(nil, sc.clientSPK[:], app.RandomNonce.bytes(), &sc.serverLPK, &sc.ID.LSK)
	if len(ct) != 48 {
		return errors.New("error encrypting client short-term public key")
	}
	copy(app.Ciphertext[:], ct[0:48])

	var appBuf bytes.Buffer
	err := binary.Write(&appBuf, binary.LittleEndian, app)
	if err != nil {
		return err
	}

	//create auth packet ciphertext
	sc.clientNonce.setCounter(1)
	apct := box.Seal(nil, appBuf.Bytes(), sc.clientNonce.bytes(), &sc.serverSPK, &sc.clientSSK)
	if len(apct) != 144 {
		panic("error encrypting payload")
	}
	copy(ap.Ciphertext[:], apct[0:144])

	return binary.Write(wr, binary.LittleEndian, ap)
}

func (sc *SessionContext) dispatchAckMsg(wr io.Writer, mp messagePacket) error {
	ackP := ackPacket{
		PktType:  clientAck,
		SenderID: mp.Sender,
		MsgID:    mp.ID}

	var serializedAckPkt bytes.Buffer
	err := binary.Write(&serializedAckPkt, binary.LittleEndian, ackP)
	if err != nil {
		return err
	}

	sc.clientNonce.increaseCounter()
	ackpCipherText := box.Seal(nil, serializedAckPkt.Bytes(), sc.clientNonce.bytes(), &sc.serverSPK, &sc.clientSSK)

	err = binary.Write(wr, binary.LittleEndian, uint16(len(ackpCipherText)))
	if err != nil {
		return err
	}

	return binary.Write(wr, binary.LittleEndian, ackpCipherText)
}

func (sc *SessionContext) dispatchEchoMsg(wr io.Writer, oldEchoPacket echoPacket) error {
	ep := echoPacket{
		Counter: oldEchoPacket.Counter + 1}
	var serializedEchoPkt bytes.Buffer
	err := binary.Write(&serializedEchoPkt, binary.LittleEndian, ep)
	if err != nil {
		return err
	}

	sc.clientNonce.increaseCounter()
	epCipherText := box.Seal(nil, serializedEchoPkt.Bytes(), sc.clientNonce.bytes(), &sc.serverSPK, &sc.clientSSK)

	err = binary.Write(wr, binary.LittleEndian, uint16(len(epCipherText)))
	if err != nil {
		return err
	}

	return binary.Write(wr, binary.LittleEndian, epCipherText)
}

func (sc *SessionContext) dispatchMessage(wr io.Writer, m Message) error {
	mh := m.header()

	randNonce := newRandomNonce()

	recipient, ok := sc.ID.Contacts.Get(mh.recipient.String())
	if !ok {
		var tr ThreemaRest
		var err error

		recipient, err = tr.GetContactByID(mh.recipient)
		if err != nil {
			return errors.New("Recipient's PublicKey could not be found!")
		}
		sc.ID.Contacts.Add(recipient)
	}

	msgData, err := m.MarshalBinary()
	if err != nil {
		return err
	}

	msgCipherText := box.Seal(nil, msgData, randNonce.bytes(), &recipient.LPK, &sc.ID.LSK)

	messagePkt := messagePacket{
		PktType:    sendingMsg,
		Sender:     mh.sender,
		Recipient:  mh.recipient,
		ID:         mh.id,
		Time:       mh.time,
		Flags:      msgFlags{PushMessage: true},
		PubNick:    mh.pubNick,
		Nonce:      randNonce,
		Ciphertext: msgCipherText,
	}

	var serializedMsgPkt bytes.Buffer
	err = binary.Write(&serializedMsgPkt, binary.LittleEndian, messagePkt)
	if err != nil {
		return err
	}

	sc.clientNonce.increaseCounter()
	serializedMsgPktCipherText := box.Seal(nil, serializedMsgPkt.Bytes(), sc.clientNonce.bytes(), &sc.serverSPK, &sc.clientSSK)

	err = binary.Write(wr, binary.LittleEndian, uint16(len(serializedMsgPktCipherText)))
	if err != nil {
		return err
	}

	return binary.Write(wr, binary.LittleEndian, serializedMsgPktCipherText)
}
