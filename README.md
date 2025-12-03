## Start do projeto (dev)
1. Build e run: `docker compose up --build`.
2. Verificar se a API está up: `http://localhost:8080`: `docker compose ps` e/ou `docker compose logs -f app`.
3. O Mongo fica em `mongodb://admin:admin@localhost:27017/auctions?authSource=admin`.

## Crir users(Mongo)
Não há endpoint para criar users como tal podemos inserir `mongodb` no docker (gera/ajuste IDs UUID válidos):
```bash
docker compose exec -T mongodb mongosh -u admin -p admin --authenticationDatabase admin <<'EOF'
  use auctions
  db.users.insertOne({ _id: "b5080cce-23c4-4a8f-93c3-2f6d5d4a0d2a", name: "Sergio" })
  db.users.insertOne({ _id: "148a8744-53d7-4b64-b94c-ec543e58908a", name: "Joana" })
EOF
```
## Criar um leilão (auction)
Payload exemplo (`condition`: 0=New, 1=Used, 2=Refurbished):
```bash
curl -X POST http://localhost:8080/auction \
  -H "Content-Type: application/json" \
  -d '{
    "product_name": "Panados",
    "category": "comida",
    "description": "Panados caseiros deliciosos",
    "condition": 0
  }'
```
Como nao e retornado o uuid do leilão deve ser consultado para se poder realizar bids 
```bash
curl http://localhost:8080/auction?status=0
```

## Criando um lance (bid)
Use um `user_id` existente e o `auction_id` do leilão que foi inserido como `Panados`:
```bash
curl -X POST http://localhost:8080/bid \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "148a8744-53d7-4b64-b94c-ec543e58908a",
    "auction_id": "18ab801f-4fdb-4f70-89c5-e1c9f6e31285",
    "amount": 50.0
  }'
```
Resposta esperada: HTTP 201 sem o body. O leilão fecha automaticamente após `AUCTION_INTERVAL`.

## Run testes
- Local (host): `go test ./...`
- Via Docker (usando a imagem Go): `docker compose exec app go test ./...`
  - Nota: antes de executar os testes primeiro tem de ser executado `docker compose up --build` e verificar se os
  containers estao em execucao `docker compose ps`.

