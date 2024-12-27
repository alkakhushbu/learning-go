up:
	docker compose up --build -d

down:
	docker compose down


generate-key-pair:
	openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in private.pem -pubout -out pubkey.pem