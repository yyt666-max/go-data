package signature

import (
	"bytes"
	"encoding/pem"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

var (
	ErrorInavlidCertFile = errors.New("invalid certificate")
)

func EncodeSign(secret string, subject map[string]any, sign []byte) ([]byte, error) {
	data, err := yaml.Marshal(subject)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:    fmt.Sprintf("%s AUTHORITY", strings.ToUpper(secret)),
		Headers: map[string]string{},
		Bytes:   data,
	})

	if err != nil {
		return nil, err
	}
	err = pem.Encode(buf, &pem.Block{
		Type:    fmt.Sprintf("%s SIGNATURE", strings.ToUpper(secret)),
		Headers: map[string]string{},
		Bytes:   sign,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func EncodePublicKey(secret string, publicKey []byte) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := pem.Encode(buf, &pem.Block{
		Type:    fmt.Sprintf("%s PUBLIC KEY", strings.ToUpper(secret)),
		Headers: map[string]string{},
		Bytes:   publicKey,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}
func DecodePublic(secret string, data []byte) ([]byte, error) {
	p, _ := pem.Decode(data)
	return p.Bytes, nil
}
func DecodePem(secret string, data []byte) (subject map[string]any, sign []byte, err error) {
	next := data
	authorityType := fmt.Sprintf("%s AUTHORITY", strings.ToUpper(secret))
	signatureType := fmt.Sprintf("%s SIGNATURE", strings.ToUpper(secret))
	for {
		if len(next) == 0 {
			break
		}
		block, rest := pem.Decode(next)
		if block == nil {
			return nil, nil, ErrorInavlidCertFile
		}
		switch block.Type {
		case signatureType:
			sign = block.Bytes
		case authorityType:
			err := yaml.Unmarshal(block.Bytes, &subject)
			if err != nil {
				return nil, nil, err
			}
		}

		next = rest
	}
	return

}
