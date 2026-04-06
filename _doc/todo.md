
#### v0.1.2



#### v0.1.1   // 04010 End
- [] DSL v0.1, 继续完善语法文档和定版
- [] chrome 脚本操作自动化功能完整，能应对主流业务场景
- [] http+host 脚本操作自动化功能完成，能应对主流业务场景
- 4~5月写blog,录制视，6月看数据最后总结


#### v0.0.14
- []chrome  ServiceWorker ： 服务任务
    []ServiceWorker.deliverPushMessage  推送消息
    []ServiceWorker.disable  禁用
    []ServiceWorker.dispatchPeriodicSyncEvent  定时同步
    []ServiceWorker.dispatchSyncEvent  推送消息
    []ServiceWorker.enable  启用
    []ServiceWorker.setForceUpdateOnPageLoad  强制更新
    []ServiceWorker.skipWaiting  跳过等待
    []ServiceWorker.startWorker  启动工作
    []ServiceWorker.stopAllWorkers  停止所有工作
    []ServiceWorker.stopWorker  停止工作
    []ServiceWorker.unregister  取消注册
    []ServiceWorker.updateRegistration  更新注册

- []chrome  WebAudio ： 此域名允许查看 Web Audio API。
    []WebAudio.disable
    []WebAudio.enable
    []WebAudio.getRealtimeData   获取实时数据

- []chrome  WebAuthn ： 该域允许配置虚拟身份验证器来测试 WebAuthn API。   // 04010 End
    []WebAuthn.addCredential  添加凭证
    []WebAuthn.addVirtualAuthenticator  添加虚拟身份验证器
    []WebAuthn.clearCredentials  清空凭证
    []WebAuthn.disable  禁用
    []WebAuthn.enable  启用
    []WebAuthn.getCredential  获取凭证
    []WebAuthn.getCredentials  获取凭证
    []WebAuthn.removeCredential  删除凭证
    []WebAuthn.removeVirtualAuthenticator  删除虚拟身份验证器
    []WebAuthn.setAutomaticPresenceSimulation  设置自动存在模拟
    []WebAuthn.setCredentialProperties  设置凭证属性
    []WebAuthn.setResponseOverrideBits  设置响应覆盖位
    []WebAuthn.setUserVerified  设置用户验证

- 执行js代码的能力
- 更多示例
  1.
  2.
  3.
  4.
  5.
  6.
  7.

// 0409 End

#### v0.0.13
- []chrome  Profiler ： 分析器域
    []Profiler.disable  禁用
    []Profiler.enable  启用
    []Profiler.getBestEffortCoverage  获取最佳效果覆盖
    []Profiler.setSamplingInterval  设置采样间隔
    []Profiler.start  
    []Profiler.startPreciseCoverage  开始精确覆盖
    []Profiler.stop  停止
    []Profiler.stopPreciseCoverage  停止精确覆盖
    []Profiler.takePreciseCoverage  获取精确覆盖

- []chrome  PWA ： 该域允许与浏览器交互以控制 PWA。
    []PWA.changeAppUserSettings  更改 PWA 用户设置
    []PWA.getOsAppState  获取操作系统应用状态
    []PWA.install  安装 PWA
    []PWA.launch  启动 PWA
    []PWA.launchFilesInApp   启动文件
    []PWA.openCurrentPageInApp  在应用中打开当前页面
    []PWA.uninstall  卸载 PWA

- []chrome  Runtime ： 运行时域通过远程求值和镜像对象公开 JavaScript 运行时环境。
    []Runtime.addBinding 如果 executionContextId 为空，则会在所有被检查上下文的全局对象（包括之后创建的上下文）上添加具有给定名称的绑定，并且绑定在重新加载后仍然存在。
    []Runtime.awaitPromise 使用给定的 Promise 对象 ID 向 Promise 添加处理程序。
    []Runtime.callFunctionOn 调用给定对象上具有给定声明的函数。结果的对象组继承自目标对象。
    []Runtime.compileScript 编译表达式。
    []Runtime.disable 禁用执行上下文创建的报告。
    []Runtime.discardConsoleEntries  丢弃收集到的异常和控制台 API 调用。
    []Runtime.enable  启用事件报告功能，用于报告执行上下文的创建情况executionContextCreated。
    []Runtime.evaluate  对全局对象求表达式的值。
    []Runtime.getProperties  返回给定对象的属性。结果的对象组继承自目标对象。
    []Runtime.globalLexicalScopeNames  返回全局作用域中的所有 let、const 和 class 变量。
    []Runtime.queryObjects  返回指定对象组中的对象。
    []Runtime.releaseObject  释放对象组中指定的对象。
    []Runtime.releaseObjectGroup  释放对象组。
    []Runtime.removeBinding  删除绑定。
    []Runtime.runIfWaitingForDebugger  运行等待的调试器。
    []Runtime.runScript  运行脚本。
    []Runtime.setAsyncCallStackDepth  设置异步调用堆栈深度。
    []Runtime.getExceptionDetails  获取异常详细信息。
    []Runtime.getHeapUsage  获取堆使用情况。
    []Runtime.getIsolateId  获取隔离 ID。
    []Runtime.setCustomObjectFormatterEnabled  设置自定义对象格式化程序是否启用。
    []Runtime.setMaxCallStackSizeToCapture  设置调用堆栈大小以捕获。
    []Runtime.terminateExecution  终止执行。

- []chrome  Security ： 安全域   // 0409 End
    []Security.disable
    []Security.enable
    []Security.setIgnoreCertificateErrors  处理触发certificateError事件的证书错误。

- []chrome  Storage ： 存储
    []Storage.clearCookies  清除所有 Cookie。
    []Storage.clearDataForOrigin  清除指定源的 Cookie。
    []Storage.clearDataForStorageKey  清除指定 StorageKey 的 Cookie。
    []Storage.getCookies  获取所有 Cookie。
    []Storage.getUsageAndQuota  获取指定源的 Cookie 使用情况。
    []Storage.setCookies  设置 Cookie。
    []Storage.setProtectedAudienceKAnonymity  设置受保护的受众的 K-匿名。
    []Storage.trackCacheStorageForOrigin  跟踪指定源的 CacheStorage。
    []Storage.trackCacheStorageForStorageKey  跟踪指定 StorageKey 的 CacheStorage。
    []Storage.trackIndexedDBForOrigin  跟踪指定源的 IndexedDB。
    []Storage.trackIndexedDBForStorageKey  跟踪指定 StorageKey 的 IndexedDB。
    []Storage.untrackCacheStorageForOrigin  停止跟踪指定源的 CacheStorage。
    []Storage.untrackCacheStorageForStorageKey  停止跟踪指定 StorageKey 的 CacheStorage。
    []Storage.untrackIndexedDBForOrigin  停止跟踪指定源的 IndexedDB。
    []Storage.untrackIndexedDBForStorageKey  停止跟踪指定 StorageKey 的 IndexedDB。
    []Storage.clearSharedStorageEntries  删除指定存储桶的共享存储条目。
    []Storage.clearTrustTokens  删除所有共享存储令牌。
    []Storage.deleteSharedStorageEntry  删除指定存储桶的共享存储条目。
    []Storage.deleteStorageBucket  删除指定存储桶。
    []Storage.getAffectedUrlsForThirdPartyCookieMetadata  获取指定源的受影响的 URL。
    []Storage.getInterestGroupDetails  获取指定兴趣组详情。
    []Storage.getRelatedWebsiteSets  获取相关网站集。
    []Storage.getSharedStorageEntries  获取指定存储桶的共享存储条目。
    []Storage.getSharedStorageMetadata  获取指定存储桶的共享存储元数据。
    []Storage.getStorageKey   获取指定存储桶的存储密钥。
    []Storage.getTrustTokens  获取所有共享存储令牌。
    []Storage.overrideQuotaForOrigin  覆盖指定源的存储配额。
    []Storage.resetSharedStorageBudget  重置共享存储的预算。
    []Storage.runBounceTrackingMitigations  运行 Bunce Tracking 拦截。
    []Storage.sendPendingAttributionReports  发送挂起的ATTRIBUTION报告。
    []Storage.setAttributionReportingLocalTestingMode  设置ATTRIBUTION报告本地测试模式。
    []Storage.setAttributionReportingTracking  设置ATTRIBUTION报告跟踪。
    []Storage.setInterestGroupAuctionTracking  设置兴趣组拍卖跟踪。
    []Storage.setInterestGroupTracking  设置兴趣组跟踪。
    []Storage.setSharedStorageEntry  设置共享存储条目。
    []Storage.setSharedStorageTracking  设置共享存储跟踪。
    []Storage.setStorageBucketTracking  设置存储ucket跟踪。

