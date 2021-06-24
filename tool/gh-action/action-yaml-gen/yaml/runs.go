package yaml

import (
	"fmt"
)

const ImageName = "ghcr.io/dena/unity-meta-check/unity-meta-check-gh-action"

func ImageWithTag(version string) DockerImage {
	return DockerImage(fmt.Sprintf("%s:%s", ImageName, version))
}

type (
	DockerImage   string
	BrandingIcon  string
	BrandingColor string

	DockerAction struct {
		// Image is the Docker image to use as the container to run the action. The value can be the Docker base image name,
		// a local Dockerfile in your repository, or a public image in Docker Hub or another registry.
		// To reference a Dockerfile local to your repository, the file must be named Dockerfile and you must use a path
		// relative to your action metadata file. The docker application will execute this file.
		// SEE: https://docs.github.com/en/actions/creating-actions/metadata-syntax-for-github-actions#runsimage
		Image DockerImage

		// Env specifies a key/value map of environment variables to set in the container environment.
		// SEE: https://docs.github.com/en/actions/creating-actions/metadata-syntax-for-github-actions#runsenv
		Env map[string]string

		// Args is an array of strings that define the inputs for a Docker container. Inputs can include hardcoded strings.
		// GitHub passes the args to the container's ENTRYPOINT when the container starts up.
		// The args are used in place of the CMD instruction in a Dockerfile. If you use CMD in your Dockerfile, use the guidelines ordered by preference:
		//   1. Document required arguments in the action's README and omit them from the CMD instruction.
		//   2. Use defaults that allow using the action without specifying any args.
		//   3. If the action exposes a --help flag, or something similar, use that to make your action self-documenting.
		// If you need to pass environment variables into an action, make sure your action runs a command shell to perform variable substitution.
		// For example, if your entrypoint attribute is set to "sh -c", args will be run in a command shell.
		// Alternatively, if your Dockerfile uses an ENTRYPOINT to run the same command ("sh -c"),
		// args will execute in a command shell.
		//
		// For more information about using the CMD instruction with GitHub Actions, see "Dockerfile support for GitHub Actions."
		// SEE: https://docs.github.com/en/actions/creating-actions/metadata-syntax-for-github-actions#runsargs
		Args []string
	}
)
