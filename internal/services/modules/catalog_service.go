package modules

import (
	"time"

	"stackyrd/config"
	"stackyrd/pkg/infrastructure"
	"stackyrd/pkg/interfaces"
	"stackyrd/pkg/logger"
	"stackyrd/pkg/registry"
	"stackyrd/pkg/request"
	"stackyrd/pkg/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	catalogCollection = "catalog"
)

type Book struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title      string             `bson:"title" json:"title"`
	Author     string             `bson:"author" json:"author"`
	ISBN       string             `bson:"isbn,omitempty" json:"isbn,omitempty"`
	Genre      string             `bson:"genre,omitempty" json:"genre,omitempty"`
	Status     string             `bson:"status" json:"status"`
	Location   string             `bson:"location,omitempty" json:"location,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type CatalogService struct {
	enabled bool
	mongo   *infrastructure.MongoConnectionManager
	logger  *logger.Logger
}

func NewCatalogService(mongo *infrastructure.MongoConnectionManager, enabled bool, logger *logger.Logger) *CatalogService {
	return &CatalogService{enabled: enabled, mongo: mongo, logger: logger}
}

func (s *CatalogService) Name() string      { return "Catalog Service" }
func (s *CatalogService) WireName() string   { return "catalog" }
func (s *CatalogService) Enabled() bool      { return s.enabled }
func (s *CatalogService) Endpoints() []string { return []string{"/catalog"} }
func (s *CatalogService) Get() interface{}   { return s }

func (s *CatalogService) RegisterRoutes(g *gin.RouterGroup) {
	sub := g.Group("/catalog")
	sub.GET("", s.list)
	sub.GET("/search", s.search)
	sub.GET("/:id", s.get)
	sub.POST("", s.create)
	sub.PUT("/:id", s.update)
	sub.DELETE("/:id", s.delete)
}

func (s *CatalogService) list(c *gin.Context) {
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	cursor, err := conn.Find(ctx, catalogCollection, bson.M{})
	if err != nil {
		s.logger.Error("Failed to list books", err)
		response.InternalServerError(c, "Failed to list books")
		return
	}
	defer cursor.Close(ctx)
	var books []Book
	if err := cursor.All(ctx, &books); err != nil {
		s.logger.Error("Failed to decode books", err)
		response.InternalServerError(c, "Failed to decode books")
		return
	}
	if books == nil {
		books = []Book{}
	}
	response.Success(c, books, "Books listed")
}

func (s *CatalogService) search(c *gin.Context) {
	query := c.Query("q")
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	filter := bson.M{}
	if query != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"title": bson.M{"$regex": query, "$options": "i"}},
				{"author": bson.M{"$regex": query, "$options": "i"}},
				{"isbn": bson.M{"$regex": query, "$options": "i"}},
			},
		}
	}
	cursor, err := conn.Find(ctx, catalogCollection, filter)
	if err != nil {
		s.logger.Error("Failed to search books", err)
		response.InternalServerError(c, "Failed to search books")
		return
	}
	defer cursor.Close(ctx)
	var books []Book
	if err := cursor.All(ctx, &books); err != nil {
		s.logger.Error("Failed to decode search results", err)
		response.InternalServerError(c, "Failed to decode search results")
		return
	}
	if books == nil {
		books = []Book{}
	}
	response.Success(c, books, "Search results")
}

func (s *CatalogService) get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Book ID is required")
		return
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		response.BadRequest(c, "Invalid book ID format")
		return
	}
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	var book Book
	err = conn.FindOne(ctx, catalogCollection, bson.M{"_id": objID}).Decode(&book)
	if err != nil {
		response.NotFound(c, "Book not found")
		return
	}
	response.Success(c, book, "Book retrieved")
}

func (s *CatalogService) create(c *gin.Context) {
	var book Book
	if err := request.Bind(c, &book); err != nil {
		return
	}
	now := time.Now()
	book.ID = primitive.NewObjectID()
	book.CreatedAt = now
	book.UpdatedAt = now
	if book.Status == "" {
		book.Status = "available"
	}
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	_, err := conn.InsertOne(ctx, catalogCollection, book)
	if err != nil {
		s.logger.Error("Failed to create book", err)
		response.InternalServerError(c, "Failed to create book")
		return
	}
	response.Created(c, book, "Book created successfully")
}

func (s *CatalogService) update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Book ID is required")
		return
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		response.BadRequest(c, "Invalid book ID format")
		return
	}
	var updates map[string]interface{}
	if err := request.Bind(c, &updates); err != nil {
		return
	}
	updates["updated_at"] = time.Now()
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	result, err := conn.UpdateOne(ctx, catalogCollection, bson.M{"_id": objID}, bson.M{"$set": updates})
	if err != nil {
		s.logger.Error("Failed to update book", err)
		response.InternalServerError(c, "Failed to update book")
		return
	}
	if result.MatchedCount == 0 {
		response.NotFound(c, "Book not found")
		return
	}
	response.Success(c, nil, "Book updated successfully")
}

func (s *CatalogService) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Book ID is required")
		return
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		response.BadRequest(c, "Invalid book ID format")
		return
	}
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	result, err := conn.DeleteOne(ctx, catalogCollection, bson.M{"_id": objID})
	if err != nil {
		s.logger.Error("Failed to delete book", err)
		response.InternalServerError(c, "Failed to delete book")
		return
	}
	if result.DeletedCount == 0 {
		response.NotFound(c, "Book not found")
		return
	}
	response.Success(c, nil, "Book deleted successfully")
}

func init() {
	registry.RegisterService("catalog_service", func(cfg *config.Config, log *logger.Logger, deps *registry.Dependencies) interfaces.Service {
		helper := registry.NewServiceHelper(cfg, log, deps)
		if !helper.IsServiceEnabled("catalog_service") {
			return nil
		}
		mongoManager, ok := registry.GetTyped[infrastructure.MongoConnectionManager](deps, "mongo")
		if !helper.RequireDependency("MongoConnectionManager", ok) {
			return nil
		}
		return NewCatalogService(&mongoManager, true, log)
	})
}