- []chrome  Tethering ： 域定义了浏览器端口绑定的方法和事件。
    []Tethering.bind
    []Tethering.unbind

- []chrome  Tracing ： 追踪
    []Tracing.end 
    []Tracing.start
    []Tracing.getCategories  获取可用的跟踪类别
    []Tracing.getTrackEventDescriptor  获取跟踪事件描述符
    []Tracing.recordClockSyncMarker  记录时钟同步标记
    []Tracing.requestMemoryDump  请求内存转储

- [] 改测试的bug和优化
- [] bug和优化验收
- [] 更新文档   // 0407 End

#### v0.0.12
- [ok]chrome  Media ： 该域允许对媒体元素进行详细检查。
    [ok]Media.disable
    [ok]Media.enable

- [ok]chrome  Memory ： 内存相关
    [ok]Memory.forciblyPurgeJavaScriptMemory  通过清除 V8 内存来模拟 OomIntervention。
    [ok]Memory.getAllTimeSamplingProfile  获取自渲染进程启动以来收集的本地内存分配概况。
    [ok]Memory.getBrowserSamplingProfile  获取自浏览器进程启动以来收集的本地内存分配概况。
    [ok]Memory.getDOMCounters  返回当前 DOM 对象计数器。
    [ok]Memory.getDOMCountersForLeakDetection 在准备渲染器进行泄漏检测后，返回 DOM 对象计数器。
    [ok]Memory.getSamplingProfile  检索自上次 startSampling调用以来收集的本地内存分配配置文件。
    [ok]Memory.prepareForLeakDetection  通过终止工作进程、停止拼写检查器、删除非必要的内部缓存、运行垃圾回收等方式，为泄漏检测做好准备。
    [ok]Memory.setPressureNotificationsSuppressed  启用/禁用所有进程中的内存压力通知抑制。
    [ok]Memory.simulatePressureNotification  模拟所有进程的内存压力通知。
    [ok]Memory.startSampling  开始收集本地内存配置文件。
    [ok]Memory.stopSampling  停止收集本地内存配置文件。

- []chrome  Network ： 网络域允许跟踪页面的网络活动。它公开有关 HTTP、文件、数据和其他请求和响应的信息，包括它们的标头、正文、时间等。
    []Network.clearBrowserCache  清除浏览器缓存。
    []Network.clearBrowserCookies  清除浏览器cookie。
    []Network.deleteCookies  删除名称和 URL 或域/路径/分区密钥对匹配的浏览器 cookie。
    []Network.disable  禁用网络域。
    []Network.enable  启用网络域。
    []Network.getCookies  返回当前 URL 的所有浏览器 Cookie。
    []Network.getRequestPostData  返回请求中发送的 POST 数据。如果请求中未发送任何数据，则返回错误。
    []Network.getResponseBody  返回针对给定请求提供的内容。
    []Network.setBypassServiceWorker  切换是否忽略每个请求中的 Service Worker。
    []Network.setCacheDisabled  切换是否忽略缓存。如果启用此选项true，则不会使用缓存。
    []Network.setCookie  设置 Cookie。
    []Network.setCookies  设置多个 Cookie。
    []Network.setExtraHTTPHeaders  指定是否始终随此页面发出的请求发送额外的 HTTP 标头。
    []Network.setUserAgentOverride  允许使用给定的字符串覆盖用户代理。
    []Network.clearAcceptedEncodingsOverride   清除 setAcceptedEncodings 设置的已接受编码
    []Network.configureDurableMessages  配置将响应体存储在渲染器外部，以便跨进程导航时响应体仍然有效。
    []Network.emulateNetworkConditionsByRule  使用 URL 匹配模式为单个请求启用网络条件模拟。
    []Network.enableDeviceBoundSessions  设置跟踪设备绑定会话并获取初始会话集。
    []Network.enableReportingApi  启用报表 API 的跟踪功能，报表 API 生成的事件现在将传递给客户端。
    []Network.fetchSchemefulSite  获取特定来源的阴谋网站
    []Network.getCertificate 返回 DER 编码的证书。
    []Network.getResponseBodyForInterception  返回针对当前拦截的请求提供的内容。
    []Network.getSecurityIsolationStatus  返回有关 COEP/COOP 隔离状态的信息。
    []Network.loadNetworkResource  获取资源并返回其内容。
    []Network.overrideNetworkState  覆盖 navigator.onLine 和 navigator.connection 的状态。
    []Network.replayXHR  此方法会发送一个与原始 XMLHttpRequest 完全相同的新请求。
    []Network.searchInResponseBody  在响应内容中搜索指定的字符串。
    []Network.setAcceptedEncodings  设置可接受的内容编码列表。空列表表示不接受任何编码。
    []Network.setAttachDebugStack  指定是否在请求中附加页面脚本堆栈 ID
    []Network.setBlockedURLs 阻止URL加载。
    []Network.setCookieControls  设置第三方 Cookie 访问控制。页面需要重新加载才能生效。
    []Network.streamResourceContent 启用对给定请求 ID 的响应进行流式传输。如果启用，dataReceived 事件将包含在流式传输期间接收到的数据。
    []Network.takeResponseBodyForInterceptionAsStream  返回指向表示响应体的流的句柄。
    
- []chrome  Page ： 与被检查页面相关的操作和事件属于页面域。  // 0404 End
    []Page.addScriptToEvaluateOnNewDocument  在创建每一帧时（在加载帧的脚本之前），对给定的脚本进行评估。
    []Page.bringToFront  将页面置于最前面（激活选项卡）。
    []Page.captureScreenshot  截取页面屏幕截图。
    []Page.close  关闭当前页面。
    []Page.createIsolatedWorld  创建一个新的 isolatedWorld 并返回其 ID。
    []Page.disable  禁用性能域。
    []Page.enable  启用性能域。
    []Page.getAppManifest  获取当前文档的已处理清单。此 API 始终等待清单加载完成。
    []Page.getFrameTree   返回当前帧树结构。
    []Page.getLayoutMetrics  返回与页面布局相关的指标，例如视口边界/缩放比例。
    []Page.getNavigationHistory  返回当前页面的导航历史记录。
    []Page.handleJavaScriptDialog 接受或关闭 JavaScript 发起的对话框（alert、confirm、prompt 或 onbeforeunload）。
    []Page.navigate 将当前页面导航到指定的URL。
    []Page.navigateToHistoryEntry  将当前页面导航到指定的历史记录条目。
    []Page.printToPDF  以PDF格式打印页面。
    []Page.reload   重新加载指定页面，可选择忽略缓存。
    []Page.removeScriptToEvaluateOnNewDocument  从列表中移除指定的脚本。
    []Page.resetNavigationHistory  重置当前页面的导航历史记录。
    []Page.setBypassCSP 启用页面内容安全策略绕过功能。
    []Page.setDocumentContent  将给定的标记设置为文档的 HTML。
    []Page.setInterceptFileChooserDialog 拦截文件选择器请求并将控制权转移给协议客户端。
    []Page.setLifecycleEventsEnabled 控制页面是否会发出生命周期事件。
    []Page.stopLoading  页面停止加载
    []Page.addCompilationCache  为给定的 URL 创建编译缓存。编译缓存不会在跨进程导航后保留。
    []Page.captureSnapshot  返回页面快照的字符串形式。对于 MHTML 格式，序列化内容包括 iframe、Shadow DOM、外部资源和元素内联样式。
    []Page.clearCompilationCache  清除已初始化的编译缓存。
    []Page.crash  IO线程上的渲染器崩溃，生成小型转储文件。
    []Page.generateTestReport 生成测试报告。
    []Page.getAdScriptAncestry  获取广告脚本的祖先。
    []Page.getAnnotatedPageContent 获取主框架的带注释页面内容。这是一个实验性命令，可能会有所更改。
    []Page.getAppId   返回唯一的（PWA）应用 ID。仅当启用功能标志“WebAppEnableManifestId”时才返回值。
    []Page.getInstallabilityErrors  获取当前页面的安装错误。
    []Page.getOriginTrials 在给定帧上获取 Origin Trials。
    []Page.getPermissionsPolicyState  获取给定帧的权限策略状态。
    []Page.getResourceContent  返回给定资源的内容。
    []Page.getResourceTree  返回当前帧/资源树结构。
    []Page.produceCompilationCache  请求后端为指定的脚本生成编译缓存。
    []Page.screencastFrameAck  确认前端已收到屏幕录制帧。
    []Page.searchInResource  在资源内容中搜索给定的字符串。
    []Page.setAdBlockingEnabled  在所有网站上启用 Chrome 的实验性广告过滤器。
    []Page.setFontFamilies  设置通用字体系列。
    []Page.setFontSizes  设置通用字体大小。
    []Page.setPrerenderingAllowed  手动启用/禁用预渲染。 
    []Page.setRPHRegistrationMode  自定义处理程序 API 的扩展
    []Page.setSPCTransactionMode  设置安全支付确认交易模式。
    []Page.setWebLifecycleState  尝试更新页面的 Web 生命周期状态。
    []Page.startScreencast  开始使用screencastFrame事件发送每一帧。
    []Page.stopScreencast  停止使用screencastFrame事件发送每一帧。
    []Page.waitForDebugger  暂停页面执行。可以使用通用的 Runtime.runIfWaitingForDebugger 恢复执行。
    
