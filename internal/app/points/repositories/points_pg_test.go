package repositories_test

import (
	"context"

	"github.com/Christian-007/fit-forge/internal/app/points/repositories"
	"github.com/Christian-007/fit-forge/internal/pkg/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PointsRepository", func() {
	var (
		repo repositories.PointsRepository
		db   *pgxpool.Pool
		ctx  context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		db = testutil.GetTestDb()

		repo = repositories.NewPointsRepositoryPg(db)
	})

	Describe("Find users for subscription deduction", func() {
		It("should not return any error", func() {
			_, err := repo.FindUsersForSubscriptionDeduction(ctx)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
