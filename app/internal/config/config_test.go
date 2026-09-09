package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var knownEnv = []string{
	"MCP_TRANSPORT", "MCP_NAME", "MCP_PATH",
	"SERVER_PORT", "SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT",
	"LOG_MODE",
	"X_API_KEY", "X_API_KEY_SECRET", "X_ACCESS_TOKEN", "X_ACCESS_TOKEN_SECRET",
	"MAX_IDLE_CONNS_PER_HOST", "X_DEBUG",
	"POSTS_MAX_QUERIES", "POSTS_MAX_USERNAMES", "POSTS_MAX_IDS", "POSTS_DEFAULT_RESULTS",
	"POSTS_MAX_RESULTS", "POSTS_DEFAULT_DAYS", "POSTS_MAX_DAYS",
	"POSTS_DEFAULT_TIMEOUT_MS", "POSTS_MIN_TIMEOUT_MS", "POSTS_MAX_TIMEOUT_MS",
	"USERS_MAX_USERNAMES", "USERS_MAX_IDS", "USERS_DEFAULT_RESULTS", "USERS_MAX_RESULTS",
	"USERS_DEFAULT_TIMEOUT_MS", "USERS_MIN_TIMEOUT_MS", "USERS_MAX_TIMEOUT_MS",
	"LISTS_DEFAULT_RESULTS", "LISTS_MAX_RESULTS",
	"LISTS_DEFAULT_TIMEOUT_MS", "LISTS_MIN_TIMEOUT_MS", "LISTS_MAX_TIMEOUT_MS",
	"MEDIA_ROOT", "MEDIA_CHUNK_BYTES", "MEDIA_POLL_INTERVAL_MS",
	"MEDIA_DEFAULT_TIMEOUT_MS", "MEDIA_MIN_TIMEOUT_MS", "MEDIA_MAX_TIMEOUT_MS",
}

var required = map[string]string{
	"X_API_KEY":             "key",
	"X_API_KEY_SECRET":      "key-secret",
	"X_ACCESS_TOKEN":        "token",
	"X_ACCESS_TOKEN_SECRET": "token-secret",
	"MEDIA_ROOT":            "/tmp",
}

func isolateEnv(t *testing.T) {
	t.Helper()

	for _, name := range knownEnv {
		t.Setenv(name, "")
		require.NoError(t, os.Unsetenv(name))
	}
}

func setRequired(t *testing.T) {
	t.Helper()

	for name, value := range required {
		t.Setenv(name, value)
	}
}

func TestLoadWithEveryRequiredVariable(t *testing.T) {
	isolateEnv(t)
	setRequired(t)

	cfg, err := Load("")

	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.NotNil(t, cfg.MCP)
	require.NotNil(t, cfg.Server)
	require.NotNil(t, cfg.Logger)
	require.NotNil(t, cfg.XClient)
	require.NotNil(t, cfg.Posts)
	require.NotNil(t, cfg.Users)
	require.NotNil(t, cfg.Lists)
	require.NotNil(t, cfg.Media)

	assert.Equal(t, "key", cfg.XClient.APIKey)
	assert.Equal(t, "key-secret", cfg.XClient.APIKeySecret)
	assert.Equal(t, "token", cfg.XClient.AccessToken)
	assert.Equal(t, "token-secret", cfg.XClient.AccessTokenSecret)
	assert.Equal(t, "/tmp", cfg.Media.Root)
}

func TestLoadMaterializesDefaults(t *testing.T) {
	isolateEnv(t)
	setRequired(t)

	cfg, err := Load("")
	require.NoError(t, err)

	assert.Equal(t, "stdio", cfg.MCP.Transport)
	assert.Equal(t, "mcp-x", cfg.MCP.Name)
	assert.Equal(t, "/mcp", cfg.MCP.Path)
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, 60*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, 60*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(t, "local", cfg.Logger.Mode)
	assert.Equal(t, 100, cfg.XClient.MaxIdleConnsPerHost)
	assert.False(t, cfg.XClient.Debug)
	assert.Equal(t, 5, cfg.Posts.MaxQueries)
	assert.Equal(t, 10, cfg.Posts.DefaultResults)
	assert.Equal(t, 100, cfg.Posts.MaxResults)
	assert.Equal(t, 7, cfg.Posts.MaxDays)
	assert.Equal(t, int64(15000), cfg.Posts.DefaultTimeoutMs)
	assert.Equal(t, 100, cfg.Users.MaxUsernames)
	assert.Equal(t, 1000, cfg.Users.MaxResults)
	assert.Equal(t, 25, cfg.Lists.DefaultResults)
	assert.Equal(t, 100, cfg.Lists.MaxResults)
	assert.Equal(t, 4194304, cfg.Media.ChunkBytes)
	assert.Equal(t, int64(600000), cfg.Media.MaxTimeoutMs)
}

