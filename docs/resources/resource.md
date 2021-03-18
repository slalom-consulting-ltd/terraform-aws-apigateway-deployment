---
page_title: "apigwdeployment_resource Resource - terraform-provider-apigwdeployment"
subcategory: ""
description: |-
  Sample resource in the Terraform provider apigwdeployment.
---

# Resource `apigwdeployment_resource`

Sample resource in the Terraform provider apigwdeployment.

## Example Usage

```terraform
resource "apigwdeployment" "example" {
  rest_api_id = "rest_api_id"
  stage_name = "stage_name"
  description = "description"
  stage_description = "stage_description"
  triggers = {"value1": "value1"}
  variables = {"value1": "value1"}
  canary_settings_percentTraffic = 20
  canary_settings_stageVariableOverrides = {"value1": "value1"}
  canary_settings_useStageCache = true
}
```

## Schema

### Optional

- **stage_name** (String, Optional)
- **stage_description** (String, Optional)
- **triggers** (Map, Optional)
- **variables** (Map, Optional)
- **canary_settings_percentTraffic** (String, Optional)
- **canary_settings_stageVariableOverrides** (Map, Optional)
- **canary_settings_useStageCache** (Boolean, Optional)


