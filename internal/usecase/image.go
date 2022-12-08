package usecase

import (
	"bytes"
	"context"
	"fmt"
	"grpc/app/internal/entity"
	"grpc/app/internal/repository"
	"grpc/app/proto/imagestorage"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"regexp"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ImageUsecase struct {
	log *logrus.Logger
	repository.ImageRepository
	imagestorage.UnimplementedImageStorageServer
}

func NewImageUsecase(repo repository.ImageRepository, log *logrus.Logger) imagestorage.ImageStorageServer {
	return &ImageUsecase{
		ImageRepository: repo,
	}
}

var (
	imagePath = os.Getenv("IMAGES_DIR")
)

func (u *ImageUsecase) SaveImage(ctx context.Context, in *imagestorage.Image) (*imagestorage.SaveImageResponse, error) {
	reader := bytes.NewReader(in.File)
	img, _, err := image.Decode(reader)
	if err != nil {
		u.log.Errorln(err)
		return nil, err
	}

	fileName := fmt.Sprintf("%s.%s", randStringRunes(10), in.FileType)
	filePath := fmt.Sprintf("%s%s", imagePath, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		u.log.Errorln(err)
		return nil, err
	}
	defer out.Close()

	var opt jpeg.Options
	opt.Quality = 1

	err = jpeg.Encode(out, img, &opt)
	if err != nil {
		u.log.Errorln(err)
		return nil, err
	}

	err = u.ImageRepository.SaveImage(fileName)
	if err != nil {
		u.log.Errorln(err)
		return nil, err
	}

	return nil, nil
}

func (u *ImageUsecase) LoadImageList(ctx context.Context, in *imagestorage.LoadImageListRequest) (*imagestorage.ImageInfoList, error) {
	data, err := u.ImageRepository.LoadImageList()
	if err != nil {
		u.log.Errorln(err)
		return nil, err
	}
	pbdata := make([]*imagestorage.ImageInfo, 0)
	for _, d := range data {
		pbdata = append(pbdata, &imagestorage.ImageInfo{
			Filename:    d.FileName,
			CreatedDate: timestamppb.New(d.CreatedDate),
			UpdatedDate: timestamppb.New(d.UpdatedDate),
		})
	}

	return &imagestorage.ImageInfoList{
		Images: pbdata,
	}, nil
}

func (u *ImageUsecase) FindImage(ctx context.Context, in *imagestorage.FindImageRequest) (*imagestorage.Image, error) {
	fileType := fileTypeRE.FindAllString(in.Filename, -1)
	if len(fileType) == 0 {
		return nil, entity.ErrIncorrectFileType
	}
	filePath := fmt.Sprintf("%s%s", imagePath, in.Filename)
	f, err := os.Open(filePath)
	if err != nil {
		u.log.Errorln(err)
		return &imagestorage.Image{}, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	if err != nil {
		return &imagestorage.Image{}, err
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, image, nil)

	return &imagestorage.Image{
		File:     buf.Bytes(),
		FileType: fileType[0],
	}, err
}

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	fileTypeRE  = regexp.MustCompile(`\.[a-zAZ]+$`)
)

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
