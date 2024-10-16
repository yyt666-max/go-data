package signature

import (
	"fmt"
	"github.com/google/uuid"
	_ "testing"
)

func ExampleSign() {
	privateKey, publicKey, err := GenerateRSAKey()
	if err != nil {
		fmt.Println(err)
		return
	}
	key := uuid.NewString()
	subjectDecode := map[string]any{
		"company": "eolink",
		"edition": "ultimate(旗舰版）",
		"begin":   "2023-01-01",
		"end":     "2024-12-31",
		"cluster": 5,
		"control": 3,
		"node":    0,
		"code":    "xxff",
	}
	subjectSign, sign, err := Sign(privateKey, key, subjectDecode)
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := EncodeSign("apinto", subjectSign, sign)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, ok := Verify("apinto", data, publicKey)
	fmt.Println("Verify :", ok)

	//output:Verify : true
}
