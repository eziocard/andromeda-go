package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/eziocard/andromeda-go/initializers"
	"github.com/eziocard/andromeda-go/internal/dto"
	"github.com/eziocard/andromeda-go/internal/models"
	"github.com/eziocard/andromeda-go/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserCreate(c *gin.Context) {
	var body struct {
		Name       string `json:"name" binding:"required"`
		LastName   string `json:"lastname" binding:"required"`
		Contact    string `json:"contact" binding:"required"`
		Email      string `json:"email" binding:"required,email"`
		Password   string `json:"password" binding:"required,min=8"`
		RoleID     uint   `json:"roleId" binding:"required"`
		BusinessID *uint  `json:"businessId"` // opcional
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	email := strings.ToLower(strings.TrimSpace(body.Email))

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "no se pudo procesar la contraseña",
		})
		return
	}

	var role models.Role
	if err := initializers.DB.First(&role, body.RoleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el rol especificado no existe"})
		return
	}

	// Solo validamos el negocio si fue enviado
	if body.BusinessID != nil {
		var business models.Business
		if err := initializers.DB.First(&business, *body.BusinessID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "el negocio especificado no existe"})
			return
		}
	}

	user := models.User{
		Name:       body.Name,
		LastName:   body.LastName,
		Contact:    body.Contact,
		Email:      email,
		Password:   string(hash),
		RoleID:     body.RoleID,
		BusinessID: body.BusinessID,
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no se pudo crear el usuario: " + err.Error(),
		})
		return
	}

	if err := initializers.DB.Preload("Role").Preload("Business").First(&user, user.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "usuario creado pero no se pudo cargar la información completa",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
}

const (
	maxFailedAttempts = 5
	lockoutDuration   = time.Minute * 15
)

// hash "dummy" para comparar aunque el usuario no exista y así
// evitar que el tiempo de respuesta delate si el email está registrado.
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("dummy-password"), bcrypt.DefaultCost)

func UserLogin(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	email := strings.ToLower(strings.TrimSpace(body.Email))

	var user models.User
	err := initializers.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		// Comparamos contra un hash dummy para no filtrar por timing si el email existe o no.
		bcrypt.CompareHashAndPassword(dummyHash, []byte(body.Password))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "cuenta bloqueada temporalmente, intenta más tarde",
		})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "cuenta inactiva"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		registerFailedAttempt(&user)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	// login correcto: reseteamos contador de intentos fallidos
	user.FailedLoginAttempts = 0
	user.LockedUntil = nil
	now := time.Now()
	user.LastLogin = &now
	initializers.DB.Save(&user)

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo generar el token"})
		return
	}

	refreshToken, err := issueRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo generar el refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         user,
	})
}

func registerFailedAttempt(user *models.User) {
	user.FailedLoginAttempts++
	if user.FailedLoginAttempts >= maxFailedAttempts {
		lockUntil := time.Now().Add(lockoutDuration)
		user.LockedUntil = &lockUntil
		user.FailedLoginAttempts = 0
	}
	initializers.DB.Save(user)
}

func issueRefreshToken(userID uint) (string, error) {
	rawToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	refreshToken := models.RefreshToken{
		UserID:    userID,
		TokenHash: utils.HashToken(rawToken),
		ExpiresAt: time.Now().Add(utils.RefreshTokenTTL),
	}

	if err := initializers.DB.Create(&refreshToken).Error; err != nil {
		return "", err
	}

	return rawToken, nil
}

func UserRefresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenHash := utils.HashToken(body.RefreshToken)

	var stored models.RefreshToken
	if err := initializers.DB.Where("token_hash = ? AND revoked = ?", tokenHash, false).First(&stored).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token inválido"})
		return
	}

	if stored.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expirado"})
		return
	}

	var user models.User
	if err := initializers.DB.First(&user, stored.UserID).Error; err != nil || !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no disponible"})
		return
	}

	// Rotación: revocamos el token usado y emitimos uno nuevo
	initializers.DB.Model(&stored).Update("revoked", true)

	newRefreshToken, err := issueRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo generar el refresh token"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo generar el token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": newRefreshToken,
	})
}

