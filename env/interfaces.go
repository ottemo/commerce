package env

import (
	"time"

	"github.com/ottemo/commerce/utils"
)

// Package global constants
const (
	ConstConfigTypeID       = utils.ConstDataTypeID
	ConstConfigTypeBoolean  = utils.ConstDataTypeBoolean
	ConstConfigTypeVarchar  = utils.ConstDataTypeVarchar
	ConstConfigTypeText     = utils.ConstDataTypeText
	ConstConfigTypeInteger  = utils.ConstDataTypeInteger
	ConstConfigTypeDecimal  = utils.ConstDataTypeDecimal
	ConstConfigTypeMoney    = utils.ConstDataTypeMoney
	ConstConfigTypeFloat    = utils.ConstDataTypeFloat
	ConstConfigTypeDatetime = utils.ConstDataTypeDatetime
	ConstConfigTypeJSON     = utils.ConstDataTypeJSON
	ConstConfigTypeHTML     = utils.ConstDataTypeHTML
	ConstConfigTypeGroup    = "group"
	ConstConfigTypeSecret   = "secret"

	ConstLogPrefixError   = "ERROR"
	ConstLogPrefixWarning = "WARNING"
	ConstLogPrefixDebug   = "DEBUG"
	ConstLogPrefixInfo    = "INFO"

	ConstErrorLevelAPI        = 10
	ConstErrorLevelModel      = 9
	ConstErrorLevelActor      = 8
	ConstErrorLevelHelper     = 7
	ConstErrorLevelService    = 4
	ConstErrorLevelServiceAct = 3
	ConstErrorLevelCritical   = 2
	ConstErrorLevelStartStop  = 1
	ConstErrorLevelExternal   = 0

	ConstErrorModule = "env"
	ConstErrorLevel  = ConstErrorLevelService
)

// InterfaceSchedule is an interface to system schedule service
type InterfaceSchedule interface {
	Execute()
	RunTask(params map[string]interface{}) error

	Enable() error
	Disable() error

	Set(param string, value interface{}) error
	Get(param string) interface{}

	GetInfo() map[string]interface{}
}

// InterfaceScheduler is an interface to system scheduler service
type InterfaceScheduler interface {
	ListTasks() []string
	RegisterTask(name string, task FuncCronTask) error

	ScheduleAtTime(scheduleTime time.Time, taskName string, taskParams map[string]interface{}) (InterfaceSchedule, error)
	ScheduleRepeat(cronExpr string, taskName string, taskParams map[string]interface{}) (InterfaceSchedule, error)
	ScheduleOnce(cronExpr string, taskName string, taskParams map[string]interface{}) (InterfaceSchedule, error)

	ListSchedules() []InterfaceSchedule
}

// InterfaceEventBus is an interface to system event processor
type InterfaceEventBus interface {
	RegisterListener(event string, listener FuncEventListener)
	New(event string, eventData map[string]interface{})
}

// InterfaceErrorBus is an interface to system error processor
type InterfaceErrorBus interface {
	GetErrorLevel(error) int
	GetErrorCode(error) string
	GetErrorMessage(error) string

	RegisterListener(FuncErrorListener)

	Dispatch(err error) error
	Modify(err error, module string, level int, code string) error

	Prepare(module string, level int, code string, message string) error
	New(module string, level int, code string, message string) error
	Raw(message string) error
}

// InterfaceLogger is an interface to system logging service
type InterfaceLogger interface {
	Log(storage string, prefix string, message string)

	LogError(err error)
	LogEvent(f LogFields, eventName string)
}

type LogFields map[string]interface{}

// InterfaceIniConfig is an interface to startup configuration predefined values service
type InterfaceIniConfig interface {
	SetWorkingSection(sectionName string) error
	SetValue(valueName string, value string) error

	GetSectionValue(sectionName string, valueName string, defaultValue string) string
	GetValue(valueName string, defaultValue string) string

	ListSections() []string
	ListItems() []string
	ListSectionItems(sectionName string) []string
}

// InterfaceConfig is an interface to configuration values managing service
type InterfaceConfig interface {
	RegisterItem(Item StructConfigItem, Validator FuncConfigValueValidator) error
	UnregisterItem(Path string) error

	ListPathes() []string
	GetValue(Path string) interface{}
	SetValue(Path string, Value interface{}) error

	GetGroupItems() []StructConfigItem
	GetItemsInfo(Path string) []StructConfigItem

	Load() error
	Reload() error
}

// InterfaceOttemoError is an interface to errors generated by error bus service
type InterfaceOttemoError interface {
	ErrorFull() string
	ErrorLevel() int
	ErrorCode() string
	ErrorMessage() string
	ErrorCallStack() string

	IsHandled() bool
	MarkHandled() bool

	IsLogged() bool
	MarkLogged() bool

	error
}

// FuncConfigValueValidator is a configuration value validator callback function prototype
type FuncConfigValueValidator func(interface{}) (interface{}, error)

// FuncEventListener is an event listener callback function prototype
//   - return value is continue flag, so listener should return false to stop event propagation
type FuncEventListener func(string, map[string]interface{}) bool

// FuncErrorListener is an error listener callback function prototype
type FuncErrorListener func(error) bool

// FuncCronTask is a callback function prototype executes by scheduler
type FuncCronTask func(params map[string]interface{}) error

// StructConfigItem is a structure to hold information about particular configuration value
type StructConfigItem struct {
	Path  string
	Value interface{}

	Type string

	Editor  string
	Options interface{}

	Label       string
	Description string

	Image string
}


type InterfaceScript interface {
	Interact() error
	Execute(code string) (interface{}, error)
	Get(name string) (interface{}, error)
	Set(name string, value interface{}) error
}

type InterfaceScriptEngine interface {
	GetScriptName() string
	GetScriptInstance() InterfaceScript
	Get(name string) (interface{}, error)
	Set(name string, value interface{}) error
}
