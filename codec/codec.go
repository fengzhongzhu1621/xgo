package codec

// Codec defines the interface of business communication protocol,
// which contains head and body. It only parses the body in binary,
// and then the business body struct will be handled by serializer.
// In common, the body's protocol is pb, json, etc. Specially,
// we can register our own serializer to handle other body type.
type Codec interface {
	// Encode pack the body into binary buffer.
	// client: Encode(msg, reqBody)(request-buffer, err)
	// server: Encode(msg, rspBody)(response-buffer, err)
	Encode(message Msg, body []byte) (buffer []byte, err error)

	// Decode unpack the body from binary buffer
	// server: Decode(msg, request-buffer)(reqBody, err)
	// client: Decode(msg, response-buffer)(rspBody, err)
	Decode(message Msg, buffer []byte) (body []byte, err error)
}
