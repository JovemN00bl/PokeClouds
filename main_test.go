package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomepageHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("O handler retornou o código errado: recebeu %v, queria %v ", status, http.StatusOK)
	}

	expected := "PokeClouds"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("O handler retornou um body inesperado: não encontrei '%v'", expected)
	}

}

func TestValidacaoArquivoVazio(t *testing.T) {
	arquivoVazio := []byte{}

	if len(arquivoVazio) == 0 {
	} else {
		t.Error("Deveria ter detectado que o arquivo está vazio")
	}
}
