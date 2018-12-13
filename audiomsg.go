package o3

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
)

//AudioMessage represents an image message as sent e2e encrypted to other threema users
type AudioMessage struct {
	messageHeader
	audioMessageBody
}

type audioMessageBody struct {
	Duration uint16 // The audio clips duration in seconds
	BlobID   [16]byte
	ServerID byte
	Size     uint32
	Key      [32]byte
}

// NewAudioMessage returns a ImageMessage ready to be encrypted
func NewAudioMessage(sc *SessionContext, recipient string, filename string) (AudioMessage, error) {
	recipientID := NewIDString(recipient)

	im := AudioMessage{
		messageHeader{
			sender:    sc.ID.ID,
			recipient: recipientID,
			id:        NewMsgID(),
			time:      time.Now(),
			pubNick:   sc.ID.Nick,
		},
		audioMessageBody{},
	}
	err := im.SetAudioData(filename, *sc)
	if err != nil {
		return AudioMessage{}, err
	}
	return im, nil
}

// GetPrintableContent returns a printable represantion of an AudioMessage
func (am AudioMessage) GetPrintableContent() string {
	return fmt.Sprintf("AudioMSG: https://%2x.blob.threema.ch/%16x, Size: %d, Nonce: %24x", am.ServerID, am.BlobID, am.Size, am.Key)
}

// GetAudioData return the decrypted audio, needs the recipients secret key
func (am AudioMessage) GetAudioData(sc SessionContext) ([]byte, error) {
	return downloadAndDecryptSym(am.BlobID, am.Key)
}

// SetAudioData encrypts and uploads the audio. Sets the blob info in the ImageMessage. Needs the recipients public key.
func (am *AudioMessage) SetAudioData(filename string, sc SessionContext) error {
	plainAudio, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.New("could not load audio")
	}

	// TODO: Should we have a whole media lib as dependency just to set this to the proper value?
	am.Duration = 0xFF

	am.Key, am.ServerID, am.Size, am.BlobID, err = encryptAndUploadSym(plainAudio)

	return err
}

//Serialize returns a fully serialized byte slice of an AudioMessage
func (am AudioMessage) MarshalBinary() ([]byte, error) {
	padding, err := genPadding()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, AUDIOMESSAGE)
	// AudioClip duration
	binary.Write(&buf, binary.LittleEndian, uint16(0xffff))
	binary.Write(&buf, binary.LittleEndian, am.BlobID)
	binary.Write(&buf, binary.LittleEndian, am.Size)
	binary.Write(&buf, binary.LittleEndian, am.Key)
	binary.Write(&buf, binary.LittleEndian, padding)

	return buf.Bytes(), nil
}
