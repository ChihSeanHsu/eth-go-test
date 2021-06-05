package models_test

import (
	"context"
	"fmt"
	"github.com/eth-go-test/pkg/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func setupDB() *models.DB {
	db := models.InitDB(10, 1)
	db.AutoMigrate(&models.Block{}, &models.Transaction{}, &models.Log{})
	return db
}

func tearDownDB(db *models.DB) {
	db.Migrator().DropTable(&models.Log{})
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.Block{})
}

func truncateTable(db *models.DB) {
	var modelArray []interface{}
	modelArray = append(modelArray, &models.Log{}, &models.Transaction{}, &models.Block{})
	sql := "TRUNCATE TABLE  %s;"
	stmt := &gorm.Statement{DB: db.DB}
	for _, m := range modelArray {
		stmt.Parse(m)
		db.Raw(fmt.Sprintf(sql, stmt.Table))
	}
}

var _ = Describe("Model test", func() {
	var db *models.DB
	BeforeEach(func() {
		db = setupDB()
	})
	AfterEach(func() {
		tearDownDB(db)
	})
	Describe("test GetBlockByID", func() {
		AfterEach(func() {
			truncateTable(db)
		})
		Context("success", func() {
			It("found", func() {
				ctx := context.Background()
				block, err := db.GetBlockByID(ctx, 1)
				Expect(block).To(Equal(models.Block{}))
				Expect(err).To(Equal(models.ErrNotFound))
			})
		})
		Context("failed", func() {
			It("not found", func() {
				ctx := context.Background()
				block, err := db.GetBlockByID(ctx, 1)
				Expect(block).To(Equal(models.Block{}))
				Expect(err).To(Equal(models.ErrNotFound))
			})
		})
	})
	Describe("test GetBlocks", func() {
		Context("failed", func() {
			It("not found", func() {
				ctx := context.Background()
				blocks, err := db.GetBlocks(ctx, 1)
				Expect(blocks).To(Equal([]models.Block{}))
				Expect(err).To(Equal(models.ErrNotFound))
			})
		})
	})
	Describe("test GetTxByHash", func() {
		Context("failed", func() {
			It("not found", func() {
				ctx := context.Background()
				tx, err := db.GetTxByHash(ctx, "test")
				Expect(tx).To(Equal(models.Transaction{}))
				Expect(err).To(Equal(models.ErrNotFound))
			})
		})
	})
})
