package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-api-rest/controllers"
	"github.com/gin-api-rest/database"
	"github.com/gin-api-rest/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var ID int

func RegistroDasRotasDeTeste() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	rotas := gin.Default()
	return rotas

}

func CriaAlunoMock() {

	aluno := models.Aluno{Nome: "Nome do aluno", CPF: "00000000000", RG: "000000000"}
	database.DB.Create(&aluno)
	ID = int(aluno.ID)

}

func DeletaAlunoMock() {

	var aluno models.Aluno
	database.DB.Delete(&aluno, ID)
}

func TestDeSatusCode(t *testing.T) {

	r := RegistroDasRotasDeTeste()
	r.GET("/:nome", controllers.Saudacao)
	req, _ := http.NewRequest("GET", "/Nia", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais!")
	mockDaResposta := `{"API diz:":"E ai Nia, tudo bem?"}`
	respostaBody, _ := ioutil.ReadAll(resposta.Body)
	assert.Equal(t, mockDaResposta, string(respostaBody))

}

func TestListandoTdsOsAlunosHandler(t *testing.T) {

	database.ConectaComBamcoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := RegistroDasRotasDeTeste()
	r.GET("/alunos", controllers.ExibetodosAlunos)
	req, _ := http.NewRequest("GET", "/alunos", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusOK, resposta.Code)

}

func TestBuscaAlunoPorCPFHandler(t *testing.T) {

	database.ConectaComBamcoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := RegistroDasRotasDeTeste()
	r.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCpf)
	req, _ := http.NewRequest("GET", "/alunos/cpf/00000000000", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusOK, resposta.Code)

}

func TestBuscoAlunoPorIDHandler(t *testing.T) {

	database.ConectaComBamcoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := RegistroDasRotasDeTeste()
	r.GET("/alunos/:id", controllers.BuscaALunoPorId)
	pathDaBusca := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("GET", pathDaBusca, nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)

	var alunoMock models.Aluno

	json.Unmarshal(resposta.Body.Bytes(), &alunoMock)
	assert.Equal(t, "Nome do aluno", alunoMock.Nome)
	assert.Equal(t, "00000000000", alunoMock.CPF)
	assert.Equal(t, "000000000", alunoMock.RG)
	assert.Equal(t, http.StatusOK, resposta.Code)

}

func TestDeletaAlunoHandler(t *testing.T) {

	database.ConectaComBamcoDeDados()
	CriaAlunoMock()
	r := RegistroDasRotasDeTeste()
	r.DELETE("/alunos/:id", controllers.DeletaAluno)
	pathDeBusca := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("DELETE", pathDeBusca, nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusOK, resposta.Code)
}

func TestEditaUmAlunoHandler(t *testing.T) {

	database.ConectaComBamcoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := RegistroDasRotasDeTeste()
	r.PATCH("/alunos/:id", controllers.EditaAluno)
	aluno := models.Aluno{Nome: "Nome do aluno", CPF: "12300000000", RG: "000000789"}
	valorJson, _ := json.Marshal(aluno)
	patchParaEditar := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("PATCH", patchParaEditar, bytes.NewBuffer(valorJson))
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)

	var alunoMockAtualizado models.Aluno

	json.Unmarshal(resposta.Body.Bytes(), &alunoMockAtualizado)
	assert.Equal(t, "Nome do aluno", alunoMockAtualizado.Nome)
	assert.Equal(t, "12300000000", alunoMockAtualizado.CPF)
	assert.Equal(t, "000000789", alunoMockAtualizado.RG)

}
