package functions

import (
	"encoding/json"
	"math/big"

	"github.com/jackc/pgtype"
	"gopkg.in/inf.v0"
)

// Helper function to convert various types to *big.Int
func ToBigInt(val interface{}) *big.Int {
	switch v := val.(type) {
	case *inf.Dec:
		// Convert *inf.Dec to *big.Int
		bigIntValue := new(big.Int)
		bigIntValue.SetString(v.String(), 10)
		return bigIntValue
	case *big.Float:
		// Convert *big.Float to *big.Int
		bigIntValue := new(big.Int)
		v.Int(bigIntValue)
		return bigIntValue
	case int:
		return big.NewInt(int64(v))
	case int16:
		return big.NewInt(int64(v))
	case int32:
		return big.NewInt(int64(v))
	case int64:
		return big.NewInt(v)
	case string:
		bigIntValue := new(big.Int)
		bigIntValue, success := bigIntValue.SetString(v, 10)
		if !success {
			var bigIntNullable *big.Int
			return bigIntNullable
		}
		return bigIntValue
	case *string:
		if v == nil {
			return nil
		}
		bigIntValue := new(big.Int)
		bigIntValue, success := bigIntValue.SetString(*v, 10)
		if !success {
			return nil
		}
		return bigIntValue
	case pgtype.Numeric:
		if v.Int == nil {
			return nil
		}
		bigInt := new(big.Int).Set(v.Int)
		if v.Exp > 0 {
			bigInt.Mul(bigInt, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(v.Exp)), nil))
		} else if v.Exp < 0 {
			return nil // fractional numbers cannot be converted to big.Int
		}

		return bigInt
	case json.Number:
		bigIntValue := new(big.Int)
		bigIntValue, success := bigIntValue.SetString(string(v), 10)
		if !success {
			return nil
		}
		return bigIntValue
	default:
		return nil
	}
}
