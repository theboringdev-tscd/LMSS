package modules

import (
	"context"
	"time"

	"stackyrd/config"
	"stackyrd/internal/middleware"
	"stackyrd/pkg/infrastructure"
	"stackyrd/pkg/interfaces"
	"stackyrd/pkg/logger"
	"stackyrd/pkg/registry"
	"stackyrd/pkg/request"
	"stackyrd/pkg/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	authCollection  = "users"
	authDefaultPass = "admin"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Name         string             `bson:"name" json:"name"`
	Role         string             `bson:"role" json:"role"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type AuthService struct {
	enabled bool
	mongo   *infrastructure.MongoConnectionManager
	logger  *logger.Logger
	secret  string
}

func NewAuthService(mongo *infrastructure.MongoConnectionManager, enabled bool, secret string, logger *logger.Logger) *AuthService {
	return &AuthService{enabled: enabled, mongo: mongo, secret: secret, logger: logger}
}

func (s *AuthService) Name() string         { return "Auth Service" }
func (s *AuthService) WireName() string      { return "auth" }
func (s *AuthService) Enabled() bool         { return s.enabled }
func (s *AuthService) Endpoints() []string   { return []string{"/auth/login", "/auth/refresh", "/auth/logout"} }
func (s *AuthService) Get() interface{}      { return s }

func (s *AuthService) RegisterRoutes(g *gin.RouterGroup) {
	sub := g.Group("/auth")
	sub.POST("/login", s.login)
	sub.POST("/refresh", s.refresh)
	sub.POST("/logout", s.logout)
}

func (s *AuthService) login(c *gin.Context) {
	var req LoginRequest
	if err := request.Bind(c, &req); err != nil {
		return
	}

	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}

	ctx := c.Request.Context()
	var user User
	err := conn.FindOne(ctx, authCollection, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		response.Unauthorized(c, "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Unauthorized(c, "Invalid email or password")
		return
	}

	token, err := middleware.GenerateToken(
		user.ID.Hex(),
		user.Name,
		user.Email,
		user.Role,
		s.secret,
		24*time.Hour,
	)
	if err != nil {
		s.logger.Error("Failed to generate token", err)
		response.InternalServerError(c, "Failed to generate token")
		return
	}

	response.Success(c, LoginResponse{Token: token, User: user}, "Login successful")
}

func (s *AuthService) refresh(c *gin.Context) {
	userID := middleware.GetUserID(c)
	username := middleware.GetUsername(c)
	email := middleware.GetUserEmail(c)
	role := middleware.GetUserRole(c)

	if userID == "" {
		response.Unauthorized(c, "No valid token")
		return
	}

	token, err := middleware.GenerateToken(userID, username, email, role, s.secret, 24*time.Hour)
	if err != nil {
		s.logger.Error("Failed to refresh token", err)
		response.InternalServerError(c, "Failed to refresh token")
		return
	}

	response.Success(c, map[string]string{"token": token}, "Token refreshed")
}

func (s *AuthService) logout(c *gin.Context) {
	response.Success(c, nil, "Logged out successfully")
}

func init() {
	registry.RegisterService("auth_service", func(cfg *config.Config, log *logger.Logger, deps *registry.Dependencies) interfaces.Service {
		helper := registry.NewServiceHelper(cfg, log, deps)
		if !helper.IsServiceEnabled("auth_service") {
			return nil
		}

		mongoManager, ok := registry.GetTyped[infrastructure.MongoConnectionManager](deps, "mongo")
		if !helper.RequireDependency("MongoConnectionManager", ok) {
			return nil
		}

		secret := cfg.Auth.Secret
		if secret == "" {
			secret = "lmss-default-secret"
		}

		svc := NewAuthService(&mongoManager, true, secret, log)
		svc.seedAdmin()
		return svc
	})
}

func (s *AuthService) seedAdmin() {
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		s.logger.Warn("Cannot seed admin user: mongo connection unavailable")
		return
	}

	ctx := newContext()
	count, err := conn.CountDocuments(ctx, authCollection, bson.M{})
	if err != nil || count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(authDefaultPass), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash admin password", err)
		return
	}

	admin := User{
		ID:           primitive.NewObjectID(),
		Email:        "admin@library.org",
		PasswordHash: string(hash),
		Name:         "Admin Librarian",
		Role:         "admin",
		CreatedAt:    time.Now(),
	}

	_, err = conn.InsertOne(ctx, authCollection, admin)
	if err != nil {
		s.logger.Error("Failed to seed admin user", err)
		return
	}
	s.logger.Info("Seeded default admin user (admin@library.org / admin)")
}

func newContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}
