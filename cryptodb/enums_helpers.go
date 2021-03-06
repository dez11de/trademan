// Code generated by "enumer -json -sql -type Direction,Side,TakeProfitStrategy,TakeProfitMoment,Status,OrderType,OrderKind,LogSource -output enums_helpers.go"; DO NOT EDIT.

package cryptodb

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const _DirectionName = "LongShort"

var _DirectionIndex = [...]uint8{0, 4, 9}

const _DirectionLowerName = "longshort"

func (i Direction) String() string {
	if i < 0 || i >= Direction(len(_DirectionIndex)-1) {
		return fmt.Sprintf("Direction(%d)", i)
	}
	return _DirectionName[_DirectionIndex[i]:_DirectionIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _DirectionNoOp() {
	var x [1]struct{}
	_ = x[Long-(0)]
	_ = x[Short-(1)]
}

var _DirectionValues = []Direction{Long, Short}

var _DirectionNameToValueMap = map[string]Direction{
	_DirectionName[0:4]:      Long,
	_DirectionLowerName[0:4]: Long,
	_DirectionName[4:9]:      Short,
	_DirectionLowerName[4:9]: Short,
}

var _DirectionNames = []string{
	_DirectionName[0:4],
	_DirectionName[4:9],
}

// DirectionString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func DirectionString(s string) (Direction, error) {
	if val, ok := _DirectionNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _DirectionNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Direction values", s)
}

// DirectionValues returns all values of the enum
func DirectionValues() []Direction {
	return _DirectionValues
}

// DirectionStrings returns a slice of all String values of the enum
func DirectionStrings() []string {
	strs := make([]string, len(_DirectionNames))
	copy(strs, _DirectionNames)
	return strs
}

// IsADirection returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Direction) IsADirection() bool {
	for _, v := range _DirectionValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Direction
func (i Direction) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Direction
func (i *Direction) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Direction should be a string, got %s", data)
	}

	var err error
	*i, err = DirectionString(s)
	return err
}

func (i Direction) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *Direction) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of Direction: %[1]T(%[1]v)", value)
	}

	val, err := DirectionString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}

const _SideName = "BuySell"

var _SideIndex = [...]uint8{0, 3, 7}

const _SideLowerName = "buysell"

func (i Side) String() string {
	if i < 0 || i >= Side(len(_SideIndex)-1) {
		return fmt.Sprintf("Side(%d)", i)
	}
	return _SideName[_SideIndex[i]:_SideIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _SideNoOp() {
	var x [1]struct{}
	_ = x[Buy-(0)]
	_ = x[Sell-(1)]
}

var _SideValues = []Side{Buy, Sell}

var _SideNameToValueMap = map[string]Side{
	_SideName[0:3]:      Buy,
	_SideLowerName[0:3]: Buy,
	_SideName[3:7]:      Sell,
	_SideLowerName[3:7]: Sell,
}

var _SideNames = []string{
	_SideName[0:3],
	_SideName[3:7],
}

// SideString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func SideString(s string) (Side, error) {
	if val, ok := _SideNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _SideNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Side values", s)
}

// SideValues returns all values of the enum
func SideValues() []Side {
	return _SideValues
}

// SideStrings returns a slice of all String values of the enum
func SideStrings() []string {
	strs := make([]string, len(_SideNames))
	copy(strs, _SideNames)
	return strs
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

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of Side: %[1]T(%[1]v)", value)
	}

	val, err := SideString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}

const _TakeProfitStrategyName = "ManualAutoLinear"

var _TakeProfitStrategyIndex = [...]uint8{0, 6, 16}

const _TakeProfitStrategyLowerName = "manualautolinear"

