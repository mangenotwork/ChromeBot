(() => {
    const result = {
        success: false,
        message: "",
        elementFound: false,
        error: null
    };

    try {
        // 1. 替换XPath参数（确保XPath正确）
        const buttonXPath = __XPATH__;

        // 2. 查找目标元素
        const button = document.evaluate(
            buttonXPath,
            document,
            null,
            XPathResult.FIRST_ORDERED_NODE_TYPE,
            null
        ).singleNodeValue;

        if (!button) {
            result.message = `未找到XPath对应的元素: ${buttonXPath}`;
            result.elementFound = false;
            return result;
        }
        result.elementFound = true;

        // 3. 前置操作：确保元素在可视区域（滚动到元素位置）
        button.scrollIntoView({
            behavior: 'auto', // 禁用平滑滚动，避免检测
            block: 'center',
            inline: 'center'
        });

        // 4. 核心：融合所有有效点击方案的函数
        function simulateUltimateClick(element) {
            // ========== 方案1：精准坐标触发完整鼠标事件序列（你的代码1核心） ==========
            const rect = element.getBoundingClientRect();
            // 取元素中心坐标（更贴近真人点击位置）
            const clickX = rect.left + rect.width / 2;
            const clickY = rect.top + rect.height / 2;

            // 定义完整的鼠标事件序列
            const mouseEvents = [
                ['mousedown', { button: 0, buttons: 1, clientX: clickX, clientY: clickY }],
                ['mouseup', { button: 0, buttons: 0, clientX: clickX, clientY: clickY }],
                ['click', { button: 0, clientX: clickX, clientY: clickY }]
            ];

            // 触发原生鼠标事件
            mouseEvents.forEach(([type, options]) => {
                const event = new MouseEvent(type, {
                    view: window,
                    bubbles: true,        // 必须冒泡，适配父元素事件监听
                    cancelable: true,
                    button: options.button,
                    buttons: options.buttons || 0,
                    clientX: options.clientX,
                    clientY: options.clientY,
                    screenX: window.screenX + options.clientX,
                    screenY: window.screenY + options.clientY
                });
                element.dispatchEvent(event);
            });

            // ========== 方案2：适配React/Vue等框架的自定义事件（你的代码2核心） ==========
            const frameworkClickEvent = new Event('click', {
                bubbles: true,
                cancelable: true
            });
            // 添加框架可能依赖的属性
            frameworkClickEvent._synthetic = true;
            frameworkClickEvent._reactName = 'onClick';
            frameworkClickEvent.nativeEvent = new MouseEvent('click', { bubbles: true });
            element.dispatchEvent(frameworkClickEvent);

            // ========== 方案3：键盘Enter键触发（你的代码3核心，兜底方案） ==========
            // 先聚焦元素
            element.focus({ preventScroll: false });
            element.dispatchEvent(new Event('focus', { bubbles: true }));

            // 触发Enter键按下+抬起
            const enterDownEvent = new KeyboardEvent('keydown', {
                key: 'Enter',
                code: 'Enter',
                keyCode: 13,
                charCode: 13,
                which: 13,
                bubbles: true,
                cancelable: true
            });
            element.dispatchEvent(enterDownEvent);

            const enterUpEvent = new KeyboardEvent('keyup', {
                key: 'Enter',
                code: 'Enter',
                keyCode: 13,
                charCode: 13,
                which: 13,
                bubbles: true,
                cancelable: true
            });
            element.dispatchEvent(enterUpEvent);

            // ========== 最终兜底：原生click ==========
            element.click();

            // ========== 针对a标签的终极兜底：主动跳转 ==========
            if (element.tagName === 'A' && element.href) {
                const cleanHref = element.href.trim();
                // 延迟100ms跳转（模拟真人操作延迟）
                setTimeout(() => {
                    window.location.href = cleanHref;
                }, 100);
                return { success: true, href: cleanHref };
            }

            return { success: true };
        }

        // 执行终极点击
        const clickResult = simulateUltimateClick(button);
        if (!clickResult.success) {
            result.message = "模拟点击函数执行失败";
            return result;
        }

        // 5. 结果返回
        if (clickResult.href) {
            result.message = `点击成功，目标链接：${clickResult.href}（已触发跳转）`;
        } else {
            result.message = "点击成功，非链接元素";
        }
        result.success = true;
        return result;

    } catch (error) {
        result.error = error.message;
        result.message = `点击失败：${error.message}`;
        return result;
    }
})();