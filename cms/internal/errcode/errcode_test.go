package errcode

import (
	"fmt"

	"testing"
)

func TestXxx(t *testing.T) {
	err := LoginFail.WithDetails(IncorrectPassword.Error())
	fmt.Println(err)
	fmt.Printf("%s\n", err)
	fmt.Printf("%v\n", err)
	fmt.Printf("%+v\n", err)
	fmt.Printf("%#v\n", err)
}
