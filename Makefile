.PHONY: deploy
deploy:
	vercel --prod

.PHONY: run
run:
	go run main.go