func (i TakeProfitStrategy) String() string {
	if i < 0 || i >= TakeProfitStrategy(len(_TakeProfitStrategyIndex)-1) {
		return fmt.Sprintf("TakeProfitStrategy(%d)", i)
	}
	return _TakeProfitStrategyName[_TakeProfitStrategyIndex[i]:_TakeProfitStrategyIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _TakeProfitStrategyNoOp() {
	var x [1]struct{}
	_ = x[Manual-(0)]
	_ = x[AutoLinear-(1)]
}

var _TakeProfitStrategyValues = []TakeProfitStrategy{Manual, AutoLinear}

var _TakeProfitStrategyNameToValueMap = map[string]TakeProfitStrategy{
	_TakeProfitStrategyName[0:6]:       Manual,
	_TakeProfitStrategyLowerName[0:6]:  Manual,
	_TakeProfitStrategyName[6:16]:      AutoLinear,
	_TakeProfitStrategyLowerName[6:16]: AutoLinear,
}

var _TakeProfitStrategyNames = []string{
	_TakeProfitStrategyName[0:6],
	_TakeProfitStrategyName[6:16],
}

// TakeProfitStrategyString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func TakeProfitStrategyString(s string) (TakeProfitStrategy, error) {
	if val, ok := _TakeProfitStrategyNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _TakeProfitStrategyNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to TakeProfitStrategy values", s)
}

// TakeProfitStrategyValues returns all values of the enum
func TakeProfitStrategyValues() []TakeProfitStrategy {
	return _TakeProfitStrategyValues
}

// TakeProfitStrategyStrings returns a slice of all String values of the enum
func TakeProfitStrategyStrings() []string {
	strs := make([]string, len(_TakeProfitStrategyNames))
	copy(strs, _TakeProfitStrategyNames)
	return strs
}

// IsATakeProfitStrategy returns "true" if the value is listed in the enum definition. "false" otherwise
func (i TakeProfitStrategy) IsATakeProfitStrategy() bool {
	for _, v := range _TakeProfitStrategyValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for TakeProfitStrategy
func (i TakeProfitStrategy) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for TakeProfitStrategy
func (i *TakeProfitStrategy) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("TakeProfitStrategy should be a string, got %s", data)
	}

	var err error
	*i, err = TakeProfitStrategyString(s)
	return err
}

func (i TakeProfitStrategy) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *TakeProfitStrategy) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of TakeProfitStrategy: %[1]T(%[1]v)", value)
	}

	val, err := TakeProfitStrategyString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}

const _TakeProfitMomentName = "ImmediatelyOnEntryFilled"

var _TakeProfitMomentIndex = [...]uint8{0, 11, 24}

const _TakeProfitMomentLowerName = "immediatelyonentryfilled"

func (i TakeProfitMoment) String() string {
	if i < 0 || i >= TakeProfitMoment(len(_TakeProfitMomentIndex)-1) {
		return fmt.Sprintf("TakeProfitMoment(%d)", i)
	}
	return _TakeProfitMomentName[_TakeProfitMomentIndex[i]:_TakeProfitMomentIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _TakeProfitMomentNoOp() {
	var x [1]struct{}
	_ = x[Immediately-(0)]
	_ = x[OnEntryFilled-(1)]
}

var _TakeProfitMomentValues = []TakeProfitMoment{Immediately, OnEntryFilled}

var _TakeProfitMomentNameToValueMap = map[string]TakeProfitMoment{
	_TakeProfitMomentName[0:11]:       Immediately,
	_TakeProfitMomentLowerName[0:11]:  Immediately,
	_TakeProfitMomentName[11:24]:      OnEntryFilled,
	_TakeProfitMomentLowerName[11:24]: OnEntryFilled,
}

var _TakeProfitMomentNames = []string{
	_TakeProfitMomentName[0:11],
	_TakeProfitMomentName[11:24],
}

// TakeProfitMomentString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func TakeProfitMomentString(s string) (TakeProfitMoment, error) {
	if val, ok := _TakeProfitMomentNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _TakeProfitMomentNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to TakeProfitMoment values", s)
}

// TakeProfitMomentValues returns all values of the enum
func TakeProfitMomentValues() []TakeProfitMoment {
	return _TakeProfitMomentValues
}

// TakeProfitMomentStrings returns a slice of all String values of the enum
func TakeProfitMomentStrings() []string {
	strs := make([]string, len(_TakeProfitMomentNames))
	copy(strs, _TakeProfitMomentNames)
	return strs
}

// IsATakeProfitMoment returns "true" if the value is listed in the enum definition. "false" otherwise
func (i TakeProfitMoment) IsATakeProfitMoment() bool {
	for _, v := range _TakeProfitMomentValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for TakeProfitMoment
func (i TakeProfitMoment) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for TakeProfitMoment
func (i *TakeProfitMoment) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("TakeProfitMoment should be a string, got %s", data)
	}

	var err error
	*i, err = TakeProfitMomentString(s)
	return err
}

