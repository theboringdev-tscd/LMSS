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
	patronsCollection = "patrons"
)

type Patron struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name" binding:"required"`
	Email     string             `bson:"email" json:"email" binding:"required"`
	Phone     string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Active    bool               `bson:"active" json:"active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type PatronsService struct {
	enabled bool
	mongo   *infrastructure.MongoConnectionManager
	logger  *logger.Logger
}

func NewPatronsService(mongo *infrastructure.MongoConnectionManager, enabled bool, logger *logger.Logger) *PatronsService {
	return &PatronsService{enabled: enabled, mongo: mongo, logger: logger}
}

func (s *PatronsService) Name() string      { return "Patrons Service" }
func (s *PatronsService) WireName() string   { return "patrons" }
func (s *PatronsService) Enabled() bool      { return s.enabled }
func (s *PatronsService) Endpoints() []string { return []string{"/patrons"} }
func (s *PatronsService) Get() interface{}   { return s }

func (s *PatronsService) RegisterRoutes(g *gin.RouterGroup) {
	sub := g.Group("/patrons")
	sub.GET("", s.list)
	sub.GET("/search", s.search)
	sub.GET("/:id", s.get)
	sub.POST("", s.create)
	sub.PUT("/:id", s.update)
	sub.DELETE("/:id", s.delete)
}

func (s *PatronsService) list(c *gin.Context) {
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	cursor, err := conn.Find(ctx, patronsCollection, bson.M{})
	if err != nil {
		s.logger.Error("Failed to list patrons", err)
		response.InternalServerError(c, "Failed to list patrons")
		return
	}
	defer cursor.Close(ctx)
	var patrons []Patron
	if err := cursor.All(ctx, &patrons); err != nil {
		s.logger.Error("Failed to decode patrons", err)
		response.InternalServerError(c, "Failed to decode patrons")
		return
	}
	if patrons == nil {
		patrons = []Patron{}
	}
	response.Success(c, patrons, "Patrons listed")
}

func (s *PatronsService) search(c *gin.Context) {
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
				{"name": bson.M{"$regex": query, "$options": "i"}},
				{"email": bson.M{"$regex": query, "$options": "i"}},
			},
		}
	}
	cursor, err := conn.Find(ctx, patronsCollection, filter)
	if err != nil {
		s.logger.Error("Failed to search patrons", err)
		response.InternalServerError(c, "Failed to search patrons")
		return
	}
	defer cursor.Close(ctx)
	var patrons []Patron
	if err := cursor.All(ctx, &patrons); err != nil {
		s.logger.Error("Failed to decode search results", err)
		response.InternalServerError(c, "Failed to decode search results")
		return
	}
	if patrons == nil {
		patrons = []Patron{}
	}
	response.Success(c, patrons, "Search results")
}

func (s *PatronsService) get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Patron ID is required")
		return
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		response.BadRequest(c, "Invalid patron ID format")
		return
	}
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	var patron Patron
	err = conn.FindOne(ctx, patronsCollection, bson.M{"_id": objID}).Decode(&patron)
	if err != nil {
		response.NotFound(c, "Patron not found")
		return
	}
	response.Success(c, patron, "Patron retrieved")
}

func (s *PatronsService) create(c *gin.Context) {
	var patron Patron
	if err := request.Bind(c, &patron); err != nil {
		return
	}
	now := time.Now()
	patron.ID = primitive.NewObjectID()
	patron.CreatedAt = now
	patron.UpdatedAt = now
	if !patron.Active {
		patron.Active = true
	}
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	_, err := conn.InsertOne(ctx, patronsCollection, patron)
	if err != nil {
		s.logger.Error("Failed to create patron", err)
		response.InternalServerError(c, "Failed to create patron")
		return
	}
	response.Created(c, patron, "Patron created successfully")
}

func (s *PatronsService) update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Patron ID is required")
		return
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		response.BadRequest(c, "Invalid patron ID format")
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
	result, err := conn.UpdateOne(ctx, patronsCollection, bson.M{"_id": objID}, bson.M{"$set": updates})
	if err != nil {
		s.logger.Error("Failed to update patron", err)
		response.InternalServerError(c, "Failed to update patron")
		return
	}
	if result.MatchedCount == 0 {
		response.NotFound(c, "Patron not found")
		return
	}
	response.Success(c, nil, "Patron updated successfully")
}

func (s *PatronsService) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Patron ID is required")
		return
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		response.BadRequest(c, "Invalid patron ID format")
		return
	}
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	result, err := conn.DeleteOne(ctx, patronsCollection, bson.M{"_id": objID})
	if err != nil {
		s.logger.Error("Failed to delete patron", err)
		response.InternalServerError(c, "Failed to delete patron")
		return
	}
	if result.DeletedCount == 0 {
		response.NotFound(c, "Patron not found")
		return
	}
	response.Success(c, nil, "Patron deleted successfully")
}

func init() {
	registry.RegisterService("patrons_service", func(cfg *config.Config, log *logger.Logger, deps *registry.Dependencies) interfaces.Service {
		helper := registry.NewServiceHelper(cfg, log, deps)
		if !helper.IsServiceEnabled("patrons_service") {
			return nil
		}
		mongoManager, ok := registry.GetTyped[infrastructure.MongoConnectionManager](deps, "mongo")
		if !helper.RequireDependency("MongoConnectionManager", ok) {
			return nil
		}
		return NewPatronsService(&mongoManager, true, log)
	})
}
