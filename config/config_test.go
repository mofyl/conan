package config

import (
	"fmt"
	"testing"
)

func TestFile(t *testing.T) {

	newFile("")
	err := Decode(CONFIG_FILE_NAME, "log", nil)
	if err != nil {
		fmt.Println("1111")
		return
	}

}
