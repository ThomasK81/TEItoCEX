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

`CTSExtract.go` is written in Go and can be easily compiled for your system. Flick me a message if you are interested (this is just a get it out before the holidays initial repo-setup).


# Convert CEX to JSON

```
cd /path/to/TEI/data
go run ~/src/go/src/TEItoCEX/CTSExtract.go ~/data/cex/foo.json -JSON
```

# Convert CEX to XML (OAI-DC)

# All nodes from opengreekandlatin.github.io/First1KGreek (186909)
```
time (cd ~/data/First1KGreek/ ; go run ~/src/go/src/TEItoCEX/CTSExtract.go ~/data/First1KGreek.xml -XML )
takes 1m45 on my laptop
```

# small data set for tests
```
time ((cd ~/data/tlg0090/ ; go run ~/src/go/src/TEItoCEX/CTSExtract.go ~/data/tlg0090.xml -XML )
```

# git handling

This repo is maintained on github *and* bitbucket.

# push to github
```
git remote add github  https://github.com/tgoerke/TEItoCEX.git
git push -u github JSON
```

# push to playground
```
git remote add playground ssh://git@code.gerdi-project.de:7999/playg/teitocex.git
git push -u playground JSON
```

