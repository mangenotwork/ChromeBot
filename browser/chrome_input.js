(function() {
    // 统一占位符
    const xpath = __XPATH__;
    const newValue = __INPUTTEXT__;
    // 用于存储最终执行结果（CDP需要显式返回）
    let finalResult = {
        success: false,
        message: '',
        elementFound: false,
        inputValue: ''
    };

    // ====================== CDP适配：同步等待元素（避免异步时序问题） ======================
    /**
     * 同步等待元素（CDP环境下异步定时器易出问题）
     * @param {string} xpath - XPath表达式
     * @param {number} timeout - 超时时间（毫秒）
     * @returns {HTMLElement|null}
     */
    function waitForElementByXPathSync(xpath, timeout = 5000) {
        const startTime = Date.now();
        while (Date.now() - startTime < timeout) {
            try {
                const el = document.evaluate(
                    xpath,
                    document,
                    null,
                    XPathResult.FIRST_ORDERED_NODE_TYPE,
                    null
                ).singleNodeValue;
                if (el && (el.tagName === 'TEXTAREA' || el.tagName === 'INPUT')) {
                    return el;
                }
            } catch (e) {
                console.error('XPath解析错误:', e);
            }
            // 模拟等待（同步阻塞，CDP环境更稳定）
            for (let i = 0; i < 1000000; i++) {}
        }
        return null;
    }

    // ====================== React适配：更新React受控组件状态 ======================
    /**
     * 查找React Fiber节点并更新受控组件状态
     * @param {HTMLElement} element - 目标输入框元素
     * @param {string} value - 要输入的文本值
     * @returns {boolean} 是否更新成功
     */
    function updateReactState(element, value) {
        // 查找React内部属性（兼容不同React版本）
        const keys = Object.keys(element);
        const reactKey = keys.find(key =>
            key.startsWith('__reactInternalInstance$') ||
            key.startsWith('__reactFiber$')
        );

        if (reactKey && element[reactKey]) {
            let fiberNode = element[reactKey];

            // 向上查找最近的state节点
            while (fiberNode) {
                if (fiberNode.memoizedProps) {
                    // 尝试找到onChange处理器
                    if (fiberNode.memoizedProps.onChange) {
                        const syntheticEvent = {
                            target: element,
                            currentTarget: element,
                            type: 'change',
                            nativeEvent: new Event('change'),
                            preventDefault: () => {},
                            stopPropagation: () => {},
                            isDefaultPrevented: () => false,
                            isPropagationStopped: () => false
                        };

                        // 更新DOM值
                        element.value = value;

                        // 调用onChange处理器更新React state
                        fiberNode.memoizedProps.onChange(syntheticEvent);
                        return true;
                    }

                    // 尝试找到value属性（受控组件）
                    if (fiberNode.memoizedProps.value !== undefined) {
                        if (fiberNode.memoizedProps.onChange) {
                            const event = {
                                target: { value: value },
                                currentTarget: element
                            };
                            fiberNode.memoizedProps.onChange(event);
                            return true;
                        }
                    }
                }

                // 继续向上查找父Fiber节点
                fiberNode = fiberNode.return;
            }
        }

        return false;
    }

    /**
     * 触发React合成事件（备用方案）
     * @param {HTMLElement} element - 目标输入框元素
     * @param {string} value - 要输入的文本值
     * @returns {boolean} 是否触发成功
     */
    function triggerReactSyntheticEvent(element, value) {
        // 设置DOM基础值
        element.value = value;

        // 构造React专用Input事件
        const reactEvent = new Event('input', { bubbles: true });
        reactEvent.simulated = true;
        reactEvent._reactName = 'onChange';
        reactEvent._targetInst = element[Object.keys(element).find(k => k.startsWith('__reactFiber'))];
        reactEvent.nativeEvent = new Event('input');

        // 触发input事件
        element.dispatchEvent(reactEvent);

        // 触发change事件
        const changeEvent = new Event('change', { bubbles: true });
        changeEvent.simulated = true;
        element.dispatchEvent(changeEvent);

        return true;
    }

    /**
     * React组件终极输入方案（兜底逻辑）
     * @param {HTMLElement} element - 目标输入框元素
     * @param {string} value - 要输入的文本值
     * @returns {boolean} 是否执行成功
     */
    function reactFallbackInput(element, value) {
        // 备份原始值
        const originalValue = element.value;

        // 强制设置DOM值
        element.value = value;
        element._value = value; // 兼容部分框架内部属性

        // 触发所有可能的事件
        const events = ['input', 'change', 'blur', 'focus'];
        for (const eventName of events) {
            const event = new Event(eventName, {
                bubbles: true,
                cancelable: true
            });

            // 为input事件补充CDP环境所需属性
            if (eventName === 'input') {
                event.data = value;
                event.inputType = 'insertText';
                event.isComposing = false;
            }

            element.dispatchEvent(event);
        }

        // 模拟用户操作时序（CDP环境下同步替代setTimeout）
        element.blur(); // 模拟失去焦点触发验证
        // 同步延迟替代setTimeout
        for (let i = 0; i < 500000; i++) {}
        element.focus(); // 重新获取焦点

        return true;
    }

    // ====================== 百度输入框专用：CDP适配版 ======================
    /**
     * CDP环境下百度输入框输入（强制触发所有必要事件）
     * @param {HTMLElement} el - 输入框元素
     * @param {string} text - 输入文本
     * @returns {boolean}
     */
    function inputToBaiduTextareaCDP(el, text) {
        if (!el || !text) {
            finalResult.message = '元素或文本为空';
            return false;
        }

        try {
            // 1. 强制聚焦（CDP环境下需主动激活）
            el.focus();
            el.scrollIntoView({ behavior: 'instant' }); // 确保元素在视口内

            // 2. 清空原有值（强制赋值，绕过框架拦截）
            el.value = '';
            // 触发原生focus事件
            el.dispatchEvent(new Event('focus', { bubbles: true, cancelable: true }));

            // 3. 逐字符输入（模拟真实用户输入，CDP环境下更易被识别）
            let currentValue = '';
            for (let i = 0; i < text.length; i++) {
                const char = text[i];
                currentValue += char;

                // 强制赋值
                el.value = currentValue;

                // 构造标准InputEvent（CDP环境下需完整参数）
                const inputEvent = new InputEvent('input', {
                    bubbles: true,
                    cancelable: true,
                    data: char,
                    inputType: 'insertText',
                    isComposing: false,
                    target: el
                });
                el.dispatchEvent(inputEvent);

                // 构造Change事件（模拟用户输入后的change）
                const changeEvent = new Event('change', { bubbles: true, cancelable: true });
                Object.defineProperty(changeEvent, 'target', {
                    get: () => ({ value: currentValue })
                });
                el.dispatchEvent(changeEvent);

                // 短延迟（模拟真实打字间隔）
                for (let i = 0; i < 500000; i++) {}
            }

            // 4. 最终失焦+验证
            el.dispatchEvent(new Event('blur', { bubbles: true, cancelable: true }));
            finalResult.success = true;
            finalResult.message = '百度专属输入成功';
            finalResult.inputValue = el.value;
            return true;
        } catch (e) {
            finalResult.message = '百度输入逻辑报错: ' + e.message;
            // 终极兜底：直接赋值
            el.value = text;
            el.dispatchEvent(new Event('input', { bubbles: true }));
            finalResult.success = true;
            finalResult.inputValue = el.value;
            return true;
        }
    }

    // ====================== 主执行逻辑（CDP适配：同步+显式返回） ======================
    try {
        // 1. 同步查找元素（CDP环境下异步定时器不可靠）
        const element = waitForElementByXPathSync(xpath);
        if (!element) {
            finalResult.message = '未找到目标元素: ' + xpath;
            finalResult.elementFound = false;
        } else {
            finalResult.elementFound = true;

            // 优先处理百度专属输入框
            if (element.id === 'chat-textarea') {
                inputToBaiduTextareaCDP(element, newValue);
            }
            // 其次尝试React组件输入逻辑
            else if (updateReactState(element, newValue)) {
                finalResult.success = true;
                finalResult.message = '通过React状态更新成功';
                finalResult.inputValue = element.value;
            }
            else if (triggerReactSyntheticEvent(element, newValue)) {
                finalResult.success = true;
                finalResult.message = '通过React合成事件更新成功';
                finalResult.inputValue = element.value;
            }
            // 最后使用React兜底方案或通用输入逻辑
            else if (reactFallbackInput(element, newValue)) {
                finalResult.success = true;
                finalResult.message = '使用React备用方案输入成功';
                finalResult.inputValue = element.value;
            }
            // 通用输入逻辑（兜底）
            else {
                element.focus();
                element.value = newValue;
                ['input', 'change', 'blur'].forEach(evt => {
                    element.dispatchEvent(new Event(evt, { bubbles: true }));
                });
                finalResult.success = true;
                finalResult.message = '通用输入逻辑执行成功';
                finalResult.inputValue = element.value;
            }
        }
    } catch (globalErr) {
        finalResult.message = '全局执行错误: ' + globalErr.message;
        finalResult.success = false;
    }

    // ====================== CDP关键：显式返回结果（避免undefined） ======================
    console.log('CDP输入执行结果:', finalResult);
    return finalResult; // 必须显式返回，CDP才能拿到结果
})()