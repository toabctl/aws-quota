package cmd

import (
	"context"
	"log"
	"os"
	"sync"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"

	"github.com/spf13/cobra"
	"github.com/jedib0t/go-pretty/v6/table"

	"aws-quota/common"
)

func init() {
	rootCmd.AddCommand(serviceQuotasListCmd)
	serviceQuotasListCmd.PersistentFlags().String("servicecode", "", "The service code (eg. ec2)")
	serviceQuotasListCmd.MarkPersistentFlagRequired("servicecode")
}

var serviceQuotasListCmd = &cobra.Command{
	Use:   "service-quotas-list",
	Short: "List available service quotas in all regions",
	Long:  `List available service quotas in all regions`,
	Run: func(cmd *cobra.Command, args []string) {
		servicecode, err := cmd.Flags().GetString("servicecode")
		if err != nil {
			log.Fatalf("can not ger servicecode: %v", err)
		}
		regions := common.AwsRegions()

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"region", "quota name", "quota code", "quota value", "adjustable"})
		var tMutex sync.Mutex

		var wg sync.WaitGroup

		for _, r := range regions {
			wg.Add(1)

			go func(reg string) {
				defer wg.Done()
				cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(reg))
				if err != nil {
					log.Fatalf("unable to load SDK config, %v", err)
				}
	
				client := servicequotas.NewFromConfig(cfg)

				params := &servicequotas.ListServiceQuotasInput{
					ServiceCode: aws.String(servicecode),
					MaxResults: aws.Int32(100),
				}
				paginator := servicequotas.NewListServiceQuotasPaginator(client, params, func(o *servicequotas.ListServiceQuotasPaginatorOptions) {
				})
				for paginator.HasMorePages() {
					output, err := paginator.NextPage(context.TODO())
					if err != nil {
						log.Fatalf("error: %v", err)
						return
					}
					for _, q := range output.Quotas {
						tMutex.Lock()
						t.AppendRows([]table.Row{
							{reg, aws.ToString(q.QuotaName), aws.ToString(q.QuotaCode), aws.ToFloat64(q.Value), q.Adjustable},
						})
						tMutex.Unlock()
					}
				}
			}(r)
		}
		wg.Wait()
		t.Render()
	},
}