func (i TakeProfitMoment) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *TakeProfitMoment) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of TakeProfitMoment: %[1]T(%[1]v)", value)
	}

	val, err := TakeProfitMomentString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}

const _StatusName = "UnplannedPlannedCreatedOrderedNewUntriggeredPartiallyFilledFilledRejectedErrorStoppedPendingCancelCancelledDeactivatedClosedLiquidatedArchived"

var _StatusIndex = [...]uint8{0, 9, 16, 23, 30, 33, 44, 59, 65, 73, 78, 85, 98, 107, 118, 124, 134, 142}

const _StatusLowerName = "unplannedplannedcreatedorderednewuntriggeredpartiallyfilledfilledrejectederrorstoppedpendingcancelcancelleddeactivatedclosedliquidatedarchived"

func (i Status) String() string {
	if i < 0 || i >= Status(len(_StatusIndex)-1) {
		return fmt.Sprintf("Status(%d)", i)
	}
	return _StatusName[_StatusIndex[i]:_StatusIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _StatusNoOp() {
	var x [1]struct{}
	_ = x[Unplanned-(0)]
	_ = x[Planned-(1)]
	_ = x[Created-(2)]
	_ = x[Ordered-(3)]
	_ = x[New-(4)]
	_ = x[Untriggered-(5)]
	_ = x[PartiallyFilled-(6)]
	_ = x[Filled-(7)]
	_ = x[Rejected-(8)]
	_ = x[Error-(9)]
	_ = x[Stopped-(10)]
	_ = x[PendingCancel-(11)]
	_ = x[Cancelled-(12)]
	_ = x[Deactivated-(13)]
	_ = x[Closed-(14)]
	_ = x[Liquidated-(15)]
	_ = x[Archived-(16)]
}

var _StatusValues = []Status{Unplanned, Planned, Created, Ordered, New, Untriggered, PartiallyFilled, Filled, Rejected, Error, Stopped, PendingCancel, Cancelled, Deactivated, Closed, Liquidated, Archived}

var _StatusNameToValueMap = map[string]Status{
	_StatusName[0:9]:          Unplanned,
	_StatusLowerName[0:9]:     Unplanned,
	_StatusName[9:16]:         Planned,
	_StatusLowerName[9:16]:    Planned,
	_StatusName[16:23]:        Created,
	_StatusLowerName[16:23]:   Created,
	_StatusName[23:30]:        Ordered,
	_StatusLowerName[23:30]:   Ordered,
	_StatusName[30:33]:        New,
	_StatusLowerName[30:33]:   New,
	_StatusName[33:44]:        Untriggered,
	_StatusLowerName[33:44]:   Untriggered,
	_StatusName[44:59]:        PartiallyFilled,
	_StatusLowerName[44:59]:   PartiallyFilled,
	_StatusName[59:65]:        Filled,
	_StatusLowerName[59:65]:   Filled,
	_StatusName[65:73]:        Rejected,
	_StatusLowerName[65:73]:   Rejected,
	_StatusName[73:78]:        Error,
	_StatusLowerName[73:78]:   Error,
	_StatusName[78:85]:        Stopped,
	_StatusLowerName[78:85]:   Stopped,
	_StatusName[85:98]:        PendingCancel,
	_StatusLowerName[85:98]:   PendingCancel,
	_StatusName[98:107]:       Cancelled,
	_StatusLowerName[98:107]:  Cancelled,
	_StatusName[107:118]:      Deactivated,
	_StatusLowerName[107:118]: Deactivated,
	_StatusName[118:124]:      Closed,
	_StatusLowerName[118:124]: Closed,
	_StatusName[124:134]:      Liquidated,
	_StatusLowerName[124:134]: Liquidated,
	_StatusName[134:142]:      Archived,
	_StatusLowerName[134:142]: Archived,
}

var _StatusNames = []string{
	_StatusName[0:9],
	_StatusName[9:16],
	_StatusName[16:23],
	_StatusName[23:30],
	_StatusName[30:33],
	_StatusName[33:44],
	_StatusName[44:59],
	_StatusName[59:65],
	_StatusName[65:73],
	_StatusName[73:78],
	_StatusName[78:85],
	_StatusName[85:98],
	_StatusName[98:107],
	_StatusName[107:118],
	_StatusName[118:124],
	_StatusName[124:134],
	_StatusName[134:142],
}

// StatusString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func StatusString(s string) (Status, error) {
	if val, ok := _StatusNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _StatusNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Status values", s)
}

// StatusValues returns all values of the enum
func StatusValues() []Status {
	return _StatusValues
}

// StatusStrings returns a slice of all String values of the enum
func StatusStrings() []string {
	strs := make([]string, len(_StatusNames))
	copy(strs, _StatusNames)
	return strs
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

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of Status: %[1]T(%[1]v)", value)
	}

	val, err := StatusString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}

