package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

var jsonOutput bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出域名下的所有 DNS 解析记录",
	Example: `  dns-cli list -d example.com
  dns-cli list -d example.com --json`,
	RunE: runList,
}

func init() {
	listCmd.Flags().BoolVar(&jsonOutput, "json", false, "以 JSON 格式输出")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	if err := requireDomain(); err != nil {
		return err
	}

	client, err := newDNSPodClient()
	if err != nil {
		return err
	}

	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = common.StringPtr(domain)

	response, err := client.DescribeRecordList(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("API 错误: %s", err)
	}
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}

	if jsonOutput {
		fmt.Println(response.ToJsonString())
		return nil
	}

	// 表格输出
	records := response.Response.RecordList
	if len(records) == 0 {
		fmt.Println("未找到任何解析记录")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "ID\t名称\t类型\t线路\t值\tTTL\t状态\n")
	fmt.Fprintf(w, "--\t----\t----\t----\t--\t---\t----\n")

	for _, r := range records {
		status := "启用"
		if r.Status != nil && *r.Status == "DISABLE" {
			status = "暂停"
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%d\t%s\n",
			safeUint64(r.RecordId),
			safeString(r.Name),
			safeString(r.Type),
			safeString(r.Line),
			safeString(r.Value),
			safeUint64(r.TTL),
			status,
		)
	}

	w.Flush()
	fmt.Printf("\n共 %d 条记录\n", len(records))
	return nil
}

func safeString(p *string) string {
	if p == nil {
		return "-"
	}
	return *p
}

func safeUint64(p *uint64) uint64 {
	if p == nil {
		return 0
	}
	return *p
}
