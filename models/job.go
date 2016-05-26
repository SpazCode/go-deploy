package models

type (
	Job struct {
		title 			string
		description		string
		repo_link 		string
		test_command	string	
		deploy_location	string
		deploy_password	string
		deploy_username	string
		current_build	string
		last_run		string
	}
)