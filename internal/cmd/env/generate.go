package env

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/pkg/env"
)

var GenerateCmd = &cli.Command{
	Name:        "generate",
	Usage:       "Generate env file for a specific app and environment",
	Description: "Generates a .env file for a specific app and environment with a custom output path. Useful for VM deployments.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "app",
			Usage:    "Name of the app to generate env file for",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "environment",
			Aliases:  []string{"env"},
			Usage:    "Target environment (development, staging, production)",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "output",
			Usage: "Output path for the .env file (defaults to {app.Path}/.env.{environment})",
		},
	},
	Action: cmdtools.Wrap(Generate),
}

// Generate creates a .env file for a specific app and environment
func Generate(ctx context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()
	printer := input.Printer()

	// Get flags
	appName := input.Command.String("app")
	environmentStr := input.Command.String("environment")
	outputPath := input.Command.String("output")

	// Parse environment
	environment, err := env.Parse(environmentStr)
	if err != nil {
		return fmt.Errorf("invalid environment: %w", err)
	}

	// Find the app
	var targetApp *appdef.App
	for _, app := range appDef.Apps {
		if app.Name == appName {
			targetApp = &app
			break
		}
	}

	if targetApp == nil {
		return fmt.Errorf("app '%s' not found in app.json", appName)
	}

	// Resolve secrets
	err = secrets.Resolve(ctx, appDef, secrets.ResolveConfig{
		SOPSClient: input.SOPSClient(),
		BaseDir:    input.BaseDir,
	})
	if err != nil {
		return err
	}

	// Merge shared and app-specific environments
	mergedApp := targetApp.MergeEnvironments(appDef.Shared.Env)

	// Get environment-specific vars
	var vars appdef.EnvVar
	switch environment {
	case env.Development:
		vars = mergedApp.Dev
	case env.Staging:
		vars = mergedApp.Staging
	case env.Production:
		vars = mergedApp.Production
	default:
		return fmt.Errorf("unsupported environment: %s", environment)
	}

	// Skip if no environment variables
	if len(vars) == 0 {
		printer.Warn(fmt.Sprintf("No environment variables defined for app '%s' in environment '%s'", appName, environment))
		return nil
	}

	// Determine output path
	if outputPath == "" {
		// Default to app directory
		fileName := ".env"
		if environment != env.Development {
			fileName = fmt.Sprintf(".env.%s", environment)
		}
		outputPath = filepath.Join(targetApp.Path, fileName)
	}

	// Convert vars to string map
	envMap := make(map[string]string)
	for k, v := range vars {
		envMap[k] = cast.ToString(v.Value)
	}

	// Marshal to dotenv format
	buf, err := godotenv.Marshal(envMap)
	if err != nil {
		return fmt.Errorf("failed to marshal env vars: %w", err)
	}

	// Ensure directory exists
	outputDir := filepath.Dir(outputPath)
	err = input.FS.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write file with webkit notice
	opts := []scaffold.Option{
		scaffold.WithTracking(manifest.SourceProject()),
	}

	err = input.Generator().Bytes(outputPath, []byte(buf), opts...)
	if err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}

	printer.Success(fmt.Sprintf("Generated env file for app '%s' (%s) at: %s", appName, environment, outputPath))

	return nil
}
