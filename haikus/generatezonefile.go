package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type haiku struct {
	Content string
	Author  string
}

func main() {
	fmt.Print(".")
	files, err := ioutil.ReadDir("./data/")
	if err != nil {
		log.Fatalf("Unable to read data dir: %s", err.Error())
	}

	haikus := make([]haiku, 0)

	for _, f := range files {
		filebytes, err := ioutil.ReadFile(fmt.Sprintf("./data/%s", f.Name()))
		if err != nil {
			log.Fatalf("wat 1 %s", err.Error())
		}

		lines := strings.Split(string(filebytes), "\n")

		insidehaiku := false
		HaikuContent := ""
		for _, line := range lines {
			if !insidehaiku && strings.HasPrefix(line, "[") {
				insidehaiku = true
				continue
			}
			if insidehaiku && strings.HasPrefix(line, "   [") {
				insidehaiku = false
				Author := line[7:]
				nhaiku := haiku{
					Author:  Author,
					Content: HaikuContent,
				}
				HaikuContent = ""
				haikus = append(haikus, nhaiku)
				continue
			}

			if insidehaiku && len(line) > 3 {
				HaikuContent = fmt.Sprintf("%s%s\n", HaikuContent, line[3:])
			}
		}
	}

	log.Printf("We found %d haikus", len(haikus))
}
