package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// ReportJSON is a struct for exporting to JSON
type ReportJSON struct {
	Nodecount   int           `json:"nodeCount"`
	Greekwords  int           `json:"greekWords"`
	Latinwords  int           `json:"latinWords"`
	Arabicwords int           `json:"arabicwords"`
	Catalog     []JSONCatalog `json:"catalog"`
}

// JSONCatalog is a struct for exporting to JSON. It is a sub-struct to ReportJSON.
type JSONCatalog struct {
	URN       string `json:"urn"`
	GroupName string `json:"group_name"`
	WorkName  string `json:"work_name"`
	Language  string `json:"language"`
	WordCount int    `json:"wordcount"`
	Scaife    string `json:"scaife"`
}

//CTSCatalog is the main container for CTS catalog data in the format expected by CEX, but in a way that it can integrated into a number of ways.
type CTSCatalog struct {
	URN            []string `json:"urn"`
	CitationScheme []string `json:"citation_scheme"`
	GroupName      []string `json:"group_name"`
	WorkTitle      []string `json:"work_title"`
	VersionLabel   []string `json:"version_label"`
	ExemplarLabel  []string `json:"exemplar_label"`
	Online         []string `json:"online"`
	Language       []string `json:"language"`
}

//ExportDocument is the old JSON implementation. Preserved for legacy.
type ExportDocument struct {
	URN            string `json:"urn"`
	CitationScheme string `json:"citation_scheme"`
	GroupName      string `json:"group_name"`
	WorkTitle      string `json:"work_title"`
	VersionLabel   string `json:"version_label"`
	ExemplarLabel  string `json:"exemplar_label"`
	Online         string `json:"online"`
	Language       string `json:"language"`
}

//OAIDCRecord is a container to produce DataCite metadata.
type OAIDCRecord struct {
	//XMLName  xml.Name `xml:"http://www.openarchives.org/OAI/2.0/oai_dc/ oai_dc:dc"`
	XMLName     xml.Name  `xml:"oai_dc:dc"`
	Xmlns1      string    `xml:"xmlns:oai_dc,attr"`
	Xmlns2      string    `xml:"xmlns:dc,attr"`
	Xmlns3      string    `xml:"xmlns:xsi,attr"`
	Xmlns4      string    `xml:"xsi:schemaLocation,attr"`
	ID          int       `xml:"id,attr"`
	Title       string    `xml:"dc:title"`
	Creator     string    `xml:"dc:creator"`
	Subject     string    `xml:"dc:subject"`
	Description [2]string `xml:"dc:description,omitempty"`
	Comment     string    `xml:",comment"`
	//Date    string    `xml:"dc:date"`
	Language  string `xml:"dc:language"`
	ViewURL   string `xml:"dc:view-url"`
	Publisher string `xml:"dc:publisher"`
}

//Metadata container for Xpath metadata
type Metadata struct {
	Xpath string
	Kind  string
}

//XPathInfo container for Xpath metadata
type XPathInfo struct {
	XPathInfo string `xml:"replacementPattern,attr"`
	XPathWhat string `xml:"n,attr"`
}

//LangInfo container for language metadata
type LangInfo struct {
	Language string `xml:"ident,attr"`
}

//RefPattern container for refpattern
type RefPattern struct {
	RefPattern []XPathInfo `xml:"teiHeader>encodingDesc>refsDecl>cRefPattern"`
	Title      []string    `xml:"teiHeader>fileDesc>titleStmt>title"`
	Author     []string    `xml:"teiHeader>fileDesc>titleStmt>author"`
	Languages  []LangInfo  `xml:"teiHeader>profileDesc>langUsage>language"`
}

// type teiHeader struct {
// 	RefPattern []XPathInfo `xml:"teiHeader>encodingDesc>refsDecl>cRefPattern"`
// }

//SmallestNode container for smallest CTS node
type SmallestNode struct {
	InnerXML string `xml:",innerxml"`
	Number   string `xml:"n,attr"`
}

//StartTEI1Direct is one of many parsing containers.
type StartTEI1Direct struct {
	Node []SmallestNode `xml:"text>body>div"`
}

//StartTEI1p is one of many parsing containers.
type StartTEI1p struct {
	Node []SmallestNode `xml:"text>body>div>p"`
}

