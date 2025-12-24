package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func AnalisarNuvem(imgData []byte) (string, error) {
	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("Erro ao criar cliente: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemma-3-12b")

	prompt := []genai.Part{
		genai.ImageData("jpeg", imgData),
		genai.Text("Atue como um meteorologista especialista. " +
			"Identifique qual tipo de nuvem é esta na imagem. " +
			"Diga o nome científico, características visuais que confirmam isso " +
			"e se ela indica chuva ou tempo bom. Seja didático."),
	}

	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		return "", fmt.Errorf("Erro ao gerar conteudo: %v", err)
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		texto := fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])
		return texto, nil
	}

	return "Não consegui identificar a nuvem.", nil
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("1. Recebi requisição! Método: %s", r.Method)

	if r.Method != http.MethodPost {
		log.Printf("⛔ ERRO: Método incorreto. Esperava POST, recebi: %s", r.Method)
		http.Error(w, "Metodo não permitido", http.StatusMethodNotAllowed)
		return
	}

	log.Println("2. Tentando ler o arquivo do formulário...")
	file, _, err := r.FormFile("foto")
	if err != nil {
		http.Error(w, "Erro ao ler a imagem", http.StatusBadRequest)
	}
	defer file.Close()

	imgBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Erro ao processar a imagem", http.StatusInternalServerError)
		return
	}
	log.Printf("3. Imagem lida com sucesso! Tamanho: %d bytes", len(imgBytes))

	log.Println("4. Enviando para o Gemini...")
	respostaIA, err := AnalisarNuvem(imgBytes)
	if err != nil {
		log.Printf("Erro na IA: %v", err)
		http.Error(w, "Erro ao analisar nuvem", http.StatusInternalServerError)
		return
	}

	log.Println("5. Sucesso! Enviando resposta para o site.")
	w.Write([]byte(respostaIA))

}

func main() {

	http.HandleFunc("/analisar", UploadHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
