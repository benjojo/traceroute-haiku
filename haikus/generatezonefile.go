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
				HaikuContent = strings.Trim(HaikuContent, "\n ")

				nhaiku := haiku{
					Author:  Author,
					Content: HaikuContent,
				}
				if len(strings.Split(HaikuContent, "\n")) == 3 {
					haikus = append(haikus, nhaiku)
				}
				HaikuContent = ""

				continue
			}

			if insidehaiku && len(line) > 3 {
				HaikuContent = fmt.Sprintf("%s%s\n", HaikuContent, line[3:])
			}
		}
	}

	log.Printf("# We found %d haikus", len(haikus))
	var hid uint16 = 1
	for _, hku := range haikus {
		// Okay here is the idea, Here is what a normal prefix address looks like:
		// 0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.c.0.0.0.0.0.5.1.7.0.a.2.ip6.arpa
		// now we have around 3078 entries to be done here, so we need about ~12 bits
		// or 3 nibbles to work with, plus the final bit to increment for the letters.
		// so I will divide address space like so:
		//
		//
		// Z in the increment bit, this goes up with the sentance. X is the haiku ID.
		// We will always end on 4'th Z nibble, since 3 sentances and 1 author credit.
		hid++
		hidhexnibbles := fmt.Sprintf("%04x", hid)
		sentances := strings.Split(hku.Content, "\n")
		addrtemplate := fmt.Sprintf(".0.%s.%s.%s.%s.0.0.0.0.0.0.0.0.0.0.0.0.0.0.c.0.0.0.0.0.5.1.7.0.a.2.ip6.arpa.",
			string(hidhexnibbles[0]), string(hidhexnibbles[1]), string(hidhexnibbles[2]), string(hidhexnibbles[3]))
		fmt.Printf("0%s\t10\tIN\tPTR\t%s\n", addrtemplate, "haiku-trace.x.benjojo.co.uk.")
		fmt.Printf("1%s\t10\tIN\tPTR\t%s\n", addrtemplate, dnsfySentance(sentances[0]))
		fmt.Printf("2%s\t10\tIN\tPTR\t%s\n", addrtemplate, dnsfySentance(sentances[1]))
		fmt.Printf("3%s\t10\tIN\tPTR\t%s\n", addrtemplate, dnsfySentance(sentances[2]))
		fmt.Printf("4%s\t10\tIN\tPTR\t%s\n", addrtemplate, dnsfySentance(fmt.Sprintf("author %s", hku.Author)))
	}
}

func dnsfySentance(in string) string {
	in = strings.Replace(in, "\t", "", -1)
	in = strings.Replace(in, "—", "", -1)
	in = strings.Replace(in, " ", ".", -1)
	return strings.Trim(in, ".") + "."
}
