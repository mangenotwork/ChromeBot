(() => {
    const result = {
        success: false,
        error: null,
        message: "",
        before: { x: 0, y: 0 },
        after: { x: 0, y: 0 },
        target: { x: 0, y: 0 },
        maxScroll: { x: 0, y: 0 }
    };

    try {
        // 独立参数
        const targetX = __SCROLL_X__;
        const targetY = __SCROLL_Y__;

        // 1. 参数校验
        if (typeof targetX !== 'number' || typeof targetY !== 'number') {
            result.error = "参数错误";
            result.message = "x/y 必须为数字类型";
            return result;
        }

        // 2. 计算页面最大滚动范围（关键：判断是否超出边界）
        const doc = document.documentElement;
        const body = document.body;
        result.maxScroll.x = Math.max(doc.scrollWidth, body.scrollWidth) - window.innerWidth;
        result.maxScroll.y = Math.max(doc.scrollHeight, body.scrollHeight) - window.innerHeight;

        // 3. 修正目标坐标（不超过最大滚动范围）
        const finalX = Math.min(Math.max(0, targetX), result.maxScroll.x);
        const finalY = Math.min(Math.max(0, targetY), result.maxScroll.y);
        result.target.x = finalX;
        result.target.y = finalY;

        // 4. 记录滚动前位置
        result.before.x = window.scrollX || doc.scrollLeft || body.scrollLeft;
        result.before.y = window.scrollY || doc.scrollTop || body.scrollTop;

        // 5. 终极滚动方案（多容器全覆盖，解决滚动不到位）
        const forceScroll = () => {
            // 方案1：标准window滚动
            window.scrollTo(finalX, finalY);
            // 方案2：直接设置html/body滚动值（兜底）
            doc.scrollLeft = finalX;
            doc.scrollTop = finalY;
            body.scrollLeft = finalX;
            body.scrollTop = finalY;
            // 方案3：针对body嵌套的特殊处理
            if (document.scrollingElement) {
                document.scrollingElement.scrollLeft = finalX;
                document.scrollingElement.scrollTop = finalY;
            }
        };

        // 立即执行滚动（多次执行确保生效）
        forceScroll();
        setTimeout(forceScroll, 50); // 延迟50ms再执行一次

        // 6. 记录滚动后位置
        result.after.x = window.scrollX || doc.scrollLeft || body.scrollLeft;
        result.after.y = window.scrollY || doc.scrollTop || body.scrollTop;

        // 7. 校验结果（允许±5px误差）
        const xDiff = Math.abs(result.after.x - finalX);
        const yDiff = Math.abs(result.after.y - finalY);
        const isReached = xDiff <= 5 && yDiff <= 5;

        if (isReached) {
            result.success = true;
            result.message = `滚动成功：目标(${finalX},${finalY})，实际(${result.after.x},${result.after.y})`;
        } else if (finalY >= result.maxScroll.y || finalX >= result.maxScroll.x) {
            result.success = true;
            result.message = `已滚动到最大范围：目标(${targetX},${targetY})，修正后(${finalX},${finalY})，最大(${result.maxScroll.x},${result.maxScroll.y})`;
        } else {
            result.error = "滚动未到位";
            result.message = `滚动失败：目标(${finalX},${finalY})，实际(${result.after.x},${result.after.y})，最大(${result.maxScroll.x},${result.maxScroll.y})`;
        }

    } catch (error) {
        result.error = "执行异常";
        result.message = error.message;
        result.stack = error.stack;
    }

    return result;
})()