func TestLoadOverridesDefaults(t *testing.T) {
	isolateEnv(t)
	setRequired(t)

	t.Setenv("MCP_TRANSPORT", "http")
	t.Setenv("MCP_PATH", "/x")
	t.Setenv("SERVER_PORT", "9999")
	t.Setenv("SERVER_READ_TIMEOUT", "5s")
	t.Setenv("LOG_MODE", "prod")
	t.Setenv("X_DEBUG", "true")
	t.Setenv("POSTS_MAX_RESULTS", "50")
	t.Setenv("MEDIA_CHUNK_BYTES", "1024")

	cfg, err := Load("")
	require.NoError(t, err)

	assert.Equal(t, "http", cfg.MCP.Transport)
	assert.Equal(t, "/x", cfg.MCP.Path)
	assert.Equal(t, "9999", cfg.Server.Port)
	assert.Equal(t, 5*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, "prod", cfg.Logger.Mode)
	assert.True(t, cfg.XClient.Debug)
	assert.Equal(t, 50, cfg.Posts.MaxResults)
	assert.Equal(t, 1024, cfg.Media.ChunkBytes)
}

func TestLoadNamesTheEnvironmentVariable(t *testing.T) {
	cases := []struct {
		name        string
		unset       string
		set         map[string]string
		wantMessage string
		wantMissing string
	}{
		{
			name:        "a missing api key",
			unset:       "X_API_KEY",
			wantMessage: "X_API_KEY",
			wantMissing: "APIKey",
		},
		{
			name:        "a missing access token secret",
			unset:       "X_ACCESS_TOKEN_SECRET",
			wantMessage: "X_ACCESS_TOKEN_SECRET",
			wantMissing: "AccessTokenSecret",
		},
		{
			name:        "a missing media root",
			unset:       "MEDIA_ROOT",
			wantMessage: "MEDIA_ROOT",
			wantMissing: "Root",
		},
		{
			name:        "an api key set to the empty string counts as missing",
			set:         map[string]string{"X_API_KEY": ""},
			wantMessage: "X_API_KEY",
			wantMissing: "APIKey",
		},
		{
			name:        "a transport outside the enum",
			set:         map[string]string{"MCP_TRANSPORT": "grpc"},
			wantMessage: "MCP_TRANSPORT",
			wantMissing: "Transport",
		},
		{
			name:        "a result limit above what the api accepts",
			set:         map[string]string{"POSTS_MAX_RESULTS": "500"},
			wantMessage: "POSTS_MAX_RESULTS",
			wantMissing: "MaxResults",
		},
		{
			name:        "a zero day window",
			set:         map[string]string{"POSTS_MAX_DAYS": "0"},
			wantMessage: "POSTS_MAX_DAYS",
			wantMissing: "MaxDays",
		},
		{
			name:        "a non-positive connection pool",
			set:         map[string]string{"MAX_IDLE_CONNS_PER_HOST": "0"},
			wantMessage: "MAX_IDLE_CONNS_PER_HOST",
			wantMissing: "MaxIdleConnsPerHost",
		},
		{
			name:        "a users page size above the api maximum",
			set:         map[string]string{"USERS_MAX_IDS": "101"},
			wantMessage: "USERS_MAX_IDS",
			wantMissing: "MaxIDs",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			isolateEnv(t)
			setRequired(t)

			if c.unset != "" {
				require.NoError(t, os.Unsetenv(c.unset))
			}

			for name, value := range c.set {
				t.Setenv(name, value)
			}

			cfg, err := Load("")

			require.Error(t, err)
			assert.Nil(t, cfg)
			assert.Contains(t, err.Error(), c.wantMessage)
			assert.NotContains(t, err.Error(), c.wantMissing)
		})
	}
}

