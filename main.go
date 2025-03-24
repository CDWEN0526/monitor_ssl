package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// 阿里云获取域名列表
func getDamainList(key string, secret string) ([]string, error) {
	var domains_list []string
	jsonData, err := DescribeDomains(key, secret)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return []string{}, err
	}

	domainsList, ok := jsonData["Domains"].(map[string]interface{})["Domain"]
	if !ok {
		fmt.Printf("Error: %v\n", err)
		return []string{}, err
	}

	for _, domain := range domainsList.([]interface{}) {
		d, _ := domain.(map[string]interface{})["DomainName"].(string)
		domains_list = append(domains_list, d)
	}
	return domains_list, nil
}

// 阿里云获取域名解析记录
func getDamainRecordsList(key string, secret string, domain string) ([]string, error) {
	var records_list []string
	jsonData, err := DescribeDomainRecords(key, secret, domain)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return []string{}, err
	}
	recordsList, ok := jsonData["DomainRecords"].(map[string]interface{})["Record"]
	if !ok {
		return []string{}, nil
	}
	for _, r := range recordsList.([]interface{}) {
		record := r.(map[string]interface{})
		status := record["Status"].(string)
		rr := record["RR"].(string)
		if status == "ENABLE" {
			records_list = append(records_list, rr+"."+domain)
		}
	}
	return records_list, nil

}

// 获取阿里云域名解析列表
func getRecordsList(key string, secret string) ([]string, error) {
	var wg sync.WaitGroup
	maxConcurrent := 3
	semaphore := make(chan struct{}, maxConcurrent)
	domain_list, err := getDamainList(key, secret)
	result := make([]string, 0)
	resultMutex := &sync.Mutex{}
	if err != nil {
		return nil, err
	}
	for _, domain := range domain_list {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			records, err := getDamainRecordsList(key, secret, domain)
			if err != nil {
				fmt.Println("获取域名解析记录失败:", err)
				return
			}
			for _, record := range records {
				resultMutex.Lock()
				result = append(result, record)
				resultMutex.Unlock()
			}
		}(domain)
	}
	wg.Wait()
	return result, nil
}

// 获取腾讯云域名列表
func txGetDamainList(id string, key string) ([]string, error) {
	var domain_list []string
	jsonData, err := txDescribeDomainList(id, key)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return []string{}, err
	}
	domainList, ok := jsonData["Response"].(map[string]interface{})["DomainList"]
	if !ok {
		fmt.Printf("Error: %v\n", err)
		return []string{}, err
	}
	for _, domain := range domainList.([]interface{}) {
		d, _ := domain.(map[string]interface{})["Punycode"].(string)
		domain_list = append(domain_list, d)
	}
	return domain_list, nil
}

// 获取腾讯云域名解析记录
func txGetDamainRecordsList(id string, key string, domain string) ([]string, error) {
	var records_list []string
	jsonData, err := txDescribeRecord(domain, id, key)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return []string{}, err
	}
	recordsList, ok := jsonData["Response"].(map[string]interface{})["RecordList"]
	if !ok {
		fmt.Printf("Error: %v\n", err)
		return []string{}, err
	}
	for _, r := range recordsList.([]interface{}) {
		status := r.(map[string]interface{})["Status"].(string)
		name := r.(map[string]interface{})["Name"].(string)
		if status == "ENABLE" {
			records_list = append(records_list, name+"."+domain)
		}
	}
	return records_list, nil
}

// 获取腾讯云域名解析列表
func txGetRecordsList(id string, key string) ([]string, error) {
	var wg sync.WaitGroup
	maxConcurrent := 3
	semaphore := make(chan struct{}, maxConcurrent)
	domain_list, err := txGetDamainList(id, key)
	result := make([]string, 0)
	resultMutex := &sync.Mutex{}
	if err != nil {
		return nil, err
	}
	for _, domain := range domain_list {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			records, err := txGetDamainRecordsList(id, key, domain)
			if err != nil {
				fmt.Println("获取域名解析记录失败:", err)
				return
			}
			for _, record := range records {
				resultMutex.Lock()
				result = append(result, record)
				resultMutex.Unlock()
			}
		}(domain)
	}
	wg.Wait()
	return result, nil
}

