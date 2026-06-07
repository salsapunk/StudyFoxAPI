package handler

// arrumar os códigos dos erros

import (
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
	ctx := c.Request.Context()
	usuario, _ := c.Get("usuario")

	user := model.Usuario(usuario.(model.Usuario))

	err := tH.service.Validate(ctx, &user)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
	}

	model.OK(c, usuario)
}

func (tH *TaskHandler) RequireAuth(c *gin.Context) {
	ctx := c.Request.Context()

	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		model.Fail(c, http.StatusUnauthorized, err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		model.Fail(c, http.StatusUnauthorized, err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			model.Fail(c, http.StatusUnauthorized, "Token expirado")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		v, ok := claims["sub"].(float64)
		if !ok {
			model.Fail(c, http.StatusUnauthorized, "Erro ao ler subject")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		matricula := int(v)

		usuario, err := tH.service.LerUsuario(ctx, matricula)
		if err != nil || usuario.Matricula == 0 {
			model.Fail(c, http.StatusUnauthorized, err.Error())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("usuario", usuario)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		model.Fail(c, http.StatusUnauthorized, "Erro ao ler claims")
		return
	}
}

func (tH *TaskHandler) SingUp(c *gin.Context) {
	ctx := c.Request.Context()
	var usuario model.Usuario

	if err := c.ShouldBindJSON(&usuario); err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if usuario.Senha == "" {
		model.Fail(c, http.StatusBadRequest, "Erro: field Senha is empty")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(usuario.Senha), 10)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	usuario.Senha = string(hash)

	id, err := tH.service.CriarUsuario(ctx, &usuario)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, id)
}

func (tH *TaskHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var usuario model.Usuario

	if err := c.ShouldBindJSON(&usuario); err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	var usuarioDB model.Usuario
	usuarioDB, err := tH.service.LerUsuario(ctx, usuario.Email)
	if err != nil {
		model.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(usuarioDB.Senha), []byte(usuario.Senha))
	if err != nil {
		model.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usuarioDB.Matricula,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", true, true)
	model.OK(c, "Login executed successfully")
}

func (tH *TaskHandler) CriarMateria(c *gin.Context) {
	ctx := c.Request.Context()

	usuario, _ := c.Get("usuario")

	var materia model.Materia
	err := c.ShouldBindJSON(&materia)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	materia.Matricula = usuario.(model.Usuario).Matricula

	codigo, err := tH.service.CriarMateria(ctx, &materia)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, codigo)
}

func (tH *TaskHandler) CriarTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	usuario, _ := c.Get("usuario")

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, err.Error())
		return
	}

	var tarefa model.Tarefa
	err = c.ShouldBindJSON(&tarefa)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, err.Error())
		return
	}

	tarefa.Codigo = codigo

	id, err := tH.service.CriarTarefa(ctx, usuario.(model.Usuario).Matricula, &tarefa)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, id)
}

// READ
func (tH *TaskHandler) ListarMaterias(c *gin.Context) {
	ctx := c.Request.Context()

	usuario, _ := c.Get("usuario")

	materias, err := tH.service.ListarMaterias(ctx, usuario.(model.Usuario).Matricula)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, materias)
}

func (tH *TaskHandler) LerMateria(c *gin.Context) {
	ctx := c.Request.Context()

	usuario, _ := c.Get("usuario")

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	materias, err := tH.service.LerMateria(ctx, usuario.(model.Usuario).Matricula, codigo)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, materias)
}

func (tH *TaskHandler) ListarTarefas(c *gin.Context) {
	ctx := c.Request.Context()

	usuario, _ := c.Get("usuario")

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	tarefas, err := tH.service.ListarTarefas(ctx, usuario.(model.Usuario).Matricula, codigo)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, tarefas)
}

func (tH *TaskHandler) LerTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	usuario, _ := c.Get("usuario")

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	param = c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	tarefa, err := tH.service.LerTarefa(ctx, usuario.(model.Usuario).Matricula, codigo, id)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	model.OK(c, tarefa)
}

// UPDATE

