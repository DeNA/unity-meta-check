package unity

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestNewMetaNecessityInUnityProject(t *testing.T) {
	cases := []struct {
		path     typedpath.SlashPath
		expected bool
	}{
		{"Assets/Foo", true},
		{"Assets/Foo/Bar", true},
		{"Packages/com.my.pkg/README.md", true},
		// https://docs.unity3d.com/2020.2/Documentation/Manual/cus-layout.html
		{"LocalPackages/com.my.local.pkg/README.md", true},
		// https://docs.unity3d.com/ja/2023.1/Manual/SpecialFolders.html
		{"Assets/Plugins/Foo.xcframework", true},
		{"Assets/Plugins/Foo.androidlib", true},
		{"Assets/Plugins/Foo.androidpack", true},
		{"Assets/Plugins/Foo.bundle", true},
		{"Assets/Plugins/Foo.framework", true},

		{"", false},
		{".git", false},
		{"somethingUserDir", false},
		{"Library", false},
		{"Assets", false},
		{"Assets/", false},
		{"Assets/Foo.meta", false},
		{"Assets/.hidden", false},
		{"Assets/hidden~", false},
		{"Assets/hidden.tmp", false},
		{"Assets/.hidden/file", false},
		{"Assets/Dir/.hidden", false},
		{"Packages", false},
		{"Packages/manifest.json", false},
		{"Packages/com.my.pkg", false},
		{"Packages/com.my.pkg/Documentation~", false},
		{"LocalPackages", false},
		{"LocalPackages/com.my.local.pkg", false},
		{"LocalPackages/com.my.local.pkg/Documentation~", false},

		// https://docs.unity3d.com/ja/2023.1/Manual/SpecialFolders.html
		{"Assets/Plugins/Foo.bundle/Contents/Info.plist", false},
		{"Assets/Plugins/Foo.framework/Versions/A/Foo", false},
		{"Assets/Plugins/Foo.xcframework/Info.plist", false},
		{"Assets/Plugins/Foo.xcframework/ios-armv7_arm64/Foo.xcframework/Info.plist", false},
		{"Assets/Plugins/Foo.androidlib/AndroidManifest.xml", false},
		{"Assets/Plugins/Foo.androidpack/src/main/assets/Bar.bundle", false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%q -> %t", c.path, c.expected), func(t *testing.T) {
			localPkgPaths := []typedpath.SlashPath{
				"Packages/com.my.pkg",
				"LocalPackages/com.my.local.pkg",
			}
			requiresMeta := NewMetaNecessityInUnityProject(localPkgPaths)

			actual := requiresMeta(c.path)

			if actual != c.expected {
				t.Errorf("want %t, got %t", c.expected, actual)
				return
			}
		})
	}
}

func TestNewMetaNecessityInUnityProjectSubDir(t *testing.T) {
	cases := []struct {
		path     typedpath.SlashPath
		expected bool
	}{
		// https://docs.unity3d.com/2019.2/Documentation/Manual/cus-layout.html
		{"README.md", true},
		{"Runtime", true},
		{"Runtime/Unity.MyPkg.asmdef", true},
		{"Runtime/RuntimeExample.cs", true},
		{"Editor", true},
		{"Editor/Unity.MyPkg.Editor.asmdef", true},
		{"Editor/EditorExample.cs", true},
		{"lib/native.a", true},
		{"Tests", true},
		{"Tests/Runtime", true},
		{"Tests/Runtime/Unity.MyPkg.Tests.asmdef", true},
		{"Tests/Runtime/RuntimeExampleTest.cs", true},
		{"Tests/Editor", true},
		{"Tests/Editor/Unity.MyPkg.Editor.Tests.asmdef", true},
		{"Tests/Editor/EditorExampleTest.cs", true},
		{"lib/native.a", true},

		// https://docs.unity3d.com/ja/2023.1/Manual/SpecialFolders.html
		{"Runtime/Foo.bundle", true},
		{"Runtime/Foo.framework", true},
		{"Runtime/Foo.xcframework", true},
		{"Runtime/Foo.androidlib", true},
		{"Runtime/Foo.androidpack", true},

		{"", false},
		{".git", false},
		{"Runtime.meta", false},
		{"Documentation~", false},
		{"Documentation~/com.my.pkg.md", false},

		// https://docs.unity3d.com/ja/2023.1/Manual/SpecialFolders.html
		{"Runtime/Foo.bundle/Contents/Info.plist", false},
		{"Runtime/Foo.framework/Versions/A/Foo", false},
		{"Runtime/Foo.xcframework/Info.plist", false},
		{"Runtime/Foo.xcframework/ios-armv7_arm64/Foo.xcframework/Info.plist", false},
		{"Runtime/Foo.androidlib/AndroidManifest.xml", false},
		{"Runtime/Foo.androidpack/src/main/assets/Bar.bundle", false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%q -> %t", c.path, c.expected), func(t *testing.T) {
			requiresMeta := NewMetaNecessityInUnityProjectSubDir()

			actual := requiresMeta(c.path)

			if actual != c.expected {
				t.Errorf("want %t, got %t", c.expected, actual)
				return
			}
		})
	}
}

func TestTrimMeta(t *testing.T) {
	cases := []struct {
		SlashPath typedpath.SlashPath
		Expected  typedpath.SlashPath
	}{
		{
			SlashPath: "path/to/foo.meta",
			Expected:  "path/to/foo",
		},
		{
			SlashPath: "path/to/test.meta",
			Expected:  "path/to/test",
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v -> %v", c.SlashPath, c.Expected), func(t *testing.T) {
			actual := TrimMetaFromSlash(c.SlashPath)

			if actual != c.Expected {
				t.Error(cmp.Diff(c.Expected, actual))
				return
			}
		})
	}
}
