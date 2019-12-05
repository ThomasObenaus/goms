DASHBOARD_DIR="./monitoring/dashboards"

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"

up: sep ## Start infra components
	@echo "--> Startup (+build) infra components"
	@echo "--> Copy dashboards from $(DASHBOARD_DIR)"
	@mkdir -p monitoring/grafana/dashboards
	@cp -R $(DASHBOARD_DIR) monitoring/grafana
	@docker-compose up --build -d
	@rm -rf monitoring/grafana/dashboards

down: sep ## Stop infra components
	@echo "--> Stop infra components"
	@docker-compose down

sep:
	@echo "----------------------------------------------------------------------------------"