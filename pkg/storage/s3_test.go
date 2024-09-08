package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockS3Client struct {
	PutObjectFunc     func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	DeleteObjectFunc  func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	ListObjectsV2Func func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	GetObjectFunc     func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	HeadObjectFunc    func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

func (m *mockS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return m.PutObjectFunc(ctx, params, optFns...)
}

func (m *mockS3Client) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	return m.DeleteObjectFunc(ctx, params, optFns...)
}

func (m *mockS3Client) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	return m.ListObjectsV2Func(ctx, params, optFns...)
}

func (m *mockS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return m.GetObjectFunc(ctx, params, optFns...)
}

func (m *mockS3Client) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	return m.HeadObjectFunc(ctx, params, optFns...)
}

func setupS3Storage(t *testing.T) *S3 {
	t.Helper()
	return &S3{
		client: &mockS3Client{
			PutObjectFunc: func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				return nil, nil
			},
			DeleteObjectFunc: func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return nil, nil
			},
			ListObjectsV2Func: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
				return nil, nil
			},
			GetObjectFunc: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				return nil, nil
			},
			HeadObjectFunc: func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return nil, nil
			},
		},
		bucket: "test-bucket",
	}
}

func setupS3StorageForPersistenceTest(t *testing.T) *S3 {
	t.Helper()

	files := make(map[string][]byte)

	mockClient := &mockS3Client{
		PutObjectFunc: func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			content, err := io.ReadAll(params.Body)
			if err != nil {
				return nil, err
			}
			files[*params.Key] = content
			return &s3.PutObjectOutput{}, nil
		},
		GetObjectFunc: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			content, ok := files[*params.Key]
			if !ok {
				return nil, &types.NoSuchKey{}
			}
			return &s3.GetObjectOutput{
				Body: io.NopCloser(bytes.NewReader(content)),
			}, nil
		},
		DeleteObjectFunc: func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
			delete(files, *params.Key)
			return &s3.DeleteObjectOutput{}, nil
		},
		ListObjectsV2Func: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			var contents []types.Object
			for key := range files {
				contents = append(contents, types.Object{Key: aws.String(key)})
			}
			return &s3.ListObjectsV2Output{
				Contents: contents,
			}, nil
		},
		HeadObjectFunc: func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
			content, ok := files[*params.Key]
			if !ok {
				return nil, &types.NoSuchKey{}
			}
			return &s3.HeadObjectOutput{
				ContentLength: aws.Int64(int64(len(content))),
				LastModified:  aws.Time(time.Now()),
			}, nil
		},
	}

	return &S3{
		client: mockClient,
		bucket: "test-bucket",
	}
}

func TestS3Config_Validate(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   S3Config
		wantErr bool
	}{
		"Valid config": {
			input: S3Config{
				Bucket:          "test-bucket",
				Region:          "us-west-2",
				AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
				SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			},
			wantErr: false,
		},
		"Missing bucket": {
			input: S3Config{
				Region:          "us-west-2",
				AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
				SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			},
			wantErr: true,
		},
		"Missing region": {
			input: S3Config{
				Bucket:          "test-bucket",
				AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
				SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			},
			wantErr: true,
		},
		"Missing access key ID": {
			input: S3Config{
				Bucket:          "test-bucket",
				Region:          "us-west-2",
				SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			},
			wantErr: true,
		},
		"Missing secret access key": {
			input: S3Config{
				Bucket:      "test-bucket",
				Region:      "us-west-2",
				AccessKeyID: "AKIAIOSFODNN7EXAMPLE",
			},
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := test.input.Validate()
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestNewS3Storage(t *testing.T) {
	t.Parallel()

	var loadDefaultConfig func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error)
	var newFromConfig func(cfg aws.Config, optFns ...func(*s3.Options)) s3ClientAPI

	originalLoadDefaultConfig := loadDefaultConfig
	defer func() { loadDefaultConfig = originalLoadDefaultConfig }()
	loadDefaultConfig = func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
		return aws.Config{}, nil
	}

	originalNewFromConfig := newFromConfig
	defer func() { newFromConfig = originalNewFromConfig }()
	newFromConfig = func(cfg aws.Config, optFns ...func(*s3.Options)) s3ClientAPI {
		return &mockS3Client{}
	}

	tt := map[string]struct {
		input S3Config
		want  bool
	}{
		"Valid": {
			input: S3Config{
				Bucket:          "test-bucket",
				Region:          "us-west-2",
				AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
				SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			},
			want: true,
		},
		"Invalid": {
			input: S3Config{
				Bucket: "test-bucket",
				Region: "us-west-2",
			},
			want: false,
		},
		"With Custom Endpoint": {
			input: S3Config{
				Bucket:          "test-bucket",
				Region:          "us-west-2",
				AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
				SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				Endpoint:        "http://localhost:9000",
			},
			want: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s3Storage, err := NewS3Storage(context.Background(), test.input)
			assert.Equal(t, test.want, err == nil && s3Storage != nil)
		})
	}
}

