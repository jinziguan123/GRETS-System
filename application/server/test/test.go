package main

import (
	"grets_server/pkg/did"
	"log"
)

// ... (Paste your HexToPublicKey and VerifySignature functions here) ...

func main() {
	publicKeyHex := "047c1383d39f5e15df8cf58c889ca6ea616c9bbd61edb1778f048dd7219a16745898ed2c228642100ea256b77e049e2f2a5ef176e443bb1e8317d80af990717d8f"
	message := "did:grets:investor:ffa2a4575771ec5c:865a19c45930c02ed694880167f5f511:ec8ac7a9c9b0cad8aed2c73a24859c6b"                                 // Use the message that causes the issue
	signatureHex := "531326f50d0fddff37f869ad3206fb799b819eefe32f6c835a806330dce658a79b93161caa004ff370164fd2661898bacbeac1382341d2f4fd46fc477f3f84cd" // Use the signature that causes the issue

	log.Println("Starting verification...")
	isValid, err := did.VerifySignature(publicKeyHex, message, signatureHex)
	if err != nil {
		log.Fatalf("VerifySignature returned error: %v", err)
	}
	log.Printf("Verification result: %t\n", isValid)
}
