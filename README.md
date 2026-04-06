# 🎬 CineFinder

## 1. Visão do Produto

Para amantes de cinema que têm dificuldades em encontrar filmes disponíveis para locação de forma rápida.

O **CineFinder** é uma ferramenta simples para pesquisar títulos em uma locadora.  
Permite o cadastro de filmes, a busca fácil por eles e checar sua disponibilidade para aluguel.

Diferente de sistemas muito complexos como a Amazon Prime, que possuem um vasto catálogo de filmes e oferecem acesso online, nosso produto foca na **simplicidade e praticidade**, sendo ideal para usuários interessados também em **mídia física**, uma prática hoje menos comum.

---

## 2. Definição do MVP

### 2.1 Funcionalidades do MVP (Dentro do Escopo)

#### 📌 Cadastro de filmes
- Inserção de:
  - Título
  - Diretor
  - Ano de lançamento
  - Gênero

#### 🔍 Busca
- Pesquisa completa ou parcial
- Ignora letras maiúsculas/minúsculas

#### 📃 Listagem
- Exibição de:
  - Filmes disponíveis
  - Filmes alugados
  - Quantidade de exemplares

#### 📅 Reserva digital
- Permite reservar um filme
- Exibe:
  - Preço
  - Prazo

#### 📚 Histórico de empréstimos (locadora)
- Registro de:
  - Empréstimos
  - Devoluções

#### 🔐 Sistema de login básico
- Autenticação com:
  - Nome
  - E-mail
- Necessário para realizar reservas

---

### 2.2 Funcionalidades Fora do Escopo

- Registro de transações financeiras  
- Sistema de autenticação avançado (senha, perfis, etc.)  
- Histórico detalhado para usuários  
- Notificações de devolução de filmes  

---

## 3. Product Backlog Inicial

### 🟢 P1 – Cadastro de Filmes
**Como** locador,  
**quero** cadastrar filmes,  
**para** manter o catálogo atualizado.

**Critérios de Aceitação:**
- Campos obrigatórios:
  - Título
  - Diretor
  - Ano
  - Gênero
- O sistema deve evitar duplicidade de registros

---

### 🟢 P2 – Autenticação de Usuário
**Como** usuário,  
**quero** me autenticar,  
**para** poder reservar filmes.

**Critérios de Aceitação:**
- Login com nome e e-mail
- Validação básica dos dados
- Apenas usuários autenticados podem reservar

---

### 🟢 P3 – Busca de Filmes
**Como** usuário,  
**quero** buscar filmes por título,  
**para** encontrar rapidamente o que desejo.

**Critérios de Aceitação:**
- Busca parcial (LIKE)
- Ignora maiúsculas/minúsculas
- Retorna lista de resultados

---

### 🟢 P4 – Listagem de Filmes
**Como** usuário,  
**quero** visualizar a lista de filmes,  
**para** saber quais estão disponíveis.

**Critérios de Aceitação:**
- Exibir status:
  - Disponível
  - Alugado
- Mostrar quantidade de exemplares
- Carregamento rápido

---

### 🟢 P5 – Reserva de Filmes
**Como** usuário,  
**quero** reservar um filme,  
**para** garantir sua disponibilidade.

**Critérios de Aceitação:**
- Vincular filme ao usuário
- Exibir:
  - Preço
  - Prazo
- Impedir reservas duplicadas

---

### 🟢 P6 – Histórico de Empréstimos
**Como** locador,  
**quero** visualizar o histórico de empréstimos,  
**para** analisar o uso do acervo.

**Critérios de Aceitação:**
- Armazenar:
  - Data de saída
  - Data de devolução
- Permitir filtragem:
  - Por filme
  - Por data

## 4. Link do Projeto escrito e vídeo

[Arquivos no Drive](https://drive.google.com/drive/folders/15M8KNpg-0WwB8Phuqfl7JoVw4xv0fqVP?usp=sharing)