package did

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// KeyPair 密钥对
type KeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

// GenerateKeyPair 生成ECDSA密钥对
func GenerateKeyPair() (*KeyPair, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("生成密钥对失败: %v", err)
	}

	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

// PublicKeyToHex 将公钥转换为十六进制字符串
func (kp *KeyPair) PublicKeyToHex() string {
	x := kp.PublicKey.X.Bytes()
	y := kp.PublicKey.Y.Bytes()

	// 确保x和y都是32字节
	xBytes := make([]byte, 32)
	yBytes := make([]byte, 32)
	copy(xBytes[32-len(x):], x)
	copy(yBytes[32-len(y):], y)

	// 04前缀表示未压缩的公钥
	publicKeyBytes := append([]byte{0x04}, append(xBytes, yBytes...)...)
	return hex.EncodeToString(publicKeyBytes)
}

// PrivateKeyToHex 将私钥转换为十六进制字符串
func (kp *KeyPair) PrivateKeyToHex() string {
	privateKeyBytes := kp.PrivateKey.D.Bytes()
	// 确保私钥是32字节
	keyBytes := make([]byte, 32)
	copy(keyBytes[32-len(privateKeyBytes):], privateKeyBytes)
	return hex.EncodeToString(keyBytes)
}

// HexToPublicKey 从十六进制字符串恢复公钥
func HexToPublicKey(hexStr string) (*ecdsa.PublicKey, error) {
	keyBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("解码公钥失败: %v", err)
	}

	if len(keyBytes) != 65 || keyBytes[0] != 0x04 {
		return nil, fmt.Errorf("无效的公钥格式")
	}

	x := new(big.Int).SetBytes(keyBytes[1:33])
	y := new(big.Int).SetBytes(keyBytes[33:65])

	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}, nil
}

// HexToPrivateKey 从十六进制字符串恢复私钥
func HexToPrivateKey(hexStr string) (*ecdsa.PrivateKey, error) {
	keyBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("解码私钥失败: %v", err)
	}

	d := new(big.Int).SetBytes(keyBytes)
	privateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
		},
		D: d,
	}

	// 计算公钥
	privateKey.PublicKey.X, privateKey.PublicKey.Y = elliptic.P256().ScalarBaseMult(keyBytes)

	return privateKey, nil
}

// SignMessage 使用私钥签名消息
func (kp *KeyPair) SignMessage(message []byte) (string, error) {
	hash := sha256.Sum256(message)
	r, s, err := ecdsa.Sign(rand.Reader, kp.PrivateKey, hash[:])
	if err != nil {
		return "", fmt.Errorf("签名失败: %v", err)
	}

	// 将r和s转换为字节并连接
	rBytes := r.Bytes()
	sBytes := s.Bytes()

	// 确保r和s都是32字节
	rPadded := make([]byte, 32)
	sPadded := make([]byte, 32)
	copy(rPadded[32-len(rBytes):], rBytes)
	copy(sPadded[32-len(sBytes):], sBytes)

	signature := append(rPadded, sPadded...)
	return hex.EncodeToString(signature), nil
}

// VerifySignature 验证签名
func VerifySignature(publicKeyHex, message, signatureHex string) (bool, error) {
	// 恢复公钥
	publicKey, err := HexToPublicKey(publicKeyHex)
	if err != nil {
		return false, fmt.Errorf("恢复公钥失败: %v", err)
	}

	// 解码签名
	signatureBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, fmt.Errorf("解码签名失败: %v", err)
	}

	if len(signatureBytes) != 64 {
		return false, fmt.Errorf("无效的签名长度")
	}

	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:])

	// 计算消息哈希
	hash := sha256.Sum256([]byte(message))

	// 验证签名
	return ecdsa.Verify(publicKey, hash[:], r, s), nil
}

// GenerateHash 生成SHA256哈希
func GenerateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
