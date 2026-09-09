package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	xclient "github.com/Role1776/mcp-x/app/internal/adapter/x/client"
	"github.com/Role1776/mcp-x/app/internal/pkg/logger"
	"github.com/Role1776/mcp-x/app/internal/pkg/mcpserver"
	"github.com/Role1776/mcp-x/app/internal/pkg/server"
	pkgvalidator "github.com/Role1776/mcp-x/app/internal/pkg/validator"
	listsUseCase "github.com/Role1776/mcp-x/app/internal/usecase/x/lists"
	mediaUseCase "github.com/Role1776/mcp-x/app/internal/usecase/x/media"
	postsUseCase "github.com/Role1776/mcp-x/app/internal/usecase/x/posts"
	usersUseCase "github.com/Role1776/mcp-x/app/internal/usecase/x/users"
)

const keysHint = "get the four keys at https://developer.x.com -> your project -> Keys and tokens; " +
	"set the app permissions to Read and Write BEFORE generating the Access Token and Secret"

type Config struct {
	MCP     *mcpserver.Config
	Server  *server.Config
	Logger  *logger.Config
	XClient *xclient.Config
	Posts   *postsUseCase.Config
	Users   *usersUseCase.Config
	Lists   *listsUseCase.Config
	Media   *mediaUseCase.Config
}

func Load(envPath string) (*Config, error) {
	cfg := &Config{
		MCP:     &mcpserver.Config{},
		Server:  &server.Config{},
		Logger:  &logger.Config{},
		XClient: &xclient.Config{},
		Posts:   &postsUseCase.Config{},
		Users:   &usersUseCase.Config{},
		Lists:   &listsUseCase.Config{},
		Media:   &mediaUseCase.Config{},
	}

	if err := loadEnvFile(envPath); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg); err != nil {
		return nil, describeParse(cfg, err)
	}

	if err := pkgvalidator.Validate(cfg); err != nil {
		return nil, describe(cfg, err)
	}

	return cfg, nil
}

// describeParse turns the struct field names caarlos0/env reports back into the
// environment variable names the operator actually sets, so a bad value names
// X_DEBUG and not Debug.
func describeParse(cfg *Config, err error) error {
	parseErrors := parseErrorsOf(err)
	if len(parseErrors) == 0 {
		return fmt.Errorf("parse env: %w", err)
	}

	names := envNames(cfg)

	problems := make([]string, 0, len(parseErrors))
	for _, parseErr := range parseErrors {
		problems = append(problems, fmt.Sprintf("%s: %v", envNameOfField(names, parseErr.Name), parseErr.Err))
	}

	return fmt.Errorf("invalid environment values: %s", strings.Join(problems, ", "))
}

func parseErrorsOf(err error) []env.ParseError {
	var aggregate env.AggregateError
	if !errors.As(err, &aggregate) {
		var single env.ParseError
		if errors.As(err, &single) {
			return []env.ParseError{single}
		}

		return nil
	}

	result := make([]env.ParseError, 0, len(aggregate.Errors))
	for _, inner := range aggregate.Errors {
		var parseErr env.ParseError
		if errors.As(inner, &parseErr) {
			result = append(result, parseErr)
		}
	}

	return result
}

// envNameOfField maps a struct field name onto an environment variable name.
// env.ParseError carries the bare field name and not its path, so a name shared
// by several sections is narrowed down to the variables actually set.
func envNameOfField(names map[string]string, field string) string {
	candidates := make([]string, 0, 1)
	for path, name := range names {
		if strings.HasSuffix(path, "."+field) {
			candidates = append(candidates, name)
		}
	}

	if len(candidates) == 0 {
		return field
	}

	present := make([]string, 0, len(candidates))
	for _, name := range candidates {
		if _, ok := os.LookupEnv(name); ok {
			present = append(present, name)
		}
	}

	if len(present) > 0 {
		candidates = present
	}

	sort.Strings(candidates)

	return strings.Join(candidates, "/")
}

// describe turns validator field paths back into the environment variable names
// the operator actually sets, so a missing key names X_API_KEY and not APIKey.
func describe(cfg *Config, err error) error {
	var invalid validator.ValidationErrors
	if !errors.As(err, &invalid) {
		return fmt.Errorf("validate config: %w", err)
	}

	names := envNames(cfg)

	missing := make([]string, 0, len(invalid))
	other := make([]string, 0, len(invalid))

	for _, fieldErr := range invalid {
		name := names[fieldErr.StructNamespace()]
		if name == "" {
			name = fieldErr.Namespace()
		}

		if fieldErr.Tag() == "required" {
			missing = append(missing, name)

			continue
		}

		other = append(other, fmt.Sprintf("%s (%s)", name, fieldErr.Tag()))
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s\n%s",
			strings.Join(missing, ", "), keysHint)
	}

	return fmt.Errorf("invalid configuration: %s", strings.Join(other, ", "))
}

func envNames(cfg *Config) map[string]string {
	names := make(map[string]string)
	collect(reflect.ValueOf(cfg).Elem(), reflect.TypeOf(cfg).Elem().Name(), names)

	return names
}

func collect(value reflect.Value, path string, names map[string]string) {
	valueType := value.Type()

	for i := range valueType.NumField() {
		field := valueType.Field(i)
		fieldPath := path + "." + field.Name

		if name, ok := field.Tag.Lookup("env"); ok {
			names[fieldPath] = name

			continue
		}

		fieldValue := value.Field(i)
		if fieldValue.Kind() == reflect.Pointer {
			if fieldValue.IsNil() {
				continue
			}
			fieldValue = fieldValue.Elem()
		}

		if fieldValue.Kind() == reflect.Struct {
			collect(fieldValue, fieldPath, names)
		}
	}
}

func loadEnvFile(envPath string) error {
	var err error
	if envPath != "" {
		err = godotenv.Load(envPath)
	}

	switch {
	case err == nil:
		return nil
	case errors.Is(err, os.ErrNotExist):
		fmt.Fprintf(os.Stderr, "%s no env file found at %q, using the ambient environment\n",
			time.Now().Format(time.RFC3339), envPath)

		return nil
	default:
		return fmt.Errorf("load env file %s: %w", envPath, err)
	}
}
