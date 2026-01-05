export function useDirectDialog(
    dialogElement, errorElement,
    // Receive
    codeBox, createCodeButton, codeField,
    // Send
    sendBox, fileInput, sendButton
) {
    function open() {
        setError('')
        clearCopySection()
        disableCreateCodeButton(false)
        clearSendSection()
        dialogElement().showModal()
    }

    function close() {
        dialogElement().close()
    }

    function showCodeBox(show = true) {
        codeBox().hidden = !show
    }

    function disableCreateCodeButton(disabled = true) {
        createCodeButton().disabled = disabled
    }

    function disableSendBox(disabled = true) {
        const el = sendBox()

        if (disabled) {
            el.classList.add('disabled')
            el.inert = true
        } else {
            el.classList.remove('disabled')
            el.inert = false
        }
    }

    function disableSendButton(disabled = true) {
        sendButton().disabled = disabled
    }

    function setCode(code) {
        codeField().textContent = code
    }

    async function copyCode() {
        const text = codeField().textContent?.trim() ?? ''
        if (!text) return

        if (navigator.clipboard?.writeText) {
            return await navigator.clipboard.writeText(codeField().textContent ?? '')
        }
        throw Error("Connection is not secure")
    }

    function clearFileInput() {
        fileInput().value = ''
    }

    function clearCopySection() {
        setCode('')
        showCodeBox(false)
    }

    function clearSendSection() {
        disableSendBox(false)
        disableSendButton(false)
        clearFileInput()
    }

    function setError(msg) {
        const el = errorElement()
        if (!msg) {
            el.hidden = true
            el.textContent = ''
            return
        }
        el.hidden = false
        el.textContent = msg
    }

    return {
        open,
        close,
        showCodeBox,
        disableCreateCodeButton,
        disableSendBox,
        disableSendButton,
        setCode,
        copyCode,
        clearFileInput,
        clearCopySection,
        clearSendSection,
        setError,
    }
}