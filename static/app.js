const { createApp } = Vue;

createApp({
    data() {
        return {
            arquivo: null,
            previewUrl: null,
            loading: false,
            resultado: null
        }
    },
    computed: {
        resultadoFormatado() {
            if (!this.resultado) return "";
            return this.resultado.replace(/\n/g, '<br>');
        }
    },
    methods: {
        selecionarImagem(event) {
            const files = event.target.files;
            if (files.length === 0) return;

            this.arquivo = files[0];
            this.previewUrl = URL.createObjectURL(this.arquivo);
        },

        async analisarNuvem() {
            if (!this.arquivo) return;

            this.loading = true;

            const formData = new FormData();
            formData.append("foto", this.arquivo);

            try {
                // Chama o backend em Go
                const response = await fetch('/analisar', {
                    method: 'POST',
                    body: formData
                });

                if (!response.ok) throw new Error("Erro no servidor");

                const texto = await response.text();
                this.resultado = texto;

            } catch (error) {
                console.error(error);
                alert("Erro ao analisar a nuvem. Tente novamente.");
            } finally {
                this.loading = false;
            }
        },

        resetar() {
            this.arquivo = null;
            this.previewUrl = null;
            this.resultado = null;
            this.loading = false;
        }
    }
}).mount('#app');