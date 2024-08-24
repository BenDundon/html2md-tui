package main

import (
	"fmt"
    "log"
    "os"
    "io"
	md "github.com/JohannesKaufmann/html-to-markdown"
)

func main() {
	converter := md.NewConverter("", true, nil)

    html := readFile("someFile.txt") 
    
	markdown, err := converter.ConvertString(html)

	if err != nil {
  		log.Fatal(err)
	}
	
	fmt.Println("md ->", markdown)
}

func readFile(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

    // Write the file text somewhere else so we can close it
    fileText, err := io.ReadAll(file)

    if err != nil {
    	log.Fatal(err)
    }

    // File *Must* be manually closed
    file.Close()
    
	return string(fileText)
}
