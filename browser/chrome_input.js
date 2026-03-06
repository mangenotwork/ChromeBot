(function() {
    // 统一占位符（和你原有代码保持一致）
    const xpath = __XPATH__;
    const newValue = __INPUTTEXT__;

    // 标记：是否为百度场景（通过XPath/selector特征判断）
    const isBaiduScene = xpath.includes('chat-textarea') || xpath === '#chat-textarea';

    // ====================== 通用工具函数 ======================
    // 1. 等待元素加载（通用版）
    function waitForElementByXPath(xpath, timeout = 3000) {
        return new Promise((resolve) => {
            const timer = setInterval(() => {
                const el = document.evaluate(
                    xpath,
                    document,
                    null,
                    XPathResult.FIRST_ORDERED_NODE_TYPE,
                    null
                ).singleNodeValue;
                if (el) {
                    clearInterval(timer);
                    resolve(el);
                }
            }, 50);
            setTimeout(() => {
                clearInterval(timer);
                resolve(document.evaluate(xpath, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue);
            }, timeout);
        });
    }

    // 2. 通用React状态更新（你原有方法1）
    function updateReactState(element, value) {
        const keys = Object.keys(element);
        const reactKey = keys.find(key =>
            key.startsWith('__reactInternalInstance$') ||
            key.startsWith('__reactFiber$')
        );

        if (reactKey && element[reactKey]) {
            let fiberNode = element[reactKey];
            while (fiberNode) {
                if (fiberNode.memoizedProps) {
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
                        element.value = value;
                        fiberNode.memoizedProps.onChange(syntheticEvent);
                        return true;
                    }
                    if (fiberNode.memoizedProps.value !== undefined && fiberNode.memoizedProps.onChange) {
                        const event = { target: { value: value }, currentTarget: element };
                        fiberNode.memoizedProps.onChange(event);
                        return true;
                    }
                }
                fiberNode = fiberNode.return;
            }
        }
        return false;
    }

    // 3. 通用React合成事件（你原有方法2）
    function triggerReactSyntheticEvent(element, value) {
        const reactEvent = new Event('input', { bubbles: true });
        reactEvent.simulated = true;
        reactEvent._reactName = 'onChange';
        reactEvent._targetInst = element[Object.keys(element).find(k => k.startsWith('__reactFiber'))];
        reactEvent.nativeEvent = new Event('input');
        element.dispatchEvent(reactEvent);

        const changeEvent = new Event('change', { bubbles: true });
        changeEvent.simulated = true;
        element.dispatchEvent(changeEvent);
        return true;
    }

    // 4. 通用备用方案（你原有最终手段）
    function fallbackInput(element, value) {
        const originalValue = element.value;
        element.value = value;
        element._value = newValue;

        const events = ['input', 'change', 'blur', 'focus'];
        for (const eventName of events) {
            const event = new Event(eventName, { bubbles: true, cancelable: true });
            if (eventName === 'input') {
                event.data = value;
                event.inputType = 'insertText';
                event.isComposing = false;
            }
            element.dispatchEvent(event);
        }

        setTimeout(() => { element.blur(); }, 100);
        setTimeout(() => { element.focus(); }, 200);
        console.log('使用备用方案完成输入');
        return true;
    }

    // ====================== 百度专属逻辑 ======================
    // 百度输入框专用输入方法（你的生效代码）
    async function inputToBaiduTextarea(el, text) {
        if (!el) {
            console.error('未找到百度输入框');
            return false;
        }

        // 聚焦+清空
        el.focus();
        el.value = '';
        el.dispatchEvent(new Event('focus', { bubbles: true }));

        // 初始change事件
        const initChangeEvent = new Event('change', { bubbles: true });
        Object.defineProperty(initChangeEvent, 'target', { get: () => ({ value: el.value }) });
        el.dispatchEvent(initChangeEvent);

        // 逐字符输入
        for (let i = 0; i < text.length; i++) {
            const char = text[i];
            el.value += char;

            // 原生InputEvent
            const inputEvent = new InputEvent('input', {
                bubbles: true,
                cancelable: true,
                data: char,
                inputType: 'insertText',
                isComposing: false,
                target: el
            });
            el.dispatchEvent(inputEvent);

            // React合成Change事件
            const reactChangeEvent = new Event('change', { bubbles: true });
            Object.defineProperties(reactChangeEvent, {
                target: { get: () => el },
                currentTarget: { get: () => el },
                value: { get: () => el.value }
            });
            el.dispatchEvent(reactChangeEvent);

            await new Promise(resolve => setTimeout(resolve, 30));
        }

        el.dispatchEvent(new Event('blur', { bubbles: true }));
        console.log('百度专属方案输入完成，最终值：', el.value);
        return true;
    }

    // ====================== 主执行逻辑 ======================
    async function main() {
        let element;
        // 1. 获取目标元素（兼容XPath和selector）
        if (xpath.startsWith('#') || xpath.startsWith('.')) {
            // 是CSS选择器（百度场景）
            element = await new Promise((resolve) => {
                const timer = setInterval(() => {
                    const el = document.querySelector(xpath);
                    if (el) { clearInterval(timer); resolve(el); }
                }, 50);
                setTimeout(() => { clearInterval(timer); resolve(document.querySelector(xpath)); }, 3000);
            });
        } else {
            // 是XPath（通用场景）
            element = await waitForElementByXPath(xpath);
        }

        if (!element) {
            console.error('未找到目标元素：', xpath);
            return false;
        }

        // 2. 分场景执行输入逻辑
        if (isBaiduScene) {
            // 百度场景：优先用专属逻辑
            await inputToBaiduTextarea(element, newValue);
        } else {
            // 通用场景：按你原有逻辑执行
            if (updateReactState(element, newValue)) {
                console.log('通过React状态更新成功');
            } else if (triggerReactSyntheticEvent(element, newValue)) {
                console.log('通过合成事件更新成功');
            } else {
                fallbackInput(element, newValue);
            }
        }
        return true;
    }

    // 执行主逻辑
    main();
})()