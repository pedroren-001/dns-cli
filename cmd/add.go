package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

var (
	recordType string
	recordLine string
	value      string
	subDomain  string
	ttl        uint64
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "添加一条 DNS 解析记录",
	Example: `  dns-cli add -d example.com --type A --sub www --value 1.2.3.4
  dns-cli add -d example.com --type CNAME --sub blog --value cdn.example.com
  dns-cli add -d example.com --type MX --sub @ --value mail.example.com --line 默认`,
	RunE: runAdd,
}

func init() {
	addCmd.Flags().StringVar(&recordType, "type", "", "记录类型 (必填，如 A, CNAME, MX, TXT 等)")
	addCmd.Flags().StringVar(&subDomain, "sub", "", "主机记录/子域名 (必填，如 www, @, mail)")
	addCmd.Flags().StringVar(&value, "value", "", "记录值 (必填，如 IP 地址或域名)")
	addCmd.Flags().StringVar(&recordLine, "line", "默认", "解析线路")
	addCmd.Flags().Uint64Var(&ttl, "ttl", 600, "TTL 值（秒）")

	addCmd.MarkFlagRequired("type")
	addCmd.MarkFlagRequired("sub")
	addCmd.MarkFlagRequired("value")

	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	if err := requireDomain(); err != nil {
		return err
	}

	client, err := newDNSPodClient()
	if err != nil {
		return err
	}

	request := dnspod.NewCreateRecordRequest()
	request.Domain = common.StringPtr(domain)
	request.RecordType = common.StringPtr(recordType)
	request.RecordLine = common.StringPtr(recordLine)
	request.Value = common.StringPtr(value)
	request.SubDomain = common.StringPtr(subDomain)
	request.TTL = common.Uint64Ptr(ttl)

	response, err := client.CreateRecord(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("API 错误: %s", err)
	}
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}

	recordID := uint64(0)
	if response.Response != nil && response.Response.RecordId != nil {
		recordID = *response.Response.RecordId
	}

	fmt.Printf("记录添加成功\n")
	fmt.Printf("  记录 ID: %d\n", recordID)
	fmt.Printf("  域名:    %s\n", domain)
	fmt.Printf("  主机记录: %s\n", subDomain)
	fmt.Printf("  类型:    %s\n", recordType)
	fmt.Printf("  值:      %s\n", value)
	fmt.Printf("  线路:    %s\n", recordLine)
	fmt.Printf("  TTL:     %d\n", ttl)

	return nil
}
