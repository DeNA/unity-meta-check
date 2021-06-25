package checker

type TargetType string

const (
	// TargetTypeIsUnityProjectRootDirectory means the target root directory point to a root directory of the Unity project.
	TargetTypeIsUnityProjectRootDirectory TargetType = "TargetTypeIsUnityProjectRootDirectory"

	// TargetTypeIsUnityProjectSubDirectory means the target root directory point to a sub directory that need meta files of the Unity project.
	// NOTE: this including UPM packages (because Packages/com.example.foo/ is a sub directory that need meta files of Unity projects).
	TargetTypeIsUnityProjectSubDirectory TargetType = "TargetTypeIsUnityProjectSubDirectory"
)

// Options for Checker.
// NOTE: We should keep this options simple because the newWorker logic is already complicated.
//       Implement on ResultFilter if you want to filter the CheckResult.
type Options struct {
	IgnoreCase                bool
	IgnoreSubmodulesAndNested bool
	TargetType                TargetType
}