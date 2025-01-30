package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"log"
	"reflect"
	"time"
)

type EllipticCurve struct {
	curve   elliptic.Curve
	pubKey  *ecdsa.PublicKey
	privKey *ecdsa.PrivateKey
}

func NewEllipticCurve(curve elliptic.Curve) *EllipticCurve {
	return &EllipticCurve{
		curve:   curve,
		privKey: new(ecdsa.PrivateKey),
	}
}

func (ec *EllipticCurve) GenerateKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privKey, err := ecdsa.GenerateKey(ec.curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	ec.privKey = privKey
	ec.pubKey = &privKey.PublicKey

	return privKey, &privKey.PublicKey, nil
}

func (ec *EllipticCurve) EncodePrivateKey(privKey *ecdsa.PrivateKey) (string, error) {
	encoded, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return "", err
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: encoded})
	return string(pemEncoded), nil
}

func (ec *EllipticCurve) EncodePublicKey(pubKey *ecdsa.PublicKey) (string, error) {
	encoded, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", err
	}

	pemEncode := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: encoded})

	return string(pemEncode), nil
}

func (ec *EllipticCurve) DecodePrivateKey(privKey string) (*ecdsa.PrivateKey, error) {
	pemBlock, _ := pem.Decode([]byte(privKey))

	pemBytes := pemBlock.Bytes
	return x509.ParseECPrivateKey(pemBytes)
}

func (ec *EllipticCurve) DecodePublicKey(pubKey string) (*ecdsa.PublicKey, error) {
	pemBlock, _ := pem.Decode([]byte(pubKey))
	pemBytes := pemBlock.Bytes

	genPubKey, err := x509.ParsePKIXPublicKey(pemBytes)
	return genPubKey.(*ecdsa.PublicKey), err
}

func (ec *EllipticCurve) VerifySignature(msg string, privKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey) ([]byte, bool, error) {
	signHash := HashMessage(msg)

	r, s, _ := ecdsa.Sign(rand.Reader, privKey, signHash)

	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)

	ok := ecdsa.Verify(pubKey, signHash, r, s)
	if !ok {
		return nil, ok, fmt.Errorf("Does not match.")
	}

	return signature, ok, nil
}

func HashMessage(message string) []byte {
	h := md5.New()

	_, err := io.WriteString(h, message)

	if err != nil {
		panic("Error while hashing message")
	}
	return h.Sum(nil)
}

// https://auth0.com/docs/secure/tokens/json-web-tokens/json-web-token-claims
func MakeJWT(userID string, timeExpire time.Duration) (string, error) {
	ec := NewEllipticCurve(elliptic.P256())
	priv, _, _ := ec.GenerateKeys()
	signingKey := priv
	newToken := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"iss": "echoServer",
			"sub": "sub",
			"exp": time.Now().Add(time.Duration(timeExpire)),
		},
	)

	return newToken.SignedString(signingKey)
}

func (ec *EllipticCurve) Test(privKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey) error {

	encPriv, err := ec.EncodePrivateKey(privKey)
	if err != nil {
		return err
	}
	encPub, err := ec.EncodePublicKey(pubKey)
	if err != nil {
		return err
	}
	priv2, err := ec.DecodePrivateKey(encPriv)
	if err != nil {
		return err
	}
	pub2, err := ec.DecodePublicKey(encPub)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(privKey, priv2) {
		return errors.New("private keys do not match")

	}
	if !reflect.DeepEqual(pubKey, pub2) {
		return errors.New("public keys do not match")
	}
	fmt.Println("Everything correct!!")
	return nil
}

func maiin() {

	s, err := MakeJWT("12345678", time.Minute)
	if err != nil {
		log.Println(err)
	}
	log.Println(s)

	/*
		ec := NewEllipticCurve(elliptic.P256())

		priv, pub, _ := ec.GenerateKeys()
		fmt.Println(ec.EncodePrivateKey(priv))
		fmt.Println(ec.EncodePublicKey(pub))

		fmt.Println(ec.VerifySignature("jamones, Jamones", priv, pub))
	*/
}
