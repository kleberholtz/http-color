# Go Color

Aplicação web simples em Go. A página mostra a cor configurada para a branch e o hostname do container.

## Executar com Docker Compose

Na branch desejada, rode:

```sh
docker compose up --build
```

Acesse <http://localhost:8080>. Para encerrar, use `Ctrl+C`.

Branches disponíveis:

- `vermelho`
- `verde`
- `azul`

Para trocar de cor, encerre o Compose, troque de branch e reconstrua a imagem:

```sh
git switch verde
docker compose up --build
```

O endpoint `/healthz` retorna `ok` para verificar se o servidor está ativo.