- [ok]chrome  Performance ： 性能域
    [ok]Performance.disable  禁用指标收集和报告功能。
    [ok]Performance.enable  启用指标收集和报告功能。
    [ok]Performance.getMetrics  获取运行时指标的当前值

- [ok]chrome  PerformanceTimeline ： 按照https://w3c.github.io/performance-timeline/#dom-performanceobserver中的规定，报告性能时间线事件
    [ok]PerformanceTimeline.enable  之前已缓冲的事件会在方法返回之前报告

- [ok]chrome  Preload ： 预加载域
    [ok]Preload.disable  
    [ok]Preload.enable  

- [] 改测试的bug和优化
    
- [] 更新文档   // 0405 End

#### v0.0.11
- [ok]chrome HeadlessExperimental : 此域提供仅在无头模式下支持的实验性命令。
    [ok]HeadlessExperimental.beginFrame  向目标发送 BeginFrame 消息，并在帧完成后返回。可选择捕获生成的帧的屏幕截图。要求创建目标时启用了 BeginFrameControl。设计用于与 `--run-all-compositor-stages-before-draw` 参数配合使用

- [ok]chrome HeapProfiler : 堆分析器域
    [ok]HeapProfiler.addInspectedHeapObject  允许控制台通过 $x 引用具有给定 id 的节点
    [ok]HeapProfiler.collectGarbage  垃圾回收
    [ok]HeapProfiler.disable  禁用堆分析器
    [ok]HeapProfiler.enable  启用堆分析器
    [ok]HeapProfiler.getHeapObjectId  返回给定对象在堆中的唯一标识符
    [ok]HeapProfiler.getObjectByHeapObjectId  返回给定对象在堆中的唯一标识符
    [ok]HeapProfiler.getSamplingProfile  获取采样配置文件
    [ok]HeapProfiler.startSampling  开始堆采样
    [ok]HeapProfiler.startTrackingHeapObjects  开始跟踪对象分配
    [ok]HeapProfiler.stopSampling   停止堆采样
    [ok]HeapProfiler.stopTrackingHeapObjects  停止跟踪对象分配
    [ok]HeapProfiler.takeHeapSnapshot   获取堆快照

- [ok]chrome  Inspector ： 检查域
    [ok]Inspector.disable 禁用检查器域通知。
    [ok]Inspector.enable 启用检查器域通知。

- [ok]chrome  IO ： 对 DevTools 生成的流进行输入/输出操作。  // 0402 End
    [ok]IO.close  关闭数据流，丢弃所有临时备份存储。
    [ok]IO.read  阅读一段流媒体内容

- [ok]chrome  LayerTree ： 层树
    [ok]LayerTree.compositingReasons  说明合成给定图层的原因。
    [ok]LayerTree.disable 禁用堆肥树检查。
    [ok]LayerTree.enable  启用堆肥树检查功能。
    [ok]LayerTree.loadSnapshot 返回快照标识符。
    [ok]LayerTree.makeSnapshot  返回图层快照标识符。
    [ok]LayerTree.profileSnapshot  获取堆肥树性能数据。
    [ok]LayerTree.releaseSnapshot  释放快照。
    [ok]LayerTree.replaySnapshot  重新播放后端捕获的层快照。
    [ok]LayerTree.snapshotCommandLog  返回指定快照的命令日志。

- [ok]chrome  Log ： 提供对日志条目的访问权限。
    [ok]Log.clear 清除日志。
    [ok]Log.disable 禁用日志。
    [ok]Log.enable 启用日志。
    [ok]Log.startViolationsReport 启动违规报告。
    [ok]Log.stopViolationsReport 停止违规报告。

- [ok]chrome  Overlay ： 叠加域 该域提供与在被检查页面上绘制图形相关的各种功能。
    [ok]Overlay.disable 禁用叠加。
    [ok]Overlay.enable 启用叠加。
    [ok]Overlay.getGridHighlightObjectsForTest 用于持久网格测试。
    [ok]Overlay.getHighlightObjectForTest 用于测试。
    [ok]Overlay.getSourceOrderHighlightObjectForTest  用于源顺序查看器测试。
    [ok]Overlay.hideHighlight 隐藏所有高亮显示。
    [ok]Overlay.highlightNode 高亮显示具有指定 ID 或指定 JavaScript 对象包装器的 DOM 节点。
    [ok]Overlay.highlightQuad 高亮显示给定四边形区域。坐标系是相对于主框架视口的绝对坐标。
    [ok]Overlay.highlightRect 高亮显示给定的矩形区域。坐标是相对于主框架视口的绝对坐标。
    [ok]Overlay.highlightSourceOrder 高亮显示具有给定 id 或给定 JavaScript 对象包装器的 DOM 节点的子节点的源顺序。
    [ok]Overlay.setInspectMode  进入“检查”模式。在此模式下，用户鼠标悬停的元素会高亮显示。
    [ok]Overlay.setPausedInDebuggerMessage  暂停 JavaScript 运行，并显示给定的消息。
    [ok]Overlay.setShowAdHighlights  设置是否显示广告高亮。
    [ok]Overlay.setShowContainerQueryOverlays  设置是否显示容器查询高亮。
    [ok]Overlay.setShowDebugBorders  设置是否显示调试边框。
    [ok]Overlay.setShowFlexOverlays  设置是否显示 Flex 高亮。
    [ok]Overlay.setShowFPSCounter  设置是否显示 FPS 计数器。
    [ok]Overlay.setShowGridOverlays  设置是否显示网格高亮。
    [ok]Overlay.setShowHinge  设置是否显示 3D 旋转轴。
    [ok]Overlay.setShowInspectedElementAnchor  设置是否显示被检查元素锚点。
    [ok]Overlay.setShowIsolatedElements  设置是否显示隔离元素。
    [ok]Overlay.setShowLayoutShiftRegions  设置是否显示布局偏移区域。
    [ok]Overlay.setShowPaintRects  设置是否显示绘制矩形。
    [ok]Overlay.setShowScrollBottleneckRects    设置是否显示滚动 bottleneck 矩形。
    [ok]Overlay.setShowScrollSnapOverlays  设置是否显示滚动 snap 覆盖。
    [ok]Overlay.setShowViewportSizeOnResize  设置是否显示视图大小。
    [ok]Overlay.setShowWindowControlsOverlay  设置是否显示窗口控制栏。

- [] 改测试的bug和优化
    [] chrome HeadlessExperimental 测试和完善
    [] chrome HeapProfiler 测试和完善
    [] chrome Inspector 测试和完善
    [] chrome IO 测试和完善
    [] chrome LayerTree 测试和完善
    [] chrome Log 测试和完善
    [] chrome Overlay 测试和完善
- [] bug和优化验收
- [] 更新文档   // 0404 End

