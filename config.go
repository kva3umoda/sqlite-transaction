package sqlite_transaction

type Config struct {
	DBName         string
	AutoVacuum     AutoVacuum
	CacheSizePages int
	JournalMode    JournalMode
	PageSizeBytes  int
	Synchronous    Synchronous
	TempStore      TempStore
	Logger         Logger
}

type Synchronous string

const (
	SynchronousNormal Synchronous = "NORMAL"
	SynchronousFull   Synchronous = "FULL"
	SynchronousExtra  Synchronous = "EXTRA"
	SynchronousOff    Synchronous = "OFF"
)

type AutoVacuum string

const (
	AutoVacuumNone        AutoVacuum = "NONE"
	AutoVacuumFull        AutoVacuum = "FULL"
	AutoVacuumIncremental AutoVacuum = "INCREMENTAL"
)

type JournalMode string

const (
	JournalModeDelete   JournalMode = "DELETE"
	JournalModeTruncate JournalMode = "TRUNCATE"
	JournalModePersist  JournalMode = "PERSIST"
	JournalModeMemory   JournalMode = "MEMORY"
	JournalModeWAL      JournalMode = "WAL"
	JournalModeOFF      JournalMode = "OFF"
)

type TempStore string

const (
	TempStoreDefault TempStore = "DEFAULT"
	TempStoreFile    TempStore = "FILE"
	TempStoreMemory  TempStore = "MEMORY"
)
