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

	contains := ContainsOp(Contains)
	opFactory[contains.Name()] = &contains
	containsSensitive := ContainsSensitiveOp(ContainsSensitive)
	opFactory[containsSensitive.Name()] = &containsSensitive
	notContains := NotContainsOp(NotContains)
	opFactory[notContains.Name()] = &notContains
	notContainsInsensitive := NotContainsInsensitiveOp(NotContainsInsensitive)
	opFactory[notContainsInsensitive.Name()] = &notContainsInsensitive

	endsWith := EndsWithOp(EndsWith)
	opFactory[endsWith.Name()] = &endsWith
	endsWithInsensitive := EndsWithInsensitiveOp(EndsWithInsensitive)
	opFactory[endsWithInsensitive.Name()] = &endsWithInsensitive
	notEndsWith := NotEndsWithOp(NotEndsWith)
	opFactory[notEndsWith.Name()] = &notEndsWith
	notEndsWithInsensitive := NotEndsWithInsensitiveOp(NotEndsWithInsensitive)
	opFactory[notEndsWithInsensitive.Name()] = &notEndsWithInsensitive

	isEmpty := IsEmptyOp(IsEmpty)
	opFactory[isEmpty.Name()] = &isEmpty
	isNotEmpty := IsNotEmptyOp(IsNotEmpty)
	opFactory[isNotEmpty.Name()] = &isNotEmpty

	size := SizeOp(Size)
	opFactory[size.Name()] = &size

	isNull := IsNullOp(IsNull)
	opFactory[isNull.Name()] = &isNull
	isNotNull := IsNotNullOp(IsNotNull)
	opFactory[isNotNull.Name()] = &isNotNull

	exist := ExistOp(Exist)
	opFactory[exist.Name()] = &exist
	notExist := NotExistOp(NotExist)
	opFactory[notExist.Name()] = &notExist

	obj := ObjectOp(Object)
	opFactory[obj.Name()] = &obj

	filterArr := ArrayOp(Array)
	opFactory[filterArr.Name()] = &filterArr
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

// MatchedData is the data to be matched with the expression's rule
type MatchedData interface {
	GetValue(field string) (interface{}, error)
}
