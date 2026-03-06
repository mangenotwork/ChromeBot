(() => {
    try {
        let buttonXPath = __XPATH__;
        let button = document.evaluate(
            buttonXPath,
            document,
            null,
            XPathResult.FIRST_ORDERED_NODE_TYPE,
            null
        ).singleNodeValue;
        if (button) {
            button.click();
            return true;
        }
        return false;
    } catch (error) {
        return false;
    }
})()