//StartTEI1pseg is one of many parsing containers.
type StartTEI1pseg struct {
	Node []SmallestNode `xml:"text>body>div>p>seg"`
}

//StartTEI1 is one of many parsing containers.
type StartTEI1 struct {
	Node []SmallestNode `xml:"text>body>div>div"`
}

//StartTEI1Late is one of many parsing containers.
type StartTEI1Late struct {
	Node []SmallestNode `xml:"text>body>div>div>div"`
}

//TEI3n2 is one of many parsing containers.
type TEI3n2 struct {
	Node   []SmallestNode `xml:"div"`
	Number string         `xml:"n,attr"`
}

//TEI3n3 is one of many parsing containers.
type TEI3n3 struct {
	Node   []TEI3n2 `xml:"div"`
	Number string   `xml:"n,attr"`
}

//StartTEI3 is one of many parsing containers.
type StartTEI3 struct {
	Node []TEI3n3 `xml:"text>body>div>div"`
}

//TEI3n2p is one of many parsing containers.
type TEI3n2p struct {
	Node   []SmallestNode `xml:"p"`
	Number string         `xml:"n,attr"`
}

//TEI3n3p is one of many parsing containers.
type TEI3n3p struct {
	Node   []TEI3n2p `xml:"div"`
	Number string    `xml:"n,attr"`
}

//StartTEI3p is one of many parsing containers.
type StartTEI3p struct {
	Node []TEI3n3p `xml:"text>body>div>div"`
}

//TEI4n2p is one of many parsing containers.
type TEI4n2p struct {
	Node   []SmallestNode `xml:"p"`
	Number string         `xml:"n,attr"`
}

//TEI4n3p is one of many parsing containers.
type TEI4n3p struct {
	Node   []TEI4n2p `xml:"div"`
	Number string    `xml:"n,attr"`
}

//TEI4n4p is one of many parsing containers.
type TEI4n4p struct {
	Node   []TEI4n3p `xml:"div"`
	Number string    `xml:"n,attr"`
}

//StartTEI4p is one of many parsing containers.
type StartTEI4p struct {
	Node []TEI4n4p `xml:"text>body>div>div"`
}

//TEI4n2div is one of many parsing containers.
type TEI4n2div struct {
	Node   []SmallestNode `xml:"div"`
	Number string         `xml:"n,attr"`
}

//TEI4n3div is one of many parsing containers.
type TEI4n3div struct {
	Node   []TEI4n2div `xml:"div"`
	Number string      `xml:"n,attr"`
}

//TEI4n4div is one of many parsing containers.
type TEI4n4div struct {
	Node   []TEI4n3div `xml:"div"`
	Number string      `xml:"n,attr"`
}

//StartTEI4div is one of many parsing containers.
type StartTEI4div struct {
	Node []TEI4n4div `xml:"text>body>div>div"`
}

//StartTEI2p is one of many parsing containers.
type StartTEI2p struct {
	Node []TEI3n2p `xml:"text>body>div>div"`
}

//TEI2n2ab is one of many parsing containers.
type TEI2n2ab struct {
	Node   []SmallestNode `xml:"ab"`
	Number string         `xml:"n,attr"`
}

//StartTEI2ab is one of many parsing containers.
type StartTEI2ab struct {
	Node []TEI2n2ab `xml:"text>body>div>div"`
}

//TEI2n2lgl is one of many parsing containers.
type TEI2n2lgl struct {
	Node   []SmallestNode `xml:"lg>l"`
	Number string         `xml:"n,attr"`
}

//StartTEI2lgl is one of many parsing containers.
type StartTEI2lgl struct {
	Node []TEI2n2lgl `xml:"text>body>div>div"`
}

//TEI3n2cit is one of many parsing containers.
type TEI3n2cit struct {
	Node   []SmallestNode `xml:"cit"`
	Number string         `xml:"n,attr"`
}

//TEI3n3cit is one of many parsing containers.
type TEI3n3cit struct {
	Node   []TEI3n2cit `xml:"p"`
	Number string      `xml:"n,attr"`
}

//StartTEI3cit is one of many parsing containers.
type StartTEI3cit struct {
	Node []TEI3n3cit `xml:"text>body>div>div"`
}

//TEI3n3divcit is one of many parsing containers.
type TEI3n3divcit struct {
	Node   []TEI3n2cit `xml:"div"`
	Number string      `xml:"n,attr"`
}

