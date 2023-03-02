output "source_agents_pool_id" {
  value = tfe_agent_pool.source.id
}

output "destination_agents_pool_id" {
  value = tfe_agent_pool.destination.id
}

output "source_ssh_id" {
  value = tfe_ssh_key.source.id
}

output "destination_ssh_id" {
  value = tfe_ssh_key.destination.id
}

output "source_team_id" {
  value = tfe_team.source.id
}

output "source_team_name" {
  value = tfe_team.source.name
}

output "source_gh_oauth_token_id" {
  value = tfe_oauth_client.source.oauth_token_id
}

output "destination_gh_oauth_token_id" {
  value = tfe_oauth_client.destination.oauth_token_id
}

output "source_workspace" {
  value = module.workspacer_source
}
