package main

import "os"

func downloadYunzaiFromGitee() {
    _, err := os.Stat("./Yunzai-bot")
    if err == nil {
        printWithEmptyLine("检测到当前目录下已存在 Yunzai-bot ，请问是否需要重新下载？(是:y 返回菜单:n)")
        userChoice := ReadChoice("y", "n")
        if userChoice == "y" {
            //删除文件夹
            os.RemoveAll("./Yunzai-bot")
        }
        if userChoice == "n" {
            return
        }
    }
    executeCmd("git clone --depth 1 -b main https://gitee.com/yoimiya-kokomi/Yunzai-Bot.git", "开始下载云崽...", "下载云崽成功！")
    //进入Yunzai-Bot文件夹
    os.Chdir("./Yunzai-Bot")
    b2 := checkCommand("pnpm -v")
    if !b2 {
        executeCmd("npm install pnpm -g --registry=https://registry.npmmirror.com", "开始安装 pnpm ...", "安装 pnpm 成功！")
    }
    executeCmd("pnpm config set registry https://registry.npmmirror.com")
    executeCmd("pnpm config set puppeteer_download_host=https://registry.npmmirror.com", "开始设置 pnpm 镜像源...", "设置 pnpm 镜像源成功！")
    executeCmd("pnpm install -P", "开始安装云崽依赖", "安装云崽依赖成功！")
    os.Chdir("..")
}
