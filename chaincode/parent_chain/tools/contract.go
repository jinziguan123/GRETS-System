package tools

import (
	"encoding/json"

	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// 根据迭代器构造结果
func ConstructResultByIterator[T interface{}](resultsIterator shim.StateQueryIteratorInterface) ([]*T, error) {
	var results []*T
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var result T
		if err := json.Unmarshal(response.Value, &result); err != nil {
			return nil, err
		}
		results = append(results, &result)
	}
	return results, nil
}

// 根据查询字符串查询
func SelectByQueryString[T interface{}](ctx contractapi.TransactionContextInterface, queryString string) ([]*T, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return ConstructResultByIterator[T](resultsIterator)
}
