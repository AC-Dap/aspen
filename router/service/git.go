package service

import (
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

type Repo struct {
	folder string
	remote string

	// The specific commit we want to clone from remote
	commitHash string
}

func NewRepo(folder string, remote string, commitHash string) Repo {
	return Repo{
		folder,
		remote,
		commitHash,
	}
}

func (gr Repo) runGitCommand(arg ...string) ([]byte, error) {
	cmd := exec.Command("git", arg...)
	cmd.Dir = gr.folder
	return cmd.Output()
}

func (gr Repo) Clone() error {
	log.Info().Str("folder", gr.folder).Str("remote", gr.remote).Str("commit", gr.commitHash).Msg("Cloning git repo")

	// Make folder if it doesn't exist
	if info, err := os.Stat(gr.folder); err != nil || !info.IsDir() {
		log.Debug().Str("folder", gr.folder).Msg("Creating folder for git repo")
		err := os.MkdirAll(gr.folder, 0755)
		if err != nil {
			return err
		}
	}

	// Init git repo
	if _, err := gr.runGitCommand("init", "."); err != nil {
		return err
	}

	// Add remote to repo
	if _, err := gr.runGitCommand("remote", "add", "origin", gr.remote); err != nil {
		return err
	}

	// Fetch the specific commit
	if _, err := gr.runGitCommand("fetch", "origin", gr.commitHash, "--depth", "1"); err != nil {
		return err
	}

	// Checkout the specific commit
	if _, err := gr.runGitCommand("reset", "--hard", "FETCH_HEAD"); err != nil {
		return err
	}

	return nil
}

// Updated checks if the repo is synced with remote.
func (gr Repo) Updated() bool {
	// Check if folder exists
	if info, err := os.Stat(gr.folder); err != nil || !info.IsDir() {
		return false
	}

	// Check if git is initialized
	if output, err := gr.runGitCommand("rev-parse", "--is-inside-work-tree"); err != nil || string(output) != "true\n" {
		return false
	}

	// Check if remote is set correctly
	if output, err := gr.runGitCommand("remote", "get-url", "origin"); err != nil || string(output) != gr.remote+"\n" {
		return false
	}

	// Check if current commit matches the commitHash
	if output, err := gr.runGitCommand("rev-parse", "HEAD"); err != nil || string(output) != gr.commitHash+"\n" {
		return false
	}

	return true
}
