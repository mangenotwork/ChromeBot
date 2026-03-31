package browser

import (
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// CDPDOMSnapshotCaptureSnapshot 捕获DOM结构快照
// 可选参数:
//   - computedStyles: 要包含的计算样式列表
//   - includeDOMRects: 是否包含DOM矩形信息
//   - includeBlendedBackgroundColors: 是否包含混合背景色
//
// 返回值:
//   - 包含DOM快照数据的JSON字符串
//   - error: 捕获过程中发生的错误
//
// DOMSnapshot.captureSnapshot用于捕获页面DOM结构的完整快照：
// 捕获当前DOM树的结构和内容, 包含节点的属性、样式、布局信息, 生成可用于离线分析的数据结构, 支持序列化和反序列化, 可以包含计算样式和布局信息
func CDPDOMSnapshotCaptureSnapshot(options ...DOMSnapshotOption) (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 DOMSnapshot.captureSnapshot")
	}

	// 默认配置
	config := &DOMSnapshotConfig{
		ComputedStyles:                 []string{},
		IncludeDOMRects:                false,
		IncludeBlendedBackgroundColors: false,
	}

	// 应用选项
	for _, option := range options {
		option(config)
	}

	chromeInstance.NextID++
	reqID := chromeInstance.NextID

	// 构建消息
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "DOMSnapshot.captureSnapshot",
		"params": {
			"computedStyles": ["%s"],
			"includeDOMRects": %t,
			"includeBlendedBackgroundColors": %t
		}
	}`, reqID, strings.Join(config.ComputedStyles, `","`),
		config.IncludeDOMRects, config.IncludeBlendedBackgroundColors)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 DOMSnapshot.captureSnapshot 失败:", err)
		return "", err
	}

	utils.Debugf("发送 CDP 消息: %s", message)
	timeout := 10 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				fmt.Println("[CDP DOMSnapshot.captureSnapshot] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("captureSnapshot 请求超时")
		}
	}
}

// DOMSnapshotConfig DOM快照配置
type DOMSnapshotConfig struct {
	ComputedStyles                 []string
	IncludeDOMRects                bool
	IncludeBlendedBackgroundColors bool
}

// DOMSnapshotOption 配置选项
type DOMSnapshotOption func(*DOMSnapshotConfig)

// WithComputedStyles 设置要包含的计算样式
func WithComputedStyles(styles []string) DOMSnapshotOption {
	return func(c *DOMSnapshotConfig) {
		c.ComputedStyles = styles
	}
}

// IncludeDOMRects 包含DOM矩形信息
func IncludeDOMRects(include bool) DOMSnapshotOption {
	return func(c *DOMSnapshotConfig) {
		c.IncludeDOMRects = include
	}
}

// IncludeBlendedBackgroundColors 包含混合背景色
func IncludeBlendedBackgroundColors(include bool) DOMSnapshotOption {
	return func(c *DOMSnapshotConfig) {
		c.IncludeBlendedBackgroundColors = include
	}
}

/*
// 示例2: 包含计算样式的详细快照
func exampleDetailedDOMSnapshot() {
	// 定义要捕获的计算样式
	computedStyles := []string{
		"display", "position", "width", "height",
		"color", "background-color", "font-size",
		"margin", "padding", "border",
		"z-index", "opacity", "visibility",
	}

	log.Println("正在捕获详细DOM快照（包含样式）...")

	response, err := CDPDOMSnapshotCaptureSnapshot(
		WithComputedStyles(computedStyles),
		IncludeDOMRects(true),
		IncludeBlendedBackgroundColors(true),
	)
	if err != nil {
		log.Fatalf("捕获详细DOM快照失败: %v", err)
	}

	log.Printf("详细DOM快照捕获成功，包含 %d 种计算样式", len(computedStyles))

	// 分析样式数据
	var data struct {
		Result struct {
			ComputedStyles []struct {
				Name string `json:"name"`
			} `json:"computedStyles"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		log.Printf("解析失败: %v", err)
		return
	}

	log.Printf("实际捕获的计算样式数量: %d", len(data.Result.ComputedStyles))
	for i, style := range data.Result.ComputedStyles {
		if i < 5 { // 只显示前5个
			log.Printf("  样式[%d]: %s", i+1, style.Name)
		}
	}
	if len(data.Result.ComputedStyles) > 5 {
		log.Printf("  ... 还有 %d 个样式", len(data.Result.ComputedStyles)-5)
	}
}


