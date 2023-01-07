package utils

import (
	"encoding/hex"
	"math/big"
	"sync"

	"github.com/dop251/goja"
)

// RsaEncrypt using rsa encrypt message
func RsaEncrypt(exponent string, modules string, str string) string {
	r := getRuntime()
	var fn func(exponent string, modules string, str string) string
	r.ExportTo(r.Get("encryptString"), &fn)
	return fn(exponent, modules, str)
}

func RsaEncrypt2(exponent string, modules string, str string) string {
	modulus := new(big.Int)
	modulus.SetString(modules, 16)
	exp := new(big.Int)
	exp.SetString(exponent, 16)
	decimals := []byte(str)
	result := new(big.Int)
	for i, d := range decimals {
		acc := new(big.Int)
		acc.SetInt64(int64(d))
		acc.Lsh(acc, uint(i*8))
		result.Add(result, acc)
	}
	result.Exp(result, exp, modulus)

	return hex.EncodeToString(result.Bytes())
}

// ***it not work now***
// RsaDecrypt using rsa decrypt message
// func RsaDecrypt(exponent string, modules string, str string) string {
// 	r := getRuntime()
// 	var fn func(exponent string, modules string, str string) string
// 	r.ExportTo(r.Get("decryptString"), &fn)
// 	return fn(exponent, modules, str)
// }

var instance *goja.Runtime
var once sync.Once

func getRuntime() *goja.Runtime {
	once.Do(func() {
		r := goja.New()
		// f, err := os.Open("./utils/security.js")
		// if err != nil {
		// 	panic(err)
		// }
		// var buff bytes.Buffer
		// _, err = buff.ReadFrom(f)
		// if err != nil {
		// 	panic(err)
		// }
		r.RunString(securityJS)
		instance = r
	})
	return instance
}