#### v0.0.10
- [ok]chrome BackgroundService ：  定义后台 Web 平台功能的事件。
    [ok]BackgroundService.clearEvents  清除该服务的所有已存储数据。
    [ok]BackgroundService.setRecording  设置服务的录制状态。
    [ok]BackgroundService.startObserving  启用服务的事件更新。
    [ok]BackgroundService.stopObserving  禁用该服务的事件更新。

- [ok]chrome Fetch : 允许客户端使用客户端代码替换浏览器网络层的域。
    [ok]Fetch.continueRequest  继续发送请求，并可选择修改其某些参数。
    [ok]Fetch.continueWithAuth  在 authRequired 事件发生后，继续提供 authChallengeResponse 的请求。
    [ok]Fetch.disable  禁用 fetch 域。
    [ok]Fetch.enable  启用 requestPaused 事件的触发。请求将被暂停，直到客户端调用 failRequest、fulfillRequest 或 continueRequest/continueWithAuth 中的一个。
    [ok]Fetch.failRequest 使请求因指定原因失败。
    [ok]Fetch.fulfillRequest  对请求做出响应。
    [ok]Fetch.getResponseBody  使服务器接收响应正文并将其作为单个字符串返回。
    [ok]Fetch.takeResponseBodyAsStream  返回指向表示响应体的流的句柄。

- [ok]chrome FileSystem : 文件系统域
    [ok]FileSystem.getDirectory 获取目录

- [ok]chrome DOM : 此域公开 DOM 读/写操作。   // 0331 End
    [ok]DOM.describeNode 根据节点 ID 描述节点，无需启用域。不会开始跟踪任何对象，可用于自动化。
    [ok]DOM.disable 禁用指定页面的 DOM 代理。
    [ok]DOM.enable 启用 DOM 代理。
    [ok]DOM.focus 聚焦指定元素。
    [ok]DOM.getAttributes 返回指定节点的属性。
    [ok]DOM.getBoxModel 返回给定节点的盒子。
    [ok]DOM.getDocument  返回根 DOM 节点（以及可选的子树）给调用者。隐式启用当前目标的 DOM 域事件。
    [ok]DOM.getNodeForLocation  返回指定位置的节点 ID。是否返回 nodeId 取决于 DOM 域是否启用。
    [ok]DOM.getOuterHTML  返回节点的 HTML 标记。
    [ok]DOM.hideHighlight  隐藏所有高亮显示。
    [ok]DOM.highlightNode  高亮显示 DOM 节点。
    [ok]DOM.highlightRect  高亮显示给定的矩形。
    [ok]DOM.moveTo  将节点移动到新容器中，并将其放置在给定锚点之前。
    [ok]DOM.querySelector  querySelector在指定节点上执行。
    [ok]DOM.querySelectorAll  querySelectorAll在指定节点上执行。
    [ok]DOM.removeAttribute  从具有给定 id 的元素中移除具有给定名称的属性。
    [ok]DOM.removeNode  删除具有给定 id 的节点。
    [ok]DOM.requestChildNodes  请求将给定 id 的节点的子节点以事件的形式返回给调用者， setChildNodes其中不仅检索直接子节点，而且检索到指定深度的所有子节点。
    [ok]DOM.requestNode  根据 JavaScript 节点对象引用，请求将节点发送给调用者。
    [ok]DOM.resolveNode  解析给定 NodeId 或 BackendNodeId 的 JavaScript 节点对象。
    [ok]DOM.scrollIntoViewIfNeeded  如果指定节点的指定矩形区域尚未可见，则将其滚动到视图中。
    [ok]DOM.setAttributesAsText  设置具有给定 ID 的元素的属性。
    [ok]DOM.setAttributeValue  设置具有给定 id 的元素的属性。
    [ok]DOM.setFileInputFiles  为给定的文件输入元素设置文件。
    [ok]DOM.setNodeName  设置具有给定 id 的节点的节点名称。
    [ok]DOM.setNodeValue  设置具有给定 id 的节点的节点值。
    [ok]DOM.setOuterHTML  设置节点 HTML 标记，返回新的节点 ID。

- [ok]chrome DOMDebugger : DOM调试允许在特定的DOM操作和事件上设置断点。
    [ok]DOMDebugger.getEventListeners  返回给定对象的事件监听器。
    [ok]DOMDebugger.removeDOMBreakpoint  移除使用 . 设置的 DOM 断点setDOMBreakpoint。
    [ok]DOMDebugger.removeEventListenerBreakpoint  移除特定 DOM 事件上的断点。
    [ok]DOMDebugger.removeXHRBreakpoint  移除 XMLHttpRequest 中的断点。
    [ok]DOMDebugger.setDOMBreakpoint  在对 DOM 进行特定操作时设置断点。
    [ok]DOMDebugger.setEventListenerBreakpoint  在特定 DOM 事件上设置断点
    [ok]DOMDebugger.setXHRBreakpoint  在 XMLHttpRequest 上设置断点。

- [ok] chrome IndexedDB : IndexedDB相关的域
    [ok]IndexedDB.clearObjectStore  清除对象存储中的所有条目。
    [ok]IndexedDB.deleteDatabase  删除数据库。
    [ok]IndexedDB.deleteObjectStoreEntries  从对象存储库中删除一系列条目
    [ok]IndexedDB.disable  禁用后端事件。
    [ok]IndexedDB.enable  启用来自后端的事件。
    [ok]IndexedDB.getMetadata  获取对象存储的元数据。
    [ok]IndexedDB.requestData  从对象存储或索引中请求数据。
    [ok]IndexedDB.requestDatabase  请求具有给定名称的数据库到给定框架。
    [ok]IndexedDB.requestDatabaseNames  请求给定安全源的数据库名称。

- [ok]chrome Input ： 输入域
    [ok]Input.cancelDragging  取消页面上所有正在进行的拖动操作。
    [ok]Input.dispatchKeyEvent  向页面发送关键事件。
    [ok]Input.dispatchMouseEvent  向页面发送鼠标事件。
    [ok]Input.dispatchTouchEvent  向页面发送触摸事件。
    [ok]Input.setIgnoreInputEvents  忽略输入事件
    [ok]Input.dispatchDragEvent  将拖拽事件发送到页面中。
    [ok]Input.emulateTouchFromMouseEvent 根据鼠标事件参数模拟触摸事件。
    [ok]Input.imeSetComposition  此方法设置输入法编辑器 (IME) 的当前候选文本。使用 `imeCommitComposition` 提交最终文本。
    [ok]Input.insertText  这种方法模拟插入非按键输入的文本，例如表情符号键盘或输入法编辑器。
    [ok]Input.setInterceptDrags  阻止默认的拖放行为，而是发出Input.dragIntercepted事件。
    
- [] 改测试的bug和优化
    [] chrome BackgroundService 测试和完善
    [] chrome Fetch 测试和完善
    [] chrome FileSystem 测试和完善
    [] chrome DOM 测试和完善
    [] chrome DOMDebugger 测试和完善
    [] chrome IndexedDB 测试和完善
    [] chrome Input 测试和完善
- [] bug和优化验收
- [] 更新文档   // 0403 End

