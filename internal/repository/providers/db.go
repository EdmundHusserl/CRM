package providers

import (
	"fmt"
	"os"
)

func getEnvOrDefault(envVarName, defaultTo string) string {
	envVar := os.Getenv(envVarName)
	if len(envVar) == 0 {
		return defaultTo
	}
	return envVar
}

// Produces connection string
func GetConnectionString() string {
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	password := os.Getenv("DB_PASSWORD")
	dbName := getEnvOrDefault("DB_NAME", "customers")
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost,	port, dbUser, password, dbName) 
	
	return connStr
	
}