export function useRouter() {
    const listeners = new Set()
    let currentPath = location.pathname

    function navigate(path) {
        if (!path.startsWith('/')) path = '/' + path
        if (path === currentPath) return

        const from = currentPath
        history.pushState({}, '', path)
        currentPath = location.pathname

        notify(from, currentPath)
    }

    function replace(path) {
        if (!path.startsWith('/')) path = '/' + path
        if (path === currentPath) return

        const from = currentPath
        history.replaceState({}, '', path)
        currentPath = location.pathname

        notify(from, currentPath)
    }

    function onRoute(fn) {
        listeners.add(fn)
        fn(currentPath, currentPath)
        return () => listeners.delete(fn)
    }

    function getRoomId() {
        const m = location.pathname.match(/^\/rooms\/([^/]+)\/?$/)
        return m ? m[1] : null
    }

    function notify(from, to) {
        for (const fn of listeners) fn(from, to)
    }

    window.addEventListener('popstate', () => {
        const from = currentPath
        currentPath = location.pathname
        notify(from, currentPath)
    })

    return {
        navigate,
        replace,
        onRoute,
        getRoomId
    }
}