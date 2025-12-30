import { useCreateDialog } from "./createDialog.js"
import { useRoomDataTable } from "./roomDataTable.js"
import { useLoginDialog } from "./loginDialog.js"
import { useRooms } from "./rooms.js"
import { useFiles } from "./files.js"
import { useRouter } from "./router.js"
import { useToast } from "./toast.js"
import { useFilesDataTable } from "./fileDataTable.js"
import { formatDate } from "./helpers.js"

const els = {
    // Others
    toast: () => document.getElementById('toastContainer'),
    backButton: () => document.getElementById('backBtn'),
    // Table
    tableBody: () => document.querySelector('#roomsTable tbody'),
    emptyState: () => document.getElementById('emptyState'),
    // Create dialog
    createDialog: () => document.getElementById('createDialog'),
    openCreateDialogBtn: () => document.getElementById('openCreateBtn'),
    //// Create dialog content
    createDialogForm: () => document.getElementById('createForm'),
    createDialogPassword: () => document.getElementById('createPassword'),
    createDialogLifespan: () => document.getElementById('createLifespan'),
    createDialogSubmitBtn: () => document.getElementById('createSubmitBtn'),
    createDialogError: () => document.getElementById('createError'),
    // Login dialog
    loginDialog: () => document.getElementById('loginDialog'),
    //// Login dialog content
    loginDialogForm: () => document.getElementById('loginForm'),
    loginDialogId: () => document.getElementById('loginId'),
    loginDialogPassword: () => document.getElementById('loginPassword'),
    loginDialogSubmitBtn: () => document.getElementById('loginSubmitBtn'),
    loginDialogError: () => document.getElementById('loginError'),
    // Room
    logoutBtn: () => document.getElementById('logoutBtn'),
    deleteRoomBtn: () => document.getElementById('deleteRoomBtn'),
    uploadForm: () => document.getElementById('uploadForm'),
    fileInput: () => document.getElementById('fileInput'),
    fileName: () => document.getElementById('fileName'),
    uploadBtn: () => document.getElementById('uploadBtn'),
    roomTitle: () => document.getElementById('roomTitle'),
    roomMeta: () => document.getElementById('roomMeta'),
    //// File table
    filesTableBody: () => document.querySelector('#filesTable tbody'),
    filesEmpty: () => document.getElementById('filesEmpty'),
}

const router = useRouter()
const toast = useToast(els.toast)
const createDialog = useCreateDialog(els.createDialog, els.createDialogPassword, els.createDialogLifespan, els.createDialogSubmitBtn, els.createDialogError)
const loginDialog = useLoginDialog(els.loginDialog, els.loginDialogId, els.loginDialogPassword, els.loginDialogSubmitBtn, els.loginDialogError)
const roomDataTable = useRoomDataTable(els.tableBody, els.emptyState)
const filesDataTable = useFilesDataTable(els.filesTableBody, els.filesEmpty)
const rooms = useRooms()
const files = useFiles()

function show(view) {
    document.getElementById('view-list').hidden = view !== 'list'
    document.getElementById('view-room').hidden = view !== 'room'
}

