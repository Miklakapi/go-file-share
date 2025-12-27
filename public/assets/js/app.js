import { useCreateDialog } from "./createDialog.js"
import { useDataTable } from "./dataTable.js"
import { useLoginDialog } from "./loginDialog.js"
import { useRooms } from "./rooms.js"
import { useRouter } from "./router.js"
import { useToast } from "./toast.js"

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
}

const router = useRouter()
const toast = useToast(els.toast)
const createDialog = useCreateDialog(els.createDialog, els.createDialogPassword, els.createDialogLifespan, els.createDialogSubmitBtn, els.createDialogError)
const loginDialog = useLoginDialog(els.loginDialog, els.loginDialogId, els.loginDialogPassword, els.loginDialogSubmitBtn, els.loginDialogError)
const dataTable = useDataTable(els.tableBody, els.emptyState)
const rooms = useRooms()

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

            dataTable.disableRowButtons(id, true)
            try {
                const removed = await rooms.remove(id)
                if (removed) {
                    dataTable.removeRow(id)
                    toast.show('Deleted successfully!', 'success')
                }
            } catch (error) {
                toast.show(error, 'error')
            }

            dataTable.disableRowButtons(id, false)

        }
        if (action === 'enter') {
            if (!await rooms.checkAccess(id)) loginDialog.open(id)
            else router.navigate(`/rooms/${id}`)
        }
    })

    els.backButton().addEventListener('click', () => router.navigate('/'))
}

router.onRoute(async (from, to) => {
    if (to === '/') {
        show('list')
        dataTable.loadData(await rooms.get())
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
        document.getElementById('roomTitle').textContent = `Room ${roomId}`
        return
    }
    router.replace('/')
})

document.addEventListener('DOMContentLoaded', () => {
    wireEvents()
})
