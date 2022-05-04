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
	rootCmd.AddCommand(serviceQuotaIncreaseCmd)
	serviceQuotaIncreaseCmd.PersistentFlags().String("servicecode", "", "The service code (eg. ec2)")
	serviceQuotaIncreaseCmd.MarkPersistentFlagRequired("servicecode")
	serviceQuotaIncreaseCmd.PersistentFlags().String("quotacode", "", "The quota code")
	serviceQuotaIncreaseCmd.MarkPersistentFlagRequired("quotacode")
	serviceQuotaIncreaseCmd.PersistentFlags().Float64("quotavalue", -1, "The quota value")
	serviceQuotaIncreaseCmd.MarkPersistentFlagRequired("quotavalue")
}

var serviceQuotaIncreaseCmd = &cobra.Command{
	Use:   "service-quota-increase",
	Short: "Increase service quota in all regions",
	Long:  `Increase service quota in all regions`,
	Run: func(cmd *cobra.Command, args []string) {
		servicecode, err := cmd.Flags().GetString("servicecode")
		if err != nil {
			log.Fatalf("can not ger servicecode: %v", err)
		}

		quotacode, err := cmd.Flags().GetString("quotacode")
		if err != nil {
			log.Fatalf("can not ger quotacode: %v", err)
		}

		quotavalue, err := cmd.Flags().GetFloat64("quotavalue")
		if err != nil {
			log.Fatalf("can not ger quotavalue: %v", err)
		}

		regions := common.AwsRegions()

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"region", "service code", "quota code", "desired value", "Case ID"})
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

				params := &servicequotas.RequestServiceQuotaIncreaseInput{
					ServiceCode: aws.String(servicecode),
					QuotaCode: aws.String(quotacode),
					DesiredValue: aws.Float64(quotavalue),
				}
				service_quota_change, err := client.RequestServiceQuotaIncrease(context.TODO(), params)
				if err != nil {
					log.Fatalf("can not create service quota request: %v", err)
					return
				}
				tMutex.Lock()
				t.AppendRows([]table.Row{
					{
						reg,
						aws.ToString(service_quota_change.RequestedQuota.ServiceCode),
						aws.ToString(service_quota_change.RequestedQuota.QuotaCode),
						aws.ToFloat64(service_quota_change.RequestedQuota.DesiredValue),
						aws.ToString(service_quota_change.RequestedQuota.CaseId),
					},
				})
				tMutex.Unlock()
				
			}(r)
		}
		wg.Wait()
		t.Render()
	},
}
