package integration

import (
	"avito-shop-service/internal/models"
	"avito-shop-service/internal/storage/postgresql"
	"context"
	"database/sql"
	"log"
	"log/slog"

	"github.com/stretchr/testify/suite"
)

type Repo interface {
	AddItemToInventory(context.Context, string, *models.Item) error
	GetItemByName(context.Context, string) (*models.Item, error)
	GetUserBalance(context.Context, string) (int, error)

	Create(string, string) (*models.User, error)
	GetUserByUsername(string) (*models.User, error)
}

type BuyItemIntegrationSuite struct {
	suite.Suite
	db *sql.DB
	r  Repo
}

func NewRepoSuite() *BuyItemIntegrationSuite {
	return &BuyItemIntegrationSuite{}
}

func (s *BuyItemIntegrationSuite) SetupSuite() {
	db, err := sql.Open("pgx", "postgres://user:psswd@localhost:5555/postgresDB")
	if err != nil {
		log.Fatal(err)
	}
	s.db = db
}

func (s *BuyItemIntegrationSuite) SetupTest() {
	s.r = postgresql.NewRep(s.db, slog.Default())
}

// func TestBuyItem
