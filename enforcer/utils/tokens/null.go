package tokens

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// NullPKI holds all PKI information
type NullPKI struct {
	PrivateKeyPEM []byte
	PublicKeyPEM  []byte
	AuthorityPEM  []byte
}

// NewNullPKI creates new secrets for PKI implementation based on compact encoding
func NewNullPKI(keyPEM, certPEM, caPEM []byte) (*NullPKI, error) {

	fmt.Println("I was called with null  ")

	p := &NullPKI{}
	fmt.Println("I finished with the keys   ")
	return p, nil
}

// Type implements the interface Secrets
func (p *NullPKI) Type() SecretsType {
	return PKINull
}

// EncodingKey returns the private key
func (p *NullPKI) EncodingKey() interface{} {
	return jwt.UnsafeAllowNoneSignatureType
}

// DecodingKey returns the public key
func (p *NullPKI) DecodingKey(server string, ackKey interface{}, prevKey interface{}) (interface{}, error) {

	return jwt.UnsafeAllowNoneSignatureType, nil
}

// VerifyPublicKey verifies if the inband public key is correct.
func (p *NullPKI) VerifyPublicKey(pkey []byte) (interface{}, error) {

	return jwt.UnsafeAllowNoneSignatureType, nil

}

// TransmittedKey returns the PEM of the public key in the case of PKI
// if there is no certificate cache configured
func (p *NullPKI) TransmittedKey() []byte {
	return []byte("none")
}

// AckSize returns the default size of an ACK packet
func (p *NullPKI) AckSize() uint32 {
	return uint32(235)
}

// AuthPEM returns the Certificate Authority PEM
func (p *NullPKI) AuthPEM() []byte {
	return p.AuthorityPEM
}

// TransmittedPEM returns the PEM certificate that is transmitted
func (p *NullPKI) TransmittedPEM() []byte {
	return p.PublicKeyPEM
}

// EncodingPEM returns the certificate PEM that is used for encoding
func (p *NullPKI) EncodingPEM() []byte {
	return p.PrivateKeyPEM
}
