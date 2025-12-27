export function useFilesDataTable(tableBodyElement, emptyElement) {
    function loadData(files) {
        const tbody = tableBodyElement()
        tbody.innerHTML = ''

        if (!files || files.length === 0) {
            showEmptyState(true)
            return
        }
        showEmptyState(false)

        for (const f of files) {
            const tr = document.createElement('tr')
            tr.dataset.id = f.ID
            tr.innerHTML = ``
            tbody.appendChild(tr)
        }
    }

    function removeRow(id) {
        const tbody = tableBodyElement()
        const tr = tbody.querySelector(`tr[data-id="${id}"]`)
        if (!tr) return

        tr.remove()

        if (tbody.children.length === 0) {
            showEmptyState(true)
        }
    }

    function disableRowButtons(id, disable = true) {
        const tbody = tableBodyElement()
        const tr = tbody.querySelector(`tr[data-id="${id}"]`)
        if (!tr) return

        const buttons = tr.querySelectorAll('button')
        buttons.forEach(btn => {
            btn.disabled = disable
        })
    }

    function showEmptyState(show) {
        emptyElement().hidden = !show
    }

    return {
        loadData,
        removeRow,
        disableRowButtons,
        showEmptyState
    }
}