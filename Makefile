build:
	go build -o crawler -v
clean:
	rm -rf crawler
install:
	go get github.com/disiqueira/gotree
	go get github.com/PuerkitoBio/goquery
	go get github.com/asaskevich/govalidator
test:
	go test -v -race ./...