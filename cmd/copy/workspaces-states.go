package copy

import (
	"crypto/md5"
	b64 "encoding/base64"
	"fmt"

	"github.com/hashicorp-services/tfm/cmd/helper"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

// 1. Get source state versions per workspace
// 2. Get dest state versions per workspace
// 3. Compare state serials between 2 workspaces
// 4. Get the WS ID of the WS Name to copy state too
// 5. Get the download URL of the source state
// 6. Download the State into memory
// 7. Create MD5 checksum
// 8. Incremenet the serial #
// 9. Lock the workspace if not locked
// 10. Use the StateVersions.Create to upload state to destination
// 11. Unlock the workspace if locked

// Get the source workspace state files from the provided workspace
func discoverSrcStates(c tfclient.ClientContexts, ws string) ([]*tfe.StateVersion, error) {
	o.AddMessageUserProvided("Getting list of states from source workspace ", ws)
	srcStates := []*tfe.StateVersion{}

	opts := tfe.StateVersionListOptions{
		ListOptions:  tfe.ListOptions{PageNumber: 1, PageSize: 100},
		Organization: c.SourceOrganizationName,
		Workspace:    ws,
	}
	for {
		items, err := c.SourceClient.StateVersions.List(c.SourceContext, &opts)
		if err != nil {
			return nil, err
		}

		srcStates = append(srcStates, items.Items...)

		o.AddFormattedMessageCalculated("Found %d Workspace states", len(srcStates))

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage

	}

	return srcStates, nil
}

// Get the destination workspace state files from the provided workspace
func discoverDestStates(c tfclient.ClientContexts, ws string) ([]*tfe.StateVersion, error) {
	o.AddMessageUserProvided("Getting list of States from destination Workspace ", ws)
	destStates := []*tfe.StateVersion{}

	opts := tfe.StateVersionListOptions{
		ListOptions:  tfe.ListOptions{PageNumber: 1, PageSize: 100},
		Organization: c.DestinationOrganizationName,
		Workspace:    ws,
	}
	for {
		items, err := c.DestinationClient.StateVersions.List(c.DestinationContext, &opts)
		if err != nil {
			return nil, err
		}

		destStates = append(destStates, items.Items...)

		o.AddFormattedMessageCalculated("Found %d Workspace states", len(destStates))

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage

	}

	return destStates, nil
}

// Check the existence of the state in the destination using the unique serial number
func doesStateExist(stateSerial int64, s []*tfe.StateVersion) bool {
	for _, state := range s {
		if stateSerial == state.Serial {
			return true
		}
	}
	return false
}

// Finds the destination Workspace ID of the workspace with a matching name as a workspace in the source
func getWorkspaceId(c tfclient.ClientContexts, ws string) (string, error) {
	w, err := c.DestinationClient.Workspaces.Read(c.DestinationContext, c.DestinationOrganizationName, ws)

	if err != nil {
		return "", err
	}

	return w.ID, nil
}

func downloadSourceState(c tfclient.ClientContexts, downloadUrl string) ([]byte, error) {

	// Takes download URL from StateVersions.List function and stores state as a []byte type
	state, err := c.SourceClient.StateVersions.Download(c.SourceContext, downloadUrl)
	if err != nil {
		return state, err
	}

	return state, nil
}

// Locks the workspace provided
func lockWorkspace(c tfclient.ClientContexts, destWorkspaceId string) error {
	message := "Uploading State"

	wsProperties, err := c.DestinationClient.Workspaces.ReadByID(c.DestinationContext, destWorkspaceId)
	if err != nil {
		return err
	}

	if wsProperties.Locked == false {
		c.DestinationClient.Workspaces.Lock(c.DestinationContext, destWorkspaceId, tfe.WorkspaceLockOptions{
			Reason: &message,
		})
		fmt.Println("Locking Workspace: ", destWorkspaceId)
	}
	return nil
}

// Unlocks the workspace provided
func unlockWorkspace(c tfclient.ClientContexts, destWorkspaceId string) error {
	wsProperties, err := c.DestinationClient.Workspaces.ReadByID(c.DestinationContext, destWorkspaceId)
	if err != nil {
		return err
	}

	if wsProperties.Locked == true {
		c.DestinationClient.Workspaces.Unlock(c.DestinationContext, destWorkspaceId)
		fmt.Println("Unlocking Workspace: ", destWorkspaceId)
	}
	return nil
}

// Main function for `--state` flag
func copyStates(c tfclient.ClientContexts) error {

	// Get the source target workspaces
	srcWorkspaces, err := getSrcWorkspacesCfg(c)
	if err != nil {
		return errors.Wrap(err, "failed to list Workspaces from source")
	}

	// Get/Check if Workspace map exists
	wsMapCfg, err := helper.ViperStringSliceMap("workspace-map")
	if err != nil {
		fmt.Println("invalid input for workspace-map")
	}

	// Get the destination target workspaces
	destWorkspaces, err := discoverDestWorkspaces(tfclient.GetClientContexts(), true)
	if err != nil {
		return errors.Wrap(err, "failed to list Workspaces from source")
	}

	for _, srcworkspace := range srcWorkspaces {
		destWorkSpaceName := srcworkspace.Name

		// Check if the destination Workspace name differs from the source name
		if len(wsMapCfg) > 0 {
			destWorkSpaceName = wsMapCfg[srcworkspace.Name]
		}

		// Check for the exsitence of the destination workspace in the destination target
		exists := doesWorkspaceExist(destWorkSpaceName, destWorkspaces)

		if exists {
			destWorkspaceId, err := getWorkspaceId(tfclient.GetClientContexts(), destWorkSpaceName)
			if err != nil {
				return errors.Wrap(err, "Failed to get the ID of the destination Workspace that matches the Name of the Source Workspace")
			}

			fmt.Printf("Source ws %v has a matching ws %v in destination with ID %v. Comparing existing States...\n", srcworkspace.Name, destWorkSpaceName, destWorkspaceId)

			// Get the source workspace states
			srcStates, err := discoverSrcStates(tfclient.GetClientContexts(), srcworkspace.Name)
			if err != nil {
				return errors.Wrap(err, "failed to list state files for workspace from source")
			}

			// Get the destination workspace states
			destStates, err := discoverDestStates(tfclient.GetClientContexts(), destWorkSpaceName)
			if err != nil {
				return errors.Wrap(err, "failed to list state files for workspace from destination")
			}

			// Loop each state for each source workspace with a matching workspace name in the destination,
			// check for the existence of that states serial in destination, upload state if serial doesnt exist
			for _, srcstate := range srcStates {
				exists := doesStateExist(srcstate.Serial, destStates)
				if exists {
					fmt.Printf("State Version %v with Serial %v exists in destination will not migrate\n", srcstate.StateVersion, srcstate.Serial)
				} else {

					// Download state from source
					state, err := downloadSourceState(tfclient.GetClientContexts(), srcstate.DownloadURL)

					// Base64 encode the state as a string
					stringState := b64.StdEncoding.EncodeToString(state)

					// Get the MD5 hash of the state
					md5String := fmt.Sprintf("%x", md5.Sum([]byte(state)))

					// Create an empty int
					newSerial := int64(0)

					// Get properties of the state
					currentState, _ := c.DestinationClient.StateVersions.ReadCurrent(c.DestinationContext, destWorkspaceId)

					// Of if there is a state file, set the newSerial variable to the current state serial + 1
					if currentState != nil {
						newSerial = currentState.Serial + 1
					}

					// Lock the destination workspace
					lockWorkspace(tfclient.GetClientContexts(), destWorkspaceId)
					fmt.Printf("Migrating state version %v serial %v for workspace Src: %v Dst: %v\n", srcstate.StateVersion, newSerial, srcworkspace.Name, destWorkSpaceName)
					srcstate, err := c.DestinationClient.StateVersions.Create(c.DestinationContext, destWorkspaceId, tfe.StateVersionCreateOptions{
						Type:             "",
						Lineage:          new(string),
						MD5:              tfe.String(md5String),
						Serial:           &newSerial,
						State:            tfe.String(stringState),
						Force:            new(bool),
						Run:              &tfe.Run{},
						JSONState:        new(string),
						JSONStateOutputs: new(string),
					})

					if err != nil {
						return err
					}

					o.AddDeferredMessageRead("Migrated State Serial # ", srcstate.Serial)
				}
			}
			unlockWorkspace(tfclient.GetClientContexts(), destWorkspaceId)
		} else {
			fmt.Printf("Source workspace (%v) does not exist in destination (%v). No states to migrate\n", srcworkspace.Name, destWorkSpaceName)
		}
	}
	return nil
}
