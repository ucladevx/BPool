package postgres

import (
	"strconv"

	"github.com/ucladevx/BPool/stores"
)

func generateWhereStatement(modifiers *[]stores.QueryModifier) (string, []interface{}) {
	var args []interface{}
	where := "WHERE "

	count := 1
	for _, modifier := range *modifiers {
		if modifier.Column == "AND" || modifier.Column == "OR" {
			where += modifier.Column + " "
			continue
		}

		if modifier.Column == "" || modifier.Value == nil {
			return "", nil
		}

		where += modifier.Column + modifier.Operator + "$" + strconv.Itoa(count) + " "
		args = append(args, modifier.Value)
		count++
	}

	return where, args
}
