#!/bin/bash

host="${host:-localhost:80}"


curl -X POST ${host}/hub/tasks -d \
'{
  "name": "TCA",
  "state": "Ready",
  "locator": "tca",
  "addon": "tca",
  "application": {
    "id": 1
  },
  "data": {
      "application_name": "TCA",
      "application_description": "test application",
      "technology_summary": "RHEL,.net,jboss"
    }
}'
