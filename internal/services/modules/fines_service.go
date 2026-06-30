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
	finesCollection = "fines"
)

type Fine struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	LoanID     primitive.ObjectID `bson:"loan_id" json:"loan_id"`
	PatronID   primitive.ObjectID `bson:"patron_id" json:"patron_id"`
	PatronName string             `bson:"patron_name" json:"patron_name"`
	BookTitle  string             `bson:"book_title" json:"book_title"`
	Amount     float64            `bson:"amount" json:"amount"`
	Paid       bool               `bson:"paid" json:"paid"`
	Status     string             `bson:"status" json:"status"`
	DueDate    time.Time          `bson:"due_date" json:"due_date"`
	PaidAt     *time.Time         `bson:"paid_at,omitempty" json:"paid_at,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type PayFineRequest struct {
	FineID string `json:"fine_id" binding:"required"`
}

type FinesService struct {
	enabled bool
	mongo   *infrastructure.MongoConnectionManager
	logger  *logger.Logger
}

func NewFinesService(mongo *infrastructure.MongoConnectionManager, enabled bool, logger *logger.Logger) *FinesService {
	return &FinesService{enabled: enabled, mongo: mongo, logger: logger}
}

func (s *FinesService) Name() string      { return "Fines Service" }
func (s *FinesService) WireName() string   { return "fines" }
func (s *FinesService) Enabled() bool      { return s.enabled }
func (s *FinesService) Endpoints() []string { return []string{"/fines", "/fines/pay", "/patrons/:id/fines"} }
func (s *FinesService) Get() interface{}   { return s }

func (s *FinesService) RegisterRoutes(g *gin.RouterGroup) {
	sub := g.Group("/fines")
	sub.GET("", s.list)
	sub.POST("/pay", s.pay)
	sub.GET("/patrons/:id", s.byPatron)
}

func (s *FinesService) list(c *gin.Context) {
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	cursor, err := conn.Find(ctx, finesCollection, bson.M{})
	if err != nil {
		s.logger.Error("Failed to list fines", err)
		response.InternalServerError(c, "Failed to list fines")
		return
	}
	defer cursor.Close(ctx)
	var fines []Fine
	if err := cursor.All(ctx, &fines); err != nil {
		s.logger.Error("Failed to decode fines", err)
		response.InternalServerError(c, "Failed to decode fines")
		return
	}
	if fines == nil {
		fines = []Fine{}
	}
	response.Success(c, fines, "Fines listed")
}

func (s *FinesService) byPatron(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Patron ID is required")
		return
	}
	patronObjID, err := primitive.ObjectIDFromHex(id)
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
	cursor, err := conn.Find(ctx, finesCollection, bson.M{"patron_id": patronObjID})
	if err != nil {
		s.logger.Error("Failed to list patron fines", err)
		response.InternalServerError(c, "Failed to list patron fines")
		return
	}
	defer cursor.Close(ctx)
	var fines []Fine
	if err := cursor.All(ctx, &fines); err != nil {
		s.logger.Error("Failed to decode fines", err)
		response.InternalServerError(c, "Failed to decode fines")
		return
	}
	if fines == nil {
		fines = []Fine{}
	}
	response.Success(c, fines, "Patron fines listed")
}

func (s *FinesService) pay(c *gin.Context) {
	var req PayFineRequest
	if err := request.Bind(c, &req); err != nil {
		return
	}
	fineObjID, err := primitive.ObjectIDFromHex(req.FineID)
	if err != nil {
		response.BadRequest(c, "Invalid fine ID format")
		return
	}
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	now := time.Now()
	result, err := conn.UpdateOne(ctx, finesCollection, bson.M{"_id": fineObjID}, bson.M{"$set": bson.M{"paid": true, "status": "paid", "paid_at": now, "updated_at": now}})
	if err != nil {
		s.logger.Error("Failed to pay fine", err)
		response.InternalServerError(c, "Failed to pay fine")
		return
	}
	if result.MatchedCount == 0 {
		response.NotFound(c, "Fine not found")
		return
	}
	var fine Fine
	if err := conn.FindOne(ctx, finesCollection, bson.M{"_id": fineObjID}).Decode(&fine); err != nil {
		response.InternalServerError(c, "Failed to load updated fine")
		return
	}
	response.Success(c, fine, "Fine paid successfully")
}

func init() {
	registry.RegisterService("fines_service", func(cfg *config.Config, log *logger.Logger, deps *registry.Dependencies) interfaces.Service {
		helper := registry.NewServiceHelper(cfg, log, deps)
		if !helper.IsServiceEnabled("fines_service") {
			return nil
		}
		mongoManager, ok := registry.GetTyped[infrastructure.MongoConnectionManager](deps, "mongo")
		if !helper.RequireDependency("MongoConnectionManager", ok) {
			return nil
		}
		return NewFinesService(&mongoManager, true, log)
	})
}
