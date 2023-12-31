run-monster:
	go run ./Monster-Slayer
run-snippetbox:
	cd ./snippetbox/cmd/web/ && go run . -dsn="web:password@tcp(127.0.0.1:3306)/snippetbox?parseTime=true"
run-test:
	cd ./snippetbox/cmd/web/ && go test -v . 
run-test-cover:
	cd ./snippetbox && go test -cover ./...

run-test-single:
	cd ./snippetbox && go test -v -run=TestSnippetCreate ./cmd/web/
