Unity meta checker
=====================================

A tool to check problems about [meta files](https://docs.unity3d.com/2021.1/Documentation/Manual/AssetWorkflow.html) of [Unity](https://unity.com/) on [Git](https://git-scm.com/) repositories, and also the tool can do limited autofix for meta files of auto-generated files.

This tool can check the following problems:

<dl>
<dt>Missing meta files</dt>
<dd><em>Missing</em> means that an asset should have a meta file but the meta file is not committed. This problem can cause broken asset references.
<dt>Dangling meta files</dt>
<dd><em>Dangling</em> means that a meta file exist, but the asset is not committed. This problem can cause annoying warning messages.
</dl>



Basic Usage
-----------

```console
$ # Change the current directory to your Unity project or UPM package.
$ cd /path/to/unity/project

$ # Execute unity-meta-check (specifying -silent make that unity-meta-check only show results or fatal errors).
$ unity-meta-check -silent
missing: Assets/Not/Added.meta
missing: Packages/com.my.pkg/README.meta
missing: LocalPackages/com.local.pkg/README.meta
...
dangling: Assets/Not/Removed.meta
dangling: Packages/com.my.pkg/MyPkg.csproj.meta
dangling: LocalPackages/com.local.pkg/LocalPkg.csproj.meta
...

$ # unity-meta-check exit with non-zero status if one or more missing/dangling .meta files exist.
$ echo $?
1
```

Typically, unity-meta-check don't need to specify the target type (Unity project or UPM package) because unity-meta-check can automatically detect it.
You can explicitly specify `-unity-project` or `-upm-package` to disable the automatic detection, if the detection result was not intended.

If you want to ignore some problems, you can use `.meta-check-ignore` (this format is very similar to `.gitignore` but `!` is not supported):

```console
$ unity-meta-check -silent
missing: Assets/Not/Added1.meta
missing: Assets/NotAdded2.meta

$ # You can ignore these problems using .meta-check-ignore:
$ cat .meta-check-ignore
Assets/Not             # All files in the directory or the sub directories get ignored.
Assets/NotAdded2.meta  # Also can specify the path to files.

$ # unity-meta-check will ignore these problems.
$ unity-meta-check -silent
```

See [more advanced usage](#Advanced-Usage) for more information.



Installation
------------
### Using Docker Image

This way is recommended to use unity-meta-check on CI.

```console
$ docker pull ghcr.io/dena/unity-meta-check/unity-meta-check:latest

$ cd path/to/your/proj 
$ docker run -v "$(pwd):/target" --rm ghcr.io/dena/unity-meta-check/unity-meta-check:latest -silent /target
missing Assets/AssetsMissing.meta
missing LocalPackages/com.example.local.pkg/LocalPkgMissing.meta
missing Packages/com.example.pkg/PkgMissing.meta
dangling Assets/AssetsDangling.meta
dangling LocalPackages/com.example.local.pkg/LocalPkgDangling.meta
dangling Packages/com.example.pkg/PkgDangling.meta

$ docker run --rm ghcr.io/dena/unity-meta-check/unity-meta-check:latest -help
usage: unity-meta-check [<options>] [<path>]
...
```


### Using GitHub Actions

To check only, the following YAML can cover almost case:

```yaml
name: Meta Check
on: pull_request

jobs:
  meta-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: DeNA/unity-meta-check@v3
```

See [`./action.yml`](./action.yml) for more detials.

<details>
<summary>Advanced Usage for JUnit report + Autofix + PR Comment report</summary>

The following YAML is the example for JUnit report + Autofix + PR Comment report:

```yaml
name: Meta Check
on: pull_request

jobs:
  unity-meta-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: DeNA/unity-meta-check@v3
        with:
          enable_autofix: true
          autofix_globs: .
          enable_junit: true
          junit_xml_path: junit.xml
          enable_pr_comment: true
          pr_comment_lang: ja
          pr_comment_send_success: true
        env:
          GITHUB_TOKEN: "${{ secrets.YOUR_GITHUB_TOKEN }}"

      - name: See how autofix did
        run: git status
        if: always()

      - uses: mikepenz/action-junit-report@v2
        with:
          report_paths: junit.xml
        if: always()
```
</details>

<details>
<summary>Advanced Usage for <code>push</code> events instead of <code>pull_request</code> events</summary>

```yaml
name: Meta Check
on: push

jobs:
  meta-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - uses: jwalton/gh-find-current-pr@v1
        id: findPr

      - uses: DeNA/unity-meta-check@v3
          enable_pr_comment: true
          pr_comment_pull_number: ${{ steps.findPr.outputs.number }}
        env:
          GITHUB_TOKEN: "${{ secrets.YOUR_GITHUB_TOKEN }}"
```
</details>



### Download binaries

Binaries are available on [releases](https://github.com/dena/unity-meta-check/releases).



Advanced Usage
--------------

Provided features are built on several individual binaries:

<dl>
<dt><a href="#unity-meta-check"><code>unity-meta-check</code></a>
<dd>Checker for missing/dangling meta files. The result print to stdout.
<dt><a href="#unity-meta-autofix"><code>unity-meta-check-autofix</code></a>
<dd>Autofix for meta files problems. It need a result of <code>unity-meta-check</code> via stdin.
<dt><a href="#unity-meta-check-junit"><code>unity-meta-check-junit</code></a>
<dd>Reporter for Jenkins compatible XML based JUnit reports. It need a result of <code>unity-meta-check</code> from stdin.
<dt><a href="#unity-meta-check-github-pr-comment"><code>unity-meta-check-github-pr-comment</code></a>
<dd>Reporter for GitHub comments of GitHub issues or pull requests. It need a result of <code>unity-meta-check</code> from stdin.
</dl>



### unity-meta-check

`unity-meta-check` checks missing/dangling meta files on the commit.

This check based on [a git tree object](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects#_tree_objects) instead of the working directory. It means you **MUST** commit adding missing meta files or removing dangling meta files to re-check by `unity-meta-check`.

Other complemental features such as autofix or reporting are provided other binaries.

```console
$ unity-meta-check -help
usage: unity-meta-check [<options>] [<path>]

Check missing or dangling .meta files.

  <path>
        root directory of your Unity project or UPM package to check (default "$(git rev-parse --show-toplevel)")

OPTIONS
  -debug
        set log level to DEBUG (default INFO)
  -ignore-file string
        path to .meta-check-ignore
  -ignore-dangling
        ignore dangling .meta
  -ignore-submodules
        ignore git submodules and nesting repositories (this is RECOMMENDED but not enabled by default because it can cause to miss problems in submodules or nesting repositories)
  -no-ignore-case
        treat case of file paths
  -silent
        set log level to WARN (default INFO)
  -unity-project
        check as Unity project
  -unity-project-sub-dir
        check as sub directory of Unity project
  -upm-package
        check as UPM package (same meaning of -unity-project-sub-dir)
  -version
        print version

EXAMPLE USAGES
  $ cd path/to/UnityProject
  $ unity-meta-check -silent

  $ cd path/to/any/dir
  $ unity-meta-check -silent -upm-package path/to/MyUPMPackage
  $ unity-meta-check -silent -unity-project-sub-dir path/to/UnityProject/Assets/Sub/Dir

EXAMPLE USAGES WITH OTHER TOOLS
  $ unity-meta-check <options> | unity-meta-check-junit path/to/junit.xml
  $ unity-meta-check <options> | unity-meta-check-github-pr-comment <options>
  $ unity-meta-check <options> | unity-meta-check-junit path/to/junit.xml | unity-meta-check-github-pr-comment <options>
```

If both `-silent` and `-debug` are specified, `-silent` win.



### unity-meta-autofix

`unity-meta-autofix` fix (very limited) problems about meta files. It needs a result of `unity-meta-check` via stdin. It can fix the following problems:

<dl>
<dt>Missing meta files for folders
<dd>
<details>
<summary>Example of auto-generated meta files</summary>

```yaml
fileFormatVersion: 2
guid: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
folderAsset: yes
DefaultImporter:
  externalObjects: {}
  userData: 
  assetBundleName: 
  assetBundleVariant: 
```

Automatic GUID generation does not depend on time, so it is safe if autofix runs parallel.
</details>

<dt>Missing meta files for some binaries
<dd>Sometimes you need to import auto-generated binary data files (like encoded as Protocol Buffer Binary Wire Format) programmatically. Then, autofix feature is useful because it can add meta files to the binaries.
<details>
<summary>Example of auto-generated meta files</summary>

```yaml
guid: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TextScriptImporter:
  externalObjects: {}
  userData: 
  assetBundleName: 
  assetBundleVariant: 
```

Automatic GUID generation does not depend on time, so it is safe if autofix runs parallel.
</details>
</dl>

```console
$ unity-meta-autofix -help
usage: unity-meta-autofix [<options>] <pattern> [<pattern>...]

Fix missing or dangling .meta. Currently autofix is only limited support.

ARGUMENTS
  <pattern>
        glob pattern to path where autofix allowed on

OPTIONS
  -debug
        set log level to DEBUG (default INFO)
  -dry-run
        dry run
  -root-dir string
        directory path to where unity-meta-check checked at (default ".")
  -silent
        set log level to WARN (default INFO)
  -version
        print version

EXAMPLE USAGES
  $ unity-meta-check <options> | unity-meta-autofix -dry-run path/to/autofix
  $ unity-meta-check <options> | unity-meta-autofix <options> | <other-unity-meta-check-tool>
```

Currently, autofix for dangling meta files is not supported, because it might be dangerous on some situations.



### unity-meta-check-junit

`unity-meta-check-junit` is a reporter for Jenkins compatible XML based JUnit reports. It needs a result of `unity-meta-check` via stdin.

```console
$ unity-meta-check-junit -help
usage: unity-meta-check-junit [<options>] [<path>]

Save a JUnit report file for the result from unity-meta-check via STDIN.

  <path>
        output path to write JUnit report

OPTIONS
  -version
        print version

EXAMPLE USAGES
  $ unity-meta-check <options> | unity-meta-check-junit path/to/junit-report.xml
  $ unity-meta-check <options> | unity-meta-check-junit path/to/junit-report.xml | <other-unity-meta-check-tool>
```



### unity-meta-check-github-pr-comment

`unity-meta-check-github-pr-comment` is a reporter for GitHub comments of GitHub issues or pull requests. It needs a result of `unity-meta-check` via stdin.

![](./docs/images/unity-meta-check-github-pr-comment-screenshot.png)

```console
$ unity-meta-check-github-pr-comment -help
usage: unity-meta-check-github-pr-comment [<options>]

Post a comment for the result from unity-meta-check via STDIN to GitHub Pull Request.

OPTIONS
  -api-endpoint string
        GitHub API endpoint URL (like https://github.example.com/api/v3) (default "https://api.github.com")
  -debug
        set log level to DEBUG (default INFO)
  -lang string
        language code (available: en, ja) (default "en")
  -owner string
        owner of the GitHub repository
  -pull uint
        pull request number
  -repo string
        name of the GitHub repository
  -silent
        set log level to WARN (default INFO)
  -template-file string
        custom template file
  -version
        print version

ENVIRONMENT
  GITHUB_TOKEN
        GitHub API token. The scope can be empty if your repository is public. Otherwise, the scope should contain "repo"

EXAMPLE USAGES
  $ export GITHUB_TOKEN="********"
  $ unity-meta-check <options> | unity-meta-check-github-pr-comment \
      -api-endpoint https://api.github.com \
      -owner example-org \
      -repo my-repo \
      -pull "$CIRCLE_PR_NUMBER"  # This is for CircleCI

  $ export GITHUB_TOKEN="********"  # This should be set via credentials().
  $ unity-meta-check <options> | unity-meta-check-github-pr-comment \
      -api-endpoint https://github.example.com/api/v3 \
      -owner example-org \
      -repo my-repo \
      -pull "$ghprbPullId"  # This is for Jenkins with GitHub PullRequest Builder plugin

  $ GITHUB_TOKEN="********" unity-meta-check <options> | unity-meta-check-junit path/to/unity-meta-check-result.xml | unity-meta-check-github-pr-comment <options> | <other-unity-meta-check-tool>

  $ export GITHUB_TOKEN="********"  # This should be set via credentials().
  $ unity-meta-check <options> | unity-meta-check-github-pr-comment \
      -api-endpoint https://github.example.com/api/v3 \
      -owner example-org \
      -repo my-repo \
      -pull "$ghprbPullId" \
      -template-file path/to/template.json  # template file can be used for localization for GitHub comments.

TEMPLATE FILE FORMAT EXAMPLE
  If a template file is like:

  {
    "success": "No missing/dangling .meta found. Perfect!",
    "failure": "Some missing or dangling .meta found. Fix commits are needed.",
    "header_status": "Status",
    "header_file_path": "File",
    "status_missing": "Not committed",
    "status_dangling": "Not removed"
  }

  then the output become:

  No missing/dangling .meta found. Perfect!

  or:

  Some missing or dangling .meta found. Fix commits are needed.

  | Status | File |
  |:--|:--|
  | Not committed | `path/to/missing.meta` |
  | Not removed | `path/to/dangling.meta` |
```
