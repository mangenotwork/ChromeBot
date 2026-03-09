(() => {
    // 定义返回结果结构
    const result = {
        success: false,
        error: null,
        message: ""
    };
    try {
        // 独立参数变量（替换原JSON解析）
        const xpath = __SCROLL_XPATH__;
        const isSmooth = __SCROLL_IS_SMOOTH__;

        // 验证参数合法性
        if (!xpath || typeof xpath !== 'string') {
            result.error = "参数错误";
            result.message = "xpath 不能为空且必须为字符串";
            return result;
        }
        if (typeof isSmooth !== 'boolean') {
            result.error = "参数错误";
            result.message = "isSmooth 必须为布尔类型";
            return result;
        }

        // 查找目标元素
        const element = document.evaluate(
            xpath,
            document,
            null,
            XPathResult.FIRST_ORDERED_NODE_TYPE,
            null
        ).singleNodeValue;

        // 检查元素是否存在
        if (!element) {
            result.error = "元素未找到";
            result.message = `XPath: ${xpath} 未匹配到任何元素`;
            return result;
        }

        // 执行元素滚动（兼容平滑滚动）
        if (isSmooth && 'scrollBehavior' in document.documentElement.style) {
            element.scrollIntoView({
                behavior: 'smooth',
                block: 'center', // 垂直居中显示元素
                inline: 'center' // 水平居中显示元素
            });
        } else {
            element.scrollIntoView(true); // 降级处理
        }

        // 滚动成功
        result.success = true;
        result.message = `成功滚动到 XPath: ${xpath} 对应的元素`;
    } catch (error) {
        // 捕获执行异常
        result.error = "执行异常";
        result.message = error.message;
        result.stack = error.stack; // 可选：保留堆栈信息用于调试
    }
    return result;
})()