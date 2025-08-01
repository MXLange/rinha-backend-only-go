# Rinha de Backend 2025

## Repositório

[https://github.com/MXLange/rinha-backend-only-go](https://github.com/MXLange/rinha-backend-only-go)

## Projeto

A escolha da linguagem **Go** foi motivada tanto pelo uso recorrente no meu dia a dia quanto pela sua eficiência em ambientes com recursos limitados — exatamente o cenário proposto pela Rinha.

A aplicação foi construída com foco em **concorrência leve** e **uso otimizado da memória**, aproveitando as ferramentas nativas da linguagem: **goroutines** e **channels**.

---

### 🧠 Arquitetura da Solução

-   **Fila de pagamentos em memória:**  
    Em vez de usar Redis ou outros intermediários, utilizei um `channel` como fila para armazenar pagamentos recebidos via API.

-   **Workers configuráveis via ENV:**  
    Um pool de workers é responsável por processar os pagamentos em paralelo. A quantidade de workers pode ser ajustada dinamicamente pela variável de ambiente `WORKERS`.

-   **Resumo distribuído dos pagamentos:**  
    O estado dos pagamentos é mantido em memória. Para calcular o total agregado, cada instância da API consulta diretamente a outra via HTTP, somando os valores locais com os remotos.

---

### ⚙️ Tecnologias e escolhas

-   **Go puro (sem frameworks pesados)**
-   Concorrência com `goroutines` e `sync.Mutex`
-   Comunicação entre instâncias via HTTP
-   Armazenamento e agregação totalmente em memória
-   Configurações por `ENV` para fácil tunning

---

### 💡 Destaques

-   Sem dependências externas como Redis ou DB
-   Baixo consumo de memória
-   Altamente paralelizável e escalável horizontalmente
-   Simples de entender, fácil de adaptar

---