func (tH *TaskHandler) MudarEmail(c *gin.Context) {
	ctx := c.Request.Context()

	usuarioReq, _ := c.Get("usuario")

	var usuario model.Usuario
	err := c.ShouldBindJSON(&usuario)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	err = tH.service.MudarEmail(ctx, usuario.Email, usuarioReq.(model.Usuario).Matricula)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, "Email do usuário atualizado")
}

func (tH *TaskHandler) MudarSenha(c *gin.Context) {
	ctx := c.Request.Context()

	usuarioReq, _ := c.Get("usuario")

	var usuario model.Usuario
	err := c.ShouldBindJSON(&usuario)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(usuario.Senha), 10)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	usuario.Senha = string(hash)

	err = tH.service.MudarSenha(ctx, usuario.Senha, usuarioReq.(model.Usuario).Matricula)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, "Senha do usuário atualizada")
}

func (tH *TaskHandler) MudarTema(c *gin.Context) {
	ctx := c.Request.Context()

	usuarioReq, _ := c.Get("usuario")

	var usuario model.Usuario
	err := c.ShouldBindJSON(&usuario)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	err = tH.service.MudarTema(ctx, usuario.Tema, usuarioReq.(model.Usuario).Matricula)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, "Tema do usuário atualizada")
}

func (tH *TaskHandler) MudarNomeMateria(c *gin.Context) {
	ctx := c.Request.Context()

	usuarioReq, _ := c.Get("usuario")

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	var materia model.Materia
	err = c.ShouldBindJSON(&materia)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	err = tH.service.MudarNomeMateria(ctx, materia.Nome, usuarioReq.(model.Usuario).Matricula, codigo)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, "Nome da matéria atualizado")
}

func (tH *TaskHandler) MudarNomeTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	usuario, _ := c.Get("usuario")

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	param = c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	var tarefa model.Tarefa
	err = c.ShouldBindJSON(&tarefa)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	err = tH.service.MudarNomeTarefa(ctx, tarefa.Nome, usuario.(model.Usuario).Matricula, codigo, id)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, "Nome da tarefa atualizada")
}

func (tH *TaskHandler) MudarPrazoTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	usuario, _ := c.Get("usuario")

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	param = c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	var tarefa model.Tarefa
	err = c.ShouldBindJSON(&tarefa)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	err = tH.service.MudarPrazoTarefa(ctx, tarefa.Prazo, usuario.(model.Usuario).Matricula, codigo, id)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, "Prazo da tarefa atualizada")
}

func (tH *TaskHandler) MudarAnotacaoTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	usuario, _ := c.Get("usuario")

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	param = c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	var tarefa model.Tarefa
	err = c.ShouldBindJSON(&tarefa)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	err = tH.service.MudarAnotacaoTarefa(ctx, tarefa.Anotacao, usuario.(model.Usuario).Matricula, codigo, id)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, "Anotação da tarefa atualizada")
}

func (tH *TaskHandler) MudarStatusTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	usuario, _ := c.Get("usuario")

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	param = c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	var tarefa model.Tarefa
	err = c.ShouldBindJSON(&tarefa)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	err = tH.service.MudarStatusTarefa(ctx, tarefa.Status, usuario.(model.Usuario).Matricula, codigo, id)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, "Status da tarefa atualizada")
}

// DELETE

func (tH *TaskHandler) DeletarUsuario(c *gin.Context) {
	ctx := c.Request.Context()

	usuarioReq, _ := c.Get("usuario")

	err := tH.service.DeletarUsuario(ctx, usuarioReq.(model.Usuario).Matricula)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, "Usuário deletado")
}

func (tH *TaskHandler) DeletarMateria(c *gin.Context) {
	ctx := c.Request.Context()

	usuarioReq, _ := c.Get("usuario")

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	err = tH.service.DeletarMateria(ctx, usuarioReq.(model.Usuario).Matricula, codigo)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, "Matéria deletada")
}

func (tH *TaskHandler) DeletarTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	usuario, _ := c.Get("usuario")

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	param = c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	err = tH.service.DeletarTarefa(ctx, usuario.(model.Usuario).Matricula, codigo, id)
	if err != nil {
		model.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	model.OK(c, "Tarefa deletada")
}
