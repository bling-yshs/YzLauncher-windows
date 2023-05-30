package main

import (
	"encoding/json"
	"io"
	"os"
)

func downloadYunzaiFromGitee() {
	_, err := os.Stat("./Yunzai-bot")
	if err == nil {
		printWithEmptyLine("检测到当前目录下已存在 Yunzai-bot ，请问是否需要重新下载？(是:y 返回菜单:n)")
		userChoice := ReadChoice("y", "n")
		if userChoice == "y" {
			printRedInfo("重新下载云崽会移除当前目录下的 Yunzai-bot 文件夹，云崽的数据将会被全部删除，且不可恢复，请再次确认是否继续？(是:y 返回菜单:n)")
			userChoice := ReadChoice("y", "n")
			if userChoice == "n" {
				return
			}
			//删除文件夹
			printWithEmptyLine("正在删除 Yunzai-bot 文件夹...")
			os.RemoveAll("./Yunzai-bot")
		}
		if userChoice == "n" {
			return
		}
	}
	printWithEmptyLine("请选择云崽类别：\n1.官方云崽\n2.喵喵云崽")
	userChoice := ReadChoice("1", "2")
	if userChoice == "1" {
		executeCmd("git clone --depth 1 -b main https://gitee.com/yoimiya-kokomi/Yunzai-Bot.git", "开始下载官方云崽...", "下载官方云崽成功！")
	}
	var installMiao = false
	if userChoice == "2" {
		executeCmd("git clone --depth 1 -b master https://gitee.com/yoimiya-kokomi/Miao-Yunzai.git", "开始下载喵喵云崽...", "下载喵喵云崽成功！")
		//将Miao-Yunzai文件夹重命名为Yunzai-Bot
		os.Rename("./Miao-Yunzai", "./Yunzai-Bot")
		installMiao = true
	}
	//进入Yunzai-Bot文件夹
	wd.changeToYunzai()
	b2 := checkCommand("pnpm -v")
	if !b2 {
		executeCmd("npm install pnpm -g --registry=https://registry.npmmirror.com", "开始安装 pnpm ...", "安装 pnpm 成功！")
	}
	executeCmd("pnpm config set registry https://registry.npmmirror.com", "开始设置 pnpm 镜像源...")
	executeCmd("pnpm config set PUPPETEER_DOWNLOAD_HOST=https://npmmirror.com/mirrors", "开始设置 puppeteer 镜像源...")
	delDep()
	if windowsVersion < 10 {
		executeCmd("pnpm install puppeteer@19.7.3 -w")
		executeCmd("node ./node_modules/puppeteer/install.js")
		printWithEmptyLine("正在下载修复文件...")
		_, err := os.Stat("./renderers")
		if err != nil {
			downloadFile("https://gitee.com/bling_yshs/YzLauncher-windows/raw/master/NonProjectRequirements/WinServer2012YunzaiFix/Official/puppeteer.js", "./lib/puppeteer")
		} else {
			downloadFile("https://gitee.com/bling_yshs/YzLauncher-windows/raw/master/NonProjectRequirements/WinServer2012YunzaiFix/Miao/config_default.yaml", "./renderers/puppeteer")
			downloadFile("https://gitee.com/bling_yshs/YzLauncher-windows/raw/master/NonProjectRequirements/WinServer2012YunzaiFix/Miao/puppeteer.js", "./renderers/puppeteer/lib")
		}
	} else {
		executeCmd("pnpm install puppeteer@19.8.3 -w")
		executeCmd("node ./node_modules/puppeteer/install.js")
		printWithEmptyLine("正在下载修复文件...")
		_, err := os.Stat("./renderers")
		if err == nil {
			downloadFile("https://gitee.com/bling_yshs/YzLauncher-windows/raw/master/NonProjectRequirements/Win10YunzaiFix/Miao/config_default.yaml", "./renderers/puppeteer")
			downloadFile("https://gitee.com/bling_yshs/YzLauncher-windows/raw/masterNonProjectRequirements/Win10YunzaiFix/Miao/puppeteer.js", "./renderers/puppeteer/lib")
		}
	}
	executeCmd("node ./node_modules/puppeteer/install.js")
	executeCmd("pnpm install", "开始安装云崽依赖", "安装云崽依赖成功！")
	if installMiao {
		printWithEmptyLine("开始下载必须的喵喵插件...")
		installMiaoPlugin()
	}
}

type PackageJSON struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Author       string            `json:"author"`
	Description  string            `json:"description"`
	Main         string            `json:"main"`
	Type         string            `json:"type"`
	Scripts      map[string]string `json:"scripts"`
	Dependencies map[string]string `json:"dependencies"`
	DevDeps      map[string]string `json:"devDependencies"`
	Imports      map[string]string `json:"imports"`
}

func delDep() {
	// 读取原始文件内容
	filePath := "./package.json"
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	// 解析 JSON
	var pkg PackageJSON
	if err := json.Unmarshal(bytes, &pkg); err != nil {
		panic(err)
	}
	// 修改内容
	delete(pkg.Dependencies, "puppeteer")
	// 重新编码为 JSON
	newBytes, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		panic(err)
	}
	// 覆写回去
	if err := os.WriteFile(filePath, newBytes, os.ModePerm); err != nil {
		panic(err)
	}
}
