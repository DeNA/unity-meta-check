package runner

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/options"
	"github.com/DeNA/unity-meta-check/resultfilter"
	"github.com/DeNA/unity-meta-check/tool/gh-action/inputs"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-autofix/autofix"
	prcomment "github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/github"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/l10n"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

type Options struct {
	RootDirAbs      typedpath.RawPath
	CheckerOpts     *checker.Options
	FilterOpts      *resultfilter.Options
	EnableJUnit     bool
	JUnitOutPath    typedpath.RawPath
	EnablePRComment bool
	PRCommentOpts   *prcomment.Options
	EnableAutofix   bool
	AutofixOpts     *autofix.Options
}

type Validator func(rootDirAbs typedpath.RawPath, i inputs.Inputs, token prcomment.Token) (*Options, error)

func NewValidateFunc(
	detectTargetType options.UnityProjectDetector,
	buildIgnoredGlobs options.IgnoredGlobsBuilder,
	buildAutofixOpts autofix.OptionsBuilder,
	readTmplFile l10n.TemplateFileReader,
) Validator {
	return func(rootDirAbs typedpath.RawPath, i inputs.Inputs, token prcomment.Token) (*Options, error) {
		inputTargetType, err := inputs.ValidateTargetType(i.TargetType)
		if err != nil {
			return nil, err
		}

		var targetType checker.TargetType
		switch inputTargetType {
		case inputs.TargetTypeUnityProj:
			targetType = checker.TargetTypeIsUnityProjectRootDirectory
		case inputs.TargetTypeUnityProjSubDir, inputs.TargetTypeUpmPackage:
			targetType = checker.TargetTypeIsUnityProjectSubDirectory
		case inputs.TargetTypeAuto:
			detected, err := detectTargetType(rootDirAbs)
			if err != nil {
				return nil, err
			}
			targetType = detected
		default:
			return nil, fmt.Errorf("unknown target type: %q", inputTargetType)
		}

		checkerOpts := &checker.Options{
			IgnoreCase:                bool(i.IgnoreCase),
			IgnoreSubmodulesAndNested: bool(i.IgnoreSubmodulesAndNested),
			TargetType:                targetType,
		}

		ignoredGlobs, err := buildIgnoredGlobs(i.IgnoredFilePath, rootDirAbs)
		if err != nil {
			return nil, err
		}

		filterOpts := &resultfilter.Options{
			IgnoreDangling: bool(i.IgnoreDangling),
			IgnoredGlobs:   ignoredGlobs,
			IgnoreCase:     bool(i.IgnoreCase),
		}

		var junitPath typedpath.RawPath
		if i.EnableJUnit {
			junitPath = i.JUnitXMLPath
		}

		var prcommentOps *prcomment.Options
		if i.EnablePRComment {
			var tmpl *l10n.Template
			if i.PRCommentTmplFilePath == "" {
				tmpl, err = l10n.GetTemplate(l10n.Lang(i.PRCommentLang))
				if err != nil {
					return nil, err
				}
			} else {
				tmpl, err = readTmplFile(i.PRCommentTmplFilePath)
				if err != nil {
					return nil, err
				}
				if err := l10n.ValidateTemplate(tmpl); err != nil {
					return nil, err
				}
			}

			owner, err := prcomment.ValidateOwner(i.PRCommentOwner)
			if err != nil {
				return nil, err
			}

			repo, err := prcomment.ValidateRepo(i.PRCommentRepo)
			if err != nil {
				return nil, err
			}

			pullNumber, err := prcomment.ValidatePullNumber(int(i.PRCommentPRNumber))
			if err != nil {
				return nil, err
			}

			apiEndpoint, err := prcomment.ValidateAPIEndpoint(i.PRCommentAPIEndpoint)
			if err != nil {
				return nil, err
			}

			prcommentOps = &prcomment.Options{
				Tmpl:          tmpl,
				SendIfSuccess: bool(i.PRCommentSendSuccess),
				Token:         token,
				APIEndpoint:   apiEndpoint,
				Owner:         owner,
				Repo:          repo,
				PullNumber:    pullNumber,
			}
		}

		var autofixOpts *autofix.Options
		if i.EnableAutofix {
			allowedGlobs := make([]globs.Glob, len(i.AutofixGlobs))
			for i, s := range i.AutofixGlobs {
				allowedGlobs[i] = globs.Glob(s)
			}
			autofixOpts, err = buildAutofixOpts(rootDirAbs, allowedGlobs)
			if err != nil {
				return nil, err
			}
		}

		return &Options{
			RootDirAbs:      rootDirAbs,
			CheckerOpts:     checkerOpts,
			FilterOpts:      filterOpts,
			EnableJUnit:     bool(i.EnableJUnit),
			JUnitOutPath:    junitPath,
			EnablePRComment: bool(i.EnablePRComment),
			PRCommentOpts:   prcommentOps,
			EnableAutofix:   bool(i.EnableAutofix),
			AutofixOpts:     autofixOpts,
		}, nil
	}
}
