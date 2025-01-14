package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
)

func getVaultClient() (*api.Client, error) {
	// Charge les variables d'environnement
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("erreur de chargement .env: %w", err)
	}

	token := os.Getenv("VAULT_TOKEN")
	addr := os.Getenv("VAULT_ADDR")

	if token == "" || addr == "" {
		return nil, fmt.Errorf("VAULT_TOKEN et VAULT_ADDR sont requis")
	}

	// Configure le client
	config := api.DefaultConfig()
	config.Address = addr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("erreur création client: %w", err)
	}

	client.SetToken(token)
	return client, nil
}

func getSecret(client *api.Client, path string, key string) (interface{}, error) {
	// Crée un contexte avec timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Lecture du secret - correction de Read en Get
	kvv2 := client.KVv2("secret")
	secret, err := kvv2.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture secret: %w", err)
	}

	value, exists := secret.Data[key]
	if !exists {
		return nil, fmt.Errorf("clé '%s' non trouvée", key)
	}

	return value, nil
}

func main() {
	// Initialise le client
	client, err := getVaultClient()
	if err != nil {
		log.Fatalf("Erreur d'initialisation: %v", err)
	}

	// Lit le secret
	value, err := getSecret(client, "test/secret", "tagada")
	if err != nil {
		log.Fatalf("Erreur de lecture: %v", err)
	}

	fmt.Printf("Valeur de tagada: %v\n", value)
}
