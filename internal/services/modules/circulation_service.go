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
	loansCollection = "loans"
	loanDays        = 14
)

type Loan struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	BookID      primitive.ObjectID  `bson:"book_id" json:"book_id"`
	BookTitle   string              `bson:"book_title" json:"book_title"`
	PatronID    primitive.ObjectID  `bson:"patron_id" json:"patron_id"`
	PatronName  string              `bson:"patron_name" json:"patron_name"`
	CheckedOutAt time.Time         `bson:"checked_out_at" json:"checked_out_at"`
	DueDate     time.Time           `bson:"due_date" json:"due_date"`
	ReturnedAt  *time.Time          `bson:"returned_at,omitempty" json:"returned_at,omitempty"`
	Status      string              `bson:"status" json:"status"`
	CreatedAt   time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at" json:"updated_at"`
}

type CheckoutRequest struct {
	BookID    string `json:"book_id" binding:"required"`
	PatronID string `json:"patron_id" binding:"required"`
}

type ReturnRequest struct {
	LoanID string `json:"loan_id" binding:"required"`
}

type CirculationService struct {
	enabled bool
	mongo   *infrastructure.MongoConnectionManager
	logger  *logger.Logger
}

func NewCirculationService(mongo *infrastructure.MongoConnectionManager, enabled bool, logger *logger.Logger) *CirculationService {
	return &CirculationService{enabled: enabled, mongo: mongo, logger: logger}
}

func (s *CirculationService) Name() string      { return "Circulation Service" }
func (s *CirculationService) WireName() string   { return "circulation" }
func (s *CirculationService) Enabled() bool      { return s.enabled }
func (s *CirculationService) Endpoints() []string { return []string{"/checkout", "/returns", "/loans/active", "/loans/history"} }
func (s *CirculationService) Get() interface{}   { return s }

func (s *CirculationService) RegisterRoutes(g *gin.RouterGroup) {
	sub := g.Group("/circulation")
	sub.POST("/checkout", s.checkout)
	sub.POST("/returns", s.returns)
	sub.GET("/loans/active", s.activeLoans)
	sub.GET("/loans/history", s.loanHistory)
}

func (s *CirculationService) checkout(c *gin.Context) {
	var req CheckoutRequest
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
	if book.Status == "checked_out" {
		response.BadRequest(c, "Book is already checked out")
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
	now := time.Now()
	loan := Loan{
		ID:           primitive.NewObjectID(),
		BookID:       bookObjID,
		BookTitle:    book.Title,
		PatronID:     patronObjID,
		PatronName:   patron.Name,
		CheckedOutAt: now,
		DueDate:      now.AddDate(0, 0, loanDays),
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	_, err = conn.InsertOne(ctx, loansCollection, loan)
	if err != nil {
		s.logger.Error("Failed to create loan", err)
		response.InternalServerError(c, "Failed to checkout book")
		return
	}
	_, err = conn.UpdateOne(ctx, catalogCollection, bson.M{"_id": bookObjID}, bson.M{"$set": bson.M{"status": "checked_out", "updated_at": now}})
	if err != nil {
		s.logger.Error("Failed to update book status", err)
		response.InternalServerError(c, "Failed to update book status")
		return
	}
	response.Created(c, loan, "Book checked out successfully")
}

func (s *CirculationService) returns(c *gin.Context) {
	var req ReturnRequest
	if err := request.Bind(c, &req); err != nil {
		return
	}
	loanObjID, err := primitive.ObjectIDFromHex(req.LoanID)
	if err != nil {
		response.BadRequest(c, "Invalid loan ID format")
		return
	}
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	var loan Loan
	if err := conn.FindOne(ctx, loansCollection, bson.M{"_id": loanObjID, "status": "active"}).Decode(&loan); err != nil {
		response.NotFound(c, "Active loan not found")
		return
	}
	now := time.Now()
	_, err = conn.UpdateOne(ctx, loansCollection, bson.M{"_id": loanObjID}, bson.M{"$set": bson.M{"returned_at": now, "status": "returned", "updated_at": now}})
	if err != nil {
		s.logger.Error("Failed to update loan", err)
		response.InternalServerError(c, "Failed to process return")
		return
	}
	_, err = conn.UpdateOne(ctx, catalogCollection, bson.M{"_id": loan.BookID}, bson.M{"$set": bson.M{"status": "available", "updated_at": now}})
	if err != nil {
		s.logger.Error("Failed to update book status", err)
		response.InternalServerError(c, "Failed to update book status")
		return
	}
	response.Success(c, nil, "Book returned successfully")
}

func (s *CirculationService) activeLoans(c *gin.Context) {
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	cursor, err := conn.Find(ctx, loansCollection, bson.M{"status": "active"})
	if err != nil {
		s.logger.Error("Failed to list active loans", err)
		response.InternalServerError(c, "Failed to list active loans")
		return
	}
	defer cursor.Close(ctx)
	var loans []Loan
	if err := cursor.All(ctx, &loans); err != nil {
		s.logger.Error("Failed to decode loans", err)
		response.InternalServerError(c, "Failed to decode loans")
		return
	}
	if loans == nil {
		loans = []Loan{}
	}
	response.Success(c, loans, "Active loans listed")
}

func (s *CirculationService) loanHistory(c *gin.Context) {
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}
	ctx := c.Request.Context()
	cursor, err := conn.Find(ctx, loansCollection, bson.M{})
	if err != nil {
		s.logger.Error("Failed to list loans", err)
		response.InternalServerError(c, "Failed to list loans")
		return
	}
	defer cursor.Close(ctx)
	var loans []Loan
	if err := cursor.All(ctx, &loans); err != nil {
		s.logger.Error("Failed to decode loans", err)
		response.InternalServerError(c, "Failed to decode loans")
		return
	}
	if loans == nil {
		loans = []Loan{}
	}
	response.Success(c, loans, "Loan history listed")
}

func init() {
	registry.RegisterService("circulation_service", func(cfg *config.Config, log *logger.Logger, deps *registry.Dependencies) interfaces.Service {
		helper := registry.NewServiceHelper(cfg, log, deps)
		if !helper.IsServiceEnabled("circulation_service") {
			return nil
		}
		mongoManager, ok := registry.GetTyped[infrastructure.MongoConnectionManager](deps, "mongo")
		if !helper.RequireDependency("MongoConnectionManager", ok) {
			return nil
		}
		return NewCirculationService(&mongoManager, true, log)
	})
}
