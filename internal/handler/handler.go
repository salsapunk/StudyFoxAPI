package handler

// arrumar os códigos dos erros

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/salsapunk/StudyFoxAPI/internal/model"
	"github.com/salsapunk/StudyFoxAPI/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHand(TaskServ *service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: TaskServ,
	}
}

func (tH *TaskHandler) CheckHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (tH *TaskHandler) Validate(c *gin.Context) {
	usuario, _ := c.Get("usuario")

	model.OK(c, usuario)
}

func (tH *TaskHandler) RequireAuth(c *gin.Context) {
	ctx := c.Request.Context()

	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		v, ok := claims["sub"].(float64)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		matricula := int(v)

		usuario, err := tH.service.LerUsuario(ctx, matricula)
		if err != nil || usuario.Matricula == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("usuario", usuario)
		c.Next()
	} else {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func (tH *TaskHandler) SingUp(c *gin.Context) {
	ctx := c.Request.Context()
	var usuario model.Usuario

	if err := c.ShouldBindJSON(&usuario); err != nil {
		model.Fail(c, http.StatusBadRequest, 0, err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(usuario.Senha), 10)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 0, err)
		return
	}

	usuario.Senha = string(hash)

	id, err := tH.service.CriarUsuario(ctx, &usuario)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 2, err)
		return
	}

	model.OK(c, id)
}

func (tH *TaskHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var usuario model.Usuario

	if err := c.ShouldBindJSON(&usuario); err != nil {
		fmt.Println("ler json")
		model.Fail(c, http.StatusBadRequest, 0, err)
		return
	}

	var usuarioDB model.Usuario
	usuarioDB, err := tH.service.LerUsuario(ctx, usuario.Email)
	if err != nil {
		fmt.Println("ler user")
		model.Fail(c, http.StatusBadRequest, 0, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(usuarioDB.Senha), []byte(usuario.Senha_hash))
	if err != nil {
		fmt.Println("hash")
		model.Fail(c, http.StatusBadRequest, 0, err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usuarioDB.Matricula,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 0, err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", true, true)
	model.OK(c, "Login executed successfully")
}

func (tH *TaskHandler) CriarMateria(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, http.StatusBadRequest, err)
		return
	}

	var materia model.Materia
	err = c.ShouldBindJSON(&materia)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 3, err)
		return
	}

	materia.Matricula = matricula

	codigo, err := tH.service.CriarMateria(ctx, &materia)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 4, err)
		return
	}

	model.OK(c, codigo)
}

func (tH *TaskHandler) CriarTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 5, err)
		return
	}

	var tarefa model.Tarefa
	err = c.ShouldBindJSON(&tarefa)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 6, err)
		return
	}

	tarefa.Codigo = codigo

	id, err := tH.service.CriarTarefa(ctx, &tarefa)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 7, err)
		return
	}

	model.OK(c, id)
}

// READ

func (tH *TaskHandler) LerUsuario(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, http.StatusBadRequest, err)
		return
	}

	usuario, err := tH.service.LerUsuario(ctx, matricula)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 8, err)
		return
	}

	model.OK(c, usuario)
}

func (tH *TaskHandler) ListarMaterias(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 9, err)
		return
	}

	materias, err := tH.service.ListarMaterias(ctx, matricula)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 10, err)
		return
	}

	model.OK(c, materias)
}

func (tH *TaskHandler) LerMateria(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, http.StatusBadRequest, err)
		return
	}

	param = c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, http.StatusBadRequest, err)
		return
	}

	usuario, err := tH.service.LerMateria(ctx, matricula, codigo)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 11, err)
		return
	}

	model.OK(c, usuario)
}

func (tH *TaskHandler) ListarTarefas(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 12, err)
		return
	}

	tarefas, err := tH.service.ListarTarefas(ctx, codigo)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 13, err)
		return
	}

	model.OK(c, tarefas)
}

func (tH *TaskHandler) LerTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 14, err)
		return
	}

	param = c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 15, err)
		return
	}

	tarefa, err := tH.service.LerTarefa(ctx, codigo, id)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 16, err)
		return
	}
	model.OK(c, tarefa)
}

