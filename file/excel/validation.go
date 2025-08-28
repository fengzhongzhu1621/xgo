package excel

import (
	"fmt"

	"github.com/fengzhongzhu1621/xgo/cast"
	"github.com/xuri/excelize/v2"
)

const (
	// fieldTypeBoolTrue bool类型为true的值
	fieldTypeBoolTrue = "true"
	// fieldTypeBoolFalse bool类型为false的值
	fieldTypeBoolFalse = "false"
	// enumRefSuffix 引用某一个sheet的第一列前缀
	enumRefSuffix = "!$A:$A"
	// errTitle 当填入excel数据不满足校验时，弹出的错误框标题
	errTitle = "警告"
	// errMessage 当填入excel数据不满足校验时，弹出的错误框内容
	errMessage = "此值与此单元格定义的数据验证限制不匹配。"
)

// ValidationParam validation parameter
type ValidationParam struct {
	Type   FieldType   // 单元格类型
	Sqref  string      // 应用数据验证的单元格范围（例如，"A1:A10"）
	Option interface{} // 下拉列表的选项 / 引用配置
}

// NewDataValidation 创建一个单元格验证器
func NewDataValidation(param *ValidationParam) (*excelize.DataValidation, error) {
	validation := excelize.NewDataValidation(true)
	validation.SetSqref(param.Sqref)

	switch param.Type {
	case Decimal:
		validation.Type = string(Decimal)
	case Bool:
		// 设置单元格的下拉列表的值
		if err := validation.SetDropList([]string{fieldTypeBoolTrue, fieldTypeBoolFalse}); err != nil {
			return nil, err
		}
	case Enum:
		// 设置单元格的下拉列表的值
		strArr, err := cast.ToStringSliceE(param.Option)
		if err != nil {
			return nil, err
		}
		if err := validation.SetDropList(strArr); err != nil {
			return nil, err
		}
	case Ref:
		// 在指定的单元格范围内设置数据验证的下拉列表，其源引用为工作表中的某个范围
		sheet := cast.ToString(param.Option)
		ref := fmt.Sprintf("'%s'%s", sheet, enumRefSuffix)
		validation.SetSqrefDropList(ref)
	}

	// 设置验证错误信息
	validation.SetError(excelize.DataValidationErrorStyleStop, errTitle, errMessage)

	return validation, nil
}
