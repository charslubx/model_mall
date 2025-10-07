package logic

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"modelmall/backend/internal/models"
	"modelmall/backend/internal/svc"
	"modelmall/backend/internal/types"
)

// UploadImageRequest 上传图片请求
type UploadImageRequest struct {
	Filename   string
	FilePath   string
	FileSize   int64
	MimeType   string
	Width      *int
	Height     *int
	FileData   []byte
	UploadedBy *int64
}

type UploadImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadImageLogic {
	return &UploadImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UploadImage 上传图片并进行分类
func (l *UploadImageLogic) UploadImage(req *UploadImageRequest) (*types.UploadImageResp, error) {
	// 创建图片记录
	image := &models.Image{
		Filename:   req.Filename,
		FilePath:   req.FilePath,
		FileSize:   req.FileSize,
		MimeType:   req.MimeType,
		Width:      req.Width,
		Height:     req.Height,
		UploadedBy: req.UploadedBy,
		Status:     models.ImageStatusProcessing,
	}

	// 保存图片记录到数据库
	imageID, err := l.svcCtx.ImageRepo.CreateImage(l.ctx, image)
	if err != nil {
		return nil, fmt.Errorf("保存图片记录失败: %w", err)
	}

	// 调用模型服务进行分类
	classifications, err := l.svcCtx.ModelService.ClassifyImageFromBytes(l.ctx, req.FileData, req.Filename)
	if err != nil {
		// 更新状态为失败
		_ = l.svcCtx.ImageRepo.UpdateImageStatus(l.ctx, imageID, models.ImageStatusFailed)
		return nil, fmt.Errorf("图片分类失败: %w", err)
	}

	// 保存分类结果到数据库
	for _, cls := range classifications {
		classification := &models.ImageClassification{
			ImageID:      imageID,
			Label:        cls.Label,
			Confidence:   cls.Confidence,
			ModelName:    l.svcCtx.ModelService.GetModelName(),
			ModelVersion: &l.svcCtx.ModelService.GetModelVersion(),
		}

		err = l.svcCtx.ImageRepo.CreateClassification(l.ctx, classification)
		if err != nil {
			l.Logger.Errorf("保存分类记录失败: %v", err)
		}
	}

	// 更新图片状态为已分类
	err = l.svcCtx.ImageRepo.UpdateImageStatus(l.ctx, imageID, models.ImageStatusClassified)
	if err != nil {
		l.Logger.Errorf("更新图片状态失败: %v", err)
	}

	// 构造响应
	resp := &types.UploadImageResp{
		ImageID:  imageID,
		Filename: req.Filename,
		FilePath: req.FilePath,
		FileSize: req.FileSize,
		Status:   models.ImageStatusClassified,
		Classifications: make([]types.ClassificationItem, 0, len(classifications)),
	}

	for _, cls := range classifications {
		resp.Classifications = append(resp.Classifications, types.ClassificationItem{
			Label:      cls.Label,
			Confidence: cls.Confidence,
		})
	}

	return resp, nil
}

type GetImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetImageLogic {
	return &GetImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetImage 获取图片信息
func (l *GetImageLogic) GetImage(req *types.GetImageReq) (*types.GetImageResp, error) {
	// 获取图片信息
	image, err := l.svcCtx.ImageRepo.GetImageByID(l.ctx, req.ImageID)
	if err != nil {
		return nil, fmt.Errorf("获取图片信息失败: %w", err)
	}

	// 获取分类信息
	classifications, err := l.svcCtx.ImageRepo.GetClassificationsByImageID(l.ctx, req.ImageID)
	if err != nil {
		l.Logger.Errorf("获取分类信息失败: %v", err)
		classifications = []*models.ImageClassification{}
	}

	// 构造响应
	resp := &types.GetImageResp{
		ID:         image.ID,
		Filename:   image.Filename,
		FilePath:   image.FilePath,
		FileSize:   image.FileSize,
		MimeType:   image.MimeType,
		Width:      image.Width,
		Height:     image.Height,
		UploadedBy: image.UploadedBy,
		Status:     image.Status,
		CreatedAt:  image.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  image.UpdatedAt.Format("2006-01-02 15:04:05"),
		Classifications: make([]types.ClassificationItem, 0, len(classifications)),
	}

	for _, cls := range classifications {
		resp.Classifications = append(resp.Classifications, types.ClassificationItem{
			Label:      cls.Label,
			Confidence: cls.Confidence,
		})
	}

	return resp, nil
}

type ListImagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListImagesLogic {
	return &ListImagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListImages 列出图片
func (l *ListImagesLogic) ListImages(req *types.ListImagesReq) (*types.ListImagesResp, error) {
	// 计算offset
	offset := (req.Page - 1) * req.PageSize

	// 获取图片列表
	images, err := l.svcCtx.ImageRepo.ListImages(l.ctx, offset, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取图片列表失败: %w", err)
	}

	// 构造响应
	resp := &types.ListImagesResp{
		Total: int64(len(images)), // 这里简化处理，实际应该查询总数
		List:  make([]types.GetImageResp, 0, len(images)),
	}

	for _, image := range images {
		// 获取每个图片的分类信息
		classifications, err := l.svcCtx.ImageRepo.GetClassificationsByImageID(l.ctx, image.ID)
		if err != nil {
			l.Logger.Errorf("获取分类信息失败: %v", err)
			classifications = []*models.ImageClassification{}
		}

		imageResp := types.GetImageResp{
			ID:         image.ID,
			Filename:   image.Filename,
			FilePath:   image.FilePath,
			FileSize:   image.FileSize,
			MimeType:   image.MimeType,
			Width:      image.Width,
			Height:     image.Height,
			UploadedBy: image.UploadedBy,
			Status:     image.Status,
			CreatedAt:  image.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  image.UpdatedAt.Format("2006-01-02 15:04:05"),
			Classifications: make([]types.ClassificationItem, 0, len(classifications)),
		}

		for _, cls := range classifications {
			imageResp.Classifications = append(imageResp.Classifications, types.ClassificationItem{
				Label:      cls.Label,
				Confidence: cls.Confidence,
			})
		}

		resp.List = append(resp.List, imageResp)
	}

	return resp, nil
}