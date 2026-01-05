export function useLoginDialog(dialogElement, idField, passwordField, submitButton, errorElement) {
    function open(id) {
        setError('')
        idField().value = id
        passwordField().value = ''
        dialogElement().showModal()
        passwordField().focus()
    }

    function close() {
        dialogElement().close()
    }

    function disableSubmitButton(disabled = true) {
        submitButton().disabled = disabled
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
        disableSubmitButton,
        setError
    }
}