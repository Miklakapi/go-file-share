export function useToast(toastContainerElement) {
    function show(message, type = 'info', ms = 2200) {
        const container = toastContainerElement()

        const toast = document.createElement('div')
        toast.className = `toast ${type}`
        toast.textContent = message

        container.appendChild(toast)

        requestAnimationFrame(() => toast.classList.add('show'))

        const hideTimer = setTimeout(() => {
            toast.classList.remove('show')

            setTimeout(() => {
                toast.remove()
            }, 200)
        }, ms)

        toast.addEventListener('click', () => {
            clearTimeout(hideTimer)
            toast.classList.remove('show')
            setTimeout(() => toast.remove(), 200)
        })
    }

    return {
        show
    }
}