package staging

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/pauldotknopf/darch/pkg/workspace"

	"github.com/pauldotknopf/darch/pkg/hooks"
	"github.com/pauldotknopf/darch/pkg/reference"
	"github.com/pauldotknopf/darch/pkg/utils"
)

// RunAllHooks Run the hooks on every image.
func (session *Session) RunAllHooks() error {
	allAssoications, err := session.imageStore.AllImages()
	if err != nil {
		return err
	}
	allHooks, err := hooks.GetHooks()
	if err != nil {
		return err
	}
	for _, association := range allAssoications {
		err = session.runHooksForAssociation(association, allHooks)
		if err != nil {
			return err
		}
	}
	return nil
}

// RunHooksForImage Run hooks for a single image.
func (session *Session) RunHooksForImage(imageRef reference.ImageRef) error {
	association, err := session.imageStore.Get(imageRef)
	if err != nil {
		return err
	}
	allHooks, err := hooks.GetHooks()
	if err != nil {
		return err
	}
	return session.runHooksForAssociation(association, allHooks)
}

func (session *Session) runHooksForAssociation(association reference.Association, hs []hooks.Hook) error {

	ws, err := workspace.NewWorkspace(DefaultStagingDirectoryTmp)
	if err != nil {
		return err
	}
	defer ws.Destroy()

	for _, hook := range hs {

		if !hooks.AppliesToImage(hook, association.Ref) {
			continue
		}

		var destinationHookDirectory = path.Join(ws.Path, "hooks", hook.NameWithOrder)

		err = os.MkdirAll(destinationHookDirectory, os.ModePerm)
		if err != nil {
			return err
		}

		hookFile := path.Join(hook.Path, "hook")
		if !utils.FileExists(hookFile) {
			return fmt.Errorf("Hook script %s doesn't exist", hookFile)
		}

		err = utils.CopyFile(hookFile, path.Join(destinationHookDirectory, "hook"))
		if err != nil {
			return err
		}

		cmd := exec.Command("/bin/bash", "-c", ". "+path.Join(destinationHookDirectory, "hook")+" && install")
		cmd.Env = append(os.Environ(), []string{
			fmt.Sprintf("DARCH_HOOKS_DIR=%s", hook.HooksPath),
			fmt.Sprintf("DARCH_HOOK_NAME=%s", hook.Name),
			fmt.Sprintf("DARCH_HOOK_DIR=%s", hook.Path),
			fmt.Sprintf("DARCH_HOOK_DEST_DIR=%s", destinationHookDirectory),
			fmt.Sprintf("DARCH_IMAGE_NAME=%s", association.Ref.FullName()),
		}...)
		cmd.Dir = destinationHookDirectory
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	// Now that we have ran the hooks, let's delete any hooks that may be in our image,
	// and move our new hooks into it.
	currentHooksDir := path.Join(DefaultStagingDirectoryImages, association.ID, "hooks")
	if utils.DirectoryExists(currentHooksDir) {
		err = os.RemoveAll(currentHooksDir)
		if err != nil {
			return err
		}
	}

	err = os.Rename(path.Join(ws.Path, "hooks"), currentHooksDir)
	if err != nil {
		return err
	}

	return nil
}
