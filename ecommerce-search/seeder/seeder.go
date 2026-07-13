package seeder

import (
	"log"

	"gorm.io/gorm"
)

func Run(db *gorm.DB, seederCommands []string) error {
	if len(seederCommands) > 0 {
		listSeeders := map[string]Seeder{
			"RoleSeeder": NewRoleSeeder(),
			"UserSeeder": NewUserSeeder(),
		}
		for _, name := range seederCommands {
			s, ok := listSeeders[name]
			if !ok {
				log.Printf("unknown seeder: %s", name)
				continue
			}
			if err := s.Handle(db); err != nil {
				return err
			}
		}
	} else {
		db.Exec(`DELETE FROM users`)
		db.Exec(`DELETE FROM roles`)

		for _, s := range []Seeder{
			NewRoleSeeder(),
			NewUserSeeder(),
		} {
			if err := s.Handle(db); err != nil {
				return err
			}
		}
	}
	log.Println("Seeding completed!")
	return nil
}
