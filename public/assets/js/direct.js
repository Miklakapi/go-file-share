import { api } from "./helpers.js"

export function useDirect() {
    let controller = null

    async function download(code) {
        if (controller) throw Error('Another operation is already in progress')
        controller = new AbortController()

        try {
            await api(`/direct/${code}/download`, {
                signal: controller.signal
            }, 0)
        } catch (err) {
            if (err.name === 'AbortError') return
            throw err
        } finally {
            controller = null
        }
    }

    async function upload(code, file) {
        if (!file) throw new Error('No file provided')

        if (controller) throw Error('Another operation is already in progress')
        controller = new AbortController()

        const form = new FormData()
        form.append('file', file, file.name)

        try {
            await api(`/direct/${code}/upload`, {
                method: 'POST',
                body: form
            })
        } catch (err) {
            if (err.name === 'AbortError') return
            throw err
        } finally {
            controller = null
        }
    }

    function abort() {
        if (controller) {
            controller.abort()
            controller = null
        }
    }

    return {
        download,
        upload,
        abort
    }
}