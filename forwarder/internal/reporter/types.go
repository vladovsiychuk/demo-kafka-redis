package reporter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type IntOrString int
type FloatOrString float64

func (i *IntOrString) UnmarshalJSON(b []byte) error {
	var intVal int
	if err := json.Unmarshal(b, &intVal); err == nil {
		*i = IntOrString(intVal)
		return nil
	}
	var strVal string
	if err := json.Unmarshal(b, &strVal); err == nil {
		strVal = strings.TrimSpace(strVal)
		if strVal == "" {
			*i = 0
			return nil
		}
		intVal, err := strconv.Atoi(strVal)
		if err != nil {
			return fmt.Errorf("invalid int string: %q", strVal)
		}
		*i = IntOrString(intVal)
		return nil
	}
	return fmt.Errorf("IntOrString: cannot unmarshal %s", string(b))
}

func (f *FloatOrString) UnmarshalJSON(b []byte) error {
	var floatVal float64
	if err := json.Unmarshal(b, &floatVal); err == nil {
		*f = FloatOrString(floatVal)
		return nil
	}
	var strVal string
	if err := json.Unmarshal(b, &strVal); err == nil {
		strVal = strings.TrimSpace(strVal)
		if strVal == "" {
			*f = 0
			return nil
		}
		floatVal, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return fmt.Errorf("invalid float string: %q", strVal)
		}
		*f = FloatOrString(floatVal)
		return nil
	}
	return fmt.Errorf("FloatOrString: cannot unmarshal %s", string(b))
}

type Data struct {
	Clicks      IntOrString   `json:"clicks"`
	Cost        FloatOrString `json:"cost"`
	Date        string        `json:"date"`
	Impressions IntOrString   `json:"impressions"`
	Installs    IntOrString   `json:"installs"`
}
