export function useCreateDialog(dialogElement, passwordField, lifespanField, submitButton, errorElement) {
    function open() {
        setError('')
        passwordField().value = ''
        lifespanField().value = '600'
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