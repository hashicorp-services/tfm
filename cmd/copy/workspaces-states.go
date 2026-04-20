// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package copy

import (
	"crypto/md5"
	b64 "encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp-services/tfm/cmd/helper"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

// 1. Get source state versions per workspace and reverse the slice
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

// Iterate backwards through the srcstate slice and append each element to a new slice
// to create a reverse ordered slice of srcStates

func rateLimitTest() {
	// Configure the rate limit to exceed 30 requests per second, set it to a higher value.
	requestsPerSecond := 1000
	requestInterval := time.Second / time.Duration(requestsPerSecond)

	// Create a wait group to wait for all goroutines to finish.
	var wg sync.WaitGroup

	// Launch multiple goroutines to make API requests.
	for i := 0; i < 1000; i++ { // Launch 100 goroutines
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Simulate making an API request.
			discoverDestTeams(tfclient.GetClientContexts())

			// Sleep for the specified interval before making the next request.
			time.Sleep(requestInterval)
		}()
	}

	// Wait for all goroutines to finish.
	wg.Wait()

	fmt.Println("All API requests completed.")
}

func reverseSlice(input []*tfe.StateVersion) []*tfe.StateVersion {
	inputLen := len(input)
	output := make([]*tfe.StateVersion, inputLen)

	for i, n := range input {
		j := inputLen - i - 1

		output[j] = n
	}

	return output
}

