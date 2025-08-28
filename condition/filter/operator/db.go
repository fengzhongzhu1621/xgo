package operator

// mongodb 支持的操作符
const (
	DBEQ = "$eq"
	DBNE = "$ne"

	DBIN           = "$in"
	DBNIN          = "$nin"
	DBMULTIPLELike = "$multilike"

	DBOR  = "$or"
	DBNOR = "$nor"
	DBAND = "$and"
	DBNot = "$not"

	DBLIKE = "$regex"

	// DBOPTIONS the db operator,used with $regex
	// detail to see https://docs.mongodb.com/manual/reference/operator/query/regex/#op._S_options
	DBOPTIONS = "$options"

	DBLT  = "$lt"
	DBLTE = "$lte"
	DBGT  = "$gt"
	DBGTE = "$gte"

	DBExists = "$exists"

	DBCount = "$count"

	DBGroup = "$group"

	DBMatch = "$match"

	DBSum = "$sum"

	DBPush = "$push"

	DBUNSET = "$unset"

	// DBAddToSet The $addToSet operator adds a value to an array unless the value is already present, in which case $addToSet does nothing to that array.
	DBAddToSet = "$addToSet"

	// DBPull The $pull operator removes from an existing array all instances of a value or values that match a specified condition.
	DBPull = "$pull"

	// DBAll matches arrays that contain all elements specified in the query.
	DBAll = "$all"

	// DBProject passes along the documents with the requested fields to the next stage in the pipeline
	DBProject = "$project"

	// DBSize counts and returns the total number of items in an array
	DBSize = "$size"

	// DBType selects documents where the value of the field is an instance of the specified BSON type(s).
	// Querying by data type is useful when dealing with highly unstructured data where data types are not predictable.
	DBType = "$type"

	// DBSort the db operator
	DBSort = "$sort"

	// DBReplaceRoot the db operator
	DBReplaceRoot = "$replaceRoot"

	// DBLimit the db operator to limit return number of doc
	DBLimit = "$limit"
	// Limit data quantity limit
	Limit = "$limit"

	// DBUnwind used to split values contained in an array field into separate doc
	DBUnwind = "$unwind"

	// DBLookup used to perform multi-table association operations
	DBLookup = "$lookup"

	// DBFrom collection to join
	DBFrom = "from"

	// DBLocalField field from the input documents
	DBLocalField = "localField"

	// DBForeignField field from the documents of the "from" collection
	DBForeignField = "foreignField"

	// DBAs output array field
	DBAs = "as"

	// DBSkip skip data index
	DBSkip = "$skip"
)

const (
	// GT TODO
	// Comparison Operator
	GT string = "$gt"
	// LT TODO
	LT string = "$lt"
	// GTE TODO
	GTE string = "$gte"
	// LTE TODO
	LTE string = "$lte"
	// IN TODO
	IN string = "$in"
	// NIN TODO
	NIN string = "$nin"
	// EQ TODO
	EQ string = "$eq"
	// NEQ TODO
	NEQ string = "$ne"
	// REGEX TODO
	REGEX string = "$regex"

	// AND TODO
	// Logic Operator
	AND string = "$and"
	// OR TODO
	OR string = "$or"
	// NOT TODO
	NOT string = "$not"
	// NOR TODO
	NOR string = "$nor"

	// EXISTS TODO
	// TODO:
	// Elements Operator
	EXISTS string = "$exists"
	// TYPE TODO
	TYPE string = "$type"

	// ALL TODO
	// Array Operator
	ALL string = "$all"
	// ELEMMATCH TODO
	ELEMMATCH string = "$elemMatch"
	// SIZE TODO
	SIZE string = "$size"
)
