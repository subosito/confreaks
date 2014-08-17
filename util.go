package confreaks

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
)

func downloadVideo(u, dir string) error {
	var stderr bytes.Buffer
	var cmd *exec.Cmd

	vu, err := url.Parse(u)
	if err != nil {
		return err
	}

	switch vu.Host {
	case "cdn.confreaks.com":
		c, err := exec.LookPath("wget")
		if err != nil {
			log.Fatal("Unable to find 'wget'. Please install it.")
		}

		out := fmt.Sprintf("%s/%s", dir, path.Base(vu.Path))
		cmd = exec.Command(c, "-N", "-c", "-O", out, u)
	default:
		c, err := exec.LookPath("youtube-dl")
		if err != nil {
			log.Fatal("Unable to find 'youtube-dl'. Please install it. See https://github.com/rg3/youtube-dl")
		}

		out := fmt.Sprintf("%s/%%(title)s-%%(id)s.%%(ext)s", dir)

		if vu.Host == "blip.tv" {
			cmd = exec.Command(c, "--ignore-config", "-o", out, u)
		} else {
			cmd = exec.Command(c, "-o", out, u)
		}
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}
