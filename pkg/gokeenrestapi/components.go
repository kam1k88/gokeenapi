package gokeenrestapi

import (
	"encoding/json"

	"github.com/noksa/gokeenapi/internal/gokeenspinner"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapimodels"
)

type keeneticComponents struct {
}

var (
	Components keeneticComponents
)

// GetAllComponents retrieves all components available on keenetic (installed or not)
func (*keeneticComponents) GetAllComponents() (gokeenrestapimodels.RciComponentsList, error) {
	var result gokeenrestapimodels.RciComponentsList
	err := gokeenspinner.WrapWithSpinner("Fetching components", func() error {
		body, err := Common.ExecutePostSubPath("/rci/components/list", "{}")
		if err != nil {
			return err
		}
		return json.Unmarshal(body, &result)
	})
	return result, err
}
