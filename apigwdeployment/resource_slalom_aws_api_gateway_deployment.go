package apigwdeployment

import (
	"context"
	"fmt"
	"log"
	"time"

	//"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSlalomAwsApiGatewayDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsApiGatewayDeploymentCreate,
		Read:   resourceAwsApiGatewayDeploymentRead,
		Update: resourceAwsApiGatewayDeploymentUpdate,
		Delete: resourceAwsApiGatewayDeploymentDelete,

		Schema: map[string]*schema.Schema{
			"rest_api_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"stage_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"stage_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"triggers": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"variables": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"invoke_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"execution_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"canary_settings_percent_traffic": {
				Type:     schema.TypeFloat,
				Optional: true,
			},

			"canary_settings_stage_variable_overrides": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"canary_settings_use_stage_cache": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAwsApiGatewayDeploymentCreate(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(aws.Config)
	client := apigateway.NewFromConfig(cfg)

	variables := make(map[string]string)
	for k, v := range d.Get("variables").(map[string]interface{}) {
		variables[k] = v.(string)
	}

	stageVariablesOverrrides := make(map[string]string)
	for k, v := range d.Get("canary_settings_stage_variable_overrides").(map[string]interface{}) {
		variables[k] = v.(string)
	}

	type CacheClusterSize string

	var err error
	deployment, err := client.CreateDeployment(context.TODO(), &apigateway.CreateDeploymentInput{
		CacheClusterEnabled: new(bool),
		CacheClusterSize:    types.CacheClusterSizeSize0Point5Gb,
		CanarySettings:      &types.DeploymentCanarySettings{PercentTraffic: d.Get("canary_settings_percent_traffic").(float64), StageVariableOverrides: stageVariablesOverrrides, UseStageCache: d.Get("canary_settings_use_stage_cache").(bool)},
		Description:         aws.String(d.Get("description").(string)),
		RestApiId:           aws.String(d.Get("rest_api_id").(string)),
		StageDescription:    aws.String(d.Get("stage_description").(string)),
		StageName:           aws.String(d.Get("stage_name").(string)),
		TracingEnabled:      new(bool),
		Variables:           variables,
	})

	if err != nil {
		return fmt.Errorf("Error creating API Gateway Deployment: %s", err)
	}
	_ = deployment
	return resourceAwsApiGatewayDeploymentRead(d, meta)
}

func resourceAwsApiGatewayDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(aws.Config)
	client := apigateway.NewFromConfig(cfg)

	log.Printf("[DEBUG] Reading API Gateway Deployment %s", d.Id())
	restApiId := d.Get("rest_api_id").(string)
	out, err := client.GetDeployment(context.TODO(), &apigateway.GetDeploymentInput{
		RestApiId:    aws.String(restApiId),
		DeploymentId: aws.String(d.Id()),
	})
	if err != nil {
		// if isAWSErr(err, apigateway.ErrCodeNotFoundException, "") {
		//  log.Printf("[WARN] API Gateway Deployment (%s) not found, removing from state", d.Id())
		//  d.SetId("")
		//  return nil
		// }
		// return err
		return nil
	}
	log.Printf("[DEBUG] Received API Gateway Deployment: %s", out)
	d.Set("description", out.Description)
	stageName := d.Get("stage_name").(string)
	hostname := fmt.Sprintf("%s.%s.%s", fmt.Sprintf("%s.execute-api", restApiId), cfg.Region, "amazonaws.com")
	d.Set("invoke_url", fmt.Sprintf("https://%s/%s", hostname, stageName))

	executionArn := arn.ARN{
		Partition: "aws", //meta.(*AWSClient).partition,
		Service:   "execute-api",
		Region:    cfg.Region,
		AccountID: d.Get("account_id").(string),
		Resource:  fmt.Sprintf("%s/%s", restApiId, stageName),
	}.String()
	d.Set("execution_arn", executionArn)
	if err := d.Set("created_date", out.CreatedDate.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Error setting created_date: %s", err)
	}
	return nil
}

// func resourceAwsApiGatewayDeploymentUpdateOperations(d *schema.ResourceData) []types.PatchOperation {
// 	operations := make([]types.PatchOperation, 0)
// 	if d.HasChange("description") {
// 		operations = append(operations, types.PatchOperation{
// 			Op:    aws.String(*apigateway.OpReplace),
// 			Path:  aws.String("/description"),
// 			Value: aws.String(d.Get("description").(string)),
// 		})
// 	}
// 	return operations
// }

func resourceAwsApiGatewayDeploymentUpdate(d *schema.ResourceData, meta interface{}) error {
	// cfg := meta.(aws.Config)
	// client := apigateway.NewFromConfig(cfg)
	// log.Printf("[DEBUG] Updating API Gateway API Key: %s", d.Id())
	// _, err := client.UpdateDeployment(context.TODO(), &apigateway.UpdateDeploymentInput{
	// 	DeploymentId:    aws.String(d.Id()),
	// 	RestApiId:       aws.String(d.Get("rest_api_id").(string)),
	// 	PatchOperations: resourceAwsApiGatewayDeploymentUpdateOperations(d),
	// })
	// if err != nil {
	// 	return err
	// }
	return resourceAwsApiGatewayDeploymentRead(d, meta)
}
func resourceAwsApiGatewayDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	// conn := meta.(*AWSClient).apigatewayconn
	// log.Printf("[DEBUG] Deleting API Gateway Deployment: %s", d.Id())
	// // If the stage has been updated to point at a different deployment, then
	// // the stage should not be removed when this deployment is deleted.
	// shouldDeleteStage := false
	// // API Gateway allows an empty state name (e.g. ""), but the AWS Go SDK
	// // now has validation for the parameter, so we must check first.
	// // InvalidParameter: 1 validation error(s) found.
	// //  - minimum field size of 1, GetStageInput.StageName.
	// stageName := d.Get("stage_name").(string)
	// if stageName != "" {
	//  stage, err := conn.GetStage(&apigateway.GetStageInput{
	//    StageName: aws.String(stageName),
	//    RestApiId: aws.String(d.Get("rest_api_id").(string)),
	//  })
	//  if err != nil && !isAWSErr(err, apigateway.ErrCodeNotFoundException, "") {
	//    return fmt.Errorf("error getting referenced stage: %s", err)
	//  }
	//  if stage != nil && aws.StringValue(stage.DeploymentId) == d.Id() {
	//    shouldDeleteStage = true
	//  }
	// }
	// if shouldDeleteStage {
	//  if _, err := conn.DeleteStage(&apigateway.DeleteStageInput{
	//    StageName: aws.String(d.Get("stage_name").(string)),
	//    RestApiId: aws.String(d.Get("rest_api_id").(string)),
	//  }); err == nil {
	//    return nil
	//  }
	// }
	// _, err := conn.DeleteDeployment(&apigateway.DeleteDeploymentInput{
	//  DeploymentId: aws.String(d.Id()),
	//  RestApiId:    aws.String(d.Get("rest_api_id").(string)),
	// })
	// if isAWSErr(err, apigateway.ErrCodeNotFoundException, "") {
	//  return nil
	// }
	// if err != nil {
	//  return fmt.Errorf("error deleting API Gateway Deployment (%s): %s", d.Id(), err)
	// }
	return nil
}
