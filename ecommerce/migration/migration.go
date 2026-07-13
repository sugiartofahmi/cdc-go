package migration

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func Run(db *gorm.DB, exec string) {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	m, err := initializeMigrator(sqlDB)
	if err != nil {
		panic(err)
	}

	handleDirtyMigration(m)

	if err := executeMigration(m, exec); err != nil {
		panic(err)
	}

	fmt.Println("migration applied successfully!")
}

func Create(_ *gorm.DB, fileName string) {
	if fileName == "" {
		panic("--fileName is required for create command")
	}

	nextNum := getNextSequentialNumber()
	prefix := fmt.Sprintf("%06d", nextNum)

	dir := "migration/files"
	upFile := fmt.Sprintf("%s/%s_%s.up.sql", dir, prefix, fileName)
	downFile := fmt.Sprintf("%s/%s_%s.down.sql", dir, prefix, fileName)

	if err := os.WriteFile(upFile, []byte(""), 0644); err != nil {
		panic(err)
	}
	if err := os.WriteFile(downFile, []byte(""), 0644); err != nil {
		panic(err)
	}

	fmt.Printf("created migration files:\n%s\n%s\n", upFile, downFile)
}

func initializeMigrator(sqlDB *sql.DB) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{
		MigrationsTable: "migrations",
	})
	if err != nil {
		return nil, fmt.Errorf("could not create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migration/files", "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("could not create migrate instance: %w", err)
	}

	return m, nil
}

func handleDirtyMigration(m *migrate.Migrate) {
	version, dirty, _ := m.Version()
	if dirty {
		fmt.Printf("migration is dirty at version %d, forcing...\n", version)
		if err := m.Force(int(version)); err != nil {
			panic(err)
		}
	}
}

func executeMigration(m *migrate.Migrate, exec string) error {
	var err error

	switch exec {
	case "down":
		err = m.Steps(-1)
	case "fresh":
		if err = m.Down(); err != nil && err != migrate.ErrNoChange {
			return err
		}
		err = m.Up()
	default:
		err = m.Up()
	}

	if err == migrate.ErrNoChange {
		fmt.Println("no migration changes")
		return nil
	}

	return err
}

func getNextSequentialNumber() int {
	entries, err := os.ReadDir("migration/files")
	if err != nil {
		return 1
	}

	re := regexp.MustCompile(`^(\d+)_`)
	var nums []int
	for _, e := range entries {
		if matches := re.FindStringSubmatch(e.Name()); len(matches) > 1 {
			if n, err := strconv.Atoi(matches[1]); err == nil {
				nums = append(nums, n)
			}
		}
	}

	if len(nums) == 0 {
		return 1
	}
	sort.Ints(nums)
	return nums[len(nums)-1] + 1
}
