// Based on https://github.com/creativeprojects/go-selfupdate/blob/v1.1.1/github_release.go
// and https://github.com/creativeprojects/go-selfupdate/blob/v1.1.1/github_source.go

package update

import (
	"fmt"

	"github.com/jeessy2/ddns-go/v5/util"
)

type Release struct {
	tagName string
	assets  []Asset
}

type Asset struct {
	name string
	url  string
}

// ReleaseResp 表示仓库中的 GitHub release 和 asset。
// 返回的数据比较多，但是我们只需要定义我们需要的就可以了
// 相比于原始的数据结构，这里精简了很多
type ReleaseResp struct {
	TagName string     `json:"tag_name,omitempty"`
	Assets  []struct { //如果不需要在外部使用此结构体，可以在内部直接定义
		Name               string `json:"name,omitempty"`
		BrowserDownloadURL string `json:"browser_download_url,omitempty"`
	} `json:"assets,omitempty"`
}

// getLatest 列出仓库的最新 release 并返回包装过的 Release
//
// GitHub API 文档：https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#get-the-latest-release
// 利用Github提供的api去查询指定仓库最新版本的软件
func getLatest(repo string) (*Release, error) {
	u := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	client := util.CreateHTTPClient()
	resp, err := client.Get(u)
	//这里不应该进行错误检查，因为错误检查已经统一放在了util.GetHTTPResponse里面了
	//这里会阻断下面的处理
	if err != nil {
		return nil, err
	}

	var result ReleaseResp
	err = util.GetHTTPResponse(resp, u, err, &result)
	if err != nil {
		return nil, err
	}

	return newRelease(&result), err
}

func newRelease(from *ReleaseResp) *Release {
	release := &Release{
		tagName: from.TagName,
		assets:  make([]Asset, len(from.Assets)),
	}
	for i, fromAsset := range from.Assets {
		release.assets[i] = Asset{
			name: fromAsset.Name,
			url:  fromAsset.BrowserDownloadURL,
		}
	}
	return release
}