func TestLoadRejectsUnparsableValues(t *testing.T) {
	cases := []struct {
		name     string
		variable string
		value    string
	}{
		{name: "a connection pool that is not a number", variable: "MAX_IDLE_CONNS_PER_HOST", value: "abc"},
		{name: "a timeout that is not a duration", variable: "SERVER_READ_TIMEOUT", value: "nonsense"},
		{name: "a debug flag that is not a boolean", variable: "X_DEBUG", value: "maybe"},
		{name: "a chunk size that is not a number", variable: "MEDIA_CHUNK_BYTES", value: "4mb"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			isolateEnv(t)
			setRequired(t)
			t.Setenv(c.variable, c.value)

			cfg, err := Load("")

			require.Error(t, err)
			assert.Nil(t, cfg)
			assert.Contains(t, err.Error(), c.variable)
		})
	}
}

func TestLoadFromEnvFile(t *testing.T) {
	t.Run("an env file supplies the required keys", func(t *testing.T) {
		isolateEnv(t)

		path := filepath.Join(t.TempDir(), ".env")
		require.NoError(t, os.WriteFile(path, []byte(
			"X_API_KEY=file-key\n"+
				"X_API_KEY_SECRET=file-key-secret\n"+
				"X_ACCESS_TOKEN=file-token\n"+
				"X_ACCESS_TOKEN_SECRET=file-token-secret\n"+
				"MEDIA_ROOT=/tmp\n"+
				"SERVER_PORT=7777\n"), 0o600))

		cfg, err := Load(path)

		require.NoError(t, err)
		assert.Equal(t, "file-key", cfg.XClient.APIKey)
		assert.Equal(t, "7777", cfg.Server.Port)
	})

	t.Run("the process environment wins over the file", func(t *testing.T) {
		isolateEnv(t)
		setRequired(t)

		path := filepath.Join(t.TempDir(), ".env")
		require.NoError(t, os.WriteFile(path, []byte("X_API_KEY=file-key\n"), 0o600))

		cfg, err := Load(path)

		require.NoError(t, err)
		assert.Equal(t, "key", cfg.XClient.APIKey)
	})

	t.Run("a path that does not exist falls back to the process environment", func(t *testing.T) {
		isolateEnv(t)
		setRequired(t)

		cfg, err := Load(filepath.Join(t.TempDir(), "missing.env"))

		require.NoError(t, err)
		assert.Equal(t, "key", cfg.XClient.APIKey)
	})

	t.Run("an empty path reads the process environment only", func(t *testing.T) {
		isolateEnv(t)
		setRequired(t)

		cfg, err := Load("")

		require.NoError(t, err)
		assert.Equal(t, "key", cfg.XClient.APIKey)
	})
}

func TestLoadRejectsAMediaRootThatIsNotAnAbsoluteDirectory(t *testing.T) {
	file := filepath.Join(t.TempDir(), "not-a-directory")
	require.NoError(t, os.WriteFile(file, []byte("x"), 0o600))

	cases := []struct {
		name        string
		root        string
		wantMessage string
	}{
		{
			name:        "a relative path",
			root:        "media",
			wantMessage: "must be an absolute path",
		},
		{
			name:        "a directory that does not exist",
			root:        filepath.Join(t.TempDir(), "missing"),
			wantMessage: "is not readable",
		},
		{
			name:        "a regular file",
			root:        file,
			wantMessage: "is not a directory",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			isolateEnv(t)
			setRequired(t)
			t.Setenv("MEDIA_ROOT", testCase.root)

			cfg, err := Load("")

			require.Error(t, err)
			assert.Nil(t, cfg)
			assert.Contains(t, err.Error(), "MEDIA_ROOT")
			assert.Contains(t, err.Error(), testCase.wantMessage)
		})
	}
}

func TestLoadAcceptsAnExistingAbsoluteMediaRoot(t *testing.T) {
	isolateEnv(t)
	setRequired(t)

	root := t.TempDir()
	t.Setenv("MEDIA_ROOT", root)

	cfg, err := Load("")

	require.NoError(t, err)
	assert.Equal(t, root, cfg.Media.Root)
}