//StartTEI3divcit is one of many parsing containers.
type StartTEI3divcit struct {
	Node []TEI3n3divcit `xml:"text>body>div>div"`
}

//StartTEI2 is one of many parsing containers.
type StartTEI2 struct {
	Node []TEI3n2 `xml:"text>body>div>div"`
}

//StartTEI2direct is one of many parsing containers.
type StartTEI2direct struct {
	Node []TEI3n2 `xml:"text>body>div"`
}

//StartTEI1Poem is one of many parsing containers.
type StartTEI1Poem struct {
	Node []SmallestNode `xml:"text>body>div>l"`
}

//TEI2Poemn2 is one of many parsing containers.
type TEI2Poemn2 struct {
	Node   []SmallestNode `xml:"l"`
	Number string         `xml:"n,attr"`
}

//StartTEI2Poem is one of many parsing containers.
type StartTEI2Poem struct {
	Node []TEI2Poemn2 `xml:"text>body>div>div"`
}

//TEI3n2DirectNumbered is one of many parsing containers.
type TEI3n2DirectNumbered struct {
	Node   []SmallestNode `xml:"div3"`
	Number string         `xml:"n,attr"`
}

//TEI3n3DirectNumbered is one of many parsing containers.
type TEI3n3DirectNumbered struct {
	Node   []TEI3n2DirectNumbered `xml:"div2"`
	Number string                 `xml:"n,attr"`
}

//StartTEI3DirectNumbered is one of many parsing containers.
type StartTEI3DirectNumbered struct {
	Node []TEI3n3DirectNumbered `xml:"text>body>div1"`
}

//TEI3n2poem is one of many parsing containers.
type TEI3n2poem struct {
	Node   []SmallestNode `xml:"l"`
	Number string         `xml:"n,attr"`
}

//TEI3n3poem is one of many parsing containers.
type TEI3n3poem struct {
	Node   []TEI3n2poem `xml:"div"`
	Number string       `xml:"n,attr"`
}

//StartTEI3poem is one of many parsing containers.
type StartTEI3poem struct {
	Node []TEI3n3poem `xml:"text>body>div>div"`
}

//QSmallestNode is one of many parsing containers.
type QSmallestNode struct {
	InnerXML string `xml:",innerxml"`
	Number   string `xml:"n,attr"`
}

//QueryTEI2 is one of many parsing containers.
type QueryTEI2 struct {
	Node []QSmallestNode `xml:"text>body>div>div"`
}

//QueryTEI3n3 is one of many parsing containers.
type QueryTEI3n3 struct {
	Node   []QSmallestNode `xml:"div"`
	Number string          `xml:"n,attr"`
}

//QueryTEI3 is one of many parsing containers.
type QueryTEI3 struct {
	Node []QueryTEI3n3 `xml:"text>body>div>div"`
}

//QueryTEI1 is one of many parsing containers.
type QueryTEI1 struct {
	Node []QSmallestNode `xml:"text>body>div"`
}

//QueryTEI1p is one of many parsing containers.
type QueryTEI1p struct {
	Node []QSmallestNode `xml:"text>body>div>p"`
}

//QSmallestNodeDiv is one of many parsing containers.
type QSmallestNodeDiv struct {
	InnerXML string `xml:",innerxml"`
}

//QueryTEI1div is one of many parsing containers.
type QueryTEI1div struct {
	Node []QSmallestNodeDiv `xml:"text>body>div>div"`
}

//QueryTEI0 is one of many parsing containers.
type QueryTEI0 struct {
	Node []QSmallestNode `xml:"text>body"`
}

//QueryInfo is one of many parsing containers.
type QueryInfo struct {
	InnerXML string `xml:",innerxml"`
	Number   string `xml:"n,attr"`
}