#### v0.0.9
- 支持执行js代码
- chrome CSS ： 此域公开 CSS 的读写操作。
  [ok] CSS.addRule  ruleText在给定样式表的styleSheetId位置插入一条具有给定值的新规则location
  [ok] CSS.collectClassNames  返回指定样式表中的所有类名。
  [ok] CSS.createStyleSheet  在给定的框架中创建一个新的特殊“via-inspector”样式表frameId。
  [ok] CSS.disable  禁用指定页面的 CSS 
  [ok] CSS.enable  为指定页面启用 CSS 
  [ok] CSS.forcePseudoState  确保给定节点在浏览器计算其样式时具有指定的伪类。
  [ok] CSS.forceStartingStyle  确保给定节点处于初始状态。
  [ok] CSS.getBackgroundColors  获取DOM.NodeId背景颜色
  [ok] CSS.getComputedStyleForNode  获取DOM.NodeId的计算样式
  [ok] CSS.getInlineStylesForNode  获取DOM.NodeId的内联样式
  [ok] CSS.getMatchedStylesForNode 获取DOM.NodeId的请求样式
  [ok] CSS.getMediaQueries  返回渲染引擎解析的所有媒体查询。
  [ok] CSS.getPlatformFontsForNode  请求有关我们在给定节点中渲染子文本节点时使用的平台字体的信息。
  [ok] CSS.getStyleSheetText  返回样式表的当前文本内容。
  [ok] CSS.setEffectivePropertyValueForNode 找到给定节点的具有给定 active 属性的规则，并设置该属性的新值。
  [ok] CSS.setKeyframeKey  修改关键帧规则的关键文本。
  [ok] CSS.setMediaText  修改规则选择器。
  [ok] CSS.setPropertyRulePropertyName 修改属性规则属性名称。
  [ok] CSS.setRuleSelector 修改规则选择器。
  [ok] CSS.setStyleSheetText  设置新的样式表文本。
  [ok] CSS.setStyleTexts  按指定顺序逐一应用指定的样式修改。
  [ok] CSS.startRuleUsageTracking  启用选择器录制。
  [ok] CSS.stopRuleUsageTracking  停止跟踪规则使用情况，并返回自上次调用（或自覆盖率检测开始）以来使用的规则列表 。
  [ok] CSS.takeCoverageDelta  获取自上次调用此方法（或自覆盖率检测开始）以来使用的规则列表。
  [ok] CSS.getEnvironmentVariables  返回 env() 函数中使用的默认 UA 定义的环境变量的值。
  [ok] CSS.setContainerQueryText  修改容器查询的表达式。

- chrome Debugger ： 调试器域公开了 JavaScript 调试功能。
  [ok] Debugger.continueToLocation  持续执行直至到达指定断点位置。
  [ok] Debugger.disable  禁用指定页面的调试器
  [ok] Debugger.enable  启用指定页面的调试器
  [ok] Debugger.evaluateOnCallFrame  对给定调用帧求表达式的值。
  [ok] Debugger.getPossibleBreakpoints 返回断点的可能位置。起始位置和结束位置的 scriptId 必须相同。
  [ok] Debugger.getScriptSource  返回指定脚本的源代码。
  [ok] Debugger.pause 执行到下一条 JavaScript 语句时停止。
  [ok] Debugger.restartFrame 从头开始重新启动特定的调用帧。
  [ok] Debugger.resume  恢复 JavaScript 运行。
  [ok] Debugger.searchInContent  在指定脚本中搜索字符串。
  [ok] Debugger.setAsyncCallStackDepth 启用或禁用异步调用堆栈跟踪。
  [ok] Debugger.setBreakpoint 在指定位置设置 JavaScript 断点。
  [ok] Debugger.setBreakpointByUrl 在指定位置（通过 URL 或 URL 正则表达式指定）设置 JavaScript 断点。
  [ok] Debugger.setBreakpointsActive  激活/停用页面上的所有断点。
  [ok] Debugger.setInstrumentationBreakpoint  设置检测断点。
  [ok] Debugger.setPauseOnExceptions  定义异常暂停状态。可以设置为在所有异常、未捕获的异常或已捕获的异常（无异常）时停止。异常暂停状态的初始值为none。
  [ok] Debugger.setScriptSource  实时编辑 JavaScript 源代码。 
  [ok] Debugger.setSkipAllPauses  使页面在任何暂停时（断点、异常、DOM 异常等）都不会中断。
  [ok] Debugger.setVariableValue  更改调用帧中变量的值。不支持基于对象的作用域，必须手动修改。
  [ok] Debugger.stepInto  进入函数调用。
  [ok] Debugger.stepOut  退出当前函数。
  [ok] Debugger.stepOver  跳过当前函数。
  [ok] Debugger.disassembleWasmModule  反汇编 Wasm 模块
  [ok] Debugger.getStackTrace  返回给定堆栈跟踪的stackTraceId.

- chrome Emulation : 该域名模拟了页面的不同环境。
  [ok] Emulation.clearDeviceMetricsOverride  清除已覆盖的设备指标。
  [ok] Emulation.clearGeolocationOverride  清除已覆盖的地理位置位置和错误。
  [ok] Emulation.clearIdleOverride  清除空闲状态覆盖。
  [ok] Emulation.setCPUThrottlingRate  启用 CPU 降频功能以模拟低速 CPU。
  [ok] Emulation.setDefaultBackgroundColorOverride  设置或清除框架默认背景颜色的覆盖值。如果内容未指定背景颜色，则使用此覆盖值。
  [ok] Emulation.setDeviceMetricsOverride  覆盖设备屏幕尺寸的值（window.screen.width、window.screen.height、window.innerWidth、window.innerHeight）。
  [ok] Emulation.setEmulatedMedia 模拟 CSS 媒体查询中给定的媒体类型或媒体特性。
  [ok] Emulation.setEmulatedOSTextScale  模拟给定操作系统的文本缩放比例。
  [ok] Emulation.setEmulatedVisionDeficiency   模拟给定的视力缺陷。
  [ok] Emulation.setGeolocationOverride  覆盖地理位置位置或误差。省略纬度、经度或精度将模拟位置不可用。
  [ok] Emulation.setIdleOverride 覆盖空闲状态。
  [ok] Emulation.setScriptExecutionDisabled  切换页面中的脚本执行方式。
  [ok] Emulation.setTimezoneOverride  使用指定的时区覆盖主机系统的默认时区。
  [ok] Emulation.setTouchEmulationEnabled  在不支持触控功能的平台上启用触控功能。
  [ok] Emulation.setUserAgentOverride  允许使用给定的字符串覆盖用户代理。 userAgentMetadata必须设置此项才能发送客户端提示标头。

- chrome EventBreakpoints : 事件断点域 允许在 JavaScript 调用的原生代码中发生的操作和事件上设置 JavaScript 断点。  // 0329 End
  [ok] EventBreakpoints.disable  移除所有断点
  [ok] EventBreakpoints.removeInstrumentationBreakpoint 移除特定原生事件上的断点。
  [ok] EventBreakpoints.setInstrumentationBreakpoint  在特定原生事件上设置断点。

- chrome CacheStorage ： 缓存存储域
  [ok] CacheStorage.deleteCache 清除缓存。
  [ok] CacheStorage.deleteEntry 删除缓存条目。
  [ok] CacheStorage.requestCachedResponse 获取缓存条目。
  [ok] CacheStorage.requestCacheNames 请求缓存名称。
  [ok] CacheStorage.requestEntries 从缓存中请求数据。

- chrome Extensions : 定义浏览器扩展的命令和事件。
  [ok] Extensions.clearStorageItems 清除给定扩展存储中的数据storageArea。
  [ok] Extensions.getExtensions  获取所有已解压扩展的列表。
  [ok] Extensions.getStorageItems 从指定位置的扩展存储中获取数据storageArea。
  [ok] Extensions.loadUnpacked 从文件系统安装已解压的扩展，类似于 `--load-extension` 命令行标志。扩展安装完成后返回扩展 ID。
                             仅当客户端使用 `--remote-debugging-pipe` 标志连接且设置了 `--enable-unsafe-extension-debugging` 标志时才可用。
  [ok] Extensions.removeStorageItems  keys从给定的扩展存储中移除storageArea
  [ok] Extensions.setStorageItems  values在给定的扩展存储中设置storageArea。提供的值values 将与存储区域中的现有值合并。
  [ok] Extensions.triggerAction  运行扩展程序的默认操作。仅当客户端使用 `--remote-debugging-pipe` 标志连接且设置了 `--enable-unsafe-extension-debugging` 标志时才可用。
  [ok] Extensions.uninstall  从配置文件中卸载已解压的扩展程序（不支持其他扩展程序）。仅当客户端使用 `--remote-debugging-pipe` 
                            标志和 `--enable-unsafe-extension-debugging` 标志连接时才可用。

- chrome FedCm : 该域允许与 FedCM 对话框进行交互。  
  [ok] FedCm.clickDialogButton 点击对话框按钮
  [ok] FedCm.disable  禁用FedCm
  [ok] FedCm.dismissDialog dismiss对话框
  [ok] FedCm.enable  启用FedCm
  [ok] FedCm.openUrl  访问一个URL并打开FedCm对话框
  [ok] FedCm.resetCooldown 重置冷却时间（如果有），以允许下一次 FedCM 调用显示对话框，即使用户最近关闭了某个对话框。
  [ok] FedCm.selectAccount 选择账户

