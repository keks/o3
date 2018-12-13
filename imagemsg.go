package o3

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
)

//ImageMessage represents an image message as sent e2e encrypted to other threema users
type ImageMessage struct {
	messageHeader
	imageMessageBody
}

type imageMessageBody struct {
	BlobID   [16]byte
	ServerID byte
	Size     uint32
	Nonce    nonce
}

// NewImageMessage returns a ImageMessage ready to be encrypted
func NewImageMessage(sc *SessionContext, recipient string, filename string) (ImageMessage, error) {
	recipientID := NewIDString(recipient)

	im := ImageMessage{
		messageHeader{
			sender:    sc.ID.ID,
			recipient: recipientID,
			id:        NewMsgID(),
			time:      time.Now(),
			pubNick:   sc.ID.Nick,
		},
		imageMessageBody{},
	}
	err := im.SetImageData(filename, *sc)
	if err != nil {
		return ImageMessage{}, err
	}
	return im, nil
}

// GetPrintableContent returns a printable represantion of a ImageMessage.
func (im ImageMessage) GetPrintableContent() string {
	return fmt.Sprintf("ImageMSG: https://%2x.blob.threema.ch/%16x, Size: %d, Nonce: %24x", im.ServerID, im.BlobID, im.Size, im.Nonce.nonce)
}

// GetImageData return the decrypted Image needs the recipients secret key
func (im ImageMessage) GetImageData(sc SessionContext) ([]byte, error) {
	return downloadAndDecryptAsym(sc, im.BlobID, im.Sender().String(), im.Nonce)
}

// SetImageData encrypts and uploads the image. Sets the blob info in the ImageMessage. Needs the recipients public key.
func (im *ImageMessage) SetImageData(filename string, sc SessionContext) error {
	plainImage, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.New("could not load image")
	}

	im.Nonce, im.ServerID, im.Size, im.BlobID, err = encryptAndUploadAsym(sc, plainImage, im.recipient.String())

	return err
}

//Serialize returns a fully serialized byte slice of an ImageMessage
func (im ImageMessage) MarshalBinary() ([]byte, error) {
	padding, err := genPadding()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, IMAGEMESSAGE)
	binary.Write(&buf, binary.LittleEndian, im.BlobID)
	binary.Write(&buf, binary.LittleEndian, im.Size)
	binary.Write(&buf, binary.LittleEndian, im.Nonce)
	binary.Write(&buf, binary.LittleEndian, padding)

	return buf.Bytes(), nil
}

//GroupImageMessage represents a group image message as sent e2e encrypted to other threema users
type GroupImageMessage struct {
	groupMessageHeader
	messageHeader
	groupImageMessageBody
}

//Serialize returns a fully serialized byte slice of a GroupImageMessage
func (im GroupImageMessage) MarshalBinary() ([]byte, error) {
	padding, err := genPadding()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, GROUPIMAGEMESSAGE)
	binary.Write(&buf, binary.LittleEndian, im.groupMessageHeader)
	binary.Write(&buf, binary.LittleEndian, im.BlobID)
	binary.Write(&buf, binary.LittleEndian, im.Size)
	binary.Write(&buf, binary.LittleEndian, im.Key)
	binary.Write(&buf, binary.LittleEndian, padding)

	return buf.Bytes(), nil
}

// GetImageData return the decrypted Image needs the recipients secret key
func (im GroupImageMessage) GetImageData(sc SessionContext) ([]byte, error) {
	return downloadAndDecryptSym(im.BlobID, im.Key)
}

// SetImageData encrypts the given image symmetrically and adds it to the message
func (im *GroupImageMessage) SetImageData(filename string) error {
	return im.groupImageMessageBody.setImageData(filename)
}

func (im *groupImageMessageBody) setImageData(filename string) error {
	plainImage, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.New("could not load image")
	}

	im.Key, im.ServerID, im.Size, im.BlobID, err = encryptAndUploadSym(plainImage)

	return err
}
