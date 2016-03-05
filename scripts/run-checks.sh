#!/bin/bash

go vet -x ./...
go test -cpu=1,2,3,4,5,6,7,8 ./...

# avoid excessive output
go test -cpu=1,2,3,4,5,6,7,8 ./... 2> bench.log
  
time ./edify download_specs 14B
time ./edify extract_specs 14B
time ./edify purge_specs 14B
time ./edify parse testdata/UNCL.14B
time ./edify full_parse --version 14B -d testdata/d14b/
time ./edify query  -d testdata/d14b/ 
time ./edify query  -d testdata/d14b/ --version 14B
time ./edify query  -d testdata/d14b/ -m testdata/messages/INVOIC_1.txt
time ./edify query  -d testdata/d14b/ -m testdata/messages/INVOIC_1.txt  -q "grp:Group_7[0]/seg:CUX[0]"
# - time ./edify  parse testdata/EDED.14B
# - time ./edify  parse testdata/EDCD.14B

./docker-cmd.sh go test -cover ./...