// 示例3: 页面对比分析
func examplePageComparison() {
	// 捕获页面A的DOM快照
	log.Println("捕获页面A的DOM快照...")
	snapshotA, err := CDPDOMSnapshotCaptureSnapshot()
	if err != nil {
		log.Printf("捕获页面A失败: %v", err)
		return
	}

	// 执行一些操作（例如：点击按钮、表单填写等）
	log.Println("执行页面操作...")
	// performPageActions()

	// 等待页面更新
	time.Sleep(2 * time.Second)

	// 捕获页面B的DOM快照
	log.Println("捕获页面B的DOM快照...")
	snapshotB, err := CDPDOMSnapshotCaptureSnapshot()
	if err != nil {
		log.Printf("捕获页面B失败: %v", err)
		return
	}

	// 比较快照
	comparison, err := CompareDOMSnapshots(snapshotA, snapshotB)
	if err != nil {
		log.Fatalf("比较快照失败: %v", err)
	}

	// 生成报告
	report := GenerateComparisonReport(comparison)
	fmt.Println(report)

	// 输出JSON格式结果
	jsonResult, _ := json.MarshalIndent(comparison, "", "  ")
	fmt.Println("\n=== JSON格式结果 ===")
	fmt.Println(string(jsonResult))
}



// 示例4: DOM结构分析工具
func exampleDOMStructureAnalysis() {
	log.Println("开始DOM结构分析...")

	response, err := CDPDOMSnapshotCaptureSnapshot(
		WithComputedStyles([]string{"display"}),
		IncludeDOMRects(true),
	)
	if err != nil {
		log.Fatalf("捕获DOM快照失败: %v", err)
	}

	// 解析DOM快照
	var data struct {
		Result struct {
			Documents []struct {
				Nodes []struct {
					NodeType   int    `json:"nodeType"`
					NodeName   string `json:"nodeName"`
					BackendID  int    `json:"backendNodeId"`
					ChildCount int    `json:"childCount,omitempty"`
				} `json:"nodes"`
			} `json:"documents"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		log.Printf("解析失败: %v", err)
		return
	}

	if len(data.Result.Documents) == 0 {
		log.Println("未找到文档")
		return
	}

	doc := data.Result.Documents[0]
	nodes := doc.Nodes

	log.Printf("DOM结构分析结果:")
	log.Printf("  总节点数: %d", len(nodes))

	// 节点类型统计
	typeCounts := make(map[string]int)
	for _, node := range nodes {
		nodeType := getNodeTypeName(node.NodeType)
		typeCounts[nodeType]++
	}

	log.Printf("  节点类型分布:")
	for nodeType, count := range typeCounts {
		percentage := float64(count) / float64(len(nodes)) * 100
		log.Printf("    %s: %d (%.1f%%)", nodeType, count, percentage)
	}

	// 查找特定类型的节点
	var divNodes, spanNodes, imgNodes []int
	for _, node := range nodes {
		switch node.NodeName {
		case "DIV":
			divNodes = append(divNodes, node.BackendID)
		case "SPAN":
			spanNodes = append(spanNodes, node.BackendID)
		case "IMG":
			imgNodes = append(imgNodes, node.BackendID)
		}
	}

	log.Printf("  特定元素数量:")
	log.Printf("    DIV元素: %d", len(divNodes))
	log.Printf("    SPAN元素: %d", len(spanNodes))
	log.Printf("    IMG元素: %d", len(imgNodes))
}


*/

// 辅助函数
func getNodeTypeName(nodeType int) string {
	switch nodeType {
	case 1:
		return "ELEMENT_NODE"
	case 3:
		return "TEXT_NODE"
	case 8:
		return "COMMENT_NODE"
	case 9:
		return "DOCUMENT_NODE"
	case 10:
		return "DOCUMENT_TYPE_NODE"
	case 11:
		return "DOCUMENT_FRAGMENT_NODE"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", nodeType)
	}
}

// ==================================================== 比较快照的实现 ====================================================
// DOMSnapshotComparison DOM快照比较结果
type DOMSnapshotComparison struct {
	NodeChanges       int           `json:"nodeChanges"`
	AttributeChanges  int           `json:"attributeChanges"`
	TextChanges       int           `json:"textChanges"`
	AddedNodes        []int         `json:"addedNodes"`
	RemovedNodes      []int         `json:"removedNodes"`
	ModifiedNodes     []NodeChange  `json:"modifiedNodes"`
	ChangedAttributes []AttrChange  `json:"changedAttributes"`
	ChangedTexts      []TextChange  `json:"changedTexts"`
	TotalNodesA       int           `json:"totalNodesA"`
	TotalNodesB       int           `json:"totalNodesB"`
	SimilarityScore   float64       `json:"similarityScore"`
	ComparisonTime    time.Duration `json:"comparisonTime"`
}

// NodeChange 节点变化信息
type NodeChange struct {
	NodeID     int      `json:"nodeId"`
	NodeName   string   `json:"nodeName"`
	ChangeType string   `json:"changeType"` // "added", "removed", "modified"
	Changes    []string `json:"changes"`    // 具体变化描述
}

// AttrChange 属性变化信息
type AttrChange struct {
	NodeID    int               `json:"nodeId"`
	NodeName  string            `json:"nodeName"`
	Added     map[string]string `json:"added"`     // 新增的属性
	Removed   []string          `json:"removed"`   // 删除的属性
	Modified  map[string]string `json:"modified"`  // 修改的属性（key: 属性名, value: 新值）
	OldValues map[string]string `json:"oldValues"` // 修改前的值
}

// TextChange 文本变化信息
type TextChange struct {
	NodeID     int    `json:"nodeId"`
	NodeName   string `json:"nodeName"`
	OldText    string `json:"oldText"`
	NewText    string `json:"newText"`
	TextLength int    `json:"textLength"`
	Diff       string `json:"diff"` // 文本差异摘要
}

// DOMNodeInfo DOM节点信息
type DOMNodeInfo struct {
	BackendID     int                    `json:"backendNodeId"`
	NodeName      string                 `json:"nodeName"`
	NodeType      int                    `json:"nodeType"`
	NodeValue     string                 `json:"nodeValue,omitempty"`
	ChildCount    int                    `json:"childCount,omitempty"`
	Attributes    map[string]string      `json:"attributes,omitempty"`
	ComputedStyle map[string]interface{} `json:"computedStyle,omitempty"`
	Layout        map[string]interface{} `json:"layout,omitempty"`
	ParentID      int                    `json:"parentId,omitempty"`
	Children      []int                  `json:"children,omitempty"`
	Index         int                    `json:"index"` // 在父节点中的位置
	Depth         int                    `json:"depth"` // 节点深度
}

// DOMSnapshotData DOM快照数据
type DOMSnapshotData struct {
	Nodes     map[int]*DOMNodeInfo `json:"nodes"`
	RootNodes []int                `json:"rootNodes"`
	NodeCount int                  `json:"nodeCount"`
	DocURL    string               `json:"documentURL"`
	Title     string               `json:"title"`
	Timestamp time.Time            `json:"timestamp"`
}

// CompareDOMSnapshots 比较两个DOM快照
func CompareDOMSnapshots(snapshotA, snapshotB string) (*DOMSnapshotComparison, error) {
	startTime := time.Now()

	// 1. 解析快照数据
	dataA, err := parseSnapshotData(snapshotA)
	if err != nil {
		return nil, fmt.Errorf("解析快照A失败: %w", err)
	}

	dataB, err := parseSnapshotData(snapshotB)
	if err != nil {
		return nil, fmt.Errorf("解析快照B失败: %w", err)
	}

	comparison := &DOMSnapshotComparison{
		AddedNodes:        []int{},
		RemovedNodes:      []int{},
		ModifiedNodes:     []NodeChange{},
		ChangedAttributes: []AttrChange{},
		ChangedTexts:      []TextChange{},
		TotalNodesA:       dataA.NodeCount,
		TotalNodesB:       dataB.NodeCount,
	}

	// 2. 比较节点结构
	compareNodeStructure(dataA, dataB, comparison)

	// 3. 比较节点属性
	compareNodeAttributes(dataA, dataB, comparison)

	// 4. 比较文本内容
	compareNodeTexts(dataA, dataB, comparison)

	// 5. 计算相似度分数
	comparison.SimilarityScore = calculateSimilarityScore(comparison, dataA.NodeCount, dataB.NodeCount)

	// 6. 设置统计信息
	comparison.NodeChanges = len(comparison.AddedNodes) + len(comparison.RemovedNodes) + len(comparison.ModifiedNodes)
	comparison.AttributeChanges = countAttributeChanges(comparison.ChangedAttributes)
	comparison.TextChanges = len(comparison.ChangedTexts)
	comparison.ComparisonTime = time.Since(startTime)

	return comparison, nil
}

// parseSnapshotData 解析DOM快照数据
func parseSnapshotData(snapshot string) (*DOMSnapshotData, error) {
	var response struct {
		Result struct {
			Documents []struct {
				DocumentURL string `json:"documentURL"`
				Title       string `json:"title"`
				Nodes       []struct {
					BackendID  int      `json:"backendNodeId"`
					NodeName   string   `json:"nodeName"`
					NodeType   int      `json:"nodeType"`
					NodeValue  string   `json:"nodeValue,omitempty"`
					ChildCount int      `json:"childCount,omitempty"`
					Attributes []string `json:"attributes,omitempty"`
					ParentID   int      `json:"parentIndex,omitempty"`
				} `json:"nodes"`
				Layout struct {
					NodeIndex []int `json:"nodeIndex"`
				} `json:"layout,omitempty"`
			} `json:"documents"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(snapshot), &response); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	if len(response.Result.Documents) == 0 {
		return nil, fmt.Errorf("快照中没有文档")
	}

	doc := response.Result.Documents[0]
	data := &DOMSnapshotData{
		Nodes:     make(map[int]*DOMNodeInfo),
		RootNodes: []int{},
		NodeCount: len(doc.Nodes),
		DocURL:    doc.DocumentURL,
		Title:     doc.Title,
		Timestamp: time.Now(),
	}

	// 构建节点映射
	for i, rawNode := range doc.Nodes {
		node := &DOMNodeInfo{
			BackendID:  rawNode.BackendID,
			NodeName:   rawNode.NodeName,
			NodeType:   rawNode.NodeType,
			NodeValue:  rawNode.NodeValue,
			ChildCount: rawNode.ChildCount,
			Index:      i,
		}

		// 解析属性
		if len(rawNode.Attributes) > 0 {
			node.Attributes = make(map[string]string)
			for j := 0; j < len(rawNode.Attributes); j += 2 {
				if j+1 < len(rawNode.Attributes) {
					key := rawNode.Attributes[j]
					value := rawNode.Attributes[j+1]
					node.Attributes[key] = value
				}
			}
		}

		// 设置父节点
		if rawNode.ParentID > 0 && rawNode.ParentID < len(doc.Nodes) {
			node.ParentID = doc.Nodes[rawNode.ParentID].BackendID
		}

		data.Nodes[node.BackendID] = node
	}

	// 构建子节点关系
	for _, node := range data.Nodes {
		if node.ParentID > 0 {
			if parent, exists := data.Nodes[node.ParentID]; exists {
				parent.Children = append(parent.Children, node.BackendID)
			}
		} else {
			// 根节点
			data.RootNodes = append(data.RootNodes, node.BackendID)
		}
	}

	// 计算节点深度
	calculateNodeDepth(data)

	return data, nil
}

// calculateNodeDepth 计算节点深度
func calculateNodeDepth(data *DOMSnapshotData) {
	var calculate func(nodeID, depth int)
	calculate = func(nodeID, depth int) {
		node, exists := data.Nodes[nodeID]
		if !exists {
			return
		}

		node.Depth = depth
		for _, childID := range node.Children {
			calculate(childID, depth+1)
		}
	}

	for _, rootID := range data.RootNodes {
		calculate(rootID, 0)
	}
}

// compareNodeStructure 比较节点结构
func compareNodeStructure(dataA, dataB *DOMSnapshotData, comparison *DOMSnapshotComparison) {
	// 查找新增的节点
	for nodeID, nodeB := range dataB.Nodes {
		if _, exists := dataA.Nodes[nodeID]; !exists {
			comparison.AddedNodes = append(comparison.AddedNodes, nodeID)

			comparison.ModifiedNodes = append(comparison.ModifiedNodes, NodeChange{
				NodeID:     nodeID,
				NodeName:   nodeB.NodeName,
				ChangeType: "added",
				Changes:    []string{fmt.Sprintf("新增%s节点", nodeB.NodeName)},
			})
		}
	}

	// 查找删除的节点
	for nodeID, nodeA := range dataA.Nodes {
		if _, exists := dataB.Nodes[nodeID]; !exists {
			comparison.RemovedNodes = append(comparison.RemovedNodes, nodeID)

			comparison.ModifiedNodes = append(comparison.ModifiedNodes, NodeChange{
				NodeID:     nodeID,
				NodeName:   nodeA.NodeName,
				ChangeType: "removed",
				Changes:    []string{fmt.Sprintf("删除%s节点", nodeA.NodeName)},
			})
		}
	}
}

// compareNodeAttributes 比较节点属性
func compareNodeAttributes(dataA, dataB *DOMSnapshotData, comparison *DOMSnapshotComparison) {
	// 只比较两个快照中都存在的节点
	for nodeID, nodeA := range dataA.Nodes {
		nodeB, exists := dataB.Nodes[nodeID]
		if !exists {
			continue
		}

		attrChange := AttrChange{
			NodeID:    nodeID,
			NodeName:  nodeA.NodeName,
			Added:     make(map[string]string),
			Removed:   []string{},
			Modified:  make(map[string]string),
			OldValues: make(map[string]string),
		}

		hasChanges := false

		// 检查新增的属性
		for attrName, attrValueB := range nodeB.Attributes {
			attrValueA, existsA := nodeA.Attributes[attrName]
			if !existsA {
				attrChange.Added[attrName] = attrValueB
				hasChanges = true
			} else if attrValueA != attrValueB {
				attrChange.Modified[attrName] = attrValueB
				attrChange.OldValues[attrName] = attrValueA
				hasChanges = true
			}
		}

		// 检查删除的属性
		for attrName := range nodeA.Attributes {
			if _, existsB := nodeB.Attributes[attrName]; !existsB {
				attrChange.Removed = append(attrChange.Removed, attrName)
				hasChanges = true
			}
		}

		if hasChanges {
			comparison.ChangedAttributes = append(comparison.ChangedAttributes, attrChange)

			// 添加到修改节点列表
			changes := []string{}
			if len(attrChange.Added) > 0 {
				changes = append(changes, fmt.Sprintf("新增%d个属性", len(attrChange.Added)))
			}
			if len(attrChange.Removed) > 0 {
				changes = append(changes, fmt.Sprintf("删除%d个属性", len(attrChange.Removed)))
			}
			if len(attrChange.Modified) > 0 {
				changes = append(changes, fmt.Sprintf("修改%d个属性", len(attrChange.Modified)))
			}

			comparison.ModifiedNodes = append(comparison.ModifiedNodes, NodeChange{
				NodeID:     nodeID,
				NodeName:   nodeA.NodeName,
				ChangeType: "modified",
				Changes:    changes,
			})
		}
	}
}

// compareNodeTexts 比较节点文本
func compareNodeTexts(dataA, dataB *DOMSnapshotData, comparison *DOMSnapshotComparison) {
	for nodeID, nodeA := range dataA.Nodes {
		nodeB, exists := dataB.Nodes[nodeID]
		if !exists {
			continue
		}

		// 只比较文本节点
		if nodeA.NodeType == 3 && nodeB.NodeType == 3 { // TEXT_NODE
			if nodeA.NodeValue != nodeB.NodeValue {
				textChange := TextChange{
					NodeID:     nodeID,
					NodeName:   "TEXT_NODE",
					OldText:    nodeA.NodeValue,
					NewText:    nodeB.NodeValue,
					TextLength: len(nodeB.NodeValue),
					Diff:       calculateTextDiff(nodeA.NodeValue, nodeB.NodeValue),
				}

				comparison.ChangedTexts = append(comparison.ChangedTexts, textChange)

				// 添加到修改节点列表
				comparison.ModifiedNodes = append(comparison.ModifiedNodes, NodeChange{
					NodeID:     nodeID,
					NodeName:   "TEXT_NODE",
					ChangeType: "modified",
					Changes:    []string{fmt.Sprintf("文本变化: 长度%d→%d", len(nodeA.NodeValue), len(nodeB.NodeValue))},
				})
			}
		}
	}
}

// calculateTextDiff 计算文本差异
func calculateTextDiff(textA, textB string) string {
	const maxDiffLength = 50
	if len(textA) <= maxDiffLength && len(textB) <= maxDiffLength {
		return fmt.Sprintf("'%s' → '%s'", textA, textB)
	}

	// 简化的差异表示
	if textA == "" {
		return fmt.Sprintf("[新增文本] 长度: %d", len(textB))
	}
	if textB == "" {
		return "[删除文本]"
	}

	// 计算相似度
	similarity := calculateTextSimilarity(textA, textB)
	return fmt.Sprintf("相似度: %.1f%%", similarity*100)
}

// calculateTextSimilarity 计算文本相似度
func calculateTextSimilarity(textA, textB string) float64 {
	if textA == textB {
		return 1.0
	}

	// 使用Levenshtein距离计算相似度
	lenA, lenB := len(textA), len(textB)
	maxLen := lenA
	if lenB > maxLen {
		maxLen = lenB
	}

	if maxLen == 0 {
		return 1.0
	}

	distance := levenshteinDistance(textA, textB)
	return 1.0 - float64(distance)/float64(maxLen)
}

// levenshteinDistance 计算Levenshtein距离
func levenshteinDistance(s1, s2 string) int {
	len1, len2 := len(s1), len(s2)
	matrix := make([][]int, len1+1)

	for i := 0; i <= len1; i++ {
		matrix[i] = make([]int, len2+1)
		matrix[i][0] = i
	}

	for j := 0; j <= len2; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // 删除
				matrix[i][j-1]+1,      // 插入
				matrix[i-1][j-1]+cost, // 替换
			)
		}
	}

	return matrix[len1][len2]
}

