package operator

// opFactory 操作符仓库
var opFactory map[OpType]IOperator

// init 初始化操作符仓库
func init() {
	opFactory = make(map[OpType]IOperator)

	eq := EqualOp(Equal)
	opFactory[eq.Name()] = &eq
	ne := NotEqualOp(NotEqual)
	opFactory[ne.Name()] = &ne

	in := InOp(In)
	opFactory[in.Name()] = &in
	nin := NotInOp(NotIn)
	opFactory[nin.Name()] = &nin

	lt := LessOp(Less)
	opFactory[lt.Name()] = &lt
	lte := LessOrEqualOp(LessOrEqual)
	opFactory[lte.Name()] = &lte

	gt := GreaterOp(Greater)
	opFactory[gt.Name()] = &gt
	gte := GreaterOrEqualOp(GreaterOrEqual)
	opFactory[gte.Name()] = &gte

	datetimeLt := DatetimeLessOp(DatetimeLess)
	opFactory[datetimeLt.Name()] = &datetimeLt
	datetimeLte := DatetimeLessOrEqualOp(DatetimeLessOrEqual)
	opFactory[datetimeLte.Name()] = &datetimeLte
	datetimeGt := DatetimeGreaterOp(DatetimeGreater)
	opFactory[datetimeGt.Name()] = &datetimeGt
	datetimeGte := DatetimeGreaterOrEqualOp(DatetimeGreaterOrEqual)
	opFactory[datetimeGte.Name()] = &datetimeGte

	beginsWith := BeginsWithOp(BeginsWith)
	opFactory[beginsWith.Name()] = &beginsWith
	beginsWithInsensitive := BeginsWithInsensitiveOp(BeginsWithInsensitive)
	opFactory[beginsWithInsensitive.Name()] = &beginsWithInsensitive
	notBeginsWith := NotBeginsWithOp(NotBeginsWith)
	opFactory[notBeginsWith.Name()] = &notBeginsWith
	notBeginsWithInsensitive := NotBeginsWithInsensitiveOp(NotBeginsWithInsensitive)
	opFactory[notBeginsWithInsensitive.Name()] = &notBeginsWithInsensitive

}

// Operator 根据操作符类型从仓库获取一个操作符对象
func GetOperator(opType OpType) IOperator {
	op, exist := opFactory[opType]
	if !exist {
		unknown := UnknownOp(Unknown)
		return &unknown
	}

	return op
}
