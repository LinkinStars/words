#!/bin/bash
(
cd books || exit
tar -czvf books.tar.gz *.json
)
rm -f assets/books.tar.gz
mv books/books.tar.gz assets/books.tar.gz

go build -ldflags "-s -w" -o words ./cmd/words
chmod +x words
mv words /usr/local/bin