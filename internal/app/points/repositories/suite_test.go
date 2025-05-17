package repositories_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Christian-007/fit-forge/internal/pkg/testutil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPointsRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PointsRepository Suite")
}

var _ = BeforeSuite(func() {
	ctx := context.Background()
	testutil.InitTestDb(ctx, os.Getenv("POSTGRES_TEST_DB_URL"))

	migrationFilePath := filepath.Join(os.Getenv("ROOT_DIR"), "/migrations")
	err := testutil.RunMigrations(ctx, os.Getenv("POSTGRES_TEST_DB_URL"), migrationFilePath)
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	db := testutil.GetTestDb()
	if db != nil {
		db.Close()
	}
})
