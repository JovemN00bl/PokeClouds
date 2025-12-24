package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func AnalisarNuvemTeste(imgData []byte) (string, error) {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")

	if apiKey == "" {
		return "", fmt.Errorf("A variavel GEMINI_API_KEY está vazia!")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-3-flash-preview")
	prompt := []genai.Part{
		genai.ImageData("jpeg", imgData),
		genai.Text("Identifique esta nuvem de forma breve."),
	}

	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0]), nil
	}
	return "Sem resposta", nil
}

func ma_in() {
	nomeArquivo := "teste.jpg"

	dadosImagem, err := os.ReadFile(nomeArquivo)
	if err != nil {
		log.Fatalf("Erro ao abrir a imagem: %v. Verifique se o arquivo '%s' existe na pasta.", err, nomeArquivo)
	}

	fmt.Println("Enviando para o Gemini... aguarde...")

	resposta, err := AnalisarNuvemTeste(dadosImagem)
	if err != nil {
		log.Fatalf("Erro na API: %v", err)
	}

	fmt.Println("--- Resposta do Gemini ---")
	fmt.Println(resposta)
}
