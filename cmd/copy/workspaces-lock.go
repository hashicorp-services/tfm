// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package copy

import (
	"fmt"
	"os"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

// All functions related to locking a workspace
func workspaceLock(c tfclient.ClientContexts) error {
		o.AddMessageUserProvided("Locking all configured workspaces on:", c.SourceHostname)
		
		// Get the source workspaces properties
		srcWorkspaces, err := getSrcWorkspacesCfg(c)
		if err != nil {
			return errors.Wrap(err, "failed to list Workspaces from source while checking lock status")
		}

		for _, ws := range srcWorkspaces {

			wsProperties, err := c.SourceClient.Workspaces.ReadByID(c.SourceContext, ws.ID)
			if err != nil {
				return err
			}
			
			if !wsProperties.Locked { 
				fmt.Println("Locking Workspace:", ws.Name)
				message := "tfm migration lock"
				lockStats, lockErr := c.SourceClient.Workspaces.Lock(c.SourceContext, ws.ID, tfe.WorkspaceLockOptions{
					Reason: &message,
				})
				if lockErr != nil {
					return lockErr
				}

				_ = lockStats
			} else {
				fmt.Println("Workspace is already locked:", ws.Name)
			}
		}
		// exiting here as we don't want to migrate the workspaces at this time
		os.Exit(0)
		return nil
	}
