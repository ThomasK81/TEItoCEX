package main

import (
	"encoding/xml"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

type CTSCatalog struct {
	URN  []string										`json:"urn"`
	CitationScheme []string					`json:"citation_scheme"`
	GroupName []string							`json:"group_name"`
	WorkTitle []string							`json:"work_title"`
	VersionLabel []string						`json:"version_label"`
	ExemplarLabel []string					`json:"exemplar_label"`
	Online []string									`json:"online"`
	Language []string								`json:"language"`
}

type ExportDocument struct {
	URN  string                   `json:"urn"`
	CitationScheme string         `json:"citation_scheme"`
	GroupName string              `json:"group_name"`
	WorkTitle string              `json:"work_title"`
	VersionLabel string           `json:"version_label"`
	ExemplarLabel string          `json:"exemplar_label"`
	Online string                 `json:"online"`
	Language string               `json:"language"`
}

type OAI_DC_Record struct {
	//XMLName  xml.Name `xml:"http://www.openarchives.org/OAI/2.0/oai_dc/ oai_dc:dc"`
	XMLName  xml.Name `xml:"oai_dc:dc"`
	Xmlns1    string  `xml:"xmlns:oai_dc,attr"`
	Xmlns2    string  `xml:"xmlns:dc,attr"`
	Xmlns3    string  `xml:"xmlns:xsi,attr"`
	Xmlns4    string  `xml:"xsi:schemaLocation,attr"`
	Id        int     `xml:"id,attr"`
	Title     string  `xml:"dc:title"`
	Creator   string  `xml:"dc:creator"`
	Subject   string  `xml:"dc:subject"`
	Description [2]string `xml:"dc:description,omitempty"`
	Comment string    `xml:",comment"`
	//Date    string    `xml:"dc:date"`
  Language string   `xml:"dc:language"`
}

type Metadata struct {
	Xpath string
	Kind  string
}

type XPathInfo struct {
	XPathInfo string `xml:"replacementPattern,attr"`
	XPathWhat string `xml:"n,attr"`
}

type LangInfo struct {
	Language string `xml:"ident,attr"`
}

type RefPattern struct {
	RefPattern []XPathInfo `xml:"teiHeader>encodingDesc>refsDecl>cRefPattern"`
	Title      []string    `xml:"teiHeader>fileDesc>titleStmt>title"`
	Author     []string    `xml:"teiHeader>fileDesc>titleStmt>author"`
	Languages  []LangInfo  `xml:"teiHeader>profileDesc>langUsage>language"`
}

type teiHeader struct {
	RefPattern []XPathInfo `xml:"teiHeader>encodingDesc>refsDecl>cRefPattern"`
}

type SmallestNode struct {
	InnerXML string `xml:",innerxml"`
	Number   string `xml:"n,attr"`
}

type StartTEI1Direct struct {
	Node []SmallestNode `xml:"text>body>div"`
}

type StartTEI1p struct {
	Node []SmallestNode `xml:"text>body>div>p"`
}

type StartTEI1pseg struct {
	Node []SmallestNode `xml:"text>body>div>p>seg"`
}

type StartTEI1 struct {
	Node []SmallestNode `xml:"text>body>div>div"`
}

type StartTEI1Late struct {
	Node []SmallestNode `xml:"text>body>div>div>div"`
}

type TEI3n2 struct {
	Node   []SmallestNode `xml:"div"`
	Number string         `xml:"n,attr"`
}

type TEI3n3 struct {
	Node   []TEI3n2 `xml:"div"`
	Number string   `xml:"n,attr"`
}

type StartTEI3 struct {
	Node []TEI3n3 `xml:"text>body>div>div"`
}

type TEI3n2p struct {
	Node   []SmallestNode `xml:"p"`
	Number string         `xml:"n,attr"`
}

type TEI3n3p struct {
	Node   []TEI3n2p `xml:"div"`
	Number string    `xml:"n,attr"`
}

type StartTEI3p struct {
	Node []TEI3n3p `xml:"text>body>div>div"`
}

type TEI4n2p struct {
	Node   []SmallestNode `xml:"p"`
	Number string         `xml:"n,attr"`
}

type TEI4n3p struct {
	Node   []TEI4n2p `xml:"div"`
	Number string    `xml:"n,attr"`
}

type TEI4n4p struct {
	Node   []TEI4n3p `xml:"div"`
	Number string    `xml:"n,attr"`
}

type StartTEI4p struct {
	Node []TEI4n4p `xml:"text>body>div>div"`
}

type StartTEI2p struct {
	Node []TEI3n2p `xml:"text>body>div>div"`
}

type TEI2n2ab struct {
	Node   []SmallestNode `xml:"ab"`
	Number string         `xml:"n,attr"`
}

type StartTEI2ab struct {
	Node []TEI2n2ab `xml:"text>body>div>div"`
}

type TEI2n2lgl struct {
	Node   []SmallestNode `xml:"lg>l"`
	Number string         `xml:"n,attr"`
}

type StartTEI2lgl struct {
	Node []TEI2n2lgl `xml:"text>body>div>div"`
}

type TEI3n2cit struct {
	Node   []SmallestNode `xml:"cit"`
	Number string         `xml:"n,attr"`
}

type TEI3n3cit struct {
	Node   []TEI3n2cit `xml:"p"`
	Number string      `xml:"n,attr"`
}

type StartTEI3cit struct {
	Node []TEI3n3cit `xml:"text>body>div>div"`
}

type StartTEI2 struct {
	Node []TEI3n2 `xml:"text>body>div>div"`
}

type StartTEI2direct struct {
	Node []TEI3n2 `xml:"text>body>div"`
}

type StartTEI1Poem struct {
	Node []SmallestNode `xml:"text>body>div>l"`
}

type TEI2Poemn2 struct {
	Node   []SmallestNode `xml:"l"`
	Number string         `xml:"n,attr"`
}
type StartTEI2Poem struct {
	Node []TEI2Poemn2 `xml:"text>body>div>div"`
}

type TEI3n2DirectNumbered struct {
	Node   []SmallestNode `xml:"div3"`
	Number string         `xml:"n,attr"`
}

type TEI3n3DirectNumbered struct {
	Node   []TEI3n2DirectNumbered `xml:"div2"`
	Number string                 `xml:"n,attr"`
}

type StartTEI3DirectNumbered struct {
	Node []TEI3n3DirectNumbered `xml:"text>body>div1"`
}

type TEI3n2poem struct {
	Node   []SmallestNode `xml:"l"`
	Number string         `xml:"n,attr"`
}

type TEI3n3poem struct {
	Node   []TEI3n2poem `xml:"div"`
	Number string       `xml:"n,attr"`
}

type StartTEI3poem struct {
	Node []TEI3n3poem `xml:"text>body>div>div"`
}

type QSmallestNode struct {
	InnerXML string `xml:",innerxml"`
	Number   string `xml:"n,attr"`
}

type QueryTEI2 struct {
	Node []QSmallestNode `xml:"text>body>div>div"`
}

type QueryTEI3n3 struct {
	Node   []QSmallestNode `xml:"div"`
	Number string          `xml:"n,attr"`
}

type QueryTEI3 struct {
	Node []QueryTEI3n3 `xml:"text>body>div>div"`
}

type QueryTEI1 struct {
	Node []QSmallestNode `xml:"text>body>div"`
}

type QueryTEI1p struct {
	Node []QSmallestNode `xml:"text>body>div>p"`
}

type QSmallestNodeDiv struct {
	InnerXML string `xml:",innerxml"`
}

type QueryTEI1div struct {
	Node []QSmallestNodeDiv `xml:"text>body>div>div"`
}

type QueryTEI0 struct {
	Node []QSmallestNode `xml:"text>body"`
}

type QueryInfo struct {
	InnerXML string `xml:",innerxml"`
	Number   string `xml:"n,attr"`
}

func main() {
	outputFile := ""
	switch len(os.Args) {
	case 1:
		fmt.Println("Usage: CTSExtract [output-filename] [optionally: -CSV|JSON|XML]")
		os.Exit(3)
	case 2,3:
		outputFile = os.Args[1]
	default:
		fmt.Println("Usage: CTSExtract [output-filename] [optionally: -CSV|JSON]|XML")
		os.Exit(3)
	}

	tagsRegExp := regexp.MustCompile(`<[/]*[^>]*>`)
	greekWordRegExp := regexp.MustCompile(`\p{Greek}+`)
	latinWordRegExp := regexp.MustCompile(`\p{Latin}+`)
	arabicWordRegExp := regexp.MustCompile(`\p{Arabic}+`)

	var querystrings []string
	var identifiers []string
	var texts []string
	var greekwordcounts []string
	var latinwordcounts []string
	var arabicwordcounts []string
	var ctscatalog CTSCatalog

	filecount := 0
	greekwords := 0
	latinwords := 0
	arabicwords := 0
	noxpath := []string{}
	xmlFiles := checkExt(".xml")
	for _, file := range xmlFiles {

		xmlFile, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
		}
    //fmt.Println(file)
		byteValue, _ := ioutil.ReadAll(xmlFile)
		var xpathinfo RefPattern
		xml.Unmarshal(byteValue, &xpathinfo)
		if len(xpathinfo.RefPattern) == 0 {
			noxpath = append(noxpath, path.Base(file))
		}
		if len(xpathinfo.RefPattern) > 0 {
			var meta []Metadata
			for i := range xpathinfo.RefPattern {
				meta = append(meta, Metadata{Xpath: xpathinfo.RefPattern[i].XPathInfo, Kind: xpathinfo.RefPattern[i].XPathWhat})
			}
			sort.Slice(meta, func(i int, j int) bool {
				return len(meta[i].Xpath) < len(meta[j].Xpath)
			})
			querystring := meta[len(meta)-1].Xpath
			whatkind := []string{}
			for i := range meta {
				whatkind = append(whatkind, meta[i].Kind)
			}
			languages := []string{}
			for i := range xpathinfo.Languages {
				languages = append(languages, xpathinfo.Languages[i].Language)
			}
			language := strings.Join(languages, ",")
			kind := strings.Join(whatkind, ",")
			querystring = strings.Replace(querystring, "#xpath(", "", -1)
			querystring = strings.Replace(querystring, ")", "", -1)
			urn := strings.Replace(path.Base(file), ".xml", "", -1)
			ctscatalog.URN = append(ctscatalog.URN, urn)
			ctscatalog.CitationScheme = append(ctscatalog.CitationScheme, kind)
			group := strings.Join(xpathinfo.Author, ",")
			group = strings.Replace(group, "\n", " ", -1)
			group = tagsRegExp.ReplaceAllString(group, "")
			group = strings.TrimSpace(group)
			ctscatalog.GroupName = append(ctscatalog.GroupName, group)
			worktitle := strings.Join(xpathinfo.Title, ",")
			worktitle = strings.Replace(worktitle, "\n", " ", -1)
			worktitle = tagsRegExp.ReplaceAllString(worktitle, "")
			worktitle = strings.TrimSpace(worktitle)
			ctscatalog.WorkTitle = append(ctscatalog.WorkTitle, worktitle)
			ctscatalog.VersionLabel = append(ctscatalog.VersionLabel, "")
			ctscatalog.ExemplarLabel = append(ctscatalog.ExemplarLabel, "")
			ctscatalog.Online = append(ctscatalog.Online, "True")
			ctscatalog.Language = append(ctscatalog.Language, language)
			switch {
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']/tei:div[@n='$2']/tei:div[@n='$3']/tei:p[@n='$4']":
				fmt.Print("1")
				filecount = filecount + 1
				var data StartTEI4p
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						for k := range data.Node[i].Node[j].Node {
							for l := range data.Node[i].Node[j].Node[k].Node {
								id := []string{data.Node[i].Number, data.Node[i].Node[j].Number, data.Node[i].Node[j].Node[k].Number, data.Node[i].Node[j].Node[k].Node[l].Number}
								identifier := strings.Join(id, ".")
								identifier = strings.Join([]string{baseIdentifier, identifier}, ":")
								text := data.Node[i].Node[j].Node[k].Node[l].InnerXML
								text = stringcleaning(text)

								words := greekWordRegExp.FindAllString(text, -1)
								latinword := latinWordRegExp.FindAllString(text, -1)
								arabicword := arabicWordRegExp.FindAllString(text, -1)
								greekwords = greekwords + len(words)
								wordcount := strconv.Itoa(len(words))
								latinwords = latinwords + len(latinword)
								arabicwords = arabicwords + len(arabicword)
								latinwordcount := strconv.Itoa(len(latinword))
								arabicwordcount := strconv.Itoa(len(arabicword))
								latinwordcounts = append(latinwordcounts, latinwordcount)
								arabicwordcounts = append(arabicwordcounts, arabicwordcount)
								identifiers = append(identifiers, identifier)
								texts = append(texts, text)
								greekwordcounts = append(greekwordcounts, wordcount)
							}
						}
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div//tei:div[@n='$1']" || querystring == "/tei:TEI/tei:text/tei:body/tei:div//tei:div[@n=\\'$1\\']":
				fmt.Print("2")
				filecount = filecount + 1
				var data QueryTEI1
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					r := strings.NewReader(data.Node[i].InnerXML)
					decoder := xml.NewDecoder(r)
					for {
						t, _ := decoder.Token()
						if t == nil {
							break
						}
						switch se := t.(type) {
						case xml.StartElement:
							if se.Name.Local == "div" {
								var info QueryInfo
								decoder.DecodeElement(&info, &se)
								identifier := strings.Join([]string{baseIdentifier, info.Number}, ":")
								text := info.InnerXML
								text = stringcleaning(text)

								words := greekWordRegExp.FindAllString(text, -1)
								latinword := latinWordRegExp.FindAllString(text, -1)
								arabicword := arabicWordRegExp.FindAllString(text, -1)
								greekwords = greekwords + len(words)
								wordcount := strconv.Itoa(len(words))
								latinwords = latinwords + len(latinword)
								arabicwords = arabicwords + len(arabicword)
								latinwordcount := strconv.Itoa(len(latinword))
								arabicwordcount := strconv.Itoa(len(arabicword))
								latinwordcounts = append(latinwordcounts, latinwordcount)
								arabicwordcounts = append(arabicwordcounts, arabicwordcount)
								identifiers = append(identifiers, identifier)
								texts = append(texts, text)
								greekwordcounts = append(greekwordcounts, wordcount)
							}
						}
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div//tei:div[@subtype='fragment'][@n='$1']":
				fmt.Print("3")
				filecount = filecount + 1
				var data QueryTEI1div
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					r := strings.NewReader(data.Node[i].InnerXML)
					decoder := xml.NewDecoder(r)
					for {
						t, _ := decoder.Token()
						if t == nil {
							break
						}
						switch se := t.(type) {
						case xml.StartElement:
							if se.Name.Local == "div" {
								var info QueryInfo
								decoder.DecodeElement(&info, &se)
								identifier := strings.Join([]string{baseIdentifier, info.Number}, ":")
								text := info.InnerXML
								text = stringcleaning(text)

								words := greekWordRegExp.FindAllString(text, -1)
								latinword := latinWordRegExp.FindAllString(text, -1)
								arabicword := arabicWordRegExp.FindAllString(text, -1)
								greekwords = greekwords + len(words)
								wordcount := strconv.Itoa(len(words))
								latinwords = latinwords + len(latinword)
								arabicwords = arabicwords + len(arabicword)
								latinwordcount := strconv.Itoa(len(latinword))
								arabicwordcount := strconv.Itoa(len(arabicword))
								latinwordcounts = append(latinwordcounts, latinwordcount)
								arabicwordcounts = append(arabicwordcounts, arabicwordcount)
								identifiers = append(identifiers, identifier)
								texts = append(texts, text)
								greekwordcounts = append(greekwordcounts, wordcount)
							}
						}
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:p//tei:l[@n='$1']":
				fmt.Print("4")
				filecount = filecount + 1
				var data QueryTEI1p
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					r := strings.NewReader(data.Node[i].InnerXML)
					decoder := xml.NewDecoder(r)
					for {
						t, _ := decoder.Token()
						if t == nil {
							break
						}
						switch se := t.(type) {
						case xml.StartElement:
							if se.Name.Local == "l" {
								var info QueryInfo
								decoder.DecodeElement(&info, &se)
								identifier := strings.Join([]string{baseIdentifier, info.Number}, ":")
								text := info.InnerXML
								text = stringcleaning(text)

								words := greekWordRegExp.FindAllString(text, -1)
								latinword := latinWordRegExp.FindAllString(text, -1)
								arabicword := arabicWordRegExp.FindAllString(text, -1)
								greekwords = greekwords + len(words)
								wordcount := strconv.Itoa(len(words))
								latinwords = latinwords + len(latinword)
								arabicwords = arabicwords + len(arabicword)
								latinwordcount := strconv.Itoa(len(latinword))
								arabicwordcount := strconv.Itoa(len(arabicword))
								latinwordcounts = append(latinwordcounts, latinwordcount)
								arabicwordcounts = append(arabicwordcounts, arabicwordcount)
								identifiers = append(identifiers, identifier)
								texts = append(texts, text)
								greekwordcounts = append(greekwordcounts, wordcount)
							}
						}
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div//tei:l[@n='$1']" || querystring == "/tei:TEI/tei:text/tei:body/tei:div//tei:l[@n=\\'$1\\']":
				fmt.Print("5")
				filecount = filecount + 1
				var data QueryTEI1
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					r := strings.NewReader(data.Node[i].InnerXML)
					decoder := xml.NewDecoder(r)
					for {
						t, _ := decoder.Token()
						if t == nil {
							break
						}
						switch se := t.(type) {
						case xml.StartElement:
							if se.Name.Local == "l" {
								var info QueryInfo
								decoder.DecodeElement(&info, &se)
								identifier := strings.Join([]string{baseIdentifier, info.Number}, ":")
								text := info.InnerXML
								text = stringcleaning(text)

								words := greekWordRegExp.FindAllString(text, -1)
								latinword := latinWordRegExp.FindAllString(text, -1)
								arabicword := arabicWordRegExp.FindAllString(text, -1)
								greekwords = greekwords + len(words)
								wordcount := strconv.Itoa(len(words))
								latinwords = latinwords + len(latinword)
								arabicwords = arabicwords + len(arabicword)
								latinwordcount := strconv.Itoa(len(latinword))
								arabicwordcount := strconv.Itoa(len(arabicword))
								latinwordcounts = append(latinwordcounts, latinwordcount)
								arabicwordcounts = append(arabicwordcounts, arabicwordcount)
								identifiers = append(identifiers, identifier)
								texts = append(texts, text)
								greekwordcounts = append(greekwordcounts, wordcount)
							}
						}
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body//tei:l[@n=\\'$1\\']":
				fmt.Print("6")
				filecount = filecount + 1
				var data QueryTEI0
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					r := strings.NewReader(data.Node[i].InnerXML)
					decoder := xml.NewDecoder(r)
					for {
						t, _ := decoder.Token()
						if t == nil {
							break
						}
						switch se := t.(type) {
						case xml.StartElement:
							if se.Name.Local == "l" {
								var info QueryInfo
								decoder.DecodeElement(&info, &se)
								identifier := strings.Join([]string{baseIdentifier, info.Number}, ":")
								text := info.InnerXML
								text = stringcleaning(text)

								words := greekWordRegExp.FindAllString(text, -1)
								latinword := latinWordRegExp.FindAllString(text, -1)
								arabicword := arabicWordRegExp.FindAllString(text, -1)
								greekwords = greekwords + len(words)
								wordcount := strconv.Itoa(len(words))
								latinwords = latinwords + len(latinword)
								arabicwords = arabicwords + len(arabicword)
								latinwordcount := strconv.Itoa(len(latinword))
								arabicwordcount := strconv.Itoa(len(arabicword))
								latinwordcounts = append(latinwordcounts, latinwordcount)
								arabicwordcounts = append(arabicwordcounts, arabicwordcount)
								identifiers = append(identifiers, identifier)
								texts = append(texts, text)
								greekwordcounts = append(greekwordcounts, wordcount)
							}
						}
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']//tei:div[@n='$2']":
				fmt.Print("7")
				filecount = filecount + 1
				var data QueryTEI2
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					ID2 := data.Node[i].Number
					r := strings.NewReader(data.Node[i].InnerXML)
					decoder := xml.NewDecoder(r)
					for {
						t, _ := decoder.Token()
						if t == nil {
							break
						}
						switch se := t.(type) {
						case xml.StartElement:
							if se.Name.Local == "div" {
								var info QueryInfo
								decoder.DecodeElement(&info, &se)
								id := []string{ID2, info.Number}
								identifier := strings.Join(id, ".")
								identifier = strings.Join([]string{baseIdentifier, identifier}, ":")
								text := info.InnerXML
								text = stringcleaning(text)

								words := greekWordRegExp.FindAllString(text, -1)
								latinword := latinWordRegExp.FindAllString(text, -1)
								arabicword := arabicWordRegExp.FindAllString(text, -1)
								greekwords = greekwords + len(words)
								wordcount := strconv.Itoa(len(words))
								latinwords = latinwords + len(latinword)
								arabicwords = arabicwords + len(arabicword)
								latinwordcount := strconv.Itoa(len(latinword))
								arabicwordcount := strconv.Itoa(len(arabicword))
								latinwordcounts = append(latinwordcounts, latinwordcount)
								arabicwordcounts = append(arabicwordcounts, arabicwordcount)
								identifiers = append(identifiers, identifier)
								texts = append(texts, text)
								greekwordcounts = append(greekwordcounts, wordcount)
							}
						}
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']/tei:div[@n='$2']//tei:div[@n='$3']":
				fmt.Print("8")
				filecount = filecount + 1
				var data QueryTEI3
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					ID2 := data.Node[i].Number
					for j := range data.Node[i].Node {
						ID3 := data.Node[i].Node[j].Number
						r := strings.NewReader(data.Node[i].Node[j].InnerXML)
						decoder := xml.NewDecoder(r)
						for {
							t, _ := decoder.Token()
							if t == nil {
								break
							}
							switch se := t.(type) {
							case xml.StartElement:
								if se.Name.Local == "div" {
									var info QueryInfo
									decoder.DecodeElement(&info, &se)
									id := []string{ID2, ID3, info.Number}
									identifier := strings.Join(id, ".")
									identifier = strings.Join([]string{baseIdentifier, identifier}, ":")
									text := info.InnerXML
									text = stringcleaning(text)

									words := greekWordRegExp.FindAllString(text, -1)
									latinword := latinWordRegExp.FindAllString(text, -1)
									arabicword := arabicWordRegExp.FindAllString(text, -1)
									greekwords = greekwords + len(words)
									wordcount := strconv.Itoa(len(words))
									latinwords = latinwords + len(latinword)
									arabicwords = arabicwords + len(arabicword)
									latinwordcount := strconv.Itoa(len(latinword))
									arabicwordcount := strconv.Itoa(len(arabicword))
									latinwordcounts = append(latinwordcounts, latinwordcount)
									arabicwordcounts = append(arabicwordcounts, arabicwordcount)
									identifiers = append(identifiers, identifier)
									texts = append(texts, text)
									greekwordcounts = append(greekwordcounts, wordcount)
								}
							}
						}
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']//tei:l[@n='$2']":
				fmt.Print("9")
				filecount = filecount + 1
				var data QueryTEI2
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					ID2 := data.Node[i].Number
					r := strings.NewReader(data.Node[i].InnerXML)
					decoder := xml.NewDecoder(r)
					for {
						t, _ := decoder.Token()
						if t == nil {
							break
						}
						switch se := t.(type) {
						case xml.StartElement:
							if se.Name.Local == "l" {
								var info QueryInfo
								decoder.DecodeElement(&info, &se)
								id := []string{ID2, info.Number}
								identifier := strings.Join(id, ".")
								identifier = strings.Join([]string{baseIdentifier, identifier}, ":")
								text := info.InnerXML
								text = stringcleaning(text)

								words := greekWordRegExp.FindAllString(text, -1)
								latinword := latinWordRegExp.FindAllString(text, -1)
								arabicword := arabicWordRegExp.FindAllString(text, -1)
								greekwords = greekwords + len(words)
								wordcount := strconv.Itoa(len(words))
								latinwords = latinwords + len(latinword)
								arabicwords = arabicwords + len(arabicword)
								latinwordcount := strconv.Itoa(len(latinword))
								arabicwordcount := strconv.Itoa(len(arabicword))
								latinwordcounts = append(latinwordcounts, latinwordcount)
								arabicwordcounts = append(arabicwordcounts, arabicwordcount)
								identifiers = append(identifiers, identifier)
								texts = append(texts, text)
								greekwordcounts = append(greekwordcounts, wordcount)
							}
						}
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']/tei:div[@n='$2']/tei:p[@n='$3']" || querystring == "/tei:TEI/tei:text/tei:body/tei:div[@type='edition']/tei:div[@n=\\'$1\\']/tei:div[@n=\\'$2\\']/tei:p[@n=\\'$3\\']":
				fmt.Print("A")
				filecount = filecount + 1
				var data StartTEI3p
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						for k := range data.Node[i].Node[j].Node {
							id := []string{data.Node[i].Number, data.Node[i].Node[j].Number, data.Node[i].Node[j].Node[k].Number}
							identifier := strings.Join(id, ".")
							identifier = strings.Join([]string{baseIdentifier, identifier}, ":")
							text := data.Node[i].Node[j].Node[k].InnerXML
							text = stringcleaning(text)

							words := greekWordRegExp.FindAllString(text, -1)
							latinword := latinWordRegExp.FindAllString(text, -1)
							arabicword := arabicWordRegExp.FindAllString(text, -1)
							greekwords = greekwords + len(words)
							wordcount := strconv.Itoa(len(words))
							latinwords = latinwords + len(latinword)
							arabicwords = arabicwords + len(arabicword)
							latinwordcount := strconv.Itoa(len(latinword))
							arabicwordcount := strconv.Itoa(len(arabicword))
							latinwordcounts = append(latinwordcounts, latinwordcount)
							arabicwordcounts = append(arabicwordcounts, arabicwordcount)
							identifiers = append(identifiers, identifier)
							texts = append(texts, text)
							greekwordcounts = append(greekwordcounts, wordcount)
						}
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']/tei:p[@n='$2']/tei:cit[@n='$3']":
				fmt.Print("B")
				filecount = filecount + 1
				var data StartTEI3cit
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						for k := range data.Node[i].Node[j].Node {
							id := []string{data.Node[i].Number, data.Node[i].Node[j].Number, data.Node[i].Node[j].Node[k].Number}
							identifier := strings.Join(id, ".")
							identifier = strings.Join([]string{baseIdentifier, identifier}, ":")
							text := data.Node[i].Node[j].Node[k].InnerXML
							text = stringcleaning(text)

							words := greekWordRegExp.FindAllString(text, -1)
							latinword := latinWordRegExp.FindAllString(text, -1)
							arabicword := arabicWordRegExp.FindAllString(text, -1)
							greekwords = greekwords + len(words)
							wordcount := strconv.Itoa(len(words))
							latinwords = latinwords + len(latinword)
							arabicwords = arabicwords + len(arabicword)
							latinwordcount := strconv.Itoa(len(latinword))
							arabicwordcount := strconv.Itoa(len(arabicword))
							latinwordcounts = append(latinwordcounts, latinwordcount)
							arabicwordcounts = append(arabicwordcounts, arabicwordcount)
							identifiers = append(identifiers, identifier)
							texts = append(texts, text)
							greekwordcounts = append(greekwordcounts, wordcount)
						}
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div[@n=\\'$1\\']" || querystring == "/tei:TEI.2/tei:text/tei:body/tei:div[@n=\\'$1\\']":
				fmt.Print("C")
				filecount = filecount + 1
				var data StartTEI1Direct
				xml.Unmarshal(byteValue, &data)
				for _, node := range data.Node {
					identifier := path.Base(file)
					identifier = strings.Replace(identifier, ".xml", "", -1)
					identifier = strings.Join([]string{identifier, node.Number}, ":")
					text := node.InnerXML
					text = stringcleaning(text)

					words := greekWordRegExp.FindAllString(text, -1)
					latinword := latinWordRegExp.FindAllString(text, -1)
					arabicword := arabicWordRegExp.FindAllString(text, -1)
					greekwords = greekwords + len(words)
					wordcount := strconv.Itoa(len(words))
					latinwords = latinwords + len(latinword)
					arabicwords = arabicwords + len(arabicword)
					latinwordcount := strconv.Itoa(len(latinword))
					arabicwordcount := strconv.Itoa(len(arabicword))
					latinwordcounts = append(latinwordcounts, latinwordcount)
					arabicwordcounts = append(arabicwordcounts, arabicwordcount)
					identifiers = append(identifiers, identifier)
					texts = append(texts, text)
					greekwordcounts = append(greekwordcounts, wordcount)
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div/tei:div[@n='$1']":
				fmt.Print("D")
				filecount = filecount + 1
				var data StartTEI1Late
				xml.Unmarshal(byteValue, &data)
				for _, node := range data.Node {
					identifier := path.Base(file)
					identifier = strings.Replace(identifier, ".xml", "", -1)
					identifier = strings.Join([]string{identifier, node.Number}, ":")
					text := node.InnerXML
					text = stringcleaning(text)

					words := greekWordRegExp.FindAllString(text, -1)
					latinword := latinWordRegExp.FindAllString(text, -1)
					arabicword := arabicWordRegExp.FindAllString(text, -1)
					greekwords = greekwords + len(words)
					wordcount := strconv.Itoa(len(words))
					latinwords = latinwords + len(latinword)
					arabicwords = arabicwords + len(arabicword)
					latinwordcount := strconv.Itoa(len(latinword))
					arabicwordcount := strconv.Itoa(len(arabicword))
					latinwordcounts = append(latinwordcounts, latinwordcount)
					arabicwordcounts = append(arabicwordcounts, arabicwordcount)
					identifiers = append(identifiers, identifier)
					texts = append(texts, text)
					greekwordcounts = append(greekwordcounts, wordcount)
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:p/tei:seg[@n='$1']":
				fmt.Print("E")
				filecount = filecount + 1
				var data StartTEI1pseg
				xml.Unmarshal(byteValue, &data)
				for _, node := range data.Node {
					identifier := path.Base(file)
					identifier = strings.Replace(identifier, ".xml", "", -1)
					identifier = strings.Join([]string{identifier, node.Number}, ":")
					text := node.InnerXML
					text = stringcleaning(text)

					words := greekWordRegExp.FindAllString(text, -1)
					latinword := latinWordRegExp.FindAllString(text, -1)
					arabicword := arabicWordRegExp.FindAllString(text, -1)
					greekwords = greekwords + len(words)
					wordcount := strconv.Itoa(len(words))
					latinwords = latinwords + len(latinword)
					arabicwords = arabicwords + len(arabicword)
					latinwordcount := strconv.Itoa(len(latinword))
					arabicwordcount := strconv.Itoa(len(arabicword))
					latinwordcounts = append(latinwordcounts, latinwordcount)
					arabicwordcounts = append(arabicwordcounts, arabicwordcount)
					identifiers = append(identifiers, identifier)
					texts = append(texts, text)
					greekwordcounts = append(greekwordcounts, wordcount)
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:p[@n='$1']":
				fmt.Print("F")
				filecount = filecount + 1
				var data StartTEI1p
				xml.Unmarshal(byteValue, &data)
				for _, node := range data.Node {
					identifier := path.Base(file)
					identifier = strings.Replace(identifier, ".xml", "", -1)
					identifier = strings.Join([]string{identifier, node.Number}, ":")
					text := node.InnerXML
					text = stringcleaning(text)

					words := greekWordRegExp.FindAllString(text, -1)
					latinword := latinWordRegExp.FindAllString(text, -1)
					arabicword := arabicWordRegExp.FindAllString(text, -1)
					greekwords = greekwords + len(words)
					wordcount := strconv.Itoa(len(words))
					latinwords = latinwords + len(latinword)
					arabicwords = arabicwords + len(arabicword)
					latinwordcount := strconv.Itoa(len(latinword))
					arabicwordcount := strconv.Itoa(len(arabicword))
					latinwordcounts = append(latinwordcounts, latinwordcount)
					arabicwordcounts = append(arabicwordcounts, arabicwordcount)
					identifiers = append(identifiers, identifier)
					texts = append(texts, text)
					greekwordcounts = append(greekwordcounts, wordcount)
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div[@type='edition']/tei:div[@n='$1']" || querystring == "/tei:TEI/tei:text/tei:body/div[@type='edition']/div[@n='$1']" || querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n=\\'$1\\']" || querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']":
				fmt.Print("G")
				filecount = filecount + 1
				var data StartTEI1
				xml.Unmarshal(byteValue, &data)
				for _, node := range data.Node {
					identifier := path.Base(file)
					identifier = strings.Replace(identifier, ".xml", "", -1)
					identifier = strings.Join([]string{identifier, node.Number}, ":")
					text := node.InnerXML
					text = stringcleaning(text)

					words := greekWordRegExp.FindAllString(text, -1)
					latinword := latinWordRegExp.FindAllString(text, -1)
					arabicword := arabicWordRegExp.FindAllString(text, -1)
					greekwords = greekwords + len(words)
					wordcount := strconv.Itoa(len(words))
					latinwords = latinwords + len(latinword)
					arabicwords = arabicwords + len(arabicword)
					latinwordcount := strconv.Itoa(len(latinword))
					arabicwordcount := strconv.Itoa(len(arabicword))
					latinwordcounts = append(latinwordcounts, latinwordcount)
					arabicwordcounts = append(arabicwordcounts, arabicwordcount)
					identifiers = append(identifiers, identifier)
					texts = append(texts, text)
					greekwordcounts = append(greekwordcounts, wordcount)
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']/tei:p[@n='$2']":
				fmt.Print("H")
				filecount = filecount + 1
				var data StartTEI2p
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						id := []string{data.Node[i].Number, data.Node[i].Node[j].Number}
						identifier := strings.Join(id, ".")
						identifier = strings.Join([]string{baseIdentifier, identifier}, ":")
						text := data.Node[i].Node[j].InnerXML
						text = stringcleaning(text)

						words := greekWordRegExp.FindAllString(text, -1)
						latinword := latinWordRegExp.FindAllString(text, -1)
						arabicword := arabicWordRegExp.FindAllString(text, -1)
						greekwords = greekwords + len(words)
						wordcount := strconv.Itoa(len(words))
						latinwords = latinwords + len(latinword)
						arabicwords = arabicwords + len(arabicword)
						latinwordcount := strconv.Itoa(len(latinword))
						arabicwordcount := strconv.Itoa(len(arabicword))
						latinwordcounts = append(latinwordcounts, latinwordcount)
						arabicwordcounts = append(arabicwordcounts, arabicwordcount)
						identifiers = append(identifiers, identifier)
						texts = append(texts, text)
						greekwordcounts = append(greekwordcounts, wordcount)
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']/tei:ab[@n='$2']":
				fmt.Print("I")
				filecount = filecount + 1
				var data StartTEI2ab
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						id := []string{data.Node[i].Number, data.Node[i].Node[j].Number}
						identifier := strings.Join(id, ".")
						identifier = strings.Join([]string{baseIdentifier, identifier}, ":")
						text := data.Node[i].Node[j].InnerXML
						text = stringcleaning(text)

						words := greekWordRegExp.FindAllString(text, -1)
						latinword := latinWordRegExp.FindAllString(text, -1)
						arabicword := arabicWordRegExp.FindAllString(text, -1)
						greekwords = greekwords + len(words)
						wordcount := strconv.Itoa(len(words))
						latinwords = latinwords + len(latinword)
						arabicwords = arabicwords + len(arabicword)
						latinwordcount := strconv.Itoa(len(latinword))
						arabicwordcount := strconv.Itoa(len(arabicword))
						latinwordcounts = append(latinwordcounts, latinwordcount)
						arabicwordcounts = append(arabicwordcounts, arabicwordcount)
						identifiers = append(identifiers, identifier)
						texts = append(texts, text)
						greekwordcounts = append(greekwordcounts, wordcount)
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']/tei:lg/tei:l[@n='$2']":
				fmt.Print("J")
				filecount = filecount + 1
				var data StartTEI2lgl
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						id := []string{data.Node[i].Number, data.Node[i].Node[j].Number}
						identifier := strings.Join(id, ".")
						identifier = strings.Join([]string{baseIdentifier, identifier}, ":")
						text := data.Node[i].Node[j].InnerXML
						text = stringcleaning(text)

						words := greekWordRegExp.FindAllString(text, -1)
						latinword := latinWordRegExp.FindAllString(text, -1)
						arabicword := arabicWordRegExp.FindAllString(text, -1)
						greekwords = greekwords + len(words)
						wordcount := strconv.Itoa(len(words))
						latinwords = latinwords + len(latinword)
						arabicwords = arabicwords + len(arabicword)
						latinwordcount := strconv.Itoa(len(latinword))
						arabicwordcount := strconv.Itoa(len(arabicword))
						latinwordcounts = append(latinwordcounts, latinwordcount)
						arabicwordcounts = append(arabicwordcounts, arabicwordcount)
						identifiers = append(identifiers, identifier)
						texts = append(texts, text)
						greekwordcounts = append(greekwordcounts, wordcount)
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']/tei:div[@n='$2']" || querystring == "/tei:TEI/tei:text/tei:body/tei:div[@type='translation']/tei:div[@n='$1']/tei:div[@n='$2']" || querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n=\\'$1\\']/tei:div[@n=\\'$2\\']" || querystring == "/tei:TEI/tei:text/tei:body/div[@type='edition']/div[@n='$1']/div[@n='$2']" || querystring == "/tei:TEI/tei:text/tei:body/tei:div[@type='edition']/tei:div[@n='$1']/tei:div[@n='$2']":
				fmt.Print("K")
				filecount = filecount + 1
				var data StartTEI2
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						id := []string{data.Node[i].Number, data.Node[i].Node[j].Number}
						identifier := strings.Join(id, ".")
						identifier = strings.Join([]string{baseIdentifier, identifier}, ":")
						text := data.Node[i].Node[j].InnerXML
						text = stringcleaning(text)

						words := greekWordRegExp.FindAllString(text, -1)
						latinword := latinWordRegExp.FindAllString(text, -1)
						arabicword := arabicWordRegExp.FindAllString(text, -1)
						greekwords = greekwords + len(words)
						wordcount := strconv.Itoa(len(words))
						latinwords = latinwords + len(latinword)
						arabicwords = arabicwords + len(arabicword)
						latinwordcount := strconv.Itoa(len(latinword))
						arabicwordcount := strconv.Itoa(len(arabicword))
						latinwordcounts = append(latinwordcounts, latinwordcount)
						arabicwordcounts = append(arabicwordcounts, arabicwordcount)
						identifiers = append(identifiers, identifier)
						texts = append(texts, text)
						greekwordcounts = append(greekwordcounts, wordcount)
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div[@type='edition']/tei:div[@n='$1']/tei:div[@n='$2']/tei:div[@n='$3']" || querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']/tei:div[@n='$2']/tei:div[@n='$3']" || querystring == "/tei:TEI/tei:text/tei:body/div[@type='edition']/div[@n='$1']/div[@n='$2']/div[@n='$3']" || querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n=\\'$1\\']/tei:div[@n=\\'$2\\']/tei:div[@n=\\'$3\\']":
				fmt.Print("L")
				filecount = filecount + 1
				var data StartTEI3
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						for k := range data.Node[i].Node[j].Node {
							id := []string{data.Node[i].Number, data.Node[i].Node[j].Number, data.Node[i].Node[j].Node[k].Number}
							identifier := strings.Join(id, ".")
							identifier = strings.Join([]string{baseIdentifier, identifier}, ":")
							text := data.Node[i].Node[j].Node[k].InnerXML
							text = stringcleaning(text)

							words := greekWordRegExp.FindAllString(text, -1)
							latinword := latinWordRegExp.FindAllString(text, -1)
							arabicword := arabicWordRegExp.FindAllString(text, -1)
							greekwords = greekwords + len(words)
							wordcount := strconv.Itoa(len(words))
							latinwords = latinwords + len(latinword)
							arabicwords = arabicwords + len(arabicword)
							latinwordcount := strconv.Itoa(len(latinword))
							arabicwordcount := strconv.Itoa(len(arabicword))
							latinwordcounts = append(latinwordcounts, latinwordcount)
							arabicwordcounts = append(arabicwordcounts, arabicwordcount)
							identifiers = append(identifiers, identifier)
							texts = append(texts, text)
							greekwordcounts = append(greekwordcounts, wordcount)
						}
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:l[@n='$1']":
				fmt.Print("M")
				filecount = filecount + 1
				var data StartTEI1Poem
				xml.Unmarshal(byteValue, &data)
				for _, node := range data.Node {
					identifier := path.Base(file)
					identifier = strings.Replace(identifier, ".xml", "", -1)
					identifier = strings.Join([]string{identifier, node.Number}, ":")
					text := node.InnerXML
					text = stringcleaning(text)

					words := greekWordRegExp.FindAllString(text, -1)
					latinword := latinWordRegExp.FindAllString(text, -1)
					arabicword := arabicWordRegExp.FindAllString(text, -1)
					greekwords = greekwords + len(words)
					wordcount := strconv.Itoa(len(words))
					latinwords = latinwords + len(latinword)
					arabicwords = arabicwords + len(arabicword)
					latinwordcount := strconv.Itoa(len(latinword))
					arabicwordcount := strconv.Itoa(len(arabicword))
					latinwordcounts = append(latinwordcounts, latinwordcount)
					arabicwordcounts = append(arabicwordcounts, arabicwordcount)
					identifiers = append(identifiers, identifier)
					texts = append(texts, text)
					greekwordcounts = append(greekwordcounts, wordcount)
				}
			case querystring == "/tei:TEI.2/tei:text/tei:body/tei:div[@n=\\'$1\\']/tei:div[@n=\\'$2\\']" || querystring == "/tei:TEI/tei:text/tei:body/tei:div[@n=\\'$1\\']/tei:div[@n=\\'$2\\']":
				fmt.Print("N")
				filecount = filecount + 1
				var data StartTEI2direct
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						id := []string{data.Node[i].Number, data.Node[i].Node[j].Number}
						identifier := strings.Join(id, ".")
						identifier = strings.Join([]string{baseIdentifier, identifier}, ":")
						text := data.Node[i].Node[j].InnerXML
						text = stringcleaning(text)

						words := greekWordRegExp.FindAllString(text, -1)
						latinword := latinWordRegExp.FindAllString(text, -1)
						arabicword := arabicWordRegExp.FindAllString(text, -1)
						greekwords = greekwords + len(words)
						wordcount := strconv.Itoa(len(words))
						latinwords = latinwords + len(latinword)
						arabicwords = arabicwords + len(arabicword)
						latinwordcount := strconv.Itoa(len(latinword))
						arabicwordcount := strconv.Itoa(len(arabicword))
						latinwordcounts = append(latinwordcounts, latinwordcount)
						arabicwordcounts = append(arabicwordcounts, arabicwordcount)
						identifiers = append(identifiers, identifier)
						texts = append(texts, text)
						greekwordcounts = append(greekwordcounts, wordcount)
					}
				}
			case querystring == "/tei:TEI.2/tei:text/tei:body/tei:div1[@n=\\'$1\\']/tei:div2[@n=\\'$2\\']/tei:div3[@n=\\'$3\\']":
				fmt.Print("O")
				filecount = filecount + 1
				var data StartTEI3DirectNumbered
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						for k := range data.Node[i].Node[j].Node {
							id := []string{data.Node[i].Number, data.Node[i].Node[j].Number, data.Node[i].Node[j].Node[k].Number}
							identifier := strings.Join(id, ".")
							identifier = strings.Join([]string{baseIdentifier, identifier}, ":")
							text := data.Node[i].Node[j].Node[k].InnerXML
							text = stringcleaning(text)

							words := greekWordRegExp.FindAllString(text, -1)
							latinword := latinWordRegExp.FindAllString(text, -1)
							arabicword := arabicWordRegExp.FindAllString(text, -1)
							greekwords = greekwords + len(words)
							wordcount := strconv.Itoa(len(words))
							latinwords = latinwords + len(latinword)
							arabicwords = arabicwords + len(arabicword)
							latinwordcount := strconv.Itoa(len(latinword))
							arabicwordcount := strconv.Itoa(len(arabicword))
							latinwordcounts = append(latinwordcounts, latinwordcount)
							arabicwordcounts = append(arabicwordcounts, arabicwordcount)
							identifiers = append(identifiers, identifier)
							texts = append(texts, text)
							greekwordcounts = append(greekwordcounts, wordcount)
						}
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']/tei:div[@n='$2']/tei:l[@n='$3']":
				fmt.Print("P")
				filecount = filecount + 1
				var data StartTEI3poem
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						for k := range data.Node[i].Node[j].Node {
							id := []string{data.Node[i].Number, data.Node[i].Node[j].Number, data.Node[i].Node[j].Node[k].Number}
							identifier := strings.Join(id, ".")
							identifier = strings.Join([]string{baseIdentifier, identifier}, ":")
							text := data.Node[i].Node[j].Node[k].InnerXML
							text = stringcleaning(text)

							words := greekWordRegExp.FindAllString(text, -1)
							latinword := latinWordRegExp.FindAllString(text, -1)
							arabicword := arabicWordRegExp.FindAllString(text, -1)
							greekwords = greekwords + len(words)
							wordcount := strconv.Itoa(len(words))
							latinwords = latinwords + len(latinword)
							arabicwords = arabicwords + len(arabicword)
							latinwordcount := strconv.Itoa(len(latinword))
							arabicwordcount := strconv.Itoa(len(arabicword))
							latinwordcounts = append(latinwordcounts, latinwordcount)
							arabicwordcounts = append(arabicwordcounts, arabicwordcount)
							identifiers = append(identifiers, identifier)
							texts = append(texts, text)
							greekwordcounts = append(greekwordcounts, wordcount)
						}
					}
				}
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']/tei:l[@n='$2']":
				fmt.Print("Q")
				filecount = filecount + 1
				var data StartTEI2Poem
				xml.Unmarshal(byteValue, &data)
				baseIdentifier := path.Base(file)
				baseIdentifier = strings.Replace(baseIdentifier, ".xml", "", -1)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						id := []string{data.Node[i].Number, data.Node[i].Node[j].Number}
						identifier := strings.Join(id, ".")
						identifier = strings.Join([]string{baseIdentifier, identifier}, ":")
						text := data.Node[i].Node[j].InnerXML
						text = stringcleaning(text)
						words := greekWordRegExp.FindAllString(text, -1)
						latinword := latinWordRegExp.FindAllString(text, -1)
						arabicword := arabicWordRegExp.FindAllString(text, -1)
						greekwords = greekwords + len(words)
						wordcount := strconv.Itoa(len(words))
						latinwords = latinwords + len(latinword)
						arabicwords = arabicwords + len(arabicword)
						latinwordcount := strconv.Itoa(len(latinword))
						arabicwordcount := strconv.Itoa(len(arabicword))
						latinwordcounts = append(latinwordcounts, latinwordcount)
						arabicwordcounts = append(arabicwordcounts, arabicwordcount)
						identifiers = append(identifiers, identifier)
						texts = append(texts, text)
						greekwordcounts = append(greekwordcounts, wordcount)
					}
				}
			default:
				querystrings = append(querystrings, querystring)
			}
		}
		xmlFile.Close()
	}
	result := removeDuplicatesUnordered(querystrings)
	if len(result) != 0 {
		fmt.Println("Not read:", len(result))
		fmt.Println("Those XPATH are unknown:", result)
	}
	fmt.Println()
	fmt.Println("Read", filecount, "of", len(xmlFiles), "files.")
	if len(noxpath) != 0 {
		fmt.Println(len(noxpath), " files have no XPATH!")
		fmt.Println("See those: ", noxpath)
	}
	fmt.Println("Write nodes to file now:")

	switch len(os.Args) {
	case 2:
		fmt.Println("Writing CEX-File")
		writeCEX(outputFile, ctscatalog, identifiers, texts)
	case 3:
		if os.Args[2] == "-CSV" {
		  fmt.Println("Writing CSV-File")
		  writeCSV(outputFile, identifiers, texts, greekwordcounts, latinwordcounts, arabicwordcounts)
     }
		if os.Args[2] == "-JSON" {
		  fmt.Println("Writing JSON-File")
		  writeJSON(outputFile, ctscatalog)
    }
		if os.Args[2] == "-XML" {
		  fmt.Println("Writing XML-File")
		  writeXML(outputFile, ctscatalog)
    }
		if os.Args[2] == "-SQL" {
		  fmt.Println("Writing SQLite DB")
		  writeSQL(outputFile, ctscatalog)
    }
	default:
		fmt.Println("Invalid number of arguments")
	}

	fmt.Println("Wrote", len(identifiers), "nodes.")
	fmt.Println(greekwords, "words written in the Greek alphabet.")
	fmt.Println(latinwords, "words written in the Latin alphabet.")
	fmt.Println(arabicwords, "words written in the Arabic alphabet.")
}

func writeCEX(outputFile string, ctscatalog CTSCatalog, identifiers, texts []string) {
	f, err := os.Create(outputFile)
	check(err)
	defer f.Close()

	// cexversion
	f.WriteString("#!cexversion")
	f.WriteString("\n\n")
	f.WriteString("3.0")
	f.WriteString("\n\n")

	// ctscatalog
	f.WriteString("#!ctscatalog")
	f.WriteString("\n\n")
	f.WriteString("urn#citationScheme#groupName#workTitle#versionLabel#exemplarLabel#online#language")
	f.WriteString("\n")
	for i := range ctscatalog.URN {
		f.WriteString(ctscatalog.URN[i])
		f.WriteString("#")
		f.WriteString(ctscatalog.CitationScheme[i])
		f.WriteString("#")
		f.WriteString(ctscatalog.GroupName[i])
		f.WriteString("#")
		f.WriteString(ctscatalog.WorkTitle[i])
		f.WriteString("#")
		f.WriteString(ctscatalog.VersionLabel[i])
		f.WriteString("#")
		f.WriteString(ctscatalog.ExemplarLabel[i])
		f.WriteString("#")
		f.WriteString(ctscatalog.Online[i])
		f.WriteString("#")
		f.WriteString(ctscatalog.Language[i])
		f.WriteString("\n")
	}
	f.WriteString("\n")

	// ctsdata
	f.WriteString("#!ctsdata")
	f.WriteString("\n\n")

	for i := range identifiers {
		newtext := strings.Replace(texts[i], "#", "", -1)
		newtext = strings.Replace(newtext, `"`, `\"`, -1)
		f.WriteString(identifiers[i])
		f.WriteString("#")
		f.WriteString(newtext)
		f.WriteString("\n")
	}
}

func getRecord(ctscatalog CTSCatalog, i int) OAI_DC_Record {
    var record OAI_DC_Record
    record = OAI_DC_Record {
      Xmlns1: "http://www.openarchives.org/OAI/2.0/oai_dc/",
      Xmlns2: "http://purl.org/dc/elements/1.1/",
      Xmlns3: "http://www.w3.org/2001/XMLSchema-instance",
      Xmlns4: "http://www.openarchives.org/OAI/2.0/oai_dc/ http://www.openarchives.org/OAI/2.0/oai_dc.xsd",
      Creator: ctscatalog.GroupName[i],
      Title: ctscatalog.WorkTitle[i],
      Subject: ctscatalog.URN[i],
      Language: ctscatalog.Language[i],
      }

    //record.Comment = "http://opengreekandlatin.github.io/First1KGreek"
    record.Description[0] = "http://cts.dh.uni-leipzig.de/text/urn:cts:greekLit:" + ctscatalog.URN[i]
    return(record)
}

func writeSQL(outputFile string, ctscatalog CTSCatalog) {

  db, err := sql.Open("sqlite3", outputFile)
  check(err)
  records, _  := db.Prepare("INSERT INTO records(id, item_id, metadata_format_id, xml, state) values(? ,?, 1, ?, 1)")
  items, _  := db.Prepare("INSERT INTO items(id, id_ext, state, timestamp) values(? ,?, 'active', '1970-01-01 00:00:00')")

  for i := range ctscatalog.URN {
    output, err := xml.MarshalIndent(getRecord(ctscatalog, i), "", " ")
    if err != nil {
      fmt.Printf("error: %v\n", err)
    }
    //os.Stdout.Write(output)
    fmt.Print(".")
    records.Exec(i, i, output)
    items.Exec(i, ctscatalog.URN[i])
    check(err)
  }
  db.Close()

}
func writeXML(outputFile string, ctscatalog CTSCatalog) {
  f, err := os.Create(outputFile)
  check(err)
  defer f.Close()

  for i := range ctscatalog.URN {
    output, err := xml.MarshalIndent(getRecord(ctscatalog, i), "", " ")
    if err != nil {
      fmt.Printf("error: %v\n", err)
    }
    //os.Stdout.Write(output)
    f.Write(output) 
  }
}

func writeJSON(outputFile string, ctscatalog CTSCatalog) {
  f, err := os.Create(outputFile)
  check(err)
  defer f.Close()

  f.WriteString(string('['))
  for i := range ctscatalog.URN {
    var d ExportDocument
    d.URN = ctscatalog.URN[i]
    d.CitationScheme = ctscatalog.CitationScheme[i]
    d.GroupName = ctscatalog.GroupName[i]
    d.WorkTitle = ctscatalog.WorkTitle[i]
    d.VersionLabel = ctscatalog.VersionLabel[i]
    d.ExemplarLabel = ctscatalog.ExemplarLabel[i]
    d.Online = ctscatalog.Online[i]
    d.Language = ctscatalog.Language[i]
    b, _ := json.Marshal(d)
    f.WriteString(string(b)) 
    if i < len(ctscatalog.URN)-1 { f.WriteString(string(',')) }
  }
  f.WriteString(string(']'))
}

func writeCSV(outputFile string, identifiers, texts, greekwordcounts, latinwordcounts, arabicwordcounts []string) {
	f, err := os.Create(outputFile)
	check(err)
	defer f.Close()

	f.WriteString("identifier")
	f.WriteString("#")
	f.WriteString("text")
	f.WriteString("#")
	f.WriteString("GreekWords")
	f.WriteString("#")
	f.WriteString("LatinWords")
	f.WriteString("#")
	f.WriteString("ArabicWords")
	f.WriteString("#")
	f.WriteString("Workgroup")
	f.WriteString("#")
	f.WriteString("Work")
	f.WriteString("#")
	f.WriteString("WorkVerbose")
	f.WriteString("\n")

	for i := range identifiers {
		newtext := strings.Replace(texts[i], "#", "", -1)
		newtext = strings.Replace(newtext, `"`, `\"`, -1)
		f.WriteString(identifiers[i])
		f.WriteString("#")
		f.WriteString(newtext)
		f.WriteString("#")
		f.WriteString(greekwordcounts[i])
		f.WriteString("#")
		f.WriteString(latinwordcounts[i])
		f.WriteString("#")
		f.WriteString(arabicwordcounts[i])
		f.WriteString("#")
		baseurn := strings.Split(identifiers[i], ":")[0]
		urnslice := strings.Split(baseurn, ".")
		workgroup := urnslice[0]
		work := strings.Join(urnslice[1:len(urnslice)], ".")
		f.WriteString(workgroup)
		f.WriteString("#")
		f.WriteString(work)
		f.WriteString("#")
		f.WriteString(baseurn)
		f.WriteString("\n")
	}
}

func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	for v := range elements {
		encountered[elements[v]] = true
	}

	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

func stringcleaning(text string) string {
	tagsRegExp := regexp.MustCompile(`<[/]*[^>]*>`)
	reInsideWhtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	result := text
	result = strings.Replace(result, "\n", " ", -1)
	result = strings.Replace(result, "#", "", -1)
	result = tagsRegExp.ReplaceAllString(result, "")
	result = strings.TrimSpace(result)
	result = reInsideWhtsp.ReplaceAllString(result, " ")
	return result
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkExt(ext string) []string {
	pathS, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var files []string
	filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				if f.Name() != "__cts__.xml" && f.Name() != "build.xml" && f.Name() != "expath-pkg.xml" && f.Name() != "repo.xml" {
					files = append(files, path)
				}
			}
		}
		return nil
	})
	return files
}
