package coding

import "github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"

type MessageCoderInterface interface {
	// Decode decodes a bytes json message object to an actual transport.Message object.
	// Returns error if there are missing values for pointers.
	//But still returns the decoded message.
	Decode(message []byte) (transport.Message, error)
	Encode(message transport.Message) ([]byte, error)
}
