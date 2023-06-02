package main

import (
	"fmt"
	"net/url"
)

func main() {
	u, err := url.Parse("https://stage.apis.avata.bianjie.ai/v2​/nft​/nfts​/0xA85a2C5e5b7f3cbaa670E6504622FAeDca7072FD​/0x83F41408f214CbEA04Bb87Cd64177248CfaF182f​/110")
	if err != nil {
		// URL 解析错误
	}

	for _, r := range u.String() {
		if r == '\u200B' {
			fmt.Println("包含一个零宽空格字符")
		}
	}

	fmt.Println("不包含")
}