func min(nums ...int) int {
	minNum := nums[0]
	for _, num := range nums[1:] {
		if num < minNum {
			minNum = num
		}
	}
	return minNum
}

// countAttributeChanges 计算属性变化总数
func countAttributeChanges(changes []AttrChange) int {
	total := 0
	for _, change := range changes {
		total += len(change.Added) + len(change.Removed) + len(change.Modified)
	}
	return total
}

// calculateSimilarityScore 计算相似度分数
func calculateSimilarityScore(comparison *DOMSnapshotComparison, totalNodesA, totalNodesB int) float64 {
	if totalNodesA == 0 && totalNodesB == 0 {
		return 1.0
	}

	totalNodes := max(totalNodesA, totalNodesB)
	if totalNodes == 0 {
		return 0.0
	}

	// 计算结构相似度
	nodeChangeScore := 1.0 - float64(comparison.NodeChanges)/float64(totalNodes)

	// 计算属性相似度
	attrChangeCount := comparison.AttributeChanges
	attrScore := 1.0
	if attrChangeCount > 0 {
		// 估算平均属性数
		estimatedAttrsPerNode := 5.0
		attrScore = 1.0 - float64(attrChangeCount)/(float64(totalNodes)*estimatedAttrsPerNode)
	}

	// 计算文本相似度
	textChangeCount := len(comparison.ChangedTexts)
	textScore := 1.0
	if textChangeCount > 0 {
		// 估算文本节点数
		textNodeCount := float64(totalNodes) * 0.3 // 假设30%是文本节点
		if textNodeCount > 0 {
			textScore = 1.0 - float64(textChangeCount)/textNodeCount
		}
	}

	// 综合分数（加权平均）
	similarity := nodeChangeScore*0.5 + attrScore*0.3 + textScore*0.2

	// 确保在0-1范围内
	if similarity < 0 {
		return 0
	}
	if similarity > 1 {
		return 1
	}

	return similarity
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 辅助函数：生成比较报告
func GenerateComparisonReport(comparison *DOMSnapshotComparison) string {
	report := strings.Builder{}

	report.WriteString("=== DOM快照比较报告 ===\n\n")
	report.WriteString(fmt.Sprintf("比较时间: %v\n", comparison.ComparisonTime))
	report.WriteString(fmt.Sprintf("快照A节点数: %d\n", comparison.TotalNodesA))
	report.WriteString(fmt.Sprintf("快照B节点数: %d\n", comparison.TotalNodesB))
	report.WriteString(fmt.Sprintf("总体相似度: %.1f%%\n\n", comparison.SimilarityScore*100))

	report.WriteString("=== 变化统计 ===\n")
	report.WriteString(fmt.Sprintf("节点变化总数: %d\n", comparison.NodeChanges))
	report.WriteString(fmt.Sprintf("  新增节点: %d\n", len(comparison.AddedNodes)))
	report.WriteString(fmt.Sprintf("  删除节点: %d\n", len(comparison.RemovedNodes)))
	report.WriteString(fmt.Sprintf("  修改节点: %d\n", len(comparison.ModifiedNodes)))
	report.WriteString(fmt.Sprintf("属性变化总数: %d\n", comparison.AttributeChanges))
	report.WriteString(fmt.Sprintf("文本变化总数: %d\n\n", comparison.TextChanges))

	// 详细变化
	if len(comparison.ModifiedNodes) > 0 {
		report.WriteString("=== 修改的节点 ===\n")
		for i, change := range comparison.ModifiedNodes {
			if i >= 10 { // 只显示前10个
				report.WriteString(fmt.Sprintf("  ... 还有 %d 个修改节点\n", len(comparison.ModifiedNodes)-10))
				break
			}
			report.WriteString(fmt.Sprintf("  [%d] 节点ID: %d, 类型: %s, 变化: %v\n",
				i+1, change.NodeID, change.NodeName, change.Changes))
		}
		report.WriteString("\n")
	}

	// 属性变化详情
	if len(comparison.ChangedAttributes) > 0 {
		report.WriteString("=== 属性变化详情 ===\n")
		for i, change := range comparison.ChangedAttributes {
			if i >= 5 { // 只显示前5个
				report.WriteString(fmt.Sprintf("  ... 还有 %d 个节点的属性变化\n", len(comparison.ChangedAttributes)-5))
				break
			}
			report.WriteString(fmt.Sprintf("  节点[%d]: %s\n", change.NodeID, change.NodeName))
			if len(change.Added) > 0 {
				for attr, value := range change.Added {
					report.WriteString(fmt.Sprintf("    + 新增属性: %s=\"%s\"\n", attr, value))
				}
			}
			if len(change.Removed) > 0 {
				report.WriteString(fmt.Sprintf("    - 删除属性: %v\n", change.Removed))
			}
			if len(change.Modified) > 0 {
				for attr, newValue := range change.Modified {
					oldValue := change.OldValues[attr]
					report.WriteString(fmt.Sprintf("    * 修改属性: %s=\"%s\" → \"%s\"\n", attr, oldValue, newValue))
				}
			}
		}
		report.WriteString("\n")
	}

	// 文本变化详情
	if len(comparison.ChangedTexts) > 0 {
		report.WriteString("=== 文本变化详情 ===\n")
		for i, change := range comparison.ChangedTexts {
			if i >= 3 { // 只显示前3个
				report.WriteString(fmt.Sprintf("  ... 还有 %d 个文本变化\n", len(comparison.ChangedTexts)-3))
				break
			}
			report.WriteString(fmt.Sprintf("  节点[%d]:\n", change.NodeID))
			report.WriteString(fmt.Sprintf("    变化: %s\n", change.Diff))
		}
	}

	return report.String()
}

// CDPDOMSnapshotDisable 禁用DOMSnapshot域
func CDPDOMSnapshotDisable() (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 DOMSnapshot.disable")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "DOMSnapshot.disable"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 DOMSnapshot.disable 失败:", err)
		return "", err
	}

	utils.Debugf("发送 CDP 消息: %s", message)
	timeout := 6 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				fmt.Println("[CDP DOMSnapshot.disable] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("DOMSnapshot.disable 请求超时")
		}
	}
}

// CDPDOMSnapshotEnable 启用DOMSnapshot域
func CDPDOMSnapshotEnable() (string, error) {
	if !DefaultBrowserWS() {
		return "", nil
	}
	if chromeInstance.BrowserWSConn == nil {
		return "", fmt.Errorf("BrowserWSConn 未连接，无法调用 DOMSnapshot.enable")
	}
	chromeInstance.NextID++
	reqID := chromeInstance.NextID
	message := fmt.Sprintf(`{
		"id": %d,
		"method": "DOMSnapshot.enable"
	}`, reqID)

	err := chromeInstance.BrowserWSConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("发送 DOMSnapshot.enable 失败:", err)
		return "", err
	}

	utils.Debugf("发送 CDP 消息: %s", message)
	timeout := 6 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case respMsg, ok := <-messageQueue:
			if !ok {
				return "", fmt.Errorf("消息队列已关闭")
			}
			if reqID == respMsg.ID {
				content := utils.JsonPrettyFormat(respMsg.Content)
				fmt.Println("[CDP DOMSnapshot.enable] 收到回复 -> ", content)
				return content, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("DOMSnapshot.enable 请求超时")
		}
	}
}
