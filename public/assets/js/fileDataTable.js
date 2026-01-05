import { formatDate } from "./helpers.js"

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
            tr.dataset.id = f.id

            const name = f.name ?? ''
            const title = escapeAttr(name)

            tr.innerHTML = `
              <td title="${title}">${name}</td>
              <td>${formatBytes(f.size ?? 0)}</td>
              <td>${formatDate(f.createdAt)}</td>
              <td class="right actions">
                <button class="btn" data-action="download" data-id="${f.id}">Download</button>
                <button class="btn danger" data-action="delete" data-id="${f.id}">Delete</button>
              </td>
            `
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

    function disableRowButtons(id, disabled = true) {
        const tbody = tableBodyElement()
        const tr = tbody.querySelector(`tr[data-id="${id}"]`)
        if (!tr) return

        const buttons = tr.querySelectorAll('button')
        buttons.forEach(btn => {
            btn.disabled = disabled
        })
    }

    function showEmptyState(show) {
        emptyElement().hidden = !show
    }

    function formatBytes(bytes) {
        const b = Number(bytes || 0)
        if (!b) return '0 B'

        const units = ['B', 'KB', 'MB', 'GB', 'TB']
        let i = 0
        let n = b

        while (n >= 1024 && i < units.length - 1) {
            n /= 1024
            i++
        }

        const val = i === 0 ? Math.round(n).toString() : n.toFixed(1)
        return `${val} ${units[i]}`
    }

    function escapeAttr(s) {
        return String(s ?? '')
            .replaceAll('&', '&amp;')
            .replaceAll('"', '&quot;')
            .replaceAll('<', '&lt;')
            .replaceAll('>', '&gt;')
    }

    return {
        loadData,
        removeRow,
        disableRowButtons,
        showEmptyState
    }
}