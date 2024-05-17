package main

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

func validate(config string) bool {
	log.Print(os.Getwd())
	tf, err := tfexec.NewTerraform(".", "terraform")
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

func deploy(config string) bool {
	tf, err := tfexec.NewTerraform(".terraform", "terraform")
	if err != nil {
		log.Printf("failed to create Terraform: %s", err)
		return false
	}
	// defer tf.Close()

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