- [] 改测试的bug和优化
    [] chrome CSS 测试和完善
    [] chrome Debugger 测试和完善
    [] chrome Emulation 测试和完善
    [] chrome EventBreakpoints 测试和完善
    [] chrome CacheStorage 测试和完善
    [] chrome Extensions 测试和完善
    [] chrome FedCm 测试和完善
- [] bug和优化验收
- [] 更新文档   // 0402 End

#### v0.0.8
- [ok] chrome device 指定多种设备启动浏览器
- [ok] 支持字符串模板语法，替代加法拼接字符串 
- [ok] chrome  SystemInfo ： SystemInfo 域定义了用于查询底层系统信息的方法和事件
    [ok] SystemInfo.getInfo 获取Chrome完整系统信息
    [ok] SystemInfo.getProcessInfo  获取Chrome所有进程信息
    [ok] SystemInfo.getFeatureState 获取Chrome所有特性状态

- [ok] chrome Browser ： 浏览器域定义了用于管理浏览器的方法和事件。
    [ok] Browser.close   关闭浏览器 
    [ok] Browser.resetPermissions  重置权限
    [ok] Browser.getWindowBounds 获取浏览器窗口的位置和大小。
    [ok] Browser.getWindowForTarget 获取目标对象对应的浏览器窗口。
    [ok] Browser.setContentsSize 设置浏览器窗口的大小。
    [ok] Browser.setWindowBounds 设置浏览器窗口的位置和/或大小。
    
- [ok] chrome Target : 目标对象
    [ok] Target.activateTarget 激活target 聚焦指定页面
    [ok] Target.attachToTarget  聚焦后返回sessionID
    [ok] Target.closeTarget   关闭指定target,如果目标是页面，则页面也会被关闭。
    [ok] Target.createBrowserContext  创建一个新的空浏览器上下文。它类似于（浏览器的）无痕模式，但你可以同时拥有多个。
                                    举个通俗的例子：普通的浏览器就像你只有一台电脑，所有人用同一个账号登录。而使用 BrowserContext，
                                    就像是你瞬间变出了 10 台全新的、互不干扰的电脑，每台电脑都可以独立登录不同的账号，互不影响。
    [ok] Target.createTarget 创建target (常用于创建页面)
    [ok] Target.detachFromTarget 分离掉指定sessionID
    [ok] Target.disposeBrowserContext 删除 BrowserContext。所有属于该 BrowserContext 的页面都将被关闭，而不会调用它们的 beforeunload 钩子函数
    [ok] Target.getBrowserContexts 返回创建的所有浏览器上下文
    [ok] Target.getTargets  获取可用目标列表。
    [待定] Target.setAutoAttach  <待定>  控制是否自动附加到与当前目标直接相关的新目标（例如 iframe 或 worker）。
    [待定] Target.setDiscoverTargets <待定> 控制是否发现可用目标并通过 targetCreated/targetInfoChanged/targetDestroyed事件通知
    [ok] Target.getTargetInfo  返回目标的相关信息。

- [ok] chrome DOMSnapshot : 该域便于获取包含 DOM、布局和样式信息的文档快照。
    [ok] DOMSnapshot.captureSnapshot <深入了解> 返回文档快照，其中包含根节点的完整 DOM 树（包括 iframe、模板内容和导入的文档），以扁平数组的形式呈现
    [ok] DOMSnapshot.disable 禁用给定页面的 DOM 快照。
    [ok] DOMSnapshot.enable  启用 DOM 快照
    
- [ok] chrome DOMStorage  : 查询和修改 DOM 存储。
    [ok] DOMStorage.clear <深入了解>
    [ok] DOMStorage.disable  禁用存储跟踪，阻止将存储事件发送到客户端。
    [ok] DOMStorage.enable  启用存储跟踪功能，存储事件现在将发送给客户端。 
    [ok] DOMStorage.getDOMStorageItems   <深入了解>
    [ok] DOMStorage.removeDOMStorageItem  <深入了解>
    [ok] DOMStorage.setDOMStorageItem   <深入了解>

- [ok] 改测试的bug和优化
    [ok] chrome Browser 测试
    [ok] chrome Target 测试
    [ok] chrome DOMSnapshot 测试
    [ok] chrome DOMStorage 测试
- [ok] bug和优化验收
- [ok] 更新文档 

#### v0.0.7
- [ok] host 方法扩展   系统文件相关交互方法   注意: path要求绝对路径，用相对路径会指定当前脚本工作目录为根
  [ok] 1. file/dir s=<search word> 搜索文件或目录
  [ok] 2. file/dir c=<path> 创建文件或目录
  [ok] 3. file/dir d=<path> 删除文件或目录
  [ok] 4. file/dir m=<path> goto=<path> 移动文件或目录
  [ok] 5. file/dir cp=<path> goto=<path> 复制文件或目录
  [ok] 6. file/dir r=<path> to=<arg> 读文件
  [ok] 7. file/dir renm=<path> goto=<path> 文件或目录改名, 路径不同则移动
  [ok] 8. ls=<path> 列出文件或目录
  [ok] 9. file/dir info=<path> 文件或目录信息
  [ok] 10. file/dir w=<path> from=<arg> 将文件内容写入文件
  [ok] 11. file/dir a=<path> from=<arg> 将文件内容追加写入文件   // 预计0325完成

- [ok] host 方法扩展 
  [ok] 1. ping 
  [ok] 2. port 已开放的端口 
  [ok] 3. zip unzip 压缩解压文件或目录   

- [ok] 改测试的bug和优化
- [ok] bug和优化验收
- [ok] 更新文档   


#### v0.0.6
- [ok] 支持excel相关操作方法 底层使用 https://github.com/qax-os/excelize 库
  [ok] 1. 读取excel
  [ok] 2. 写入excel
  [ok] 3. 最小粒度的单元格操作, 读取指定的单元格数据，写指定的单元格数据
  [ok] 4. 按行读写，按列读写
  [ok] 5. 图片插入
  [ok] 6. 单元格样式设置（字体、颜色、对齐）
  [ok] 7. 合并单元格
  [ok] 9. 设置单元格公式  
  
- [ok] 支持json与字典互转方法  
  [ok] 1. json转字典
  [ok] 2. 字典转json
  [ok] 3. json查找元素方法
  [ok] 4. 读写json文件 

- [ok] excel转json
- [ok] json转excel
- [ok] 将html table 直接转存到excel文件

- [ok]示例
  [ok] 1. eastmoney_2.cbs 采集东方财经网排名数据到excel文件
  [ok] 2. gaokaochsi_2.cbs 采集阳光高考的学校数据到excel文件
  [fail] 3. lottery_1.cbs 采集中国体彩的大乐透历史开奖数据到excel文件
  [ok] 4. douban_1.cbs 采集豆瓣电影排行榜到excel文件

- [ok] 改测试的bug和优化
- [ok] bug和优化验收
- [ok] 更新文档

#### v0.0.5
- [ok] 设计脚本配置 以 @ 开头, 全局的程序执行前就需要处理，优先解析并保存在全局常量，只有脚本模式下才有
  方案，在脚本解析前，优先判断@开头的行，如果符合以下全局配置则记录到全局配置常量，否则进行报错，脚本过完后将每行的第一个 @改为#，然后执行后面的语法检查
  [ok] 1. @cron 设置定时执行脚本,语法参考 cron 核心定时参数总览
  [ok] 2. @conf_json 设置外部配置文件json,读取json文件后将值存储到 as到指定的全局字典常量,以供脚本全局使用  @conf_json path="" as=conf
  [ok] 3. @conf_yaml 设置外部配置文件yaml,读取yaml文件后将值存储到 as到指定的全局字典常量,以供脚本全局使用  @conf_ini path="" as=conf
  [ok] 4. @conf_ini 设置外部配置文件ini,读取ini文件后将值存储到 as到指定的全局字典常量,以供脚本全局使用   @conf_ini path="" as=conf
  [ok] 5. @chrome_check 全局检查是否支持chrome浏览器，提取检查，如果宿主机未安装会提前检查出来 语法 : 直接使用  @chrome_check
  [ok] 6. @network_check 全局检查网络是否OK, 如果当前宿主机未网络会提前检查出来 语法: 请求地址可以是ip也可以是域名 如: @network_check "www.baidu.com"   @network_check "254.254.254.254" 

