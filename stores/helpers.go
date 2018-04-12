package stores

// QueryModifierOp is a enum type to indicate the modifer operation
type QueryModifierOp = string

const (
	EQ  = "="
	LT  = "<"
	LTE = "<="
	NE  = "!="
	GT  = ">"
	GTE = ">="
)

var (
	// And is a query modifier
	And = QueryModifier{"AND", EQ, nil}
	// Or is a query modifier
	Or = QueryModifier{"OR", EQ, nil}
)

// QueryModifier is used in where queries to add selection criteria
type QueryModifier struct {
	Column   string
	Operator QueryModifierOp
	Value    interface{}
}

func QueryMod(col string, operator QueryModifierOp, value interface{}) QueryModifier {
	return QueryModifier{
		Column:   col,
		Operator: operator,
		Value:    value,
	}
}
