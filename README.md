# TEItoCEX
Turn CTS TEI corpora like CHS and OGL's [First1KGreek](http://opengreekandlatin.github.io/First1KGreek/) into [CEX collection files](https://github.com/cite-architecture/citedx/blob/master/docs/CEX-spec-3.0.md). 

# USAGE OSX

1. Copy `CTSExtract` into the unpacked data folder of e.g. First1Greek. 
2. Open a terminal in that folder and type: `./CTSExtract 1kGreek.cex ` (you might have to chmod +x the executable before you can use it)
3. Enjoy your new CEX collection file!

Sample output:

```
.................................................................
Read 860 of 860 files.
Write nodes to file now:
Wrote 186909 nodes.
21160286 words written in the Greek alphabet.
3755493 words written in the Latin alphabet.
5990 words written in the Arabic alphabet.
```

# Linux and Windows

CTSExtract.go` is written in Go and can be easily compiled for your system. Flick me a message if you are interested (this is just a get it out before the holidays initial repo-setup).

# Extract OAI-PMH compliant metadata
CTSExtract can be used to extract metadta fields of TEI-XML annotated
input. Currently export to CSV, JSON and XML is possible. 
The XML format complies to OAI-DC format (DataCite). Please see
OAI-PMH.md for information on OAI-PMH compliant hosting.
