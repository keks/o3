/*Functions to covert packets from go structs to byte buffers.
 *These functions will only be called from packetdispatcher and
 *their task is only conversion (inversion of the parser).
 *Errors are passed up the chain in the form of panics and will
 *be converted to go errors further up the chain.
 */

package o3

import (
	"crypto/rand"
	"math/big"
)

func genPadding() ([]byte, error) {
	paddingValueBig, err := rand.Int(rand.Reader, big.NewInt(255))
	if err != nil {
		return nil, err
	}

	paddingValue := byte(paddingValueBig.Int64())
	padding := make([]byte, paddingValue)
	for i := range padding {
		padding[i] = paddingValue
	}

	return padding, nil
}

/*
func serializeMsgPkt(mp messagePacket) *bytes.Buffer {

	buf := new(bytes.Buffer)
	serializePktType(buf, mp.PktType)
	serializeIDString(buf, mp.Sender)
	serializeIDString(buf, mp.Recipient)
	serializeMsgID(buf, mp.ID)
	serializeTime(buf, mp.Time)
	serializeMsgFlags(buf, mp.Flags)
	// The three following bytes are unused
	serializeUnusedBytes(buf)
	serializePubNick(buf, mp.PubNick)
	serializeNonce(buf, mp.Nonce)
	serializeCiphertext(buf, mp.Ciphertext)

	return buf
}
*/

// func serializeGroupAudioMsg(gam GroupAudioMessage) *bytes.Buffer {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			panic(serializerPanicHandler("AudioMsg", r))
// 		}
// 	}()
// 	buf := new(bytes.Buffer)

// 	serializeMsgType(buf, AUDIOMESSAGE)

// 	// AudioClip duration, fixed value for now
// 	serializeUint16(buf, 0x0000)

// 	serializeBlobID(buf, gam.BlobID)
// 	serializeUint32(buf, gam.Size)
// 	serializeKey(buf, gam.Key)
// 	serializePadding(buf)
// 	return buf
// }

/*
func serializerPanicHandler(context string, i interface{}) error {
	if _, ok := i.(string); ok {
		return fmt.Errorf("%s: error occurred serializing %s", context, i)
	}
	return fmt.Errorf("%s: unknown serializing error occurred: %#v", context, i)
}

func serializeHelper(buf *bytes.Buffer, i interface{}) error {
	return errors.Wrapf(binary.Write(buf, binary.LittleEndian, i), "%T", i)
}

func contextualSerializeHelper(context string, buf *bytes.Buffer, i interface{}) (*bytes.Buffer, error) {
	err := binary.Write(buf, binary.LittleEndian, i)
	if err != nil {
		return errors.New(context)
	}
	return nil
}

// serializePadding returns a byte slice filled with n repetitions of the byte value n
func serializePadding(buf *bytes.Buffer) error {
	paddingValueBig, err := rand.Int(rand.Reader, big.NewInt(255))
	if err != nil {
		return err
	}

	paddingValue := byte(paddingValueBig.Int64())
	padding := make([]byte, paddingValue)
	for i := range padding {
		padding[i] = paddingValue
	}

	err = binary.Write(buf, binary.LittleEndian, padding)
	return errors.Wrapf(err, "error writing padding of length %d", paddingValue)
}

// TODO: clean this up!
func serializeUint8(num uint8, buf *bytes.Buffer) error {
	err := binary.Write(buf, binary.LittleEndian, num)
	return errors.Wrapf(err, "error writing uint8 %v", num)
}

func serializeUnusedBytes(buf *bytes.Buffer) *bytes.Buffer {
	return serializeArbitraryData(buf, []byte{0x00, 0x00, 0x00})
}

func serializeKey(buf *bytes.Buffer, key [32]byte) *bytes.Buffer {
	return contextualSerializeHelper("key", buf, key)
}

func serializeNoncePrefix(buf *bytes.Buffer, np [16]byte) *bytes.Buffer {
	return contextualSerializeHelper("nonce prefix", buf, np)
}

func serializeIDString(buf *bytes.Buffer, is IDString) *bytes.Buffer {
	return contextualSerializeHelper("id string", buf, is)
}

func serializePubNick(buf *bytes.Buffer, pn PubNick) *bytes.Buffer {
	return contextualSerializeHelper("public nickname", buf, pn)
}

func serializeMsgStatus(buf *bytes.Buffer, msgStatus MsgStatus) *bytes.Buffer {
	return serializeByte(buf, byte(msgStatus))
}

func serializeMsgID(buf *bytes.Buffer, id uint64) *bytes.Buffer {
	return serializeUint64(buf, id)
}

func serializeTime(buf *bytes.Buffer, t time.Time) *bytes.Buffer {
	//TODO time sanity checks
	return contextualSerializeHelper("time", buf, uint32(t.Unix()))
}

func serializeNonce(buf *bytes.Buffer, n nonce) *bytes.Buffer {
	return contextualSerializeHelper("nonce", buf, n.nonce)
}

func serializeCiphertext(buf *bytes.Buffer, bts []byte) *bytes.Buffer {
	//TODO error handling written bytes vs. len(bts)?
	return contextualSerializeHelper("ciphertext", buf, bts)
}

func serializeSysData(buf *bytes.Buffer, sysData [32]byte) *bytes.Buffer {
	return contextualSerializeHelper("system data", buf, sysData)
}

func serializeText(buf *bytes.Buffer, text string) *bytes.Buffer {
	// TODO: sanatize?
	return serializeHelper(buf, []byte(text))
}

func serializeBlobID(buf *bytes.Buffer, blobID [16]byte) *bytes.Buffer {
	return serializeHelper(buf, []byte(blobID[:]))
}

func serializeArbitraryData(buf *bytes.Buffer, i interface{}) *bytes.Buffer {
	//TODO type assertions for error handling?
	//TODO what to do about context?
	return serializeHelper(buf, i)
}
*/
