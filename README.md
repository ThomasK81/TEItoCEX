# TEItoCEX
Turn CTS TEI corpora like CHS and OGL's [First1KGreek](http://opengreekandlatin.github.io/First1KGreek/) into [CEX collection files](https://github.com/cite-architecture/citedx/blob/master/docs/CEX-spec-3.0.md). 

# USAGE OSX

1. Copy `CTSExtract` into the unpacked data folder of e.g. First1Greek. 
2. Open a terminal in that folder and type: `./CTSExtract 1kGreek.cex ` (you might have to chmod +x the executable before you can use it)
3. Enjoy your new CEX collection file!

## Alternatively convert to CSV (or JSON, a flat XML, or SQL) 

1. Copy `CTSExtract` into the unpacked data folder of e.g. First1Greek. 
2. Open a terminal in that folder and type: `./CTSExtract 1kGreek.csv -CSV `
3. Enjoy your new CSV collection file!

# Sample Terminal Output

The numbers and letters shows the scheme that has been used in the original XML file:

```
KKGGGKGG58GKGGGGGGGGGGGGGGGGKGGKGKGGGGGGG7IGGGGGGGKKKKKGGGGGKKGGGKGGGGGGGGGKGGGGKGGGGGGKKGGGGGGGGGGKKGKKKKGKKKKGGGGGGKGLKKKGGGGGKGKLKGGKGKKKGGGGGKGGGGGGGGGGKGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGDKSJJKGKKGGGGKGM5555555LLKKGGGGGGGGKGGKKKGGGGGGLGKKLKGGGKGGKGGGGGKKGGGRGKGGGGGGGKKGGGKGGGGGGKKKKKKKKKKKKKKLKKKGKGGGGKGGGGGKLGKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKGGGKGKKKKKGKKKKKKKGGGKGGKGGKKKKKGGGKGGKGGGGGKKGKKGGGKKKGKGLKGKKGLKGGGEEGKKLGGGKGKLLGKGGGGGGLGGGGGKGGGKGGGGGGGGGGGGGGGKGGGKGGGGGGGGGGGGGGGGGKGGGQGGKKGKGGGKKLKKKKKKGGGGGGKKLKGGGGGGGGGKKLK4GKKLKGGGLLKKKKKKKKKKKKKGGGGKKKGGGGKKLGGGGGKGGGGGGGGKGLGGGGGKGLGGGGKLGKLLKGGGLKKLLK9GGGGGGKKGKKGKGGGGKKKKKGGGGGGGGGKKGGKKGGGGGGG3LKGKKKKKGGGGGGGKKKKKKKKKKKKKKLKLGKKLKKGGGGKGKGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGKGGKGGGGGGGGGGGGGGGGGKKGGGGGGGGKGGGGGKGGKGKKGGKGGGKGGPKKKKKKLKKKKLKKKGGGKKKGGKKGGLLLGGGGGGGKGGGGKLGGGGGGGGGGGKKGGKGGGGGGGGGKKGLKLGLGGGGGLGGLLGGLGGGGGGGGGGGGGGGGGGGKGGKGGGGGGGKGGGKGGKKKK88KKKKKGGLG
Read 974 of 974 files.
Write nodes to file now:
Writing CSV-File
Wrote 227668 nodes.
23340077 words written in the Greek alphabet.
4331600 words written in the Latin alphabet.
5996 words written in the Arabic alphabet.
The following schemes were used:
K 310
8 3
M 1
R 1
Q 1
L 47
D 1
S 1
J 2
5 8
7 1
I 1
4 1
P 1
G 591
E 2
9 1
3 1
```

# Linux and Windows

CTSExtract.go` is written in Go and can be easily compiled for your system. Flick me a message if you are interested.

# Extract OAI-PMH compliant metadata
CTSExtract can be used to extract metadta fields of TEI-XML annotated
input. Currently export to CSV, JSON and XML (and SQL) is possible. 
The XML format complies to OAI-DC format (DataCite). Please see
OAI-PMH.md for information on OAI-PMH compliant hosting.

# Producing First1kGreek JSON Catalog

```
./TEItoCEX catalog.json -Cat
```
The catalog can then replace the `catalog.json` in the gh-pages branch of the [First1KGreek](http://opengreekandlatin.github.io/First1KGreek/) repo.

