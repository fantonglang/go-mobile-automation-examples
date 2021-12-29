package main

import (
	"log"
	"runtime"
	"time"

	"github.com/fantonglang/go-mobile-automation/apis"
)

func getDevice() *apis.Device {
	if runtime.GOARCH == "arm" {
		return apis.NewNativeDevice()
	}
	//此处 c574dd45 替换成你自己的手机的device id
	_d, err := apis.NewHostDevice("c574dd45")
	if err != nil {
		log.Println("101: failed connecting to device")
		return nil
	}
	return _d
}

/*
 * 这个栗子🌰依赖以下app以及工具：
 * 1. 抖音 18.9.0，在一加手机，河马云手机荣耀2有完整跑过，运行速度快且稳定
 *		apk下载链接为 (https://caiyun.139.com/m/i?014Mchni6ugFK) 提取码 (xeHG) 注意将后缀改回.apk
 *    $ adb install douyin.apk
 * 2. atx-agent(https://github.com/openatx/atx-agent),
 * 		工具下载地址:(https://github.com/openatx/atx-agent/releases) 选择armv7，一般手机都兼容armv7，注意如果是PC模拟器比如雷电虚拟化出来的手机，使用相应的cpu arch
 *    解压tar包，将可执行文件推到手机中，并执行可执行文件 - 后台启动API服务
 *		$ adb push atx-agent /data/local/tmp
 *		$ adb shell chmod 755 /data/local/tmp/atx-agent
 *		# launch atx-agent in daemon mode
 *		$ adb shell /data/local/tmp/atx-agent server -d
 *		# stop already running atx-agent and start daemon
 *		$ adb shell /data/local/tmp/atx-agent server -d --stop
 * 3. android-uiautomator-server(https://github.com/openatx/android-uiautomator-server)
 *		工具下载地址:(https://github.com/openatx/android-uiautomator-server/releases)
 *		下载app-uiautomator-test.apk, app-uiautomator.apk, 并且安装到手机
 *		$ adb install	app-uiautomator-test.apk
 *		$ adb install	app-uiautomator.apk
 * 		在你的应用列表中会出现一个应用ATX（出租车图标），打开应用，点击“启动UIAUTOMATOR”按钮，点击“开启悬浮窗”按钮
 * 4. 开发工具 - xpath/属性拾取
 *		安装python3(version 3.6+)，并执行下面命令安装weditor
 *		$ pip3 install -U weditor
 * 5. xpath/属性拾取 - 启动weditor
 *		在命令行输入weditor，chrome上会自动弹出页面，inspector跑在17310端口，注意使用的时候把代理(比如SS/V2RAY)关了
 *		点击页面顶部的“Dump Hierarchy”按钮来刷新手机UI的控件树，网页上点击任何想自动化的控件，获取它的信息
 * 6. 开发
 *		# 先在项目文件夹中执行下面指令来初始化go module
 *		$ go mod init [模块名 - 自己取 - 任意]
 *		# 拉取SDK依赖
 *		$ go get github.com/fantonglang/go-mobile-automation
 *
 *		新建文件main.go在里面写我们要自动化的流程，参考本栗子🌰
 *		执行go run main.go
 * 7. 部署
 *		# linux/arm方式cross build流程脚本
 *		$ GOOS=linux GOARCH=arm go build
 *		# douyin-luo-live是上一步编译出来的可执行文件，替换成你编译出来可执行文件的名字
 *		$ adb push douyin-luo-live /data/local/tmp
 *		$ adb shell chmod 755 /data/local/tmp/douyin-luo-live
 * 8. 执行
 *		# douyin-luo-live是上一步编译出来的可执行文件，替换成你编译出来可执行文件的名字
 *		$ adb shell /data/local/tmp/douyin-luo-live
 *
 *		开始执行之后就不需要ADB连接了
 */
func main() {
	log.Println("001: Process Running")
	d := getDevice()
	if d == nil {
		return
	}

	// 1. start douyin
	// 如果抖音已经在运行，确保它已经关闭，然后打开抖音
	d.Shell(`am force-stop com.ss.android.ugc.aweme`)
	d.Shell(`am start -n "com.ss.android.ugc.aweme/.main.MainActivity"`)
	// 抖音界面顶部有三个tab按钮，本地（你所在的城区名称）关注 推荐，等待本地出现，然后点击
	localBtn := d.XPath(`//*[@resource-id="com.ss.android.ugc.aweme:id/nd3"]`).Wait(time.Minute)
	if localBtn == nil {
		log.Println("102: failed can't find local button")
		return
	}
	localBtn.Click()
	time.Sleep(time.Second) // 慢慢来，等1s，你也可以不等
	// 2. find search button and click
	// 抖音界面顶部右上角有个搜索按钮，等待按钮出现并点击
	searchBtn := d.XPath(`//*[@resource-id="com.ss.android.ugc.aweme:id/fad"]`).Wait(time.Minute)
	if searchBtn == nil {
		log.Println("103: failed find search button")
		return
	}
	searchBtn.Click()
	// 3. find search box and click
	// 点击搜索按钮之后，进入搜索页面，等待找到搜索输入框，并点击它
	inputBox := d.XPath(`//*[@resource-id="com.ss.android.ugc.aweme:id/fl_intput_hint_container"]`).Wait(time.Minute)
	if inputBox == nil {
		log.Println("104: failed search box doesn't exist")
		return
	}
	inputBox.Click()
	// 4. type live room name
	// 在输入框中输入直播间名字，比如“罗永浩”，希望你在看到这行代码的时候他还没因为逃税被封号
	time.Sleep(time.Second)
	d.SendKeys("罗永浩", true)
	// 5. in the dropdown list, find the live broadcasting item, and click
	// 输入直播间名字之后，下面会弹一个列表，如果是正在直播的账号，会有个红色的跳动的按钮，等待并点击按钮进入直播列表页
	liveBroadingSign := d.XPath(`//*[@resource-id="com.ss.android.ugc.aweme:id/k-b"]`).Wait(time.Minute)
	if liveBroadingSign == nil {
		log.Println("105: failed not live broadcasting")
		return
	}
	liveBroadingSign.Click()
	// 6. find the broadcasting small screen
	// 找到正在直播的小窗口并点击，但是有时候xpath并不是总是相同的，它们又个共性就是文字中带有“直播中”
	liveBroadingSmallScreen := d.XPath(`//*[@resource-id="com.ss.android.ugc.aweme:id/f-1"]`).Wait(5 * time.Second)
	if liveBroadingSmallScreen == nil {
		liveBroadingSmallScreen = d.XPath(`//*[contains(@text, "直播中")]`).Wait(5 * time.Second)
	}
	if liveBroadingSmallScreen == nil {
		log.Println("106: failed get broadcasting small screen")
		return
	}
	liveBroadingSmallScreen.Click()

	// 7. 现在你就到了直播页面，开心了吧～
	log.Println("002: Process Finished")
}
