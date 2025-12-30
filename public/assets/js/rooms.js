import { api } from "./helpers.js"

export function useRooms() {
    async function get() {
        return (await api('/rooms')).data
    }

    async function getById(id) {
        return (await api(`/rooms/${id}`)).data
    }

    async function checkAccess(id) {
        let status = true
        try {
            await api(`/rooms/${id}/access`)
        } catch {
            status = false
        }
        return status
    }

    async function auth(id, password) {
        if (!password) {
            throw Error('Password is required')
        }

        await api(`/rooms/${id}/auth`, {
            method: 'POST',
            body: JSON.stringify({ password }),
        })
    }

    async function logout(id) {
        await api(`/rooms/${id}/logout`, {
            method: 'POST'
        })
    }

    async function create(password, lifespan) {
        if (!password) {
            throw Error('Password is required')
        }
        if (!Number.isFinite(lifespan) || lifespan <= 0) {
            throw Error('Lifespan must be a positive number (seconds)')
        }

        const response = await api('/rooms', {
            method: 'POST',
            body: JSON.stringify({ password, lifespan }),
        })
        return response.data.id
    }

    async function remove(id) {
        if (!confirm('Delete this room?')) return false
        await api(`/rooms/${id}`, { method: 'DELETE' })
        return true
    }

    return {
        get,
        getById,
        checkAccess,
        auth,
        logout,
        create,
        remove
    }
}
