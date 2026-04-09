package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

var recordID uint64

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "删除一条 DNS 解析记录",
	Long: `删除指定的 DNS 解析记录。需要提供记录 ID，
可通过 dns-cli list 命令查看所有记录及其 ID。`,
	Example: `  # 先查看记录列表获取 ID
  dns-cli list -d example.com

  # 再根据 ID 删除
  dns-cli rm -d example.com --record-id 12345`,
	RunE: runRm,
}

func init() {
	rmCmd.Flags().Uint64Var(&recordID, "record-id", 0, "要删除的记录 ID (必填，可通过 list 命令查看)")
	rmCmd.MarkFlagRequired("record-id")

	rootCmd.AddCommand(rmCmd)
}

func runRm(cmd *cobra.Command, args []string) error {
	if err := requireDomain(); err != nil {
		return err
	}

	client, err := newDNSPodClient()
	if err != nil {
		return err
	}

	request := dnspod.NewDeleteRecordRequest()
	request.Domain = common.StringPtr(domain)
	request.RecordId = common.Uint64Ptr(recordID)

	_, err = client.DeleteRecord(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("API 错误: %s", err)
	}
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}

	fmt.Printf("记录删除成功 (ID: %d, 域名: %s)\n", recordID, domain)
	return nil
}
