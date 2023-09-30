// package service

// import (
// 	"context"
// 	"reflect"
// 	"testing"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/service/s3"
// 	"github.com/ccallazans/filedrop/internal/config"
// 	"github.com/ccallazans/filedrop/internal/domain/repository"
// )

// func TestFileService_DownloadFile(t *testing.T) {
// 	container := config.NewTestContainerStruct()
// 	ctx := context.Background()

// 	implFileStore := repository.NewPostgresFileStore(container.DB)
// 	implUserStore := repository.NewPostgresUserStore(container.DB)
// 	implS3Client := config.NewS3Client(aws.NewConfig())

// 	type fields struct {
// 		fileStore repository.FileStore
// 		userStore repository.UserStore
// 		s3Client  *s3.Client
// 	}
// 	type args struct {
// 		ctx      context.Context
// 		hash     string
// 		password string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    *s3.GetObjectOutput
// 		want1   string
// 		wantErr bool
// 	}{
// 		{
// 			name:    "Should Error hash file do not exist",
// 			fields:  fields{fileStore: implFileStore, userStore: implUserStore, s3Client: implS3Client},
// 			args:    args{ctx: ctx, hash: "abcde", password: ""},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &FileService{
// 				fileStore: tt.fields.fileStore,
// 				userStore: tt.fields.userStore,
// 				s3Client:  tt.fields.s3Client,
// 			}
// 			got, got1, err := s.DownloadFile(tt.args.ctx, tt.args.hash, tt.args.password)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FileService.DownloadFile() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if ((err != nil) == tt.wantErr) && (err.Error() != "file not found") {
// 				t.Errorf("FileService.DownloadFile() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("FileService.DownloadFile() got = %v, want %v", got, tt.want)
// 			}
// 			if got1 != tt.want1 {
// 				t.Errorf("FileService.DownloadFile() got1 = %v, want %v", got1, tt.want1)
// 			}
// 		})
// 	}
// }
