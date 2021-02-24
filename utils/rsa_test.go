package utils

import (
	"fmt"
	"testing"
)

func TestRsaEncrypt(t *testing.T) {
	eStr := RsaEncrypt(
		"010001",
		"00bf857a6a1e652bde0876c9a3aa75ba2384d03d11c9b0c3d2031794ff7e23e88119da7a16ee6c6206bd39c8b4710b2fd0fa58c3adf3a68428c9dfe210241bd5bcc7d273e46387e47a5b6b04d34969ddb567581dc63c45aacfa0e14277b1d9566bf44981eafb6c973c157c08b961e0e7b23b8a712f339f230a8f076fb293facf45",
		"EcardStr0xx",
	)
	fmt.Println(eStr)
}
