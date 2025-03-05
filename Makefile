run-server:
	@cd server
	@air

run-client:
	@air

run:
	@make run-server &
	@sleep 2
	@make run-client