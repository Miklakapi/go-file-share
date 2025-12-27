import { api } from "./helpers.js"

export function useRooms() {
    async function get() {
        const rooms = (await api('/rooms')).rooms
        return rooms
    }

    async function create(password, lifespan) {
        if (!password) {
            throw Error('Password is required')
        }
        if (!Number.isFinite(lifespan) || lifespan <= 0) {
            throw Error('Lifespan must be a positive number (seconds)')
        }

        try {
            await api('/rooms', {
                method: 'POST',
                body: JSON.stringify({ password, lifespan }),
            })
        } catch (e) {
            throw Error(e.message)
        }
    }

    async function remove(id) {
        if (!confirm('Delete this room?')) return false
        await api(`/rooms/${id}`, { method: 'DELETE' })
        return null
    }

    return {
        get,
        create,
        remove
    }
}
