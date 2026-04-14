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
