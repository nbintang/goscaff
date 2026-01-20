package scaffold


type DBType string

const (
	DBTypePostgres DBType = "postgres"
	DBTypeMySQL    DBType = "mysql"
)

type PresetType string

const (
	PresetBase   PresetType = "base"
	PresetFull   PresetType = "full"
)

type Overlay struct {
	Src string  
	Dst string 
}