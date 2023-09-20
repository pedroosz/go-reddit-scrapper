package utils

import (
	"bufio"
	"fmt"
	"os"
)

func ReadEntireFile(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		Fatal(fmt.Sprintf("Não foi possível ler o arquivo %s", filepath), err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	content := ""

	for scanner.Scan() {
		content += scanner.Text()
	}

	return content
}
