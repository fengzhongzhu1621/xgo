package cast

import "errors"

func TransSlice2Interface(old interface{}) ([]interface{}, error) {
	switch trans := old.(type) {
	case []string:
		new := make([]interface{}, len(trans))
		for i, v := range trans {
			new[i] = v
		}
		return new, nil
	default:
		return nil, errors.New("illegal type")
	}
}

func TransSlice2Interface2(slice []string) ([]interface{}, error) {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice, nil
}
