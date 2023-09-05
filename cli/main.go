// Command flies provides a webserver which produces detailed logs of each
// request it receives.
package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strings"

	"github.com/tnyeanderson/flies"
)

const defaultBanner = `
   __ _ _           
  / _| (_)          
 | |_| |_  ___  ___ 
 |  _| | |/ _ \/ __|
 | | | | |  __/\__ \
 |_| |_|_|\___||___/ `

func getRequestWriter(out io.Writer) flies.RequestWriter {
	switch os.Getenv("FLIES_FORMAT") {
	case "json":
		return flies.NewRequestWriterJSON(out)
	case "pipe":
		return flies.NewRequestWriter("\n", out)
	case "template":
		return getTemplateWriter(out)
	default:
		fmt.Fprintln(out, defaultBanner)
		fmt.Fprintln(out, strings.Repeat("+", 80))
		return flies.NewRequestWriterPretty(out)
	}
}

func getTemplateWriter(out io.Writer) *flies.RequestWriterTemplate {
	templateFile := os.Getenv("FLIES_TEMPLATE")
	if templateFile == "" {
		log.Fatal("FLIES_TEMPLATE_FILE must be set if FLIES_FORMAT=template")
	}
	templateText, err := os.ReadFile(templateFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	tmpl, err := template.New("").Parse(string(templateText))
	if err != nil {
		log.Fatal(err.Error())
	}
	return flies.NewRequestWriterTemplate(out, tmpl)
}

func main() {
	reqOut := os.Stdout
	s := &flies.Server{
		ErrWriter: os.Stderr,
		RawWriter: io.Discard,
	}
	s.Init()

	if os.Getenv("FLIES_RAW") != "" {
		s.RawWriter = os.Stdout
		s.ReqWriter = &flies.RequestWriterDiscard{}
	} else {
		s.ReqWriter = flies.NewMultiRequestWriter(
			s.ErrWriter,
			getRequestWriter(reqOut),
		)
	}

	err := s.Listen()
	if err != nil {
		log.Fatal(err.Error())
	}
}
