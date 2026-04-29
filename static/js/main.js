const copyButton = document.getElementById("copy-content")
if (copyButton) {
    copyButton.addEventListener("click", copyContent)
}

async function copyContent() {
    const content = document.getElementById("content").innerText
    try {
        await navigator.clipboard.writeText(content)
    } catch (err) {
        alert("Failed to copy content: " + err)
    }
}

const pasteURL = document.querySelector(".paste-url")
if (pasteURL) {
    const original = pasteURL.innerText
    pasteURL.addEventListener("mouseenter", () => pasteURL.innerHTML = "<em>click to copy</em>")
    pasteURL.addEventListener("mouseleave", () => pasteURL.innerText = original)
    pasteURL.addEventListener("click", () => navigator.clipboard.writeText(original))
}

const textArea = document.querySelector('textarea')
if (textArea) {
    textArea.addEventListener('keydown', (e) => {
        if (e.key === 'Tab') {
            // prevent default focus-jump behavior
            e.preventDefault()
            const cursorStart = textArea.selectionStart
            const cursorEnd = textArea.selectionEnd
            // insert 4 spaces at cursor position
            textArea.value = textArea.value.substring(0, cursorStart) + '    ' + textArea.value.substring(cursorEnd)
            // move cursor to after the inserted spaces
            textArea.selectionStart = textArea.selectionEnd = cursorStart + 4
        }
    })
}
