run-monster:
	go run ./Monster-Slayer
run-snippetbox:
	cd ./snippetbox/cmd/web/ && go run . -dsn="root:password@tcp(127.0.0.1:3306)/snippetbox?parseTime=true"
