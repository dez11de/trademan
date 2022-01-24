// Code generated by "enumer -sql -json -type=Side,Status,OrderType,LogSource -output enums_helpers.go"; DO NOT EDIT.

//
package cryptodb

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const _SideName = "SideLongSideShort"

var _SideIndex = [...]uint8{0, 8, 17}

func (i Side) String() string {
	if i < 0 || i >= Side(len(_SideIndex)-1) {
		return fmt.Sprintf("Side(%d)", i)
	}
	return _SideName[_SideIndex[i]:_SideIndex[i+1]]
}

var _SideValues = []Side{0, 1}

var _SideNameToValueMap = map[string]Side{
	_SideName[0:8]:  0,
	_SideName[8:17]: 1,
}

// SideString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func SideString(s string) (Side, error) {
	if val, ok := _SideNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Side values", s)
}

// SideValues returns all values of the enum
func SideValues() []Side {
	return _SideValues
}

// IsASide returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Side) IsASide() bool {
	for _, v := range _SideValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Side
func (i Side) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Side
func (i *Side) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Side should be a string, got %s", data)
	}

	var err error
	*i, err = SideString(s)
	return err
}

func (i Side) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *Side) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("value is not a byte slice")
		}

		str = string(bytes[:])
	}

	val, err := SideString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}

const _StatusName = "StatusPlannedStatusOrderedStatusFilledStatusStoppedStatusClosedStatusCancelledStatusLiquidatedStatusLogged"

var _StatusIndex = [...]uint8{0, 13, 26, 38, 51, 63, 78, 94, 106}

func (i Status) String() string {
	if i < 0 || i >= Status(len(_StatusIndex)-1) {
		return fmt.Sprintf("Status(%d)", i)
	}
	return _StatusName[_StatusIndex[i]:_StatusIndex[i+1]]
}

var _StatusValues = []Status{0, 1, 2, 3, 4, 5, 6, 7}

var _StatusNameToValueMap = map[string]Status{
	_StatusName[0:13]:   0,
	_StatusName[13:26]:  1,
	_StatusName[26:38]:  2,
	_StatusName[38:51]:  3,
	_StatusName[51:63]:  4,
	_StatusName[63:78]:  5,
	_StatusName[78:94]:  6,
	_StatusName[94:106]: 7,
}

// StatusString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func StatusString(s string) (Status, error) {
	if val, ok := _StatusNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Status values", s)
}

// StatusValues returns all values of the enum
func StatusValues() []Status {
	return _StatusValues
}

// IsAStatus returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Status) IsAStatus() bool {
	for _, v := range _StatusValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Status
func (i Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Status
func (i *Status) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Status should be a string, got %s", data)
	}

	var err error
	*i, err = StatusString(s)
	return err
}

func (i Status) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *Status) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("value is not a byte slice")
		}

		str = string(bytes[:])
	}

	val, err := StatusString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}

const _OrderTypeName = "TypeHardStopLossTypeSoftStopLossTypeEntryTypeTakeProfit"

var _OrderTypeIndex = [...]uint8{0, 16, 32, 41, 55}

func (i OrderType) String() string {
	if i < 0 || i >= OrderType(len(_OrderTypeIndex)-1) {
		return fmt.Sprintf("OrderType(%d)", i)
	}
	return _OrderTypeName[_OrderTypeIndex[i]:_OrderTypeIndex[i+1]]
}

var _OrderTypeValues = []OrderType{0, 1, 2, 3}

var _OrderTypeNameToValueMap = map[string]OrderType{
	_OrderTypeName[0:16]:  0,
	_OrderTypeName[16:32]: 1,
	_OrderTypeName[32:41]: 2,
	_OrderTypeName[41:55]: 3,
}

// OrderTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func OrderTypeString(s string) (OrderType, error) {
	if val, ok := _OrderTypeNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to OrderType values", s)
}

// OrderTypeValues returns all values of the enum
func OrderTypeValues() []OrderType {
	return _OrderTypeValues
}

// IsAOrderType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i OrderType) IsAOrderType() bool {
	for _, v := range _OrderTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for OrderType
func (i OrderType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for OrderType
func (i *OrderType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("OrderType should be a string, got %s", data)
	}

	var err error
	*i, err = OrderTypeString(s)
	return err
}

func (i OrderType) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *OrderType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("value is not a byte slice")
		}

		str = string(bytes[:])
	}

	val, err := OrderTypeString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}

const _LogSourceName = "SourceTriggerSourceSoftwareSourceUser"

var _LogSourceIndex = [...]uint8{0, 13, 27, 37}

func (i LogSource) String() string {
	if i < 0 || i >= LogSource(len(_LogSourceIndex)-1) {
		return fmt.Sprintf("LogSource(%d)", i)
	}
	return _LogSourceName[_LogSourceIndex[i]:_LogSourceIndex[i+1]]
}

var _LogSourceValues = []LogSource{0, 1, 2}

var _LogSourceNameToValueMap = map[string]LogSource{
	_LogSourceName[0:13]:  0,
	_LogSourceName[13:27]: 1,
	_LogSourceName[27:37]: 2,
}

// LogSourceString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func LogSourceString(s string) (LogSource, error) {
	if val, ok := _LogSourceNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to LogSource values", s)
}

// LogSourceValues returns all values of the enum
func LogSourceValues() []LogSource {
	return _LogSourceValues
}

// IsALogSource returns "true" if the value is listed in the enum definition. "false" otherwise
func (i LogSource) IsALogSource() bool {
	for _, v := range _LogSourceValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for LogSource
func (i LogSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for LogSource
func (i *LogSource) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("LogSource should be a string, got %s", data)
	}

	var err error
	*i, err = LogSourceString(s)
	return err
}

func (i LogSource) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *LogSource) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("value is not a byte slice")
		}

		str = string(bytes[:])
	}

	val, err := LogSourceString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}
