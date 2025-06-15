package imaging

import (
	"fmt"

	"github.com/disintegration/imaging"
)

// RemoveImageExif 移除图像的元信息
// EXIF 信息是可交换图像文件格式（Exchangeable Image File Format）的缩写，专门用于存储数码照片的属性信息和拍摄数据。它通常包含以下内容：
//
// 拍摄日期和时间：记录照片拍摄的具体时间。
// 相机型号：拍摄照片的相机品牌和型号。
// 光圈大小：拍摄时相机的光圈设置。
// 快门速度：拍摄时相机的快门速度。
// ISO 感光度：拍摄时相机的感光度设置。
// 白平衡设置：拍摄时相机的白平衡设置。
// GPS 位置信息：照片拍摄的具体地理位置。
func RemoveImageExif(src string, dst string) error {
	// 打开源图片
	img, err := imaging.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open image %s: %w", src, err)
	}

	// 调整图片大小
	resized := imaging.Resize(img, 800, 0, imaging.Lanczos)

	// 移除EXIF信息
	processed := imaging.Clone(resized)

	// 保存处理后的图片
	if err := imaging.Save(processed, dst); err != nil {
		return fmt.Errorf("failed to save image to %s: %w", dst, err)
	}

	return nil
}
