(function() {
    try {
        // 获取完整的HTML
        const html = document.documentElement.outerHTML;

        // 如果需要包含DOCTYPE
        const doctype = new XMLSerializer().serializeToString(document.doctype);
        const fullHTML = doctype + html;

        return {
            html: fullHTML,
        };
    } catch (error) {
        return {
            success: false,
            error: error.message
        };
    }
})()