function wireEvents() {
    els.openCreateDialogBtn().addEventListener('click', createDialog.open)

    els.createDialogForm().addEventListener('submit', async (e) => {
        e.preventDefault()
        if (e.submitter.value === 'cancel') {
            createDialog.close()
            return
        }
        try {
            createDialog.disableSubmitButton(true)
            const id = await rooms.create((els.createDialogPassword().value || '').trim(), Number(els.createDialogLifespan().value))
            createDialog.close()
            toast.show('Created successfully!', 'success')
            router.navigate(`/rooms/${id}`)
        } catch (error) {
            createDialog.setError(`${error}`.replace("Error:", ""))
        } finally {
            createDialog.disableSubmitButton(false)
        }
    })

    els.loginDialogForm().addEventListener('submit', async (e) => {
        e.preventDefault()
        if (e.submitter.value === 'cancel') {
            loginDialog.close()
            router.navigate('/')
            return
        }
        try {
            loginDialog.disableSubmitButton(true)
            await rooms.auth(els.loginDialogId().value, (els.loginDialogPassword().value || '').trim())
            loginDialog.close()
            toast.show('Login successful!', 'success')
            router.navigate(`/rooms/${els.loginDialogId().value}`)
        } catch (error) {
            loginDialog.setError(`${error}`.replace("Error:", ""))
        } finally {
            loginDialog.disableSubmitButton(false)
        }
    })

    els.tableBody().addEventListener('click', async (e) => {
        const btn = e.target.closest('button[data-action]')
        if (!btn) return

        const action = btn.dataset.action
        const id = btn.dataset.id
        if (action === 'delete') {
            roomDataTable.disableRowButtons(id, true)
            try {
                const removed = await rooms.remove(id)
                if (removed) {
                    roomDataTable.removeRow(id)
                    toast.show('Deleted successfully!', 'success')
                }
            } catch (error) {
                toast.show(error, 'error')
            }
            roomDataTable.disableRowButtons(id, false)
        }
        if (action === 'enter') {
            router.navigate(`/rooms/${id}`)
        }
    })

    els.logoutBtn().addEventListener('click', async () => {
        try {
            const id = router.getRoomId()
            if (!id) throw Error('Unable to get room ID')
            await rooms.logout(id)
            toast.show('Logged out!', 'success')
            router.navigate('/')
        } catch (error) {
            toast.show(error, 'error')
        }
    })

    els.deleteRoomBtn().addEventListener('click', async () => {
        try {
            const id = router.getRoomId()
            if (!id) throw Error('Unable to get room ID')
            const removed = await rooms.remove(id)
            if (removed) {
                toast.show('Deleted successfully!', 'success')
                router.navigate('/')
            }
        } catch (error) {
            toast.show(error, 'error')
        }
    })

    els.uploadForm().addEventListener('submit', async (e) => {
        e.preventDefault()
        try {
            els.uploadBtn().disabled = true
            const input = els.fileInput()

            const id = router.getRoomId()
            if (!id) throw Error('Unable to get room ID')

            const filesToUpload = input?.files ? Array.from(input.files) : []
            if (filesToUpload.length === 0) {
                throw Error('No file selected')
            }

            const results = await Promise.allSettled(filesToUpload.map(file => files.upload(id, file)))

            const failed = results.filter(r => r.status === 'rejected')
            if (failed.length > 0) {
                toast.show(failed[0].reason, 'error')
            } else {
                toast.show('Upload completed', 'success')
            }

            input.value = ''
            els.fileName().textContent = 'No file selected'

            filesDataTable.loadData(await files.get(id))
        } catch (error) {
            console.log(error)
            toast.show(error, 'error')
        } finally {
            els.uploadBtn().disabled = false
        }
    })

    els.fileInput().addEventListener('change', () => {
        const files = els.fileInput().files
        els.fileName().textContent = files?.length
            ? `${files.length} file(s) selected`
            : 'No file selected'
        els.uploadBtn().disabled = !(files && files.length > 0)
    })

    els.filesTableBody().addEventListener('click', async (e) => {
        const btn = e.target.closest('button[data-action]')
        if (!btn) return

        let roomId = router.getRoomId()
        if (!roomId) {
            toast.show('Unable to get room ID', 'error')
            return
        }

        const action = btn.dataset.action
        const fileId = btn.dataset.id
        if (action === 'delete') {
            filesDataTable.disableRowButtons(fileId, true)
            try {
                const removed = await files.remove(roomId, fileId)
                if (removed) {
                    filesDataTable.removeRow(fileId)
                    toast.show('Deleted successfully!', 'success')
                }
            } catch (error) {
                toast.show(error, 'error')
            }
            filesDataTable.disableRowButtons(fileId, false)
        }
        if (action === 'download') {
            try {
                await files.download(roomId, fileId)
            } catch (error) {
                toast.show(error, 'error')
            }
        }
    })

    els.backButton().addEventListener('click', () => router.navigate('/'))
}

router.onRoute(async (from, to) => {
    if (to === '/') {
        show('list')
        try {
            roomDataTable.loadData(await rooms.get())
        } catch (error) {
            toast.show(error, 'error')
        }
        return
    }

    const roomId = router.getRoomId()
    if (roomId) {
        if (!await rooms.checkAccess(roomId)) {
            router.navigate('/')
            loginDialog.open(roomId)
            return
        }
        show('room')
        try {
            const [filesData, roomData] = await Promise.all([files.get(roomId), rooms.getById(roomId)])

            filesDataTable.loadData(filesData)

            const roomMeta = els.roomMeta()
            els.roomTitle().textContent = `Room ${roomData.id}`
            roomMeta.querySelector("#expiresMeta").textContent = formatDate(roomData.expiresAt)
            roomMeta.querySelector("#tokensMeta").textContent = roomData.tokens
        } catch (error) {
            toast.show(error, 'error')
        }
        return
    }
    router.replace('/')
})

document.addEventListener('DOMContentLoaded', () => {
    wireEvents()
})
