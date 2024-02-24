package helper

import (
	"fmt"
	"os"
)

func AppendToFile(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("utils open file error: %+v", err)
	}
	defer file.Close()
	if _, err = file.WriteString(content + "\n"); err != nil {
		return fmt.Errorf("utils write to file error: %+v", err)
	}
	return nil
}