const _OrderTypeName = "MarketLimit"

var _OrderTypeIndex = [...]uint8{0, 6, 11}

const _OrderTypeLowerName = "marketlimit"

func (i OrderType) String() string {
	if i < 0 || i >= OrderType(len(_OrderTypeIndex)-1) {
		return fmt.Sprintf("OrderType(%d)", i)
	}
	return _OrderTypeName[_OrderTypeIndex[i]:_OrderTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _OrderTypeNoOp() {
	var x [1]struct{}
	_ = x[Market-(0)]
	_ = x[Limit-(1)]
}

var _OrderTypeValues = []OrderType{Market, Limit}

var _OrderTypeNameToValueMap = map[string]OrderType{
	_OrderTypeName[0:6]:       Market,
	_OrderTypeLowerName[0:6]:  Market,
	_OrderTypeName[6:11]:      Limit,
	_OrderTypeLowerName[6:11]: Limit,
}

var _OrderTypeNames = []string{
	_OrderTypeName[0:6],
	_OrderTypeName[6:11],
}

// OrderTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func OrderTypeString(s string) (OrderType, error) {
	if val, ok := _OrderTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _OrderTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to OrderType values", s)
}

// OrderTypeValues returns all values of the enum
func OrderTypeValues() []OrderType {
	return _OrderTypeValues
}

// OrderTypeStrings returns a slice of all String values of the enum
func OrderTypeStrings() []string {
	strs := make([]string, len(_OrderTypeNames))
	copy(strs, _OrderTypeNames)
	return strs
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

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of OrderType: %[1]T(%[1]v)", value)
	}

	val, err := OrderTypeString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}

const _OrderKindName = "MarketStopLossLimitStopLossEntryTakeProfit"

var _OrderKindIndex = [...]uint8{0, 14, 27, 32, 42}

const _OrderKindLowerName = "marketstoplosslimitstoplossentrytakeprofit"

