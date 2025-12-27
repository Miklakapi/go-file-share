export function useRouter() {
    const listeners = new Set()

    function navigate(path) {
        history.pushState({}, '', path)
        notify()
    }

    function replace(path) {
        history.replaceState({}, '', path)
        notify()
    }

    function onRoute(fn) {
        listeners.add(fn)
        fn(location.pathname)
        return () => listeners.delete(fn)
    }

    function getRoomId() {
        const m = location.pathname.match(/^\/rooms\/([^/]+)\/?$/)
        return m ? m[1] : null
    }

    function notify() {
        for (const fn of listeners) fn(location.pathname)
    }

    window.addEventListener('popstate', notify)

    return {
        navigate,
        replace,
        onRoute,
        getRoomId
    }
}