// 并发获取证书有效期
func getCertValidTime(domainList []string) (map[string]int, error) {
	var wg sync.WaitGroup
	maxConcurrent := concurrent
	semaphore := make(chan struct{}, maxConcurrent)
	result := make(map[string]int)
	resultMutex := &sync.Mutex{}

	for _, domain := range domainList {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			doamin_status, err := getHttps(domain, days_due)
			if err != nil {
				return
			}
			resultMutex.Lock()
			for k, v := range doamin_status {
				result[k] = v
			}
			resultMutex.Unlock()
		}(domain)
	}
	wg.Wait()
	return result, nil
}

// 读取config.json文件
func readConfig() (map[string]interface{}, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var data map[string]interface{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func send(subject string, body string) {
	// 邮件服务器配置
	email := config["email"].(map[string]interface{})
	smtpHost := email["host"].(string)
	smtpPort := email["port"].(string)
	from := email["user"].(string)
	password := email["pass"].(string)
	to := make([]string, 0)
	emailTo, ok := email["to"].([]interface{})
	if !ok {
		fmt.Println("配置文件中的'to'字段格式不正确")
		return
	}
	for _, v := range emailTo {
		to = append(to, v.(string))
	}

	// 调用发送邮件函数
	err := SendEmail(smtpHost, smtpPort, from, password, to, subject, body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("邮件发送成功！")

}

// 发送告警邮件
func sendAlarm(alarm_domain_name []map[string]int, title string) {
	var content string
	for _, v := range alarm_domain_name {
		for k, v := range v {
			content += fmt.Sprintf("域名：%s 剩余：%d天\n", k, v)
		}
	}
	fmt.Println(content)
	if len(content) == 0 {
		send(title+"-SSL证书即将到期", "无临期域名...")
		return
	}
	send(title+"-SSL证书即将到期", content)
}

var config map[string]interface{}
var concurrent int
var days_due int

func main() {
	startTime := time.Now()
	// 初始化全局变量config
	var err error
	config, err = readConfig()
	if err != nil {
		fmt.Println("读取配置文件失败:", err)
		return
	}
	// 并发设置，到期天数设置
	concurrent = int(config["concurrence"].(float64))
	days_due = int(config["alarm_days"].(float64))

	// 阿里云检测部分
	var records_list []string
	var alarm_domain_name []map[string]int
	// 遍历所有阿里云账户，获取域名解析记录
	for _, v := range config["aliyun"].([]interface{}) {
		var rl []string
		aliyun := v.(map[string]interface{})
		key := aliyun["key"].(string)
		secret := aliyun["secret"].(string)
		rl, err = getRecordsList(key, secret)
		if err != nil {
			fmt.Println(err)
			continue
		}
		records_list = append(records_list, rl...)
	}

	domain_validity, err := getCertValidTime(records_list)
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range domain_validity {
		alarm_domain_name = append(alarm_domain_name, map[string]int{k: v})
	}
	if len(alarm_domain_name) != 0 {
		sendAlarm(alarm_domain_name, "阿里云")
	}

	// 腾讯云检测部分
	records_list = []string{}
	alarm_domain_name = []map[string]int{}
	// 遍历所有腾讯云账户，获取域名解析记录
	for _, v := range config["tencent"].([]interface{}) {
		var rl []string
		tencent := v.(map[string]interface{})
		id := tencent["SecretId"].(string)
		key := tencent["SecretKey"].(string)
		rl, err = txGetRecordsList(id, key)
		if err != nil {
			fmt.Println(err)
			continue
		}
		records_list = append(records_list, rl...)
	}
	domain_validity, err = getCertValidTime(records_list)
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range domain_validity {
		alarm_domain_name = append(alarm_domain_name, map[string]int{k: v})
	}
	if len(alarm_domain_name) != 0 {
		sendAlarm(alarm_domain_name, "腾讯云")
	}

	// 其他域名检测部分
	records_list = []string{}
	otherDomains, ok := config["other"].([]interface{})
	alarm_domain_name = []map[string]int{}
	if !ok {
		fmt.Println("配置文件中的'other'字段格式不正确")
		return
	}
	for _, d := range otherDomains {
		domain := d.(string)
		records_list = append(records_list, domain)
	}
	domain_validity, err = getCertValidTime(records_list)
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range domain_validity {
		alarm_domain_name = append(alarm_domain_name, map[string]int{k: v})
	}
	if len(alarm_domain_name) != 0 {
		sendAlarm(alarm_domain_name, "其他")
	}

	endTime := time.Now()
	fmt.Printf("执行时间：%s\n", endTime.Sub(startTime))
}
