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
	reservationsCollection = "reservations"
)

type Reservation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookID     primitive.ObjectID `bson:"book_id" json:"book_id"`
	BookTitle  string             `bson:"book_title" json:"book_title"`
	PatronID   primitive.ObjectID `bson:"patron_id" json:"patron_id"`
	PatronName string             `bson:"patron_name" json:"patron_name"`
	Status     string             `bson:"status" json:"status"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type ReservationRequest struct {
	BookID    string `json:"book_id" binding:"required"`
	PatronID string `json:"patron_id" binding:"required"`
}

type ReservationsService struct {
	enabled bool
	mongo   *infrastructure.MongoConnectionManager
	logger  *logger.Logger
}

func NewReservationsService(mongo *infrastructure.MongoConnectionManager, enabled bool, logger *logger.Logger) *ReservationsService {
	return &ReservationsService{enabled: enabled, mongo: mongo, logger: logger}
}

func (s *ReservationsService) Name() string      { return "Reservations Service" }
func (s *ReservationsService) WireName() string   { return "reservations" }
func (s *ReservationsService) Enabled() bool      { return s.enabled }
func (s *ReservationsService) Endpoints() []string { return []string{"/reservations"} }
func (s *ReservationsService) Get() interface{}   { return s }

func (s *ReservationsService) RegisterRoutes(g *gin.RouterGroup) {
	sub := g.Group("/reservations")
	sub.GET("", s.list)
	sub.POST("", s.create)
	sub.DELETE("/:id", s.delete)
}

func (s *ReservationsService) list(c *gin.Context) {
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	cursor, err := conn.Find(ctx, reservationsCollection, bson.M{})
	if err != nil {
		s.logger.Error("Failed to list reservations", err)
		response.InternalServerError(c, "Failed to list reservations")
		return
	}
	defer cursor.Close(ctx)
	var reservations []Reservation
	if err := cursor.All(ctx, &reservations); err != nil {
		s.logger.Error("Failed to decode reservations", err)
		response.InternalServerError(c, "Failed to decode reservations")
		return
	}
	if reservations == nil {
		reservations = []Reservation{}
	}
	response.Success(c, reservations, "Reservations listed")
}

func (s *ReservationsService) create(c *gin.Context) {
	var req ReservationRequest
	if err := request.Bind(c, &req); err != nil {
		return
	}
	bookObjID, err := primitive.ObjectIDFromHex(req.BookID)
	if err != nil {
		response.BadRequest(c, "Invalid book ID format")
		return
	}
	patronObjID, err := primitive.ObjectIDFromHex(req.PatronID)
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
	var book Book
	if err := conn.FindOne(ctx, catalogCollection, bson.M{"_id": bookObjID}).Decode(&book); err != nil {
		response.NotFound(c, "Book not found")
		return
	}
	var patron Patron
	if err := conn.FindOne(ctx, patronsCollection, bson.M{"_id": patronObjID}).Decode(&patron); err != nil {
		response.NotFound(c, "Patron not found")
		return
	}
	if !patron.Active {
		response.BadRequest(c, "Patron account is not active")
		return
	}
	if exists, _ := conn.CountDocuments(ctx, reservationsCollection, bson.M{"book_id": bookObjID, "status": "pending"}); exists > 0 {
		response.BadRequest(c, "Book is already reserved")
		return
	}
	now := time.Now()
	reservation := Reservation{
		ID:         primitive.NewObjectID(),
		BookID:     bookObjID,
		BookTitle:  book.Title,
		PatronID:   patronObjID,
		PatronName: patron.Name,
		Status:     "pending",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	_, err = conn.InsertOne(ctx, reservationsCollection, reservation)
	if err != nil {
		s.logger.Error("Failed to create reservation", err)
		response.InternalServerError(c, "Failed to create reservation")
		return
	}
	response.Created(c, reservation, "Reservation created successfully")
}

func (s *ReservationsService) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Reservation ID is required")
		return
	}
	resObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		response.BadRequest(c, "Invalid reservation ID format")
		return
	}
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	result, err := conn.DeleteOne(ctx, reservationsCollection, bson.M{"_id": resObjID})
	if err != nil {
		s.logger.Error("Failed to delete reservation", err)
		response.InternalServerError(c, "Failed to delete reservation")
		return
	}
	if result.DeletedCount == 0 {
		response.NotFound(c, "Reservation not found")
		return
	}
	response.Success(c, nil, "Reservation deleted successfully")
}

func init() {
	registry.RegisterService("reservations_service", func(cfg *config.Config, log *logger.Logger, deps *registry.Dependencies) interfaces.Service {
		helper := registry.NewServiceHelper(cfg, log, deps)
		if !helper.IsServiceEnabled("reservations_service") {
			return nil
		}
		mongoManager, ok := registry.GetTyped[infrastructure.MongoConnectionManager](deps, "mongo")
		if !helper.RequireDependency("MongoConnectionManager", ok) {
			return nil
		}
		return NewReservationsService(&mongoManager, true, log)
})
}