- [ok] 实现全局定时任务
- [ok] 实现全局配置 @conf_json, @conf_yaml, @conf_ini
- [ok] 实现全局 @chrome_check
- [ok] 实现全局 @network_check
- [ok] host 宿主机的相关方法 增加 host 关键字
  1. host name
  2. host ip
  3. host info
  4. host to
- [ok] DNS查询方法
- [ok] Whois查询方法
- [ok] 端口扫描方法
- [ok] 网站死链检查
- [ok] 网站证书信息检查
- [ok] 扫描网站的url, 扫描网站的外链
- [ok] 示例
  1. check_domain_ssl.cbs 定期监控网站证书到期时间
  2. check_domain_badlink.cbs 定期对网站进行死链检查
  3. search_port.cbs 定期扫描端口
  4. case_cron.cbs 定期chrome交互网站
- [ok] 改测试的bug和优化
- [ok] bug和优化验收
- [ok] 更新文档 (主要是host,和几个新增的方法)  

#### v0.0.4
- [ok] 系统级别的交互确认弹窗
- [ok] <bug转需求> 本地系统记录chrome userpath进程，如果存在默认隔离环境的chrome进程提供交互选择关闭还是新建隔离环境
- [ok] chrome init 新增new指令初始化隔离环境
- [ok] <bug转需求> chrome连不上的时候重试
- [ok] Xpath格式检查相关的函数
- [ok] 语法: 支持 \ 作为代码的换行连接符，与命令行的换行一样
- [ok] 语法: 命令式值能绑定上函数并执行函数使用函数的返回值 如: chrome click=NowTabMatchDemoContentOP("百度一下")
- [ok] chrome 非法参数应该提示
- [ok] 支持输出chrome信息
- [ok] 如果本地未安装chrome提示安装
- [ok] <bug转需求> tab失焦后进行提示, 提供交互确认，如果是继续则默认选择第一个tab,如果是结束就结束脚本
- [ok] 测试并新增示例
  1. [ok] doubao_1.cbs 交互豆包问豆包问题
  2. [ok] zzttop_1.cbs 交互网站手机号随机生成，采集将生成的手机号
  3. [ok] toutiao_1.cbs 头条采集今日要闻访问后截图保存到本地
  4. [ok] douyin_1.cbs 抖音交互下滑视频
  
- [ok] 改测试的bug和优化
- [ok] bug和优化验收 
- [ok] 更新文档 

#### v0.0.3
- [ok] 谷歌浏览器打开方法以及语句 增加 chrome 关键字
  1. 打开窗口以及页签
  2. 访问网址
  3. 新开一个页签
  4. 查看当前页签
  5. 关闭页签
  6. 切换页签
 
- [ok] 输出html以及语句
  1. 输出页面的html
  2. 解析HTML demo树
  
- [ok] 操作点击方法以及语句
  1. 点击xpath
  2. 输入xpath

- [ok] 检查操作，检查页面是否存在指定xpath
- [ok] 滚动操作，滚动页面
- [ok] 截图操作，浏览器截图操作

- [ok] 新增xpath相关的函数 
  1. 打印当前页面的demo数(含xpath)
  2. 匹配页面标签属性值获取对应xpath  
  3. 匹配页面内容获取对应的xpath

- [ok] 等待页面完全响应完
- [ok] 整理日志打印
- [ok] 多增加实例的测试用例以及测试 
  1. [ok] eastmoney_1.cbs  东方财经网，先进官网再访问排名，获取排名列表并保存到本地文件
  2. [ok] xinhuanet_1.cbs  新华网获取要闻聚焦，点击每个要闻访问文章内容并截图保存在本地
  3. [ok] newscctv_1.cbs  央视新闻网，获取每个分类标签下的新闻列表打印新闻的标题出来
  4. [ok] gaokaochsi_1.cbs 阳光高考获取全部的学校信息
  
- [ok] 改测试的bug和优化
- [ok] bug和优化验收
- [ok] 更新文档

#### v0.0.2
- [ok] http请求方法 增加 http 关键字
  1. http参数解析
  2. http请求的功能
  3. 请求后返回值的存储
  4. 关闭gt的日志
- [ok] 字符处理方法
- [ok] 正则方法
- [ok] 时间相关方法
- [ok] http增加save方法，将请求返回的body保存到本地
- [ok] 测试内置函数，语法层面的
- [ok] 除了基础测试的代码，多增加实例的测试用例
- [ok] 改测试的bug和优化
- [ok] 更新文档

#### v0.0.1
- [ok] 项目目录结构
- [ok] 开始函数与REPL
- [ok] DSL v0.1, 基础语法
- [ok] DSL v0.1, 测试语法，更多的示例
- [ok] DSL v0.1, 语法文档
- [ok] 改测试的bug和优化
- [ok] bug和优化验收
- [ok] 更新文档  支持了 elif, switch, for in, while in, 自增自减

## 里程碑
1. 第一阶段实现 DSL 实现脚本编写操作浏览器自动化任务
2. 第二阶段大量产出测试用例给AI进行训练
3. 训练定向模型，最终做到 自热语言 -> DSL脚本 -> 自动化操作浏览器

## 需求池
- 脚本命令式应该支持 output参数来指定数据输出到哪里（支持文件，数据库，远端服务等等）
- [] host 方法扩展  
  [] disk 磁盘信息 使用率
  [] mem 内存信息 使用率
  [] cpu cpu信息 使用率
  [] pid 进程信息
  [] 获取窗口信息
  [] 操作鼠标键盘
  [] 系统层面操作软件等等...
- 实现全局 @mysql  全局声明并连接mysql as到指定对象 todo 只设计语法
- 实现全局 @redis  全局声明并连接redis as到指定对象 todo 只设计语法
- http代理
- websocket客户端
- 研究和设计系统级别的变量和常量，多脚本可用，脚本运行中断但数据变量不还在，采用磁盘持久存储不会被内存释放 | 或者研究一下windows的注册表来存储变量
- 脚本运行中ctrl+c这些来中断需要弹出框再次确认
- 支持三目运算   参考 vba 的IIF 如: b = IIf(a > 10, 1, 2)
- 是否要支持goto语法
- 多脚本场景下 是否设计 CALL语法来支持调用运行其他脚本，用CALL语句可将程序执行控制权转移到一个脚本过程中，在其他脚本结束后，再将控制权返回到调用脚本的下一行。
  （设计的时候需考虑同步调用，异步调用，调用时候值的传递，调用之间的通讯机制等等）
- excel 加解密，隐藏，锁定功能 
- excel copy sheets在新的sheets上操作，保证原有的shhets不变
- excel 获取指定区域的数据，按照矩形 x0,y1 的坐标矩形， 如 A1:C20
- excel 多个sheets之间的交互
- 识别html页面组件或控件的方法，列举可操作的节点
- host box 支持输入框
- 支持邮箱发送来实现通知的效果
- Browser.close 关闭浏览器 在现有的关闭浏览器加这个方法，如果连接还在的时候就可以先用这个关闭浏览器
- 页面DOM变化的监听函数，监听页面DOM变化，并执行回调函数（参考现有的DOMSnapshot.captureSnapshot）
- 

----
 
浏览器自动化相关资料

## 浏览器自动化操作
点击操作	左键单击 / 双击、右键单击、坐标点击、模拟点击（避免被反爬检测）	触发按钮 / 链接 / 复选框等交互（如登录按钮、勾选同意协议）
输入操作	输入文本、清空输入框、模拟键盘输入（回车 / 退格 / 快捷键）	表单填写（用户名 / 密码 / 搜索框）、触发搜索（输入后按回车）
选择操作	下拉框选择（按索引 / 值 / 文本）、勾选 / 取消勾选复选框 / 单选框	表单中的下拉选择、同意条款勾选
拖拽操作	元素拖拽（如滑块验证）、文件拖拽上传	滑块验证码、文件上传、拖拽排序
元素状态验证	检查元素是否可见 / 可点击 / 存在 / 被选中、获取元素属性 / 文本 / 尺寸 / 位置	验证操作是否生效（如输入框是否有值、按钮是否可点击）
元素滚动	滚动到元素可见位置、页面滚动（向上 / 向下 / 到顶部 / 到底部）	操作不在可视区域的元素（如页面底部的提交按钮）
元素截图	单独截取元素区域、高亮元素后截图	验证元素显示是否正确、留存测试证据


