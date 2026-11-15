# Full Cycle - Pós Go Expert - Stress Test

CLI em Go para realizar testes de carga HTTP. O usuário informa a URL de destino, o total de requisições e o número máximo de chamadas simultâneas. Ao final, o programa apresenta o relatório no terminal.

## Objetivos

- Receber os parâmetros `--url`, `--requests` e `--concurrency` pela linha de comando.
- Distribuir as chamadas entre goroutines e executar o total solicitado, inclusive quando ele não é múltiplo da concorrência.
- Informar o tempo total, a quantidade de requisições, a quantidade de respostas HTTP 200 e a distribuição dos demais status.
- Permitir a execução da CLI com Docker.

---

## Tecnologias utilizadas

- Go 1.26
- Pacotes da biblioteca padrão `flag`, `net/http` e `sync`
- Goroutines, canais e `sync.WaitGroup`
- Docker com build em dois estágios

---

## Como executar

### 1. Clonar o repositório

```bash
git clone https://github.com/danielsantello/fc-pos-desafio-stress-test.git
cd fc-pos-desafio-stress-test
```

### 2. Construir a imagem

Com o Docker disponível, execute na raiz do projeto:

```bash
docker build -t stress-test .
```

O Docker compila o executável dentro do primeiro estágio. A imagem final contém o binário e os certificados usados para verificar conexões HTTPS.

### 3. Rodar o teste de carga

Substitua a URL de exemplo pelo endereço do serviço que você quer testar:

```bash
docker run --rm stress-test \
  --url=https://seu-servico.example/ \
  --requests=1000 \
  --concurrency=10
```

Passe as flags da CLI depois do nome da imagem. O `--rm` remove o container quando o teste termina.

### 4. Experimentar com um servidor local

Com Python 3 disponível, inicie um servidor em um terminal:

```bash
python3 -m http.server 8081 --bind 127.0.0.1
```

Em outro terminal, execute a CLI. Use `--network host` se o Docker estiver configurado para esse modo de rede:

```bash
docker run --rm --network host stress-test \
  --url=http://127.0.0.1:8081/ \
  --requests=7 \
  --concurrency=3
```

Para testar a distribuição de outros status, troque o caminho `/` por um caminho inexistente, como `/__stress_test_404__`. Finalize o servidor local com `Ctrl+C`.

---

## Parâmetros

| Parâmetro | Finalidade | Exemplo |
| --- | --- | --- |
| `--url` | URL absoluta do serviço, usando HTTP ou HTTPS. | `--url=http://127.0.0.1:8081/` |
| `--requests` | Total de requisições a realizar. Deve ser maior que zero. | `--requests=23` |
| `--concurrency` | Máximo de chamadas simultâneas. Deve ser maior que zero. | `--concurrency=5` |

Os três parâmetros são obrigatórios. Para consultar a ajuda, use `docker run --rm stress-test --help`.

---

## Funcionamento

O [main.go](main.go) valida os parâmetros e imprime o relatório. Em [load.go](load.go), o programa envia uma tarefa por requisição a um canal. Cada goroutine pega uma tarefa, faz um `GET` e envia o resultado a outro canal. Se a concorrência configurada for maior que o total de requisições, o programa cria apenas as goroutines necessárias para esse total.

Uma única goroutine recebe e contabiliza os resultados. Assim, `--requests=23 --concurrency=5` produz 23 tentativas, com até cinco chamadas em andamento. O cliente usa um timeout de 10 segundos por chamada. Se receber um redirecionamento, registra o status correspondente e encerra aquela chamada ali.

## Relatório

Uma execução local contra um caminho inexistente pode produzir esta saída:

```text
Tempo total: 20.8016ms
Total de requisições realizadas: 7
HTTP 200: 0
Outros status HTTP:
  HTTP 404: 7
Falhas sem resposta HTTP: 0
```

O tempo total é o intervalo decorrido durante a execução concorrente. Os outros status aparecem em ordem numérica. Respostas como `404` e `500` entram na distribuição. Falhas sem resposta HTTP, como conexão recusada ou timeout antes de receber o status, aparecem separadamente e fazem parte do total.

---

## Autor

Daniel Santello
