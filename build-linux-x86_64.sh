program=TimedClearTask

go clean&&go fmt

go env -w GOOS=linux GOARCH=amd64

go build -v

go env -w GOOS=darwin GOARCH=arm64

mv ./$program ./$program-linux-x86-64

chmod 640 ./$program-linux-x86-64/readme.txt

tar -czvf ./$program-linux-x86-64.tar ./$program-linux-x86-64

#cp ./$program-linux-x86-64.tar /Users/mac/Documents
