package utils

import (
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
