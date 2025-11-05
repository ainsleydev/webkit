package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/infra"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestTransformOutputs(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input       infra.OutputResult
		environment env.Environment
		want        TerraformOutputProvider
	}{
		"Empty Result": {
			input: infra.OutputResult{
				Resources: map[string]map[string]any{},
			},
			environment: env.Production,
			want:        TerraformOutputProvider{},
		},
		"Single Resource Single Output": {
			input: infra.OutputResult{
				Resources: map[string]map[string]any{
					"postgres": {
						"connection_url": "postgresql://user:pass@host:5432/db",
					},
				},
			},
			environment: env.Production,
			want: TerraformOutputProvider{
				OutputKey{
					Environment:  env.Production,
					ResourceName: "postgres",
					OutputName:   "connection_url",
				}: "postgresql://user:pass@host:5432/db",
			},
		},
		"Single Resource Multiple Outputs": {
			input: infra.OutputResult{
				Resources: map[string]map[string]any{
					"postgres": {
						"connection_url": "postgresql://user:pass@host:5432/db",
						"host":           "db.example.com",
						"port":           5432,
					},
				},
			},
			environment: env.Staging,
			want: TerraformOutputProvider{
				OutputKey{
					Environment:  env.Staging,
					ResourceName: "postgres",
					OutputName:   "connection_url",
				}: "postgresql://user:pass@host:5432/db",
				OutputKey{
					Environment:  env.Staging,
					ResourceName: "postgres",
					OutputName:   "host",
				}: "db.example.com",
				OutputKey{
					Environment:  env.Staging,
					ResourceName: "postgres",
					OutputName:   "port",
				}: 5432,
			},
		},
		"Multiple Resources Multiple Outputs": {
			input: infra.OutputResult{
				Resources: map[string]map[string]any{
					"postgres": {
						"connection_url": "postgresql://user:pass@host:5432/db",
						"host":           "db.example.com",
					},
					"s3": {
						"bucket_name": "my-bucket",
						"region":      "us-east-1",
					},
				},
			},
			environment: env.Development,
			want: TerraformOutputProvider{
				OutputKey{
					Environment:  env.Development,
					ResourceName: "postgres",
					OutputName:   "connection_url",
				}: "postgresql://user:pass@host:5432/db",
				OutputKey{
					Environment:  env.Development,
					ResourceName: "postgres",
					OutputName:   "host",
				}: "db.example.com",
				OutputKey{
					Environment:  env.Development,
					ResourceName: "s3",
					OutputName:   "bucket_name",
				}: "my-bucket",
				OutputKey{
					Environment:  env.Development,
					ResourceName: "s3",
					OutputName:   "region",
				}: "us-east-1",
			},
		},
		"Various Value Types": {
			input: infra.OutputResult{
				Resources: map[string]map[string]any{
					"resource": {
						"string_val": "value",
						"int_val":    42,
						"bool_val":   true,
						"float_val":  3.14,
						"nil_val":    nil,
					},
				},
			},
			environment: env.Production,
			want: TerraformOutputProvider{
				OutputKey{
					Environment:  env.Production,
					ResourceName: "resource",
					OutputName:   "string_val",
				}: "value",
				OutputKey{
					Environment:  env.Production,
					ResourceName: "resource",
					OutputName:   "int_val",
				}: 42,
				OutputKey{
					Environment:  env.Production,
					ResourceName: "resource",
					OutputName:   "bool_val",
				}: true,
				OutputKey{
					Environment:  env.Production,
					ResourceName: "resource",
					OutputName:   "float_val",
				}: 3.14,
				OutputKey{
					Environment:  env.Production,
					ResourceName: "resource",
					OutputName:   "nil_val",
				}: nil,
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := TransformOutputs(test.input, test.environment)
			assert.Equal(t, test.want, got)
		})
	}
}
