package repositories_test

import (
	"context"

	pointsrepositories "github.com/Christian-007/fit-forge/internal/app/points/repositories"
	"github.com/Christian-007/fit-forge/internal/app/users/domains"
	usersrepositories "github.com/Christian-007/fit-forge/internal/app/users/repositories"
	sharedmodel "github.com/Christian-007/fit-forge/internal/pkg/model"
	"github.com/Christian-007/fit-forge/internal/pkg/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PointsRepository", func() {
	var (
		pointsRepo pointsrepositories.PointsRepository
		usersRepo  usersrepositories.UserRepository
		db         *pgxpool.Pool
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		db = testutil.GetTestDb()

		pointsRepo = pointsrepositories.NewPointsRepositoryPg(db)
		usersRepo = usersrepositories.NewUserRepositoryPg(db)

		_, err := db.Exec(ctx, `TRUNCATE users RESTART IDENTITY CASCADE`)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Find users for subscription deduction", func() {
		When("there is no error", func() {
			It("should return a list of users that have enough points and not", func() {
				userWithPoints, err := usersRepo.CreateWithInitialPoints(ctx, domains.UserModel{Name: "John Test", Email: "johntest@gmail.com", Password: []byte("test")})
				Expect(err).NotTo(HaveOccurred())

				users, err := pointsRepo.FindUsersForSubscriptionDeduction(ctx, userWithPoints.CreatedAt.AddDate(0, 1, 0).Format("2006-01-02"))
				Expect(err).ToNot(HaveOccurred())

				mockEligibleUsersResponse := []sharedmodel.UserWithPoints{
					{Id: userWithPoints.Id, Email: userWithPoints.Email, TotalPoints: userWithPoints.Point.TotalPoints},
				}

				Expect(users.EligibleForDeduction).To(Equal(mockEligibleUsersResponse))
				Expect(users.InsufficientPoints).To(BeNil())
			})
		})
	})
})
