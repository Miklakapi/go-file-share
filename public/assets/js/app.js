import { useCreateDialog } from "./createDialog.js"
import { useDataTable } from "./dataTable.js"
import { useRooms } from "./rooms.js"
import { useRouter } from "./router.js"

const els = {
    // Others
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
}

const router = useRouter()

const {
    open: openCreateDialog,
    close: closeCreateDialog,
    disableSubmitButton,
    setError: setCreateDialogError
} = useCreateDialog(els.createDialog, els.createDialogPassword, els.createDialogLifespan, els.createDialogSubmitBtn, els.createDialogError)

const {
    get: getRooms,
    create: createRoom,
    remove: removeRoom
} = useRooms()

const {
    loadData,
    removeRow,
    disableRowButtons,
} = useDataTable(els.tableBody, els.emptyState)

function wireEvents() {
    els.openCreateDialogBtn().addEventListener('click', openCreateDialog)

    els.createDialogForm().addEventListener('submit', async (e) => {
        e.preventDefault()
        if (e.submitter.value === 'cancel') {
            closeCreateDialog()
            return
        }
        try {
            disableSubmitButton(true)
            await createRoom((els.createDialogPassword().value || '').trim(), Number(els.createDialogLifespan().value))
            closeCreateDialog()
            loadData(await getRooms())
        } catch (error) {
            setCreateDialogError(`${error}`.replace("Error:", ""))
        } finally {
            disableSubmitButton(false)
        }
    })

    els.tableBody().addEventListener('click', async (e) => {
        const btn = e.target.closest('button[data-action]')
        if (!btn) return

        const action = btn.dataset.action
        const id = btn.dataset.id
        if (action === 'delete') {
            try {
                disableRowButtons(id, true)
                if (await removeRoom(id) === true) removeRow(id)
            } finally {
                disableRowButtons(id, false)
            }
        }
        if (action === 'enter') router.navigate(`/rooms/${id}`)
    })

    els.backButton().addEventListener('click', () => router.navigate('/'))
}

function show(view) {
    document.getElementById('view-list').hidden = view !== 'list'
    document.getElementById('view-room').hidden = view !== 'room'
}

document.addEventListener('DOMContentLoaded', () => {
    wireEvents()
})

router.onRoute(async () => {
    if (location.pathname === '/' || location.pathname === '') {
        show('list')
        loadData(await getRooms())
        return
    }

    const roomId = router.getRoomId()

    if (roomId) {
        show('room')
        document.getElementById('roomTitle').textContent = `Room ${roomId}`
        document.getElementById('roomMeta').textContent = `ID: ${roomId}`
        return
    }

    router.replace('/')
})
