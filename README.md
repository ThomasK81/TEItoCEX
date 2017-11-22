# TEItoCEX
Turn CTS TEI corpora like CHS and OGL's [First1KGreek](http://opengreekandlatin.github.io/First1KGreek/) into [CEX collection files](https://github.com/cite-architecture/citedx/blob/master/docs/CEX-spec-3.0.md). 

# USAGE OSX

1. Copy `CTSExtract` into the unpacked data folder of e.g. First1Greek. 
2. Open a terminal in that folder and type: `./CTSExtract 1kGreek.cex ` (you might have to chmod +x the executable before you can use it)
3. Enjoy your new CEX collection file.

# Linux and Windows

`CTSExtract.go` is written in Go and can be easily compiled for your system. Flick me a message if you are interested (this is just a get it out before the holidays initial repo-setup).
