//hmacs are used to encrypt user input and also validate them
//output will always be the same as long as input is the same
/*can be used to store values on users machine
for example on cookies to make sure the values are not tampered with*/

package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
)

func HashString(s string) string {
	hash := hmac.New(sha256.New, []byte("hashkey"))
	io.WriteString(hash, s)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func main() {
	input1 := "Input1"
	output1 := HashString(input1)
	fmt.Println(output1)

	input2 := "Input2"
	output2 := HashString(input2)
	fmt.Println(output2)
	fmt.Println("----------")

	fmt.Println(input1 == input2)
	fmt.Println(output1 == output2)
}
