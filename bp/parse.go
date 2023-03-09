package bp

import "net"

/*
	调用net.ParseCIDR函数，将字符串的CIDR地址转换为net.IP和net.IPNet类型的地址。
	在循环中，使用ip.Mask函数，将起始IP地址与子网掩码进行按位与运算，得到当前IP地址。
	使用ipnet.Contains函数判断当前IP地址是否在子网中，如果在，则将其加入到ips切片中。使用inc函数，将当前IP地址加1，以生成下一个IP地址。
*/
func Parseip(inip string) []string {
	ip, ipnet, err := net.ParseCIDR(inip)
	if err != nil {
		return nil
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	ips = ips[1 : len(ips)-1]
	return ips
}
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
 