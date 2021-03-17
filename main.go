package main

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
)

var link = "http://www.cbr.ru/scripts/XML_daily.asp"
var name = "Японских иен"

type ValCurs struct {
	Names  xml.Name `xml:"ValCurs"`
	Valute []Valute `xml:"Valute"`
}
type Valute struct {
	Valute   xml.Name `xml:"Valute"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Nominal  string   `xml:"Nominal"`
	Value    string   `xml:"Value"`
	Name     string   `xml:"Name"`
}

func getXML(url string, xmlUn *ValCurs) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status error: %v", resp.StatusCode)
	}
	d := xml.NewDecoder(resp.Body)

	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}
	err = d.Decode(xmlUn)
	if err != nil {
		return err
	}

	return nil
}

func makeCurList(valcurs *ValCurs) map[string]string {
	curList := make(map[string]string)
	for _, val := range valcurs.Valute {
		curList[val.Name] = val.Value
	}
	return curList
}

func main() {
	XmlS := &ValCurs{}
	err := getXML(link, XmlS)
	if err != nil {
		log.Fatal(err)
	}
	curList := makeCurList(XmlS)
	fmt.Println("курс ", name, " = ", curList[name], " рублей")
}
