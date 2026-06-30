package modules

import (
	"time"

	"stackyrd/config"
	"stackyrd/internal/middleware"
	"stackyrd/pkg/infrastructure"
	"stackyrd/pkg/interfaces"
	"stackyrd/pkg/logger"
	"stackyrd/pkg/registry"
	"stackyrd/pkg/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OverviewStats struct {
	BooksCheckedOutToday int64           `json:"books_checked_out_today"`
	NewPatronsThisWeek   int64           `json:"new_patrons_this_week"`
	OverdueReturns       int64           `json:"overdue_returns"`
	RecentActivity       []ActivityEntry `json:"recent_activity"`
}

type ActivityEntry struct {
	BookTitle  string `json:"book_title"`
	PatronName string `json:"patron_name"`
	Action     string `json:"action"`
	Timestamp  string `json:"timestamp"`
}

type ReportsService struct {
	enabled bool
	mongo   *infrastructure.MongoConnectionManager
	logger  *logger.Logger
	secret  string
}

func NewReportsService(mongo *infrastructure.MongoConnectionManager, enabled bool, secret string, logger *logger.Logger) *ReportsService {
	return &ReportsService{enabled: enabled, mongo: mongo, secret: secret, logger: logger}
}

func (s *ReportsService) Name() string       { return "Reports Service" }
func (s *ReportsService) WireName() string    { return "reports" }
func (s *ReportsService) Enabled() bool       { return s.enabled }
func (s *ReportsService) Endpoints() []string { return []string{"/reports/overview"} }
func (s *ReportsService) Get() interface{}    { return s }

func (s *ReportsService) RegisterRoutes(g *gin.RouterGroup) {
	sub := g.Group("/reports")
	sub.Use(middleware.JWTRequired(s.secret))
	sub.GET("/overview", s.overview)
}

func (s *ReportsService) overview(c *gin.Context) {
	conn, ok := s.mongo.GetConnection("primary")
	if !ok {
		response.InternalServerError(c, "Database connection unavailable")
		return
	}

	ctx := c.Request.Context()
	today := time.Now().Truncate(24 * time.Hour)
	weekAgo := today.AddDate(0, 0, -7)

	checkedOutToday, _ := conn.CountDocuments(ctx, "circulation", bson.M{
		"checked_out_at": bson.M{"$gte": today},
	})

	newPatrons, _ := conn.CountDocuments(ctx, "patrons", bson.M{
		"created_at": bson.M{"$gte": weekAgo},
	})

	overdue, _ := conn.CountDocuments(ctx, "circulation", bson.M{
		"returned": false,
		"due_date": bson.M{"$lt": time.Now()},
	})

	col := conn.Collection("circulation")
	findOpts := options.Find().SetSort(bson.M{"checked_out_at": -1}).SetLimit(10)
	cursor, err := col.Find(ctx, bson.M{}, findOpts)
	recent := []ActivityEntry{}
	if err == nil {
		defer cursor.Close(ctx)
		type circEntry struct {
			BookTitle  string    `bson:"book_title"`
			PatronName string    `bson:"patron_name"`
			Action     string    `bson:"action"`
			Timestamp  time.Time `bson:"checked_out_at"`
		}
		for cursor.Next(ctx) {
			var e circEntry
			if cursor.Decode(&e) == nil {
				action := e.Action
				if action == "" {
					action = "borrowed"
				}
				recent = append(recent, ActivityEntry{
					BookTitle:  e.BookTitle,
					PatronName: e.PatronName,
					Action:     action,
					Timestamp:  e.Timestamp.Format("2006-01-02 15:04"),
				})
			}
		}
	}

	response.Success(c, OverviewStats{
		BooksCheckedOutToday: checkedOutToday,
		NewPatronsThisWeek:   newPatrons,
		OverdueReturns:       overdue,
		RecentActivity:       recent,
	}, "Dashboard overview")
}

func init() {
	registry.RegisterService("reports_service", func(cfg *config.Config, log *logger.Logger, deps *registry.Dependencies) interfaces.Service {
		helper := registry.NewServiceHelper(cfg, log, deps)
		if !helper.IsServiceEnabled("reports_service") {
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
		return NewReportsService(&mongoManager, true, secret, log)
	})
}