// UPDATE

func (tH *TaskHandler) MudarEmail(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 17, err)
		return
	}

	var usuario model.Usuario
	err = c.ShouldBindJSON(&usuario)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 18, err)
		return
	}

	err = tH.service.MudarEmail(ctx, usuario.Email, matricula)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 19, err)
		return
	}

	model.OK(c, "Email do usuário atualizado")
}

func (tH *TaskHandler) MudarSenha(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 20, err)
		return
	}

	var usuario model.Usuario
	err = c.ShouldBindJSON(&usuario)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 21, err)
		return
	}

	err = tH.service.MudarSenha(ctx, usuario.Senha, matricula)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 22, err)
		return
	}

	model.OK(c, "Senha do usuário atualizada")
}

func (tH *TaskHandler) MudarTema(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 0, err)
		return
	}

	var usuario model.Usuario
	err = c.ShouldBindJSON(&usuario)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 0, err)
		return
	}

	err = tH.service.MudarTema(ctx, usuario.Tema, matricula)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 0, err)
		return
	}

	model.OK(c, "Tema do usuário atualizada")
}

func (tH *TaskHandler) MudarNomeMateria(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 23, err)
		return
	}

	param = c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 24, err)
		return
	}

	var materia model.Materia
	err = c.ShouldBindJSON(&materia)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 25, err)
		return
	}

	err = tH.service.MudarNomeMateria(ctx, materia.Nome, matricula, codigo)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 26, err)
		return
	}

	model.OK(c, "Nome da matéria atualizado")
}

func (tH *TaskHandler) MudarNomeTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 27, err)
		return
	}

	param = c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 28, err)
		return
	}

	var tarefa model.Tarefa
	err = c.ShouldBindJSON(&tarefa)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 29, err)
		return
	}

	err = tH.service.MudarNomeTarefa(ctx, tarefa.Nome, codigo, id)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 30, err)
		return
	}

	model.OK(c, "Nome da tarefa atualizada")
}

func (tH *TaskHandler) MudarPrazoTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 31, err)
		return
	}

	param = c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 32, err)
		return
	}

	var tarefa model.Tarefa
	err = c.ShouldBindJSON(&tarefa)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 33, err)
		return
	}

	err = tH.service.MudarPrazoTarefa(ctx, tarefa.Prazo, codigo, id)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 34, err)
		return
	}

	model.OK(c, "Prazo da tarefa atualizada")
}

func (tH *TaskHandler) MudarAnotacaoTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 35, err)
		return
	}

	param = c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 36, err)
		return
	}

	var tarefa model.Tarefa
	err = c.ShouldBindJSON(&tarefa)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 37, err)
		return
	}

	err = tH.service.MudarAnotacaoTarefa(ctx, tarefa.Anotacao, codigo, id)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 38, err)
		return
	}

	model.OK(c, "Anotação da tarefa atualizada")
}

func (tH *TaskHandler) MudarStatusTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 0, err)
		return
	}

	param = c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 0, err)
		return
	}

	var tarefa model.Tarefa
	err = c.ShouldBindJSON(&tarefa)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 0, err)
		return
	}

	err = tH.service.MudarStatusTarefa(ctx, tarefa.Status, codigo, id)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 0, err)
		return
	}

	model.OK(c, "Status da tarefa atualizada")
}

// DELETE

func (tH *TaskHandler) DeletarUsuario(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 39, err)
		return
	}

	err = tH.service.DeletarUsuario(ctx, matricula)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 40, err)
		return
	}

	model.OK(c, "Usuário deletado")
}

func (tH *TaskHandler) DeletarMateria(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 41, err)
		return
	}

	param = c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 42, err)
		return
	}

	err = tH.service.DeletarMateria(ctx, matricula, codigo)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 43, err)
		return
	}

	model.OK(c, "Matéria deletada")
}

func (tH *TaskHandler) DeletarTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 44, err)
		return
	}

	param = c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 45, err)
		return
	}

	err = tH.service.DeletarTarefa(ctx, codigo, id)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 46, err)
		return
	}

	model.OK(c, "Tarefa deletada")
}
