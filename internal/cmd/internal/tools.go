package cmdtools

import (
	"fmt"
	"os"
)

// TODO: Make sexy af
func Exit(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
	os.Exit(0)
}
