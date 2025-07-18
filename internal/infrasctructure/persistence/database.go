package persistence

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetDB(database *gorm.DB) {
	db = database
}

func GetDB() *gorm.DB {
	return db
}

func InitDB() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("⚠️  .env não encontrado, usando variáveis de ambiente padrão")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL não definida no ambiente")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Falha ao conectar no banco de dados: %v", err)
	}

	err = database.AutoMigrate(
		&UsuarioModel{},
		&Pedido{},
	)
	if err != nil {
		log.Fatalf("❌ Falha na migração do banco: %v", err)
	}

	fmt.Println("✅ Banco de dados conectado com sucesso")
	SetDB(database)
}
