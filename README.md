# Projeto de Servidor de Chat com Cliente e Bot em Go

Este é um projeto de um servidor de chat simples com suporte a múltiplos clientes e bot, desenvolvido em Go. O servidor permite a troca de mensagens públicas e privadas, além de permitir que cada usuário altere seu apelido. Um bot chamado `Bot@ReverseBot` responde a todas as mensagens públicas invertendo o texto.

## Funcionalidades

- Conexão de múltiplos usuários simultaneamente.
- Troca de mensagens públicas entre os usuários.
- Envio de mensagens privadas para usuários específicos.
- Alteração do apelido durante a sessão.
- Bot que responde a mensagens públicas invertendo o texto.

## Como Funciona

### Servidor

O servidor de chat aceita conexões TCP na porta **3000**. Ele gerencia a entrada e saída de usuários, além de transmitir mensagens públicas para todos os clientes conectados.

### Cliente

Cada cliente pode se conectar ao servidor e escolher um apelido. Os clientes podem enviar mensagens públicas, trocar seu apelido e enviar mensagens privadas para outros usuários.

### Bot

O `Bot@ReverseBot` responde automaticamente a qualquer mensagem pública invertendo o texto recebido.

## Requisitos

- Go 1.16 ou superior instalado.
- Acesso ao terminal para executar o servidor e o cliente.

## Como Executar

### Executar o Servidor

1. Compile e execute o servidor:
   ```bash
   go run server.go

2. Abra um segundo terminal e execute o cliente:
   ```bash
   go run client.go

3. Abra um terceiro terminal e execute o bot:
   ```bash
   go run bot.go

## Come usar os comandos?

1. Troca de Nome
   ```bash
   \changenick novo_apelido

2. Mensagem Privada
   ```bash
   \msg @nome_do_usuario Mensagem aqui