func TestS3_Upload(t *testing.T) {
	t.Parallel()

	t.Run("Successful Upload", func(t *testing.T) {
		t.Parallel()
		s := setupS3Storage(t)
		s.client.(*mockS3Client).PutObjectFunc = func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			assert.Equal(t, "test-bucket", *params.Bucket)
			assert.Equal(t, "test.txt", *params.Key)
			content, _ := io.ReadAll(params.Body)
			assert.Equal(t, []byte("test content"), content)
			return &s3.PutObjectOutput{}, nil
		}

		content := bytes.NewBufferString("test content")
		err := s.Upload(context.Background(), "test.txt", content)
		require.NoError(t, err)
	})

	t.Run("Upload Error", func(t *testing.T) {
		t.Parallel()
		s := setupS3Storage(t)
		s.client.(*mockS3Client).PutObjectFunc = func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			return nil, errors.New("upload error")
		}

		content := bytes.NewBufferString("error content")
		err := s.Upload(context.Background(), "error.txt", content)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "upload error")
	})
}

func TestS3_Delete(t *testing.T) {
	t.Parallel()

	t.Run("Successful Delete", func(t *testing.T) {
		t.Parallel()
		s := setupS3Storage(t)
		s.client.(*mockS3Client).DeleteObjectFunc = func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
			assert.Equal(t, "test-bucket", *params.Bucket)
			assert.Equal(t, "delete_me.txt", *params.Key)
			return &s3.DeleteObjectOutput{}, nil
		}

		err := s.Delete(context.Background(), "delete_me.txt")
		require.NoError(t, err)
	})

	t.Run("Delete Error", func(t *testing.T) {
		t.Parallel()
		s := setupS3Storage(t)
		s.client.(*mockS3Client).DeleteObjectFunc = func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
			return nil, errors.New("delete error")
		}

		err := s.Delete(context.Background(), "non_existent.txt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete error")
	})
}

func TestS3_List(t *testing.T) {
	t.Parallel()

	t.Run("Successful List", func(t *testing.T) {
		t.Parallel()
		s := setupS3Storage(t)
		s.client.(*mockS3Client).ListObjectsV2Func = func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			assert.Equal(t, "test-bucket", *params.Bucket)
			assert.Equal(t, "", *params.Prefix)
			return &s3.ListObjectsV2Output{
				Contents: []types.Object{
					{Key: new(string)},
					{Key: new(string)},
				},
			}, nil
		}

		files, err := s.List(context.Background(), "")
		require.NoError(t, err)
		assert.Len(t, files, 2)
	})

	t.Run("List Error", func(t *testing.T) {
		t.Parallel()
		s := setupS3Storage(t)
		s.client.(*mockS3Client).ListObjectsV2Func = func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			return nil, errors.New("list error")
		}

		files, err := s.List(context.Background(), "")
		assert.Error(t, err)
		assert.Nil(t, files)
		assert.Contains(t, err.Error(), "list error")
	})
}

func TestS3_Download(t *testing.T) {
	t.Parallel()

	t.Run("Successful Download", func(t *testing.T) {
		t.Parallel()
		s := setupS3Storage(t)
		s.client.(*mockS3Client).GetObjectFunc = func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			assert.Equal(t, "test-bucket", *params.Bucket)
			assert.Equal(t, "download.txt", *params.Key)
			return &s3.GetObjectOutput{
				Body: io.NopCloser(bytes.NewReader([]byte("test content"))),
			}, nil
		}

		reader, err := s.Download(context.Background(), "download.txt")
		require.NoError(t, err)
		defer reader.Close()

		content, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, []byte("test content"), content)
	})

	t.Run("Download Error", func(t *testing.T) {
		t.Parallel()
		s := setupS3Storage(t)
		s.client.(*mockS3Client).GetObjectFunc = func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			return nil, errors.New("download error")
		}

		_, err := s.Download(context.Background(), "non_existent.txt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "download error")
	})
}

func TestS3_Exists(t *testing.T) {
	t.Parallel()

	t.Run("File Exists", func(t *testing.T) {
		t.Parallel()
		s := setupS3Storage(t)
		s.client.(*mockS3Client).HeadObjectFunc = func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
			assert.Equal(t, "test-bucket", *params.Bucket)
			assert.Equal(t, "exists.txt", *params.Key)
			return &s3.HeadObjectOutput{}, nil
		}

		exists, err := s.Exists(context.Background(), "exists.txt")
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("File Does Not Exist", func(t *testing.T) {
		t.Parallel()
		s := setupS3Storage(t)
		s.client.(*mockS3Client).HeadObjectFunc = func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
			return nil, errors.New("not found")
		}

		exists, err := s.Exists(context.Background(), "non_existent.txt")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestS3_Stat(t *testing.T) {
	t.Parallel()

	t.Run("Stat File", func(t *testing.T) {
		t.Parallel()
		s := setupS3Storage(t)
		s.client.(*mockS3Client).HeadObjectFunc = func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
			assert.Equal(t, "test-bucket", *params.Bucket)
			assert.Equal(t, "stat.txt", *params.Key)
			return &s3.HeadObjectOutput{
				ContentLength: new(int64),
				LastModified:  new(time.Time),
				ContentType:   new(string),
			}, nil
		}

		info, err := s.Stat(context.Background(), "stat.txt")
		require.NoError(t, err)
		assert.NotNil(t, info.Size)
		assert.NotNil(t, info.LastModified)
		assert.NotNil(t, info.ContentType)
		assert.False(t, info.IsDir)
	})

	t.Run("Stat Error", func(t *testing.T) {
		t.Parallel()
		s := setupS3Storage(t)
		s.client.(*mockS3Client).HeadObjectFunc = func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
			return nil, errors.New("stat error")
		}

		_, err := s.Stat(context.Background(), "non_existent.txt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "stat error")
	})
}
