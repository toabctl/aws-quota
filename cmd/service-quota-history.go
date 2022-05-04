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
	rootCmd.AddCommand(serviceQuotaHistoryCmd)
}

var serviceQuotaHistoryCmd = &cobra.Command{
	Use:   "service-quota-history",
	Short: "service quota history for all regions",
	Long:  `service history for all regions`,
	Run: func(cmd *cobra.Command, args []string) {
		regions := common.AwsRegions()

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"region", "service code", "quota code", "desired value", "status", "case ID"})
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

				params := &servicequotas.ListRequestedServiceQuotaChangeHistoryInput{
					MaxResults: aws.Int32(100),
				}
				paginator := servicequotas.NewListRequestedServiceQuotaChangeHistoryPaginator(client, params, func(o *servicequotas.ListRequestedServiceQuotaChangeHistoryPaginatorOptions) {
					o.Limit = 100
				})

				for paginator.HasMorePages() {
					output, err := paginator.NextPage(context.TODO())
					if err != nil {
						log.Printf("error: %v", err)
						return
					}
					for _, rq := range output.RequestedQuotas {
						tMutex.Lock()
						t.AppendRows([]table.Row{
							{
								reg,
								aws.ToString(rq.ServiceCode),
								aws.ToString(rq.QuotaCode),
								aws.ToFloat64(rq.DesiredValue),
								rq.Status,
								aws.ToString(rq.CaseId),
							},
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