func (i OrderKind) String() string {
	if i < 0 || i >= OrderKind(len(_OrderKindIndex)-1) {
		return fmt.Sprintf("OrderKind(%d)", i)
	}
	return _OrderKindName[_OrderKindIndex[i]:_OrderKindIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _OrderKindNoOp() {
	var x [1]struct{}
	_ = x[MarketStopLoss-(0)]
	_ = x[LimitStopLoss-(1)]
	_ = x[Entry-(2)]
	_ = x[TakeProfit-(3)]
}

var _OrderKindValues = []OrderKind{MarketStopLoss, LimitStopLoss, Entry, TakeProfit}

var _OrderKindNameToValueMap = map[string]OrderKind{
	_OrderKindName[0:14]:       MarketStopLoss,
	_OrderKindLowerName[0:14]:  MarketStopLoss,
	_OrderKindName[14:27]:      LimitStopLoss,
	_OrderKindLowerName[14:27]: LimitStopLoss,
	_OrderKindName[27:32]:      Entry,
	_OrderKindLowerName[27:32]: Entry,
	_OrderKindName[32:42]:      TakeProfit,
	_OrderKindLowerName[32:42]: TakeProfit,
}

var _OrderKindNames = []string{
	_OrderKindName[0:14],
	_OrderKindName[14:27],
	_OrderKindName[27:32],
	_OrderKindName[32:42],
}

// OrderKindString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func OrderKindString(s string) (OrderKind, error) {
	if val, ok := _OrderKindNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _OrderKindNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to OrderKind values", s)
}

// OrderKindValues returns all values of the enum
func OrderKindValues() []OrderKind {
	return _OrderKindValues
}

// OrderKindStrings returns a slice of all String values of the enum
func OrderKindStrings() []string {
	strs := make([]string, len(_OrderKindNames))
	copy(strs, _OrderKindNames)
	return strs
}

// IsAOrderKind returns "true" if the value is listed in the enum definition. "false" otherwise
func (i OrderKind) IsAOrderKind() bool {
	for _, v := range _OrderKindValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for OrderKind
func (i OrderKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for OrderKind
func (i *OrderKind) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("OrderKind should be a string, got %s", data)
	}

	var err error
	*i, err = OrderKindString(s)
	return err
}

func (i OrderKind) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *OrderKind) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of OrderKind: %[1]T(%[1]v)", value)
	}

	val, err := OrderKindString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}

const _LogSourceName = "ExchangeServerUser"

var _LogSourceIndex = [...]uint8{0, 8, 14, 18}

const _LogSourceLowerName = "exchangeserveruser"

func (i LogSource) String() string {
	if i < 0 || i >= LogSource(len(_LogSourceIndex)-1) {
		return fmt.Sprintf("LogSource(%d)", i)
	}
	return _LogSourceName[_LogSourceIndex[i]:_LogSourceIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _LogSourceNoOp() {
	var x [1]struct{}
	_ = x[Exchange-(0)]
	_ = x[Server-(1)]
	_ = x[User-(2)]
}

var _LogSourceValues = []LogSource{Exchange, Server, User}

var _LogSourceNameToValueMap = map[string]LogSource{
	_LogSourceName[0:8]:        Exchange,
	_LogSourceLowerName[0:8]:   Exchange,
	_LogSourceName[8:14]:       Server,
	_LogSourceLowerName[8:14]:  Server,
	_LogSourceName[14:18]:      User,
	_LogSourceLowerName[14:18]: User,
}

var _LogSourceNames = []string{
	_LogSourceName[0:8],
	_LogSourceName[8:14],
	_LogSourceName[14:18],
}

// LogSourceString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func LogSourceString(s string) (LogSource, error) {
	if val, ok := _LogSourceNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _LogSourceNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to LogSource values", s)
}

// LogSourceValues returns all values of the enum
func LogSourceValues() []LogSource {
	return _LogSourceValues
}

// LogSourceStrings returns a slice of all String values of the enum
func LogSourceStrings() []string {
	strs := make([]string, len(_LogSourceNames))
	copy(strs, _LogSourceNames)
	return strs
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

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of LogSource: %[1]T(%[1]v)", value)
	}

	val, err := LogSourceString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}
