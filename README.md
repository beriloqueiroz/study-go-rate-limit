# Rate Limit

## introdução

- No arquivo todo.md há uma contextualização do desafio desse projeto.

- Este projeto é dividido em 3 camadas, entity, infra, usecase, que coincidem com diretórios na pasta internal e está organizado da seguinte forma:

  ```tree
  ├── cmd
  │   └── main.go
  ├── configs
  │   └── config.go
  ├── docker-compose.yaml
  ├── Dockerfile
  ├── go.mod
  ├── go.sum
  ├── internal
  │   ├── entity
  │   │   ├── limiter.go
  │   │   └── limiter_test.go
  │   ├── infra
  │   │   ├── repository
  │   │   │   ├── config_limit_repository_impl.go
  │   │   │   └── rate_limit_repository_impl.go
  │   │   └── web
  │   │       ├── routes
  │   │       │   └── api
  │   │       │       └── test-simple-route.go
  │   │       └── server
  │   │           └── webserver.go
  │   └── usecase
  │       ├── config_limit_repository.go
  │       ├── rate_limit_repository.go
  │       └── rate_limit.usecase.go
  ├── README.md
  └── todo.md
  ```

  Onde as regras do "ratelimit" estão contidas na camada de entity, o use case representa a intensão do usuário e orquestra os gateways, que são basicamente o repositório de config limit e o repositório de rate limit, abstraídos pelas interfaces contidas na camada de "usecase", então para substiuição destes, basta implementar essas interfaces e injetar no cmd/main.go. Atualmente o repositório de config limit está implementado usando as variáveis de ambiente e o repositório de rate limit usando o redis no caminho: internal/infra/repository/rate_limit_repository_impl.go

## funcionamento

### dependências

- A injeção de dependências acontece de forma resumida da seguinte forma:
  - main (entrypoint)
    - webserver
      - route
      - rateLimitUsecase (usado na middleware do server)
        - configLimitRepository
          - interface: internal/usecase/config_limit_repository.go
          - implementação: internal/infra/repository/config_limit_repository_impl.go
            - ConfigEnvironment
        - rateLimitRepository
          - interface: internal/usecase/rate_limit_repository.go
          - implementação: internal/infra/repository/rate_limit_repository_impl.go

- Ao receber uma requisição no webserver, a mesma passa por um middleware que executa o usecase, que:
  - captura informações do "rate limit repository" e "config limit repository" monta entidade de limiter que processa as regras, por final o usecase salva no "rate limit repository", com uso de dtos, esse limiter e retorna se a operação é permitida ou não. Qualquer erro no execute do usecase na middleware não deve impactar a resposta da intensão do usuário que executou a request.

## para rodar

- testes:

  ```bash
    go test ./...
  ```

- configuração das variáveis de ambiente
  - usar arquivo: .env.config

    ```json
    {
      "web_server_port": ":8080",
      "default_limit_per_ip_per_second": 3,
      "default_expiration_time_in_minutes": 2,
      "default_api_keys_limit_per_second": [
          {
            "api_key": "1",
            "limit_per_second": 3
          },
          {
            "api_key": "2",
            "limit_per_second": 1
          }
        ]
    }
    ```

- aplicação na porta 8080, e redis:

  ```bash
    docker compose up -d
  ```

- request:
  - ip (o ip deve ser passado via header):

   ```bash
    curl  -X GET \
    'http://localhost:8080/' \
    --header 'Accept: */*' \
    --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
    --header 'X-Forwarded-For: 55.235.1.55'
   ```

  - key (a key deve ser passado via header):

   ```bash
    curl  -X GET \
  'http://localhost:8080/' \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
  --header 'API_KEY: assdbasksdbas=aa'
   ```
