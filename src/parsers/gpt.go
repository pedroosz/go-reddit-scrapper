package parsers

import (
	"context"
	"os"

	"github.com/ayush6624/go-chatgpt"
	"github.com/pedroosz/go-reddit-scrapper/src/utils"
)

func buildClient() *chatgpt.Client {
	client, err := chatgpt.NewClient(os.Getenv("GPT_SECRET"))
	if err != nil {
		utils.Fatal("Não foi possível criar um cliente do CHAT GPT", err)
	}
	return client
}

func GPTParse(message string) string {
	client := buildClient()

	res, err := client.Send(context.Background(), &chatgpt.ChatCompletionRequest{
		Model: chatgpt.GPT35Turbo16k,
		Messages: []chatgpt.ChatMessage{
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: "Responda em português",
			},
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: "Traduza o texto, caso ele esteja em outra lingua que não seja português.",
			},
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: "Utilize o sistema de codificação de caracteres UTF-8",
			},
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: "Qualquer palavra de baixo calão, como sexo e assédio, troque para s*x* e ass*d**",
			},
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: "Utilize as regras gramaticais e sintáticas mais recentes que você conhece",
			},
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: "Mantenha as quebras de linha e não as altere",
			},
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: "Você é um robô de correção gramatical e sintática.",
			},
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: "Qualquer texto recebido, você irá censurar palavras de baixo calão, reescrever o texto de uma forma que faço sentido gramaticalmente e sintaticamente sem alterar muito o conteúdo e mantendo o nível de formalidade original",
			},
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: "Você deve devolver como resposta somente o texto tratado sem nenhuma mensagem adicional",
			},
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: "Você irá receber histórias do reddit. Essas histórias possuem indicativos próprios de sexo, gênero e idade. Um exemplo claro é H20, que quer dizer homem de 20 anos, e F20, que quer dizer mulher de 20 anos. Caso esse indicativo exista, você deve substituílos para o formato por extenso, como indicado.",
			},
			{
				Role:    chatgpt.ChatGPTModelRoleUser,
				Content: message,
			},
		},
	})

	if err != nil {
		utils.Fatal("Não foi possível realizar a consulta ao ChatGPT", err)
	}

	return res.Choices[0].Message.Content
}