// Get the source workspace state files from the provided workspace
func discoverSrcStates(c tfclient.ClientContexts, ws string, NumberOfStates int) ([]*tfe.StateVersion, error) {
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

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage

	}
	o.AddFormattedMessageCalculated("Found %d Workspace states", len(srcStates))

	if NumberOfStates != 0 {
		o.AddFormattedMessageCalculated("Only the %d newest workspace states will be migrated", NumberOfStates)

		// If a last X amount of states is given, remove all previous states except for the last X amount.
		// If there are fewer states to keep than there are states, set the number to keep the same amount of states
		// there are for the workspace
		if NumberOfStates > len(srcStates) {
			NumberOfStates = len(srcStates)
		}

		srcStates = srcStates[:len(srcStates)-(len(srcStates)-NumberOfStates)]
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

	if !wsProperties.Locked {

		fmt.Println("Locking Workspace: ", destWorkspaceId)
		lockStats, lockErr := c.DestinationClient.Workspaces.Lock(c.DestinationContext, destWorkspaceId, tfe.WorkspaceLockOptions{
			Reason: &message,
		})
		if lockErr != nil {
			return lockErr
		}

		_ = lockStats

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
func copyStates(c tfclient.ClientContexts, NumberOfStates int) error {

	// Get the source target workspaces
	srcWorkspaces, err := getSrcWorkspacesCfg(c)
	if err != nil {
		return errors.Wrap(err, "failed to list Workspaces from source")
	}

	if NumberOfStates > 1 {
		// fmt.Printf("\n\n**** Operation will migrate last %v states per workspace **** \n\n", NumberOfStates)
		o.AddMessageUserProvided2("\n\n", fmt.Sprint(NumberOfStates), "states per workspace will be copied over.\n\nWarning:\n\n**** THIS OPERATION SHOULD NOT BE RAN MORE THAN ONCE ***")

		if !confirm() {
			fmt.Println("\n\n**** Canceling tfm run **** ")
			os.Exit(1)
		}
	}

	// Get/Check if Workspace map exists
	wsMapCfg, err := helper.ViperStringSliceMap("workspaces-map")
	if err != nil {
		fmt.Println("invalid input for workspaces-map")
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

		// Check if the workspace name prefix and suffix are set
		if len(wsNamePrefix) > 0 || len(wsNameSuffix) > 0 {
			srcworkspaceSlice := []*tfe.Workspace{{Name: destWorkSpaceName}}
			newDestWorkspaceName := standardizeNamingConvention(srcworkspaceSlice, wsNamePrefix, wsNameSuffix)
			destWorkSpaceName = newDestWorkspaceName[0].Name
		}

		// Check for the existence of the destination workspace in the destination target
		exists := doesWorkspaceExist(destWorkSpaceName, destWorkspaces)

		if exists {
			destWorkspaceId, err := getWorkspaceId(tfclient.GetClientContexts(), destWorkSpaceName)
			if err != nil {
				return errors.Wrap(err, "Failed to get the ID of the destination Workspace that matches the Name of the Source Workspace")
			}

			fmt.Printf("Source ws %v has a matching ws %v in destination with ID %v. Comparing existing States...\n", srcworkspace.Name, destWorkSpaceName, destWorkspaceId)

			// Get the source workspace states
			srcStates, err := discoverSrcStates(tfclient.GetClientContexts(), srcworkspace.Name, NumberOfStates)
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

			for _, srcstate := range reverseSlice(srcStates) {

				exists := doesStateExist(srcstate.Serial, destStates)
				if exists {
					fmt.Printf("State Version %v with Serial %v exists in destination will not migrate\n", srcstate.StateVersion, srcstate.Serial)
				} else {

					// Download state from source
					state, err := downloadSourceState(tfclient.GetClientContexts(), srcstate.DownloadURL)

					// Create an empty int
					//newSerial := int64(1)

					// Get properties of the state
					// currentState, _ := c.DestinationClient.StateVersions.ReadCurrent(c.DestinationContext, destWorkspaceId)

					// Of if there is a state file, set the newSerial variable to the current state serial + 1
					// if currentState != nil {
					// 	newSerial = currentState.Serial + 1
					// }

					// Get state file as string
					plainTextState := string(state)

					// Define a regular expression pattern to match the "serial" value
					serialPattern := `"serial":\s*(\d+)`

					// Compile the regular expression
					serialRe := regexp.MustCompile(serialPattern)

					// Find the match
					serialMatch := serialRe.FindStringSubmatch(plainTextState)

					if len(serialMatch) != 2 {
						fmt.Println("Serial not found in JSON")
						return err
					}

					newSerialConversion, err := strconv.ParseInt(serialMatch[1], 10, 64)
					if err != nil {
						fmt.Printf("The source state wasn't a int64")
					}

					// Define a regular expression pattern to match the "lineage" value
					lineagePattern := `"lineage":\s*"([^"]+)"`

					// Compile the regular expression
					lineageRe := regexp.MustCompile(lineagePattern)

					// Find the match
					lineageMatch := lineageRe.FindStringSubmatch(plainTextState)

					if len(lineageMatch) != 2 {
						fmt.Println("Lineage not found in JSON")
						return err
					}

					lineage := lineageMatch[1]
					//fmt.Println("State JSON Lineage:", lineage)

					// Base64 encode the state as a string
					stringState := b64.StdEncoding.EncodeToString(state)

					// Get the MD5 hash of the state
					md5String := fmt.Sprintf("%x", md5.Sum([]byte(state)))

					// Lock the destination workspace
					lockWorkspace(tfclient.GetClientContexts(), destWorkspaceId)
					fmt.Printf("Migrating state version %v serial %v for workspace Src: %v Dst: %v\n", srcstate.StateVersion, newSerialConversion, srcworkspace.Name, destWorkSpaceName)
					// // ------------------------------------------------------------------------------
					// // --- START rate limiting testing code ------------------------------------------
					// // --- Comment out when not testing ------------------------------------------
					// // ------------------------------------------------------------------------------
					// Configure the rate limit to exceed 30 requests per second, set it to a higher value.
					// requestsPerSecond := 1000
					// requestInterval := time.Second / time.Duration(requestsPerSecond)

					// // Create a wait group to wait for all goroutines to finish.
					// var wg sync.WaitGroup

					// // Launch multiple goroutines to make API requests.
					// for i := 0; i < 1000; i++ { // Launch 100 goroutines
					// 	wg.Add(1)
					// 	go func() {
					// 		defer wg.Done()

					// 		// Simulate making an API request.
					// 		resp, err := discoverDestTeams(tfclient.GetClientContexts())
					// 		if err != nil {
					// 			// Handle other errors here.
					// 			fmt.Println("Error:", err)
					// 			return
					// 		}
					// 		_ = resp
					// 		// Sleep for the specified interval before making the next request.
					// 		time.Sleep(requestInterval)

					// 	}()
					// }

					// // Wait for all goroutines to finish.
					// wg.Wait()

					// fmt.Println("All API requests completed.")
					// // ------------------------------------------------------------------------------
					// // --- end rate limiting testing code ------------------------------------------
					// // ------------------------------------------------------------------------------
					srcstate, err := c.DestinationClient.StateVersions.Create(c.DestinationContext, destWorkspaceId, tfe.StateVersionCreateOptions{
						Type:    "",
						Lineage: &lineage,
						MD5:     tfe.String(md5String),
						Serial:  &newSerialConversion,
						//Serial:           &newSerial,
						State:            tfe.String(stringState),
						Force:            new(bool),
						Run:              &tfe.Run{},
						JSONState:        new(string),
						JSONStateOutputs: new(string),
					})

					if err != nil {
						// Get the current timestamp and format it as a string
						timestamp := time.Now().Format(time.RFC850)

						// Replace colons with a different character to make the filename Windows-compatible
						safeTimestamp := strings.ReplaceAll(timestamp, ":", "-")

						// Create a file to store workspace names with errors
						errorLogFileName := fmt.Sprintf("workspace_error_log_%s.txt", safeTimestamp)
						errorLogFile, err := os.Create(errorLogFileName)
						if err != nil {
							fmt.Printf("Failed to create error log file: %v\n", err)
							return err
						}
						defer errorLogFile.Close()

						// If there is an error output the error, log it, and move onto the next workspace.
						fmt.Println("failed to migrate state file. Moving onto next workspace.", err)
						errorLogFile.WriteString(fmt.Sprintf("Failed to migrate state file for source workspace: %v\n", srcworkspace.Name))
						break

					}

					_ = srcstate

				}
			}

			unlockWorkspace(tfclient.GetClientContexts(), destWorkspaceId)
		} else {
			fmt.Printf("Source workspace (%v) does not exist in destination (%v). No states to migrate\n", srcworkspace.Name, destWorkSpaceName)
		}
	}
	return nil
}