func main() {
	scheme := make(map[string]int)
	tempscheme := ""
	outputFile := ""
	switch len(os.Args) {
	case 1:
		fmt.Println("Usage: CTSExtract [output-filename] [optionally: -CSV|JSON|XML|HTML]")
		os.Exit(3)
	case 2, 3:
		outputFile = os.Args[1]
	default:
		fmt.Println("Usage: CTSExtract [output-filename] [optionally: -CSV|JSON|XML|HTML]")
		os.Exit(3)
	}
	basereg := regexp.MustCompile(`urn:\p{L}+:\p{L}+:`)
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
		basestr := "urn:cts:greekLit:"
		xmlFile, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(file)
		byteValue, _ := ioutil.ReadAll(xmlFile)
		match := basereg.FindStringSubmatch(string(byteValue))
		if len(match) > 0 {
			basestr = match[0]
		}
		var xpathinfo RefPattern
		err = xml.Unmarshal(byteValue, &xpathinfo)
		check(err)
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
			urn = basestr + urn
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
				tempscheme = "1"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI4p
				err = xml.Unmarshal(byteValue, &data)
				check(err)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						for k := range data.Node[i].Node[j].Node {
							for l := range data.Node[i].Node[j].Node[k].Node {
								id := []string{data.Node[i].Number, data.Node[i].Node[j].Number, data.Node[i].Node[j].Node[k].Number, data.Node[i].Node[j].Node[k].Node[l].Number}
								identifier := strings.Join(id, ".")
								identifier = strings.Join([]string{urn, identifier}, ":")
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
				tempscheme = "2"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data QueryTEI1
				err = xml.Unmarshal(byteValue, &data)
				check(err)
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
								err = decoder.DecodeElement(&info, &se)
								check(err)
								identifier := strings.Join([]string{urn, info.Number}, ":")
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
				tempscheme = "3"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data QueryTEI1div
				err = xml.Unmarshal(byteValue, &data)
				check(err)
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
								err = decoder.DecodeElement(&info, &se)
								check(err)
								identifier := strings.Join([]string{urn, info.Number}, ":")
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
				tempscheme = "4"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data QueryTEI1p
				xml.Unmarshal(byteValue, &data)
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
								err = decoder.DecodeElement(&info, &se)
								check(err)
								identifier := strings.Join([]string{urn, info.Number}, ":")
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
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div//tei:l[@n='$1']" || querystring == "/tei:TEI/tei:text/tei:body/tei:div//tei:l[@n=\\'$1\\']" || querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:sp/tei:l[@n='$1']":
				tempscheme = "5"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data QueryTEI1
				xml.Unmarshal(byteValue, &data)
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
								err = decoder.DecodeElement(&info, &se)
								check(err)
								identifier := strings.Join([]string{urn, info.Number}, ":")
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
				tempscheme = "6"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data QueryTEI0
				xml.Unmarshal(byteValue, &data)
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
								err = decoder.DecodeElement(&info, &se)
								check(err)
								identifier := strings.Join([]string{urn, info.Number}, ":")
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
				tempscheme = "7"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data QueryTEI2
				xml.Unmarshal(byteValue, &data)
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
								err = decoder.DecodeElement(&info, &se)
								check(err)
								id := []string{ID2, info.Number}
								identifier := strings.Join(id, ".")
								identifier = strings.Join([]string{urn, identifier}, ":")
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
				tempscheme = "8"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data QueryTEI3
				xml.Unmarshal(byteValue, &data)
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
									identifier = strings.Join([]string{urn, identifier}, ":")
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
				tempscheme = "9"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data QueryTEI2
				xml.Unmarshal(byteValue, &data)
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
								identifier = strings.Join([]string{urn, identifier}, ":")
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
				tempscheme = "A"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI3p
				xml.Unmarshal(byteValue, &data)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						for k := range data.Node[i].Node[j].Node {
							id := []string{data.Node[i].Number, data.Node[i].Node[j].Number, data.Node[i].Node[j].Node[k].Number}
							identifier := strings.Join(id, ".")
							identifier = strings.Join([]string{urn, identifier}, ":")
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
				tempscheme = "B"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI3cit
				xml.Unmarshal(byteValue, &data)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						for k := range data.Node[i].Node[j].Node {
							id := []string{data.Node[i].Number, data.Node[i].Node[j].Number, data.Node[i].Node[j].Node[k].Number}
							identifier := strings.Join(id, ".")
							identifier = strings.Join([]string{urn, identifier}, ":")
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
				tempscheme = "C"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI1Direct
				xml.Unmarshal(byteValue, &data)
				for _, node := range data.Node {
					identifier := strings.Join([]string{urn, node.Number}, ":")
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
				tempscheme = "D"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI1Late
				xml.Unmarshal(byteValue, &data)
				for _, node := range data.Node {
					identifier := strings.Join([]string{urn, node.Number}, ":")
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
				tempscheme = "E"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI1pseg
				xml.Unmarshal(byteValue, &data)
				for _, node := range data.Node {
					identifier := strings.Join([]string{urn, node.Number}, ":")
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
				tempscheme = "F"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI1p
				xml.Unmarshal(byteValue, &data)
				for _, node := range data.Node {
					identifier := strings.Join([]string{urn, node.Number}, ":")
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
				tempscheme = "G"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI1
				xml.Unmarshal(byteValue, &data)
				for _, node := range data.Node {
					identifier := strings.Join([]string{urn, node.Number}, ":")
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
				tempscheme = "H"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI2p
				xml.Unmarshal(byteValue, &data)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						id := []string{data.Node[i].Number, data.Node[i].Node[j].Number}
						identifier := strings.Join(id, ".")
						identifier = strings.Join([]string{urn, identifier}, ":")
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
				tempscheme = "I"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI2ab
				xml.Unmarshal(byteValue, &data)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						id := []string{data.Node[i].Number, data.Node[i].Node[j].Number}
						identifier := strings.Join(id, ".")
						identifier = strings.Join([]string{urn, identifier}, ":")
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
				tempscheme = "J"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI2lgl
				xml.Unmarshal(byteValue, &data)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						id := []string{data.Node[i].Number, data.Node[i].Node[j].Number}
						identifier := strings.Join(id, ".")
						identifier = strings.Join([]string{urn, identifier}, ":")
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
				tempscheme = "K"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI2
				xml.Unmarshal(byteValue, &data)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						id := []string{data.Node[i].Number, data.Node[i].Node[j].Number}
						identifier := strings.Join(id, ".")
						identifier = strings.Join([]string{urn, identifier}, ":")
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
				tempscheme = "L"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI3
				xml.Unmarshal(byteValue, &data)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						for k := range data.Node[i].Node[j].Node {
							id := []string{data.Node[i].Number, data.Node[i].Node[j].Number, data.Node[i].Node[j].Node[k].Number}
							identifier := strings.Join(id, ".")
							identifier = strings.Join([]string{urn, identifier}, ":")
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
				tempscheme = "M"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI1Poem
				xml.Unmarshal(byteValue, &data)
				for _, node := range data.Node {
					identifier := strings.Join([]string{urn, node.Number}, ":")
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
				tempscheme = "N"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI2direct
				xml.Unmarshal(byteValue, &data)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						id := []string{data.Node[i].Number, data.Node[i].Node[j].Number}
						identifier := strings.Join(id, ".")
						identifier = strings.Join([]string{urn, identifier}, ":")
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
				tempscheme = "O"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI3DirectNumbered
				xml.Unmarshal(byteValue, &data)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						for k := range data.Node[i].Node[j].Node {
							id := []string{data.Node[i].Number, data.Node[i].Node[j].Number, data.Node[i].Node[j].Node[k].Number}
							identifier := strings.Join(id, ".")
							identifier = strings.Join([]string{urn, identifier}, ":")
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
				tempscheme = "P"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI3poem
				xml.Unmarshal(byteValue, &data)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						for k := range data.Node[i].Node[j].Node {
							id := []string{data.Node[i].Number, data.Node[i].Node[j].Number, data.Node[i].Node[j].Node[k].Number}
							identifier := strings.Join(id, ".")
							identifier = strings.Join([]string{urn, identifier}, ":")
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
				tempscheme = "Q"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI2Poem
				xml.Unmarshal(byteValue, &data)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						id := []string{data.Node[i].Number, data.Node[i].Node[j].Number}
						identifier := strings.Join(id, ".")
						identifier = strings.Join([]string{urn, identifier}, ":")
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
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']/tei:div[@n='$2']/tei:div[@n='$3']/tei:div[@n='$4']":
				tempscheme = "R"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI4div
				xml.Unmarshal(byteValue, &data)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						for k := range data.Node[i].Node[j].Node {
							for l := range data.Node[i].Node[j].Node[k].Node {
								id := []string{data.Node[i].Number, data.Node[i].Node[j].Number, data.Node[i].Node[j].Node[k].Number, data.Node[i].Node[j].Node[k].Node[l].Number}
								identifier := strings.Join(id, ".")
								identifier = strings.Join([]string{urn, identifier}, ":")
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
			case querystring == "/tei:TEI/tei:text/tei:body/tei:div/tei:div[@n='$1']/tei:div[@n='$2']/tei:cit[@n='$3']":
				tempscheme = "S"
				fmt.Print(tempscheme)
				scheme[tempscheme] = scheme[tempscheme] + 1
				filecount = filecount + 1
				var data StartTEI3divcit
				xml.Unmarshal(byteValue, &data)
				for i := range data.Node {
					for j := range data.Node[i].Node {
						for k := range data.Node[i].Node[j].Node {
							id := []string{data.Node[i].Number, data.Node[i].Node[j].Number, data.Node[i].Node[j].Node[k].Number}
							identifier := strings.Join(id, ".")
							identifier = strings.Join([]string{urn, identifier}, ":")
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
		if os.Args[2] == "-HTML" {
			fmt.Println("Writing HTML Report")
			writeHTML(outputFile, ctscatalog, identifiers, texts, greekwordcounts, latinwordcounts, arabicwordcounts, greekwords, latinwords, arabicwords)
		}
		if os.Args[2] == "-Cat" {
			fmt.Println("Writing JSON Catalog")
			var jsoncat = []JSONCatalog{}
			for i := range ctscatalog.URN {
				scaifestring := ""
				itemwords := 0
				found := false
				for j, v := range identifiers {
					if !found {
						if strings.Contains(v, ctscatalog.URN[i]) {
							scaifestring = "https://scaife.perseus.org/reader/" + v
							found = true
						}
					}
					if strings.Contains(v, ctscatalog.URN[i]) {
						greek, _ := strconv.Atoi(greekwordcounts[j])
						latin, _ := strconv.Atoi(latinwordcounts[j])
						arabic, _ := strconv.Atoi(arabicwordcounts[j])
						itemwords = itemwords + greek + latin + arabic
					}
				}
				catitem := JSONCatalog{
					URN:       ctscatalog.URN[i],
					GroupName: ctscatalog.GroupName[i],
					WorkName:  ctscatalog.WorkTitle[i],
					Language:  ctscatalog.Language[i],
					WordCount: itemwords,
					Scaife:    scaifestring,
				}
				jsoncat = append(jsoncat, catitem)
			}
			var report = ReportJSON{Nodecount: len(identifiers),
				Greekwords:  greekwords,
				Latinwords:  latinwords,
				Arabicwords: arabicwords,
				Catalog:     jsoncat}
			writeCatalog(outputFile, report)
		}
	default:
		fmt.Println("Invalid number of arguments")
	}

	fmt.Println("Wrote", len(identifiers), "nodes.")
	fmt.Println(greekwords, "words written in the Greek alphabet.")
	fmt.Println(latinwords, "words written in the Latin alphabet.")
	fmt.Println(arabicwords, "words written in the Arabic alphabet.")
	fmt.Println("The following schemes were used:")
	for i, v := range scheme {
		fmt.Println(i, v)
	}
}

func writeCatalog(outputFile string, report ReportJSON) {
	jsonreport, err1 := json.Marshal(report)
	check(err1)
	f, err2 := os.Create(outputFile)
	check(err2)
	defer f.Close()
	_, err := f.WriteString(string(jsonreport))
	check(err)
}

type fileConnection struct {
	*os.File
}

func (f fileConnection) writeToFile(input string) {
	_, err := f.WriteString(input)
	check(err)
}

func writeHTML(outputFile string, ctscatalog CTSCatalog, identifiers, texts, greekwordcounts, latinwordcounts, arabicwordcounts []string, greekwords, latinwords, arabicwords int) {
	f, err := os.Create(outputFile)
	check(err)
	defer f.Close()
	fconnection := fileConnection{f}
	// HTML Header
	// HTML BODY
	fconnection.writeToFile("<div>\n")
	fconnection.writeToFile("<p>")
	fconnection.writeToFile("Greek words:" + strconv.Itoa(greekwords))
	fconnection.writeToFile("</p>\n")
	fconnection.writeToFile("<p>")
	fconnection.writeToFile("Latin words:" + strconv.Itoa(latinwords))
	fconnection.writeToFile("</p>\n")
	fconnection.writeToFile("<p>")
	fconnection.writeToFile("Arabic words:" + strconv.Itoa(arabicwords))
	fconnection.writeToFile("</p>\n")
	fconnection.writeToFile("</div>\n")
	fconnection.writeToFile("</hr>\n")
	for i := range ctscatalog.URN {
		fconnection.writeToFile("<div>\n")
		fconnection.writeToFile("<h3>")
		fconnection.writeToFile("URN:" + ctscatalog.URN[i])
		fconnection.writeToFile("</h3>\n")
		fconnection.writeToFile("<p>")
		fconnection.writeToFile("CitationScheme:" + ctscatalog.CitationScheme[i])
		fconnection.writeToFile("</p>\n")
		fconnection.writeToFile("<p>")
		fconnection.writeToFile("GroupName:" + ctscatalog.GroupName[i])
		fconnection.writeToFile("</p>\n")
		fconnection.writeToFile("<p>")
		fconnection.writeToFile("WorkTitle:" + ctscatalog.WorkTitle[i])
		fconnection.writeToFile("</p>\n")
		fconnection.writeToFile("<p>")
		fconnection.writeToFile("VersionLabel:" + ctscatalog.VersionLabel[i])
		fconnection.writeToFile("</p>\n")
		fconnection.writeToFile("<p>")
		fconnection.writeToFile("ExemplarLabel:" + ctscatalog.ExemplarLabel[i])
		fconnection.writeToFile("</p>\n")
		fconnection.writeToFile("<p>")
		fconnection.writeToFile("Language:" + ctscatalog.Language[i])
		fconnection.writeToFile("</p>\n")
		found := false
		wordcount := 0
		for j, v := range identifiers {
			if !found {
				if strings.Contains(v, ctscatalog.URN[i]) {
					fconnection.writeToFile("<p>")
					fconnection.writeToFile("First URN:" + v)
					fconnection.writeToFile("</p>\n")
					fconnection.writeToFile("<p>")
					fconnection.writeToFile("<a href=\"https://scaife.perseus.org/reader/" + v + "\">Read Online</a>")
					fconnection.writeToFile("</p>\n")
					found = true
				}
			}
			if strings.Contains(v, ctscatalog.URN[i]) {
				greek, _ := strconv.Atoi(greekwordcounts[j])
				latin, _ := strconv.Atoi(latinwordcounts[j])
				arabic, _ := strconv.Atoi(arabicwordcounts[j])
				wordcount = wordcount + greek + latin + arabic
			}
		}
		fconnection.writeToFile("<p>")
		fconnection.writeToFile("Words:" + strconv.Itoa(wordcount))
		fconnection.writeToFile("</p>\n")
		fconnection.writeToFile("</div>\n")
		fconnection.writeToFile("</hr>\n")
	}
}

func writeCEX(outputFile string, ctscatalog CTSCatalog, identifiers, texts []string) {
	f, err := os.Create(outputFile)
	check(err)
	fconnection := fileConnection{f}
	defer f.Close()

	// cexversion
	fconnection.writeToFile("#!cexversion")
	fconnection.writeToFile("\n\n")
	fconnection.writeToFile("3.0")
	fconnection.writeToFile("\n\n")

	// ctscatalog
	fconnection.writeToFile("#!ctscatalog")
	fconnection.writeToFile("\n\n")
	fconnection.writeToFile("urn#citationScheme#groupName#workTitle#versionLabel#exemplarLabel#online#language")
	fconnection.writeToFile("\n")
	for i := range ctscatalog.URN {
		fconnection.writeToFile(ctscatalog.URN[i])
		fconnection.writeToFile("#")
		fconnection.writeToFile(ctscatalog.CitationScheme[i])
		fconnection.writeToFile("#")
		fconnection.writeToFile(ctscatalog.GroupName[i])
		fconnection.writeToFile("#")
		fconnection.writeToFile(ctscatalog.WorkTitle[i])
		fconnection.writeToFile("#")
		fconnection.writeToFile(ctscatalog.VersionLabel[i])
		fconnection.writeToFile("#")
		fconnection.writeToFile(ctscatalog.ExemplarLabel[i])
		fconnection.writeToFile("#")
		fconnection.writeToFile(ctscatalog.Online[i])
		fconnection.writeToFile("#")
		fconnection.writeToFile(ctscatalog.Language[i])
		fconnection.writeToFile("\n")
	}
	fconnection.writeToFile("\n")

	// ctsdata
	fconnection.writeToFile("#!ctsdata")
	fconnection.writeToFile("\n\n")

	for i := range identifiers {
		newtext := strings.Replace(texts[i], "#", "", -1)
		newtext = strings.Replace(newtext, `"`, `\"`, -1)
		fconnection.writeToFile(identifiers[i])
		fconnection.writeToFile("#")
		fconnection.writeToFile(newtext)
		fconnection.writeToFile("\n")
	}
}

func getRecord(ctscatalog CTSCatalog, i int) (record OAIDCRecord) {
	record = OAIDCRecord{
		Xmlns1:   "http://www.openarchives.org/OAI/2.0/oai_dc/",
		Xmlns2:   "http://purl.org/dc/elements/1.1/",
		Xmlns3:   "http://www.w3.org/2001/XMLSchema-instance",
		Xmlns4:   "http://www.openarchives.org/OAI/2.0/oai_dc/ http://www.openarchives.org/OAI/2.0/oai_dc.xsd",
		Creator:  ctscatalog.GroupName[i],
		Title:    ctscatalog.WorkTitle[i],
		Subject:  ctscatalog.URN[i],
		Language: ctscatalog.Language[i],
	}

	//record.Comment = "http://opengreekandlatin.github.io/First1KGreek"
	record.Publisher = "OGLP"
	record.ViewURL = "http://cts.dh.uni-leipzig.de/text/urn:cts:greekLit:" + ctscatalog.URN[i]
	record.Description[0] = "http://cts.dh.uni-leipzig.de/text/urn:cts:greekLit:" + ctscatalog.URN[i]
	return (record)
}

func writeSQL(outputFile string, ctscatalog CTSCatalog) {

	var record OAIDCRecord
	db, err := sql.Open("sqlite3", outputFile)
	check(err)
	records, _ := db.Prepare("INSERT INTO records(id, item_id, metadata_format_id, xml, state) values(? ,?, 1, ?, 1)")
	items, _ := db.Prepare("INSERT INTO items(id, id_ext, state, timestamp) values(? ,?, 'active', '1970-01-01 00:00:00')")

	for i := range ctscatalog.URN {
		record = getRecord(ctscatalog, i)
		if record.Creator != "" {
			output, err := xml.MarshalIndent(record, "", " ")
			check(err)
			fmt.Print(".")
			_, err = records.Exec(i, i, output)
			check(err)
			_, err = items.Exec(i, ctscatalog.URN[i])
			check(err)
		}
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
		_, err = f.Write(output)
		check(err)
	}
}

func writeJSON(outputFile string, ctscatalog CTSCatalog) {
	f, err := os.Create(outputFile)
	check(err)
	defer f.Close()
	_, err = f.WriteString(string('['))
	check(err)
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
		_, err = f.WriteString(string(b))
		check(err)
		if i < len(ctscatalog.URN)-1 {
			_, err = f.WriteString(string(','))
			check(err)
		}
	}
	_, err = f.WriteString(string(']'))
	check(err)
}

func writeCSV(outputFile string, identifiers, texts, greekwordcounts, latinwordcounts, arabicwordcounts []string) {
	f, err := os.Create(outputFile)
	check(err)
	defer f.Close()
	fconnection := fileConnection{f}

	fconnection.writeToFile("identifier#text#GreekWords#LatinWords#ArabicWords#Workgroup#Work#WorkVerbose\n")

	for i := range identifiers {
		newtext := strings.Replace(texts[i], "#", "", -1)
		newtext = strings.Replace(newtext, `"`, `\"`, -1)
		fconnection.writeToFile(identifiers[i])
		fconnection.writeToFile("#")
		fconnection.writeToFile(newtext)
		fconnection.writeToFile("#")
		fconnection.writeToFile(greekwordcounts[i])
		fconnection.writeToFile("#")
		fconnection.writeToFile(latinwordcounts[i])
		fconnection.writeToFile("#")
		fconnection.writeToFile(arabicwordcounts[i])
		fconnection.writeToFile("#")
		baseurn := strings.Split(identifiers[i], ":")[3]
		urnslice := strings.Split(baseurn, ".")
		workgroup := urnslice[0]
		work := strings.Join(urnslice[1:], ".")
		fconnection.writeToFile(workgroup)
		fconnection.writeToFile("#")
		fconnection.writeToFile(work)
		fconnection.writeToFile("#")
		fconnection.writeToFile(baseurn)
		fconnection.writeToFile("\n")
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
	err2 := filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
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
	check(err2)
	return files
}
