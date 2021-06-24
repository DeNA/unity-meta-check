package yaml

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

const tagNamePrefix = "# image-name: "

func TestImage(t *testing.T) {
	expected, err := getActualImageName("../../../../.github/images/Dockerfile")
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	if ImageName != expected {
		t.Errorf("want %q, got %q", expected, ImageName)
	}
}

func getActualImageName(dockerfilePath string) (DockerImage, error) {
	dockerfile, err := os.OpenFile(dockerfilePath, os.O_RDONLY, 0)
	if err != nil {
		return "", err
	}
	defer dockerfile.Close()

	scanner := bufio.NewScanner(dockerfile)
	if !scanner.Scan() {
		return "", errors.New("want read a first line, got no lines")
	}

	firstLine := scanner.Text()

	if !strings.HasPrefix(firstLine, tagNamePrefix) {
		return "", fmt.Errorf("want a tag name prefix comment at first line, got no prefix: %q", firstLine)
	}

	return DockerImage(strings.TrimPrefix(firstLine, tagNamePrefix)), nil
}
