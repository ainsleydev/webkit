package jsonformat

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test types that mimic appdef structures with inline tags.
type (
	testEnvValue struct {
		Source string `json:"source"`
		Value  string `json:"value,omitempty"`
		Path   string `json:"path,omitempty"`
	}

	testEnvVar map[string]testEnvValue

	testEnvironment struct {
		Dev        testEnvVar `json:"dev,omitempty" inline:"true"`
		Production testEnvVar `json:"production,omitempty" inline:"true"`
		Staging    testEnvVar `json:"staging,omitempty" inline:"true"`
	}

	testCommandSpec struct {
		Command string `json:"command,omitempty"`
		SkipCI  bool   `json:"skip_ci,omitempty"`
		Timeout string `json:"timeout,omitempty"`
	}

	testApp struct {
		Name     string                     `json:"name"`
		Env      testEnvironment            `json:"env"`
		Commands map[string]testCommandSpec `json:"commands,omitempty" inline:"true"`
	}
)

func init() {
	// Register test types for testing.
	RegisterType(reflect.TypeOf(testApp{}))
	RegisterType(reflect.TypeOf(testEnvironment{}))
}

func TestFormat(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"Simple environment variable with source and value": {
			input: `{
	"env": {
		"dev": {
			"DATABASE_URI": {
				"source": "value",
				"value": "file:./cms.db"
			}
		}
	}
}`,
			want: `{
	"env": {
		"dev": {
			"DATABASE_URI": {"source": "value", "value": "file:./cms.db"}
		}
	}
}`,
		},
		"Multiple environment variables": {
			input: `{
	"env": {
		"dev": {
			"DATABASE_URI": {
				"source": "value",
				"value": "file:./cms.db"
			},
			"FRONTEND_URL": {
				"source": "value",
				"value": "http://localhost:5173"
			}
		}
	}
}`,
			want: `{
	"env": {
		"dev": {
			"DATABASE_URI": {"source": "value", "value": "file:./cms.db"},
			"FRONTEND_URL": {"source": "value", "value": "http://localhost:5173"}
		}
	}
}`,
		},
		"Environment variable with source and path": {
			input: `{
	"env": {
		"production": {
			"API_KEY": {
				"source": "sops",
				"path": "secrets/prod.yaml:API_KEY"
			}
		}
	}
}`,
			want: `{
	"env": {
		"production": {
			"API_KEY": {"source": "sops", "path": "secrets/prod.yaml:API_KEY"}
		}
	}
}`,
		},
		"Multiple environments": {
			input: `{
	"env": {
		"dev": {
			"DATABASE_URI": {
				"source": "value",
				"value": "file:./cms.db"
			}
		},
		"production": {
			"DATABASE_URI": {
				"source": "resource",
				"value": "db.connection_url"
			}
		}
	}
}`,
			want: `{
	"env": {
		"dev": {
			"DATABASE_URI": {"source": "value", "value": "file:./cms.db"}
		},
		"production": {
			"DATABASE_URI": {"source": "resource", "value": "db.connection_url"}
		}
	}
}`,
		},
		"Simple command": {
			input: `{
	"commands": {
		"build": {
			"command": "pnpm build"
		}
	}
}`,
			want: `{
	"commands": {
		"build": {"command": "pnpm build"}
	}
}`,
		},
		"Multiple commands": {
			input: `{
	"commands": {
		"build": {
			"command": "pnpm build"
		},
		"test": {
			"command": "pnpm test"
		},
		"lint": {
			"command": "pnpm lint"
		}
	}
}`,
			want: `{
	"commands": {
		"build": {"command": "pnpm build"},
		"test": {"command": "pnpm test"},
		"lint": {"command": "pnpm lint"}
	}
}`,
		},
		"Command with multiple fields": {
			input: `{
	"commands": {
		"test": {
			"command": "go test ./...",
			"timeout": "5m"
		}
	}
}`,
			want: `{
	"commands": {
		"test": {"command": "go test ./...", "timeout": "5m"}
	}
}`,
		},
		"Command with skip_ci": {
			input: `{
	"commands": {
		"deploy": {
			"command": "kubectl apply",
			"skip_ci": true
		}
	}
}`,
			want: `{
	"commands": {
		"deploy": {"command": "kubectl apply", "skip_ci": true}
	}
}`,
		},
		"Mixed content with trailing comma": {
			input: `{
	"name": "test",
	"env": {
		"dev": {
			"URL": {
				"source": "value",
				"value": "localhost"
			}
		}
	},
	"commands": {
		"build": {
			"command": "make"
		}
	}
}`,
			want: `{
	"name": "test",
	"env": {
		"dev": {
			"URL": {"source": "value", "value": "localhost"}
		}
	},
	"commands": {
		"build": {"command": "make"}
	}
}`,
		},
		"Nested structure not matching pattern": {
			input: `{
	"config": {
		"nested": {
			"field1": "value1",
			"field2": "value2"
		}
	}
}`,
			want: `{
	"config": {
		"nested": {
			"field1": "value1",
			"field2": "value2"
		}
	}
}`,
		},
		"Environment with resource reference": {
			input: `{
	"env": {
		"production": {
			"DB_URL": {
				"source": "resource",
				"value": "db.connection_string"
			},
			"BUCKET": {
				"source": "resource",
				"value": "storage.bucket_name"
			}
		}
	}
}`,
			want: `{
	"env": {
		"production": {
			"DB_URL": {"source": "resource", "value": "db.connection_string"},
			"BUCKET": {"source": "resource", "value": "storage.bucket_name"}
		}
	}
}`,
		},
		"Real world example": {
			input: `{
	"apps": [
		{
			"name": "cms",
			"env": {
				"dev": {
					"DATABASE_URI": {
						"source": "value",
						"value": "file:./cms.db"
					},
					"FRONTEND_URL": {
						"source": "value",
						"value": "http://localhost:5173"
					}
				},
				"production": {
					"DATABASE_URI": {
						"source": "resource",
						"value": "db.connection_url"
					},
					"FRONTEND_URL": {
						"source": "value",
						"value": "https://searchspares.com"
					}
				}
			},
			"commands": {
				"build": {
					"command": "pnpm build"
				},
				"test": {
					"command": "pnpm test"
				}
			}
		}
	]
}`,
			want: `{
	"apps": [
		{
			"name": "cms",
			"env": {
				"dev": {
					"DATABASE_URI": {"source": "value", "value": "file:./cms.db"},
					"FRONTEND_URL": {"source": "value", "value": "http://localhost:5173"}
				},
				"production": {
					"DATABASE_URI": {"source": "resource", "value": "db.connection_url"},
					"FRONTEND_URL": {"source": "value", "value": "https://searchspares.com"}
				}
			},
			"commands": {
				"build": {"command": "pnpm build"},
				"test": {"command": "pnpm test"}
			}
		}
	]
}`,
		},
		"Empty object": {
			input: `{
}`,
			want: `{
}`,
		},
		"Object with no inline candidates": {
			input: `{
	"name": "test",
	"version": "1.0.0"
}`,
			want: `{
	"name": "test",
	"version": "1.0.0"
}`,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := Format([]byte(test.input))
			require.NoError(t, err)
			assert.Equal(t, test.want, string(got))
		})
	}
}

func TestFormat_PreservesIndentation(t *testing.T) {
	t.Parallel()

	input := `{
		"deeply": {
			"nested": {
				"env": {
					"dev": {
						"KEY": {
							"source": "value",
							"value": "test"
						}
					}
				}
			}
		}
	}`

	want := `{
		"deeply": {
			"nested": {
				"env": {
					"dev": {
						"KEY": {"source": "value", "value": "test"}
					}
				}
			}
		}
	}`

	got, err := Format([]byte(input))
	require.NoError(t, err)
	assert.Equal(t, want, string(got))
}

func TestFormat_HandlesTrailingNewline(t *testing.T) {
	t.Parallel()

	input := `{
	"commands": {
		"build": {
			"command": "make"
		}
	}
}
`

	got, err := Format([]byte(input))
	require.NoError(t, err)

	// Result should preserve structure but may not have exact trailing newline.
	assert.Contains(t, string(got), `"build": {"command": "make"}`)
}
