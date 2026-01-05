import { useCreateDialog } from "./createDialog.js"
import { useRoomDataTable } from "./roomDataTable.js"
import { useLoginDialog } from "./loginDialog.js"
import { useRooms } from "./rooms.js"
import { useFiles } from "./files.js"
import { useRouter } from "./router.js"
import { useToast } from "./toast.js"
import { useFilesDataTable } from "./fileDataTable.js"
import { formatDate, generateNumericCode, sleep } from "./helpers.js"
import { useSSE } from "./sse.js"
import { useDirectDialog } from "./directDialog.js"
import { useDirect } from "./direct.js"

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
    //// Direct upload
    directDialog: () => document.getElementById('directDialog'),
    openDirectBtn: () => document.getElementById('openDirectBtn'),
    //// Direct upload dialog content
    directForm: () => document.getElementById('directForm'),
    directReceiveBtn: () => document.getElementById('directReceiveBtn'),
    directCodeBox: () => document.getElementById('directCodeBox'),
    directCodeValue: () => document.getElementById('directCodeValue'),
    directCopyBtn: () => document.getElementById('directCopyBtn'),
    directSendBox: () => document.getElementById('directSendBox'),
    directCodeInput: () => document.getElementById('directCodeInput'),
    directFileInput: () => document.getElementById('directFileInput'),
    directSendBtn: () => document.getElementById('directSendBtn'),
    directDialogError: () => document.getElementById('directError'),
}

let suppressRoomsRefresh = false

const router = useRouter()
const toast = useToast(els.toast)
const createDialog = useCreateDialog(els.createDialog, els.createDialogPassword, els.createDialogLifespan, els.createDialogSubmitBtn, els.createDialogError)
const directDialog = useDirectDialog(
    els.directDialog, els.directDialogError,
    els.directCodeBox, els.directReceiveBtn, els.directCodeValue,
    els.directSendBox, els.directFileInput, els.directSendBtn
)
const loginDialog = useLoginDialog(els.loginDialog, els.loginDialogId, els.loginDialogPassword, els.loginDialogSubmitBtn, els.loginDialogError)
const roomDataTable = useRoomDataTable(els.tableBody, els.emptyState)
const filesDataTable = useFilesDataTable(els.filesTableBody, els.filesEmpty)
const rooms = useRooms()
const files = useFiles()
const sse = useSSE()
const direct = useDirect()

function show(view) {
    document.getElementById('view-list').hidden = view !== 'list'
    document.getElementById('view-room').hidden = view !== 'room'
}

function wireEvents() {
    els.openCreateDialogBtn().addEventListener('click', createDialog.open)
    els.openDirectBtn().addEventListener('click', directDialog.open)
    els.backButton().addEventListener('click', () => router.navigate('/'))

    els.createDialogForm().addEventListener('submit', async (e) => {
        e.preventDefault()
        if (e.submitter.value === 'cancel') {
            createDialog.close()
            return
        }
        try {
            createDialog.disableSubmitButton(true)
            suppressRoomsRefresh = true
            const id = await rooms.create((els.createDialogPassword().value || '').trim(), Number(els.createDialogLifespan().value))
            createDialog.close()
            toast.show('Created successfully!', 'success')
            router.navigate(`/rooms/${id}`)
        } catch (error) {
            createDialog.setError(`${error}`.replace("Error:", ""))
        } finally {
            createDialog.disableSubmitButton(false)
            suppressRoomsRefresh = false
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

    els.directForm().addEventListener('submit', async (e) => {
        e.preventDefault()
        if (e.submitter.value === 'cancel') {
            directDialog.close()
            return
        }

        if (e.submitter.value === 'generateCode') {
            const code = generateNumericCode(16)
            directDialog.setError("")
            directDialog.setCode(code)
            directDialog.disableCreateCodeButton(true)
            directDialog.disableSendBox(true)
            directDialog.showCodeBox(true)
            try {
                await direct.download(code)
                directDialog.clearCopySection()
                directDialog.disableCreateCodeButton(false)
                directDialog.disableSendBox(false)
            } catch (error) {
                directDialog.setError(`${error}`.replace("Error:", ""))
            }
            return
        }

        if (e.submitter.value === 'copyCode') {
            try {
                await directDialog.copyCode()
            } catch (error) {
                toast.show(error, 'error')
                return
            }
            toast.show('Copied!', 'success')
            return
        }

        if (e.submitter.value === 'cancelCode') {
            direct.abort()
            directDialog.clearCopySection()
            directDialog.disableCreateCodeButton(false)
            directDialog.disableSendBox(false)
            return
        }

        if (e.submitter.value === 'sendFile') {
            directDialog.setError("")
            const fileInput = els.directFileInput()
            let fileToUpload = fileInput?.files ? Array.from(fileInput.files) : []
            if (fileToUpload.length === 0) {
                directDialog.setError("No file to upload")
                toast.show("No file to upload", 'error')
                return
            }
            fileToUpload = fileToUpload[0]

            directDialog.disableCreateCodeButton(true)
            directDialog.disableSendButton(true)
            try {
                await direct.upload(els.directCodeInput().value?.trim() ?? '', fileToUpload)
                directDialog.clearFileInput()
            } catch (error) {
                directDialog.setError(`${error}`.replace("Error:", ""))
            } finally {
                directDialog.disableCreateCodeButton(false)
                directDialog.disableSendButton(false)
            }
            return
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

    sse.onMessage(async e => console.info(e.data))
    sse.onEvent("RoomsChange", async e => {
        if (suppressRoomsRefresh) return
        if (router.getLocation() !== '/') return
        roomDataTable.loadData(await rooms.get())
    })
    sse.onEvent("Message", e => toast.show(e.data, 'success'))
    sse.onError(err => console.error("EventSource failed:", err))
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
