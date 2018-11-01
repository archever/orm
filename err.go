package orm

import "errors"

var (
	ErrNotFund         = errors.New("资源未找到")
	ErrDestType        = errors.New("接受类型不匹配")
	ErrNotAssignable   = errors.New("自定义类型不能赋值")
	ErrCreateEmptyData = errors.New("插入数据不能为空")
	ErrTableNotSet     = errors.New("未指定表名")
)