#### 元素定位

定位方式	语法示例（Selenium/Playwright）	适用场景
ID 定位	find_element(By.ID, "username") / page.locator("#username")	元素有唯一 ID（最稳定）
名称（Name）定位	find_element(By.NAME, "password") / page.locator("[name='password']")	表单元素（输入框 / 按钮）常用 name 属性
类名（Class）定位	find_element(By.CLASS_NAME, "btn-submit") / page.locator(".btn-submit")	按样式类定位，注意类名可能重复
标签名定位	find_element(By.TAG_NAME, "input") / page.locator("input")	批量定位同类型元素（如所有输入框）
XPath 定位	find_element(By.XPATH, "//input[@id='username']") / page.locator("//input[@id='username']")	万能定位，支持复杂路径 / 逻辑（如父节点 / 兄弟节点 / 文本匹配）
CSS 选择器定位	find_element(By.CSS_SELECTOR, "#username") / page.locator("#username")	高效定位，支持层级 / 属性 / 伪类（如 input[type='text']:first-child）
链接文本定位	find_element(By.LINK_TEXT, "登录") / page.locator("text=登录")	定位超链接（<a>标签）
部分链接文本	find_element(By.PARTIAL_LINK_TEXT, "登")	链接文本过长时模糊匹配
自定义定位	基于文本 / 正则 / 相对位置（Playwright 的locator.filter()）	复杂场景（如 “包含‘提交’且 class 为 btn 的按钮”）

#### 页面导航与加载

操作类型	具体操作	用途
页面跳转	访问指定 URL（get(url)/goto(url)）、后退 / 前进 / 刷新页面	核心导航，如打开目标网站、返回上一页
加载控制	等待页面加载完成（等待 DOM / 网络 / 渲染完成）、中断页面加载	避免操作过早执行（如页面未加载完就点击按钮）
页面信息获取	获取页面 URL / 标题 / 源代码 / 截图 / HTML、获取页面性能数据（加载时间）	验证页面是否正确跳转、保存页面快照、分析性能
页面状态控制	最大化 / 最小化 / 调整窗口大小、全屏模式、设置页面缩放比例	适配不同分辨率、模拟移动端窗口（如 375x667）


#### 会话 / 驱动管理

操作类型	具体操作	用途
驱动初始化	启动浏览器（Chrome/Firefox/Edge）、配置启动参数（无头模式 / 窗口大小 / 代理）	初始化自动化环境，如 puppeteer.launch({headless: true})、webdriver.Chrome()
会话控制	新建标签页 / 窗口、关闭标签页 / 窗口、切换标签页 / 窗口、获取当前窗口句柄	多标签 / 多窗口场景（如同时操作多个页面）
浏览器配置	设置 Cookie / 本地存储 / 会话存储、清除缓存 / 历史记录、设置用户代理（UA）	模拟不同用户环境、绕过反爬、持久化登录状态
退出驱动	关闭浏览器、释放驱动资源	自动化结束后清理环境，避免进程残留


#### 键盘与鼠标模拟

操作类型	具体操作	用途
键盘模拟	按下 / 释放单个按键（KeyDown/KeyUp）、组合键（Ctrl+C/Ctrl+V）、输入特殊字符	快捷键操作（复制 / 粘贴）、模拟用户打字节奏、触发键盘事件
鼠标模拟	鼠标移动（到元素 / 坐标）、鼠标按下 / 释放、滚轮滚动（向上 / 向下）	模拟用户鼠标移动、悬浮提示框、滚轮翻页
悬浮操作	鼠标悬浮到元素上（触发 hover 事件）	显示悬浮菜单、验证悬浮提示文本


#### 弹窗处理

弹窗类型	具体操作	用途
警告弹窗（Alert）	接受弹窗（OK）、获取弹窗文本	处理系统警告（如 “操作成功” 提示）
确认弹窗（Confirm）	接受 / 取消弹窗、获取弹窗文本	处理确认操作（如 “是否删除”）
提示弹窗（Prompt）	输入文本后确认、取消弹窗、获取弹窗提示文本	处理需要输入的弹窗（如 “请输入验证码”）
浏览器弹窗	处理下载弹窗、文件上传弹窗、权限请求弹窗（摄像头 / 麦克风 / 通知）	下载文件、上传文件、允许 / 拒绝权限请求
自定义弹窗	定位并关闭自定义弹窗（如广告弹窗、登录弹窗）	关闭干扰操作的弹窗（如页面广告）


#### 网络请求与响应

操作类型	具体操作	用途
网络监听	监听所有网络请求 / 响应、过滤指定 URL / 方法（GET/POST）	抓包分析、获取接口数据、验证接口返回
网络拦截	拦截指定请求、修改请求参数 / 响应内容、模拟接口返回（Mock）	测试异常场景（如接口返回 500）、绕过接口限制、加速测试（Mock 数据）
网络等待	等待指定请求完成、等待所有请求完成	确保接口数据加载完成后再操作（如列表数据加载）
Cookie 操作	添加 / 获取 / 删除 Cookie、设置 Cookie 有效期 / 域名 / 路径	持久化登录状态、模拟不同用户登录
本地存储操作	操作 LocalStorage/SessionStorage（添加 / 获取 / 删除 / 清空）	模拟前端存储数据、验证存储逻辑


#### JavaScript 执行

操作类型	具体操作	用途
执行 JS 代码	执行任意 JavaScript 代码、获取 JS 执行结果	操作前端无法通过常规 API 实现的功能（如修改元素样式、触发自定义事件）
注入 JS 脚本	向页面注入外部脚本（如 jQuery）、修改页面全局变量	增强页面交互能力、绕过前端限制
异步 JS 执行	执行异步 JS 代码（如等待 Promise 完成）、获取异步执行结果	处理前端异步逻辑（如接口请求完成后获取数据）


#### 文件操作

操作类型	具体操作	用途
文件上传	定位文件上传输入框、传入本地文件路径	表单中的文件上传（如头像上传、文档上传）
文件下载	设置下载路径、等待文件下载完成、验证下载文件（存在 / 大小 / 内容）	下载报表、文件后验证完整性
截图 / 录屏	页面全屏截图、指定区域截图、录制页面操作视频	测试结果留存、问题复现、自动化报告附件


#### 高级控制

操作类型	具体操作	用途
框架 /iframe 切换	切换到指定 iframe、切回主文档、嵌套 iframe 切换	操作 iframe 内的元素（如支付页面、内嵌表单）
多窗口 / 标签页	新建窗口 / 标签页、切换窗口 / 标签页、关闭窗口 / 标签页、获取所有窗口句柄	多页面并行操作（如同时登录多个账号）
超时设置	设置全局超时、元素定位超时、页面加载超时、脚本执行超时	避免自动化卡死（如元素长时间未出现时超时退出）
异常处理	捕获操作异常、重试失败操作、忽略指定异常	提高自动化稳定性（如网络波动导致的点击失败重试）
代理设置	配置 HTTP/HTTPS/SOCKS 代理、切换代理、验证代理有效性	模拟不同地区访问、绕过 IP 限制
认证处理	处理 HTTP 基本认证、OAuth 认证、验证码识别（对接打码平台）	登录需要认证的网站、处理验证码（如滑块 / 图片验证码）


## 其他文案

- ChromeBot 的作用： 专注与自动化操作浏览器脚本的编写，未减少心智脚本语言cbs语法设计简单为首要目标，然后满足所有的脚本语法，支持命令式语法，只要是有编程语言基础的人能快速上手使用，降低学习成本。
- ChromeBot, 最有帮助的是什么? 最有价值的是什么？最有趣的是什么？最值得推广的是什么？
- ChromeBot核心用户群体: 测试，爬虫，自动化任务，爱好学习者，炫技者，日常电脑办公
- ChromeBot 核心是编写自动化脚本来完成自动化任务，首要支持是Chrome浏览器的自动化，其次是支持http相关方法(如:curl的功能)，系统的脚本功能，读写多存储应用(excel,mysql,redis....)

