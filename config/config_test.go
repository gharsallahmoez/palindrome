package config_test

import (
	"github.com/gharsallahmoez/palindrome/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	// Test default config.
	t.Run("default config", func(t *testing.T) {
		conf := config.New()
		require.Equal(t, "localhost", conf.Server.Host)
		require.Equal(t, "8080", conf.Server.Port)
		require.Equal(t, "in-memory", conf.Database.Type)
	})

	// Test with custom config.
	t.Run("server config set from env", func(t *testing.T) {
		t.Setenv("SERVER_HOST", "1.1.1.1")
		t.Setenv("SERVER_PORT", "8080")
		t.Setenv("DATABASE_TYPE", "POSTGRES")
		conf := config.New()
		require.Equal(t, "1.1.1.1", conf.Server.Host)
		require.Equal(t, "8080", conf.Server.Port)
		require.Equal(t, "POSTGRES", conf.Database.Type)
	})
}

func TestInitLogger(t *testing.T) {
	// Call the InitLogger function
	config.InitLogger()

	assert.IsType(t, &logrus.JSONFormatter{}, logrus.StandardLogger().Formatter, "Formatter should be a JSONFormatter")
	jsonFormatter := logrus.StandardLogger().Formatter.(*logrus.JSONFormatter)
	assert.True(t, jsonFormatter.PrettyPrint, "PrettyPrint should be true")
	assert.True(t, jsonFormatter.DisableHTMLEscape, "DisableHTMLEscape should be true")
	assert.Equal(t, time.RFC3339, jsonFormatter.TimestampFormat, "TimestampFormat should be RFC3339")
	assert.True(t, logrus.StandardLogger().ReportCaller, "ReportCaller should be true")
	assert.Equal(t, os.Stdout, logrus.StandardLogger().Out, "Output should be os.Stdout")
	assert.Equal(t, logrus.DebugLevel, logrus.StandardLogger().Level, "Level should be DebugLevel")
}
