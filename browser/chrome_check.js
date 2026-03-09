(() => {
    try {
        const xpath = __XPATH__;
        const result = document.evaluate(
            xpath,
            document,
            null,
            XPathResult.FIRST_ORDERED_NODE_TYPE,
            null
        );
        return result.singleNodeValue !== null;
    } catch (error) {
        return false;
    }
})()