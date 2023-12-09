# Como rodar o projeto

Requisitos
Certifique-se de ter o Docker instalado no seu sistema antes de seguir as instruções abaixo.

## Docker
### Construção da Imagem
Abra um terminal e navegue até o diretório raiz do projeto.

Execute o seguinte comando para construir a imagem Docker:

```
docker build -t tusk-app .
```

Certifique-se de incluir o ponto " . " no final do comando para indicar que o Dockerfile está no diretório atual.

### Execução do Contêiner
Após a construção da imagem, execute o seguinte comando para iniciar o contêiner:
```
docker run -p 3000:3000 tusk-app
```

Este comando inicia o contêiner e mapeia a porta 3000 do host para a porta 3000 do contêiner.

Agora, você deve ser capaz de acessar o aplicativo no seu navegador em http://localhost:3000.

#### Observações
Consulte o arquivo ```.env.example``` para fornecer as variáveis de ambiente necessárias.
