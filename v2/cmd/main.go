package main

import (
	"binai.net/v2/config"
	"binai.net/v2/internal/repository"
	"binai.net/v2/internal/router"

	"log"
)

func main() {
	//migrate := flag.Bool("migrate", false, "Apply migrations")
	//flag.Parse()

	cfg, err := config.InitConfig(".env")
	if err != nil {
		log.Fatal(err)
	}
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	//if !*migrate {
	//	err := migrations.ApplyMigrationsForAuth(db, cfg)
	//	if err != nil {
	//		log.Fatalf("Failed to apply migrations: %v", err)
	//	}
	//	log.Println("Migrations applied successfully")
	//}

	repo := repository.InitRepositories(db)

	router := router.SetupRouter(cfg, repo)

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start course-certificate_service: %v\n", err)
	}
}
