package serverapp

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
)

// Run loads configuration, initializes infrastructure, and starts the HTTP server.
// The command package remains a thin compatibility entry point for Air, Docker, and CI.
func Run() error {
	configPath := resolveConfigPath()
	log.Printf("[Main] Final config path: %s", configPath)

	if _, err := config.LoadConfig(configPath); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	log.Println("[Main] Configuration loaded successfully")

	config.RegisterReloadHandler("database", func() {
		if err := core.InitDB(); err != nil {
			log.Printf("Failed to reload database configuration: %v", err)
		}
	})
	config.RegisterReloadHandler("server", func() {
		log.Println("Server configuration reloaded")
	})
	config.EnableHotReload()

	if err := core.InitDB(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	r, err := core.InitServer()
	if err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}

	if err := core.RunServer(r); err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}
	return nil
}

func resolveConfigPath() string {
	wd, _ := os.Getwd()
	log.Printf("[Main] Working directory: %s", wd)
	log.Printf("[Main] CONFIG_FILE env: %s", os.Getenv("CONFIG_FILE"))
	log.Printf("[Main] CONFIG_PATH env: %s", os.Getenv("CONFIG_PATH"))

	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}
	if configPath == "" {
		if _, err := os.Stat("/app/configs/config.yaml"); err == nil {
			configPath = "/app/configs/config.yaml"
			log.Printf("[Main] Found Docker config at: %s", configPath)
		} else if _, err := os.Stat("/app/config/config.yaml"); err == nil {
			configPath = "/app/config/config.yaml"
			log.Printf("[Main] Found Docker config at: %s", configPath)
		} else {
			configPath = "."
			log.Printf("[Main] Using current directory for config: %s", configPath)
		}
	}
	return filepath.Clean(configPath)
}
