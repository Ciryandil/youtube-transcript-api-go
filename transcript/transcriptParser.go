package transcript

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

var FORMATTING_TAGS []string = []string{
	"strong", // important
	"em",     // emphasized
	"b",      // bold
	"i",      // italic
	"mark",   // marked
	"small",  // smaller
	"del",    // deleted
	"ins",    // inserted
	"sub",    // subscript
	"sup",    // superscript
}

type TranscriptXML struct {
	Text []struct {
		Start    float64 `xml:"start,attr"`
		Duration float64 `xml:"dur,attr"`
		Body     string  `xml:",chardata"`
	} `xml:"text"`
}

func getHtmlRegex(preserveFormatting bool) (*regexp.Regexp, error) {
	if preserveFormatting {
		formatsRegex := strings.Join(FORMATTING_TAGS, "|")
		formatsRegex = fmt.Sprintf("(?i)</?(?!(%s)\b).*?\b>", formatsRegex)
		return regexp.Compile(formatsRegex)
	} else {
		return regexp.Compile("<[^>]*>")
	}
}

func parseTranscriptXml(rawData string, preserveFormatting bool) ([]TranscriptSnippet, error) {
	var transcriptXmlStruct TranscriptXML
	err := xml.Unmarshal([]byte(rawData), &transcriptXmlStruct)
	if err != nil {
		return nil, fmt.Errorf("error parsing raw data to xml: %v", err)
	}
	snippetList := make([]TranscriptSnippet, 0)
	regex, err := getHtmlRegex(preserveFormatting)
	if err != nil {
		return nil, fmt.Errorf("error getting parsing regex: %v", err)
	}
	for _, textSnippet := range transcriptXmlStruct.Text {
		snippet := TranscriptSnippet{
			Start:    textSnippet.Start,
			Duration: textSnippet.Duration,
			Text:     regex.ReplaceAllString(textSnippet.Body, ""),
		}
		snippetList = append(snippetList, snippet)
	}
	return snippetList, nil
}
