package cmd

import (
	"context"
	"log"
	"os"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"

	"github.com/spf13/cobra"
	"github.com/jedib0t/go-pretty/v6/table"
)

func init() {
	rootCmd.AddCommand(servicesListCmd)
}

var servicesListCmd = &cobra.Command{
	Use:   "services-list",
	Short: "List available services",
	Long:  `List available services`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}
	
		client := servicequotas.NewFromConfig(cfg)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"service code", "service name"})
	

		params := &servicequotas.ListServicesInput{
			MaxResults: aws.Int32(100),
		}
		paginator := servicequotas.NewListServicesPaginator(client, params, func(o *servicequotas.ListServicesPaginatorOptions) {
			o.Limit = 100
		})
	
		for paginator.HasMorePages() {
			output, err := paginator.NextPage(context.TODO())
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
			for _, s := range output.Services {
				t.AppendRows([]table.Row{
					{aws.ToString(s.ServiceCode), aws.ToString(s.ServiceName)},
				})
			}
		}
		t.Render()
	},
}
