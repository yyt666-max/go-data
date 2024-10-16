package signature

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

// GenerateRSAKey 构造密钥对
func GenerateRSAKey() (x509PrivateKey []byte, X509PublicKey []byte, err error) {
	privateKey, errGen := rsa.GenerateKey(rand.Reader, 2048)
	if errGen != nil {
		err = errGen
		return
	}
	// 通过x509标准将得到的ras私钥序列化为 ASN.1 的 DER 编码字符串
	x509PrivateKey = x509.MarshalPKCS1PrivateKey(privateKey)
	// X509对公钥编码
	X509PublicKey, err = x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return
	}
	return
}
