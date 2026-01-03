package services

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService(url string) *CloudinaryService {
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		return nil
	}
	return &CloudinaryService{cld: cld}
}
func (s *CloudinaryService) UploadAvatar(ctx context.Context, file multipart.File, userID string) (string, error) {
	resp, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:       userID,
		Folder:         "usof/avatars/",
		Overwrite:      api.Bool(true),
		Transformation: "w_400,h_400,c_fill,g_face",
	})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}

func (s *CloudinaryService) DeleteAvatar(ctx context.Context, userID string) error {
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: userID})
	return err
}