func UserLogout(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenHash := utils.HashToken(body.RefreshToken)
	initializers.DB.Model(&struct {
	}{}).Table("refresh_tokens").
		Where("token_hash = ?", tokenHash).
		Update("revoked", true)

	c.JSON(http.StatusOK, gin.H{"message": "sesión cerrada"})
}

func UsersIndex(c *gin.Context) {
	var users []models.User

	if err := initializers.DB.
		Preload("Business").
		Preload("Role").
		Find(&users).Error; err != nil {

		c.JSON(500, gin.H{
			"error": "No se pudieron obtener los usuarios",
		})
		return
	}

	var response []dto.UserResponse

	for _, user := range users {
		userResponse := dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			LastName:  user.LastName,
			Contact:   user.Contact,
			Email:     user.Email,
			LastLogin: user.LastLogin,
		}

		if user.Business != nil {
			userResponse.Business = &dto.BusinessResponse{
				ID:   user.Business.ID,
				Name: user.Business.Name,
			}
		}

		userResponse.Role = &dto.RoleResponse{
			ID:   user.Role.ID,
			Name: user.Role.Name,
		}

		response = append(response, userResponse)
	}

	c.JSON(200, gin.H{
		"users": response,
	})
}

func UserMe(c *gin.Context) {
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}

	user, ok := userValue.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno de autenticación"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func UserDelete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de usuario requerido"})
		return
	}

	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}

	if currentUserValue, exists := c.Get("user"); exists {
		if currentUser, ok := currentUserValue.(models.User); ok {
			if currentUser.ID == user.ID {
				c.JSON(http.StatusBadRequest, gin.H{"error": "no podés eliminar tu propia cuenta desde aquí"})
				return
			}
		}
	}

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo iniciar la transacción"})
		return
	}

	if err := tx.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", user.ID, false).
		Update("revoked", true).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron revocar las sesiones del usuario"})
		return
	}

	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo eliminar el usuario"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo confirmar la eliminación"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "usuario eliminado correctamente"})
}

func UserPatch(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de usuario requerido"})
		return
	}

	var body struct {
		Name       *string `json:"name"`
		LastName   *string `json:"lastname"`
		Contact    *string `json:"contact"`
		Email      *string `json:"email" binding:"omitempty,email"`
		Password   *string `json:"password"`
		RoleID     *uint   `json:"roleId"`
		BusinessID *uint   `json:"businessId"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}

	if body.RoleID != nil {
		var role models.Role
		if err := initializers.DB.First(&role, *body.RoleID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "el rol especificado no existe"})
			return
		}
		user.RoleID = *body.RoleID
	}

	if body.BusinessID != nil {
		var business models.Business
		if err := initializers.DB.First(&business, *body.BusinessID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "el negocio especificado no existe"})
			return
		}
		user.BusinessID = body.BusinessID
	}

	if body.Name != nil {
		user.Name = *body.Name
	}
	if body.LastName != nil {
		user.LastName = *body.LastName
	}
	if body.Contact != nil {
		user.Contact = *body.Contact
	}

	if body.Email != nil {
		newEmail := strings.ToLower(strings.TrimSpace(*body.Email))

		// Evitamos colisión con el email de otro usuario
		var existing models.User
		if err := initializers.DB.Where("email = ? AND id != ?", newEmail, user.ID).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "ya existe otro usuario con ese email"})
			return
		}

		user.Email = newEmail
	}

	// Contraseña: solo se actualiza si viene no vacía
	if body.Password != nil && *body.Password != "" {
		if len(*body.Password) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "la contraseña debe tener al menos 8 caracteres"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(*body.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo procesar la contraseña"})
			return
		}
		user.Password = string(hash)
	}

	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo actualizar el usuario: " + err.Error()})
		return
	}

	if err := initializers.DB.Preload("Role").Preload("Business").First(&user, user.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "usuario actualizado pero no se pudo cargar la información completa"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"user": user})
}
