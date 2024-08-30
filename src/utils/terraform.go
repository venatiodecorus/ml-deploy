package utils

import (
	"bytes"
	"context"
	"log"
	"os"
	"sync"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

var (
	tfClient *tfexec.Terraform
	tfClientOnce sync.Once
	tfClientInitErr error
)

func GetTFClient() (*tfexec.Terraform, error) {
	tfClientOnce.Do(func() {
		tfClient, tfClientInitErr = tfexec.NewTerraform(".terraform", "terraform")
	})

	if tfClientInitErr != nil {
		return nil, tfClientInitErr
	}
	return tfClient, nil
}


// Validates, formats and saves the config file
func Validate(config string) bool {
	tf, err := GetTFClient()
	if err != nil {
		log.Printf("failed to create Terraform: %s", err)
		return false
	}
	// defer tf.Close()

	var output *tfjson.ValidateOutput
	output, err = tf.Validate(context.Background())
	if err != nil {
		log.Printf("failed to validate Terraform: %s", err)
		return false
	}

	if !output.Valid {
		return false
	}

	// Format and save the config to ./terraform/main.tf
	formatted, err := tf.FormatString(context.TODO(), config)
	if err != nil {
		log.Printf("failed to format config: %s", err)
		return false
	}

	// err = os.MkdirAll("./terraform", 0755)
	// if err != nil {
	// 	log.Printf("failed to create directory: %s", err)
	// 	return false
	// }

	err = os.WriteFile("./.terraform/main.tf", []byte(formatted), 0644)
	if err != nil {
		log.Printf("failed to write config: %s", err)
		return false
	}

	return true
}

func Deploy() bool {
	tf, err := GetTFClient()
	if err != nil {
		log.Printf("failed to create Terraform: %s", err)
		return false
	}

	err = tf.Init(context.Background())
	if err != nil {
		log.Printf("failed to init Terraform: %s", err)
		return false
	}

	err = tf.Apply(context.Background())
	if err != nil {
		log.Printf("failed to apply Terraform: %s", err)
		return false
	}

	return true
}

func GetState() *tfjson.State {
	tf, err := GetTFClient()
	if err != nil {
		log.Printf("failed to create Terraform: %s", err)
		return nil
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		log.Printf("failed to show state: %s", err)
		return nil
	}

	return state
}

func GetPlan() (string,error) {
	tf, err := GetTFClient()
	if err != nil {
		log.Printf("failed to create Terraform: %s", err)
		return "",err
	}

	// Store the plan in a variable
	var buf bytes.Buffer

	plan, err := tf.PlanJSON(context.Background(), &buf)
	if err != nil {
		log.Printf("failed to show plan: %s", err)
		return "",err
	}

	if !plan {
		return "",err
	}
	return buf.String(),nil
}

func Destroy() bool {
	tf, err := GetTFClient()
	if err != nil {
		log.Printf("failed to create Terraform: %s", err)
		return false
	}

	err = tf.Destroy(context.Background())
	if err != nil {
		log.Printf("failed to destroy Terraform: %s", err)
		return false
	}

	return true
}