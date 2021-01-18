package generate

import "encoding/json"

// CriarRegisto -
func CriarRegisto(conteudo interface{}) ([]byte, error) {
	res, err := json.MarshalIndent(conteudo, "", "\t")
	if err != nil {
		return nil, err
	}
	return res, nil
}
