# emotions-ai

Projeto de estudo: detecção de expressões faciais em tempo real pela webcam.

- A detecção do rosto roda **no navegador** (MediaPipe, via WASM). O vídeo nunca sai da página.
- Apenas os **rótulos de emoção** (texto) são enviados ao servidor por WebSocket.
- A cada 8 segundos, o servidor (Go + [Genkit](https://firebase.google.com/docs/genkit)) agrupa os rótulos recebidos e pede ao Gemini um resumo curto e cauteloso da sessão.

> Isto mede **expressão facial aparente**, não o que a pessoa realmente sente. A classificação é uma heurística didática, não um classificador validado.

## Requisitos

- [Go](https://go.dev/dl/) 1.25 ou superior
- Uma chave de API do Gemini, que pode ser criada no [Google AI Studio](https://aistudio.google.com/apikey)
- Um navegador moderno (Chrome, Edge ou Firefox) com webcam
- Conexão com a internet: o navegador baixa o MediaPipe e o modelo facial de CDNs, e o servidor chama a API do Gemini

## Instalação

1. Clone o repositório e entre na pasta:

   ```bash
   git clone https://github.com/adrianodevfullstack/emotions-ai.git
   cd emotions-ai
   ```

2. Baixe as dependências:

   ```bash
   go mod tidy
   ```

3. Crie o arquivo `.env` a partir do exemplo:

   ```bash
   # Linux / macOS
   cp .env.example .env
   ```

   ```powershell
   # Windows (PowerShell)
   Copy-Item .env.example .env
   ```

4. Edite o `.env` e coloque a sua chave:

   ```env
   GEMINI_API_KEY=sua_chave_aqui
   PORT=8080
   ```

   | Variável         | Obrigatória | Descrição                               |
   | ---------------- | ----------- | --------------------------------------- |
   | `GEMINI_API_KEY` | Sim         | Chave da API do Gemini                  |
   | `PORT`           | Não         | Porta do servidor HTTP (padrão: `8080`) |

## Como executar

Na **raiz do projeto**, rode:

```bash
go run .
```

O servidor imprime algo como:

```
servindo em http://localhost:8080 (webcam) — janelas de 8s enviadas ao Genkit
```

Depois:

1. Abra <http://localhost:8080> no navegador.
2. Clique em **Connect webcam and start**.
3. Permita o acesso à webcam quando o navegador pedir.
4. Aguarde o modelo carregar (a primeira vez demora alguns segundos).
5. Faça expressões diante da câmera: a emoção detectada aparece na tela e, a cada 8 segundos, surge a interpretação da sessão.

Para gerar um executável:

```bash
go build -o emotions-ai .
./emotions-ai
```

No Windows, use `go build -o emotions-ai.exe .` e execute `.\emotions-ai.exe`. Em ambos os casos, rode o binário a partir da raiz do projeto, porque ele lê o `.env` e a pasta `web/` por caminho relativo.

## Estrutura do projeto

```
main.go                      ponto de entrada (só chama emotions.Run)
internal/emotions/
  server.go                  configuração (.env), Genkit e servidor HTTP
  handler.go                 WebSocket: recebe amostras e envia os resumos
  aggregator.go              agrupa amostras em janelas de tempo
  flow.go                    fluxo Genkit que gera a interpretação
  types.go                   tipos compartilhados
web/
  index.html                 estrutura da página
  style.css                  estilos
  app.js                     webcam, MediaPipe, heurística de emoção e WebSocket
```
