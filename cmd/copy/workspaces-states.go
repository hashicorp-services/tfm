package copy

import (
	"fmt"
	"os"

	"crypto/md5"

	"github.com/hashicorp-services/tfe-mig/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

// 1. Get source state versions per workspace
// 2. Get dest state versions per workspace
// 3. Compare state serials between 2 workspaces
// 4. Get the WS ID of the WS Name to copy state too
// 5. Get the download URL of the source state
// 6. Download the State
// 7. Create MD5 checksum
// 8. Use the StateVersions.Create to upload state tot destination
func discoverSrcStates(c tfclient.ClientContexts, ws string) ([]*tfe.StateVersion, error) {
	o.AddMessageUserProvided("Getting list of workspaces states: ", c.SourceHostname)
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

		o.AddFormattedMessageCalculated("Found %d Workspaces states", len(srcStates))

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage

	}

	return srcStates, nil
}

func discoverDestStates(c tfclient.ClientContexts, ws string) ([]*tfe.StateVersion, error) {
	o.AddMessageUserProvided("Getting list of workspaces states: ", c.DestinationHostname)
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

		o.AddFormattedMessageCalculated("Found %d Workspaces states", len(destStates))

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
func findDestWorkspaceId() (string, error) {
	// Get source workspace details
	srcWorkspaces, err := discoverSrcWorkspaces(tfclient.GetClientContexts())

	// Get Dest Workspace details
	destWorkspaces, err := discoverDestWorkspaces(tfclient.GetClientContexts())

	var srcWsName string

	for _, srcworkspace := range srcWorkspaces {
		srcWsName = srcworkspace.Name
	}

	var stateVersionsWsId string
	var destWsName string
	var destWsId string

	for _, destworkspace := range destWorkspaces {
		destWsName = destworkspace.Name
		destWsId = destworkspace.ID

		if srcWsName == destWsName {
			stateVersionsWsId = destWsId
		}

	}
	return stateVersionsWsId, err
}

func downloadSourceState(c tfclient.ClientContexts, downloadUrl string) ([]byte, error) {
	o.AddMessageUserProvided("Creating temp dir to download states from: ", c.SourceHostname)

	o.AddMessageUserProvided("Downloading State file to local host from: ", c.SourceHostname)

	// Takes download URL from StateVersions.List function and stores state as a []byte type
	state, err := c.SourceClient.StateVersions.Download(c.SourceContext, downloadUrl)
	if err != nil {
		return state, err
	}

	// Need to generate the file name based on workspacename
	//
	//

	// Need to store state files in /tmp for unix or /temp for windows
	//
	//

	// Writes the state to a provided filename. If file does not exist, os.WriteFile will create it
	if err := os.WriteFile("/Users/joshuatracy/temp-git-edits/go/tfe-migrate/file", state, 0644); err != nil {
		panic(err)
	}

	return state, nil
	//defer os.RemoveAll(dir)
}

// Calculate MD5 from the state
func calculateMD5(state []byte) []byte {
	data := state
	md5 := md5.Sum(data)

	fmt.Println("MD5 Calcualted:", md5)

	return md5[:]

}

func copyStates(c tfclient.ClientContexts) error {
	// Get the source workspaces properties
	srcWorkspaces, err := discoverSrcWorkspaces(tfclient.GetClientContexts())
	if err != nil {
		return errors.Wrap(err, "failed to list Workspaces from source")
	}

	for _, srcworkspace := range srcWorkspaces {

		destWorkspaceId, err := findDestWorkspaceId()
		if err != nil {
			return errors.Wrap(err, "Failed to get the ID of the destination Workspace that matches the Name of the Source Workspace")
		}

		// Get the source teams properties
		srcStates, err := discoverSrcStates(tfclient.GetClientContexts(), srcworkspace.Name)
		if err != nil {
			return errors.Wrap(err, "failed to list state files for workspace from source")
		}

		// Get the destination teams properties
		destStates, err := discoverDestStates(tfclient.GetClientContexts(), srcworkspace.Name)
		if err != nil {
			return errors.Wrap(err, "failed to list state files for workspace from destination")
		}

		// Loop each team in the srcTeams slice, check for the team existence in the destination,
		// and if a team exists in the destination, then do nothing, else create team in destination.
		for _, srcstate := range srcStates {
			exists := doesStateExist(srcstate.Serial, destStates)
			if exists {
				fmt.Println("State Exists in destination will not migrate", srcstate.Serial)
			} else {
				state, err := downloadSourceState(tfclient.GetClientContexts(), srcstate.DownloadURL)
				calculateMD5(state)
				// Turn MD5 from []byte to string
				md5String := string(state)
				srcstate, err := c.DestinationClient.StateVersions.Create(c.DestinationContext, destWorkspaceId, tfe.StateVersionCreateOptions{
					Type:             "",
					Lineage:          new(string),
					MD5:              &md5String,
					Serial:           &srcstate.Serial,
					State:            new(string),
					Force:            new(bool),
					Run:              &tfe.Run{},
					JSONState:        new(string),
					JSONStateOutputs: new(string),
				})
				if err != nil {
					return err
				}
				o.AddDeferredMessageRead("Migrated", srcstate.Serial)
			}
		}

	}
	return nil
}
