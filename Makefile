.PHONY: install tidy run-dev

install:
	cd frontend && npm install

tidy:
	cd backend && go mod tidy

run-dev:
	cd frontend && npm run dev