package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	path := filepath.Join("testdata", "config_test.yml")
	cfg, err := read(path)
	if err != nil {
		t.Fail()
	}
	expected := Config{Node: Node{Host: "https://eth_node_test", Confirmations: 2, StartBlockHeight: 70800},
		Db: DB{Name: "block_db_test", Host: "localhost:27017"},
	}
	assert.Equal(t, expected, *cfg)
}

func TestReadWithEnvVariables(t *testing.T) {
	_ = os.Setenv("NODE_HOST", "https://eth_node_test_from_env")
	_ = os.Setenv("NODE_STARTBLOCK", "7239023")
	_ = os.Setenv("NODE_CONFIRMATIONS", "6")
	_ = os.Setenv("DB_NAME", "block_db_test_env")
	path := filepath.Join("testdata", "config_test.yml")
	cfg, err := read(path)
	if err != nil {
		t.Fail()
	}
	expected := Config{Node: Node{Host: "https://eth_node_test_from_env", Confirmations: 6, StartBlockHeight: 7239023},
		Db: DB{Name: "block_db_test_env", Host: "localhost:27017"},
	}
	assert.Equal(t, expected, *cfg)
}
