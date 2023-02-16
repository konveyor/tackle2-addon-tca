#!/bin/bash

host="${host:-localhost:80}"


curl -X POST ${host}/hub/tasks -d \
'{
  "name": "TCA",
  "state": "Ready",
  "locator": "tcafeb",
  "addon": "tcafeb",
  "application": {
    "id": 5
  },
  "data": {
      "application_name": "TCA",
      "application_description": "test application",
      "technology_summary": "R,sap,.net,jboss"
    }
}'
