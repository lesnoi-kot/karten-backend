package entityservices_test

import (
	"os"
	"testing"
	"time"

	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

var (
	testDate = time.Unix(0, 0).UTC()
	testUser = store.User{
		ID:       33,
		SocialID: "ae4f869a3f59db15cd8173b451c97025",
		Name:     "john_doe",
		Login:    "John Doe",
		Email:    "john_doe@hotmail.com",
		URL:      "www.google.com",
	}
	fixturesFS   = os.DirFS("./testdata")
	storeService *store.Store
)

func TestIntegrationEntityservices(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") == "" {
		t.Skip("entityservices intergation tests are skipped")
	}

	storeService = store.NewStore(store.StoreConfig{
		DSN:    os.Getenv("STORE_DSN"),
		Logger: zap.NewNop().Sugar(),
		Debug:  false,
	})

	suite.Run(t, new(projectserviceSuite))
	suite.Run(t, new(userserviceSuite))
}
