OAI-PMH server for CTS TEI metadata
===================================

Clone TEI corpora files from github
-----------------------------------

cd ~/data 
git clone https://github.com/OpenGreekAndLatin/First1KGreek.git

Convert CEX to CSV
------------------

    cd /path/to/TEI/data
    go run ~/src/go/src/TEItoCEX/CTSExtract.go ~/data/cex/oglp.csv -CSV

Convert CEX to JSON
-------------------

    time (cd ~/data/First1KGreek/ ; go run ~/src/go/src/TEItoCEX/CTSExtract.go ~/data/First1KGreek.xml - JSON)

Convert CEX to XML (OAI-DC)
---------------------------

    time (cd ~/data/First1KGreek/ ; go run ~/src/go/src/TEItoCEX/CTSExtract.go ~/data/First1KGreek.xml -XML )

populate OAI-PMH server
-----------------------

git repo handling
=================

This repo is maintained on github *and* bitbucket.

push to github
--------------

    git remote add github https://github.com/tgoerke/TEItoCEX.git
    git push -u github master

push to playground
------------------

    git remote add playground ssh://git@code.gerdi-project.de:7999/playg/teitocex.git
    git push -u playground master
