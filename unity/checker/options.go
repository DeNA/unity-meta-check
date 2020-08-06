package checker

type TargetType string

const (
	TargetTypeIsUnityProjectRootDirectory TargetType = "TargetTypeIsUnityProjectRootDirectory"
	// NOTE: Unity proj sub dirs including UPM packages.
	TargetTypeIsUnityProjectSubDirectory TargetType = "TargetTypeIsUnityProjectSubDirectory"
)

// NOTE: We should keep this options simple because the newWorker logic is already complicated.
//       Implement on ResultFilter if you want to filter the CheckResult.
type Options struct {
	IgnoreCase                bool
	IgnoreSubmodulesAndNested bool
	TargetType                TargetType
}
