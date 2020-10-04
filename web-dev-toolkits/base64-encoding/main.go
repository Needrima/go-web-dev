//Base63 enoding can be used to encode strings and just like hashes, the output
//will always be the same so fargo fmt the input is same
package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	//to encode
	//var StdEncoding = NewEncoding(encodeStd)
	encstd := base64.StdEncoding // encoding standard
	//func (enc *Encoding) EncodeToString(src []byte) string
	encodedstring := encstd.EncodeToString([]byte("ourstring")) //
	fmt.Println(encodedstring)

	//encodedstring2 := encstd.EncodeToString([]byte("ourstring"))
	//fmt.Println(encodedstring2)

	//to decode
	//func (enc *Encoding) DecodeString(s string) ([]byte, error)
	bs, err := encstd.DecodeString(encodedstring)
	if err != nil {
		fmt.Println(err)
		return
	}
	str := string(bs)
	fmt.Println(str)
}
