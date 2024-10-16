package signature

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/eolinker/go-common/utils"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const HASH = crypto.SHA256

func Sign(privateKey []byte, key string, subject map[string]any) (map[string]any, []byte, error) {
	pkcs1PrivateKey, err := x509.ParsePKCS1PrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	m, data := format(key, subject)
	hash := crypto.SHA256.New()
	hash.Write(data)
	sign, err := rsa.SignPKCS1v15(rand.Reader, pkcs1PrivateKey, crypto.SHA256, hash.Sum(nil))
	if err != nil {
		return nil, nil, err
	}
	return m, sign, nil
}
func Verify(secret string, pemData []byte, publicKey []byte) (map[string]interface{}, bool) {
	rsaPublicKey, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, false
	}
	subjects, sign, err := DecodePem(secret, pemData)
	if err != nil {
		return nil, false
	}
	hash, err := readHash(subjects)
	if err != nil {
		return nil, false
	}

	data := subjectEncode(subjects)
	h := hash.New()
	h.Write(data)
	err = rsa.VerifyPKCS1v15(rsaPublicKey.(*rsa.PublicKey), hash, h.Sum(nil), sign)
	if err != nil {
		return nil, false
	}
	return subjects, true
}
func readHash(subjects map[string]any) (crypto.Hash, error) {
	if subjects == nil {
		return 0, errors.New("invalid subject")
	}
	v, has := subjects["hash"]
	if !has {
		return 0, errors.New("invalid subject")
	}
	switch value := v.(type) {
	case string:
		hash, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		}
		return crypto.Hash(hash), nil
	default:
		hash, err := strconv.Atoi(fmt.Sprint(value))
		if err != nil {
			return 0, err
		}
		return crypto.Hash(hash), nil
	}
	//return 0, errors.New("invalid subject")
}
func format(key string, subject map[string]any) (map[string]any, []byte) {
	subjectFormat := formatForSign(key, subject)

	return subjectFormat, subjectEncode(subjectFormat)
}
func subjectEncode(subject map[string]any) []byte {
	list := utils.MapToSlice(subject, func(k string, v any) *_KV {
		return &_KV{
			Name:  k,
			Value: v,
		}
	})
	sort.Sort(_KVS(list))
	data := utils.SliceToSlice(list, func(s *_KV) string {
		return url.QueryEscape(fmt.Sprint(s.Value))
	})
	return []byte(strings.Join(data, "&"))
}

func formatForSign(key string, subject map[string]any) map[string]any {
	mp := utils.CopyMaps(subject)
	mp["id"] = key
	mp["hash"] = strconv.Itoa(int(HASH))
	mp["sign_time"] = time.Now().Format(time.RFC3339)
	return mp
}

type _KV struct {
	Name  string
	Value any
}
type _KVS []*_KV

func (kvs _KVS) Len() int {
	return len(kvs)
}

func (kvs _KVS) Less(i, j int) bool {
	return kvs[i].Name < kvs[j].Name
}

func (kvs _KVS) Swap(i, j int) {
	kvs[i], kvs[j] = kvs[j], kvs[i]
}
