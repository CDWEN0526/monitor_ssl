package main

import (
	"encoding/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

// 获取域名列表接口
func DescribeDomains(key string, secret string) (map[string]interface{}, error) {
	config := sdk.NewConfig()

	// Please ensure that the environment variables ALIBABA_CLOUD_ACCESS_KEY_ID and ALIBABA_CLOUD_ACCESS_KEY_SECRET are set.
	credential := credentials.NewAccessKeyCredential(key, secret)
	/* use STS Token
	credential := credentials.NewStsTokenCredential(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"), os.Getenv("ALIBABA_CLOUD_SECURITY_TOKEN"))
	*/
	client, err := sdk.NewClientWithOptions("cn-qingdao", config, credential)
	if err != nil {
		return nil, err
	}

	request := requests.NewCommonRequest()

	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dns.aliyuncs.com"
	request.Version = "2015-01-09"
	request.ApiName = "DescribeDomains"
	request.QueryParams["PageSize"] = "100"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(response.GetHttpContentBytes(), &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

// 获取域名解析记录列表
func DescribeDomainRecords(key string, secret string, domain string) (map[string]interface{}, error) {
	config := sdk.NewConfig()

	// Please ensure that the environment variables ALIBABA_CLOUD_ACCESS_KEY_ID and ALIBABA_CLOUD_ACCESS_KEY_SECRET are set.
	credential := credentials.NewAccessKeyCredential(key, secret)
	/* use STS Token
	credential := credentials.NewStsTokenCredential(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"), os.Getenv("ALIBABA_CLOUD_SECURITY_TOKEN"))
	*/
	client, err := sdk.NewClientWithOptions("cn-qingdao", config, credential)
	if err != nil {
		return nil, err
	}

	request := requests.NewCommonRequest()

	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dns.aliyuncs.com"
	request.Version = "2015-01-09"
	request.ApiName = "DescribeDomainRecords"
	request.QueryParams["DomainName"] = domain
	request.QueryParams["PageSize"] = "500"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(response.GetHttpContentBytes(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
