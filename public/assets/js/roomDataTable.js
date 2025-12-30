import { formatDate, shortId } from "./helpers.js"

export function useRoomDataTable(tableBodyElement, emptyElement) {
    function loadData(rooms) {
        const tbody = tableBodyElement()
        tbody.innerHTML = ''

        if (!rooms || rooms.length === 0) {
            showEmptyState(true)
            return
        }
        showEmptyState(false)

        for (const r of rooms) {
            const tr = document.createElement('tr')
            tr.dataset.id = r.id
            tr.innerHTML = `
              <td title="${r.id}">${shortId(r.id)}</td>
              <td>${formatDate(r.expiresAt)}</td>
              <td>${r.files ?? 0}</td>
              <td>${r.tokens ?? 0}</td>
              <td class="right actions">
                <button class="btn" data-action="enter" data-id="${r.id}">Enter</button>
                <button class="btn danger" data-action="delete" data-id="${r.id}">Delete</button>
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