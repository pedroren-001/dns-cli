package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

// 全局 flag
var domain string

var rootCmd = &cobra.Command{
	Use:   "dns-cli",
	Short: "腾讯云 DNSPod DNS 记录管理工具",
	Long: `dns-cli 是一个基于腾讯云 DNSPod API 的命令行工具，
用于管理域名的 DNS 解析记录。

需要设置以下环境变量进行认证：
  TENCENTCLOUD_SECRET_ID   - 腾讯云 SecretId
  TENCENTCLOUD_SECRET_KEY  - 腾讯云 SecretKey

密钥可前往控制台获取：https://console.cloud.tencent.com/cam/capi`,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "", "域名 (必填，如 example.com)")
}

// Execute 执行根命令
func Execute() error {
	return rootCmd.Execute()
}

// newDNSPodClient 创建 DNSPod 客户端（所有子命令共用）
func newDNSPodClient() (*dnspod.Client, error) {
	secretID := os.Getenv("TENCENTCLOUD_SECRET_ID")
	secretKey := os.Getenv("TENCENTCLOUD_SECRET_KEY")

	if secretID == "" || secretKey == "" {
		return nil, fmt.Errorf("请设置环境变量 TENCENTCLOUD_SECRET_ID 和 TENCENTCLOUD_SECRET_KEY")
	}

	credential := common.NewCredential(secretID, secretKey)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"

	client, err := dnspod.NewClient(credential, "", cpf)
	if err != nil {
		return nil, fmt.Errorf("创建 DNSPod 客户端失败: %w", err)
	}

	return client, nil
}

// requireDomain 校验 domain flag 是否已设置
func requireDomain() error {
	if domain == "" {
		return fmt.Errorf("请通过 --domain / -d 指定域名")
	}
	return nil
}
