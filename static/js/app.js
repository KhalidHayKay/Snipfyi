// ── Toast ──
function showToast(msg, type = 'success') {
	document.querySelector('.toast')?.remove();
	const t = document.createElement('div');
	t.className = `toast ${type}`;
	t.textContent = msg;
	document.body.appendChild(t);
	setTimeout(() => {
		t.style.opacity = '0';
		t.style.transition = '.25s';
		setTimeout(() => t.remove(), 250);
	}, 2800);
}

// ── Copy ──
function copyText(text) {
	navigator.clipboard
		.writeText(text)
		.then(() => {
			showToast('Copied to clipboard');
			// Mark as copied for beforeunload on API key pages
			if (window.__keyCopied !== undefined) {
				window.__keyCopied = true;
			}
		})
		.catch((err) => {
			console.error('Copy failed:', err);
			showToast('Failed to copy', 'error');
		});
}

// ── Copy code block (for API docs) ──
function copyCode(button, codeId) {
	const codeEl = document.getElementById(codeId);
	if (!codeEl) return;
	const text = codeEl.textContent;
	navigator.clipboard
		.writeText(text)
		.then(() => {
			const prev = button.textContent;
			button.textContent = 'Copied!';
			setTimeout(() => {
				button.textContent = prev;
			}, 2000);
		})
		.catch((err) => {
			console.error('Copy failed:', err);
			showToast('Failed to copy', 'error');
		});
}

// ── QR code generation ──
// Requires qrcode.min.js loaded in layout <head> (CDN or local)
// <script src="https://cdnjs.cloudflare.com/ajax/libs/qrcodejs/1.0.0/qrcode.min.js"></script>
function generateQR(containerId, url) {
	const el = document.getElementById(containerId);
	if (!el || typeof QRCode === 'undefined') return;
	el.innerHTML = '';
	new QRCode(el, {
		text: url,
		width: 112,
		height: 112,
		colorDark: '#0d0f12',
		colorLight: '#ffffff',
		correctLevel: QRCode.CorrectLevel.M,
	});
}

function downloadQR(containerId, filename) {
	const el = document.getElementById(containerId);
	if (!el) return;
	const canvas = el.querySelector('canvas');
	const img = el.querySelector('img');
	if (canvas) {
		const a = document.createElement('a');
		a.download = filename || 'qr.png';
		a.href = canvas.toDataURL('image/png');
		a.click();
	} else if (img) {
		// QRCode.js falls back to img on some browsers
		const a = document.createElement('a');
		a.download = filename || 'qr.png';
		a.href = img.src;
		a.click();
	}
}

// ── API Key Page ──
function initApiKeyPage() {
	const keyValueEl = document.getElementById('key-value');
	const keyCopiedEl = document.getElementById('key-copied');
	if (!keyValueEl) return;

	window.doCopy = function () {
		const key = keyValueEl.textContent;
		copyText(key);
		// Show copied feedback on this page
		if (keyCopiedEl) {
			keyCopiedEl.classList.add('show');
			setTimeout(() => {
				keyCopiedEl.classList.remove('show');
			}, 3000);
		}
	};

	// Track if key was copied before leave
	window.__keyCopied = false;
	const copyBtn = document.querySelector('.key-copy');
	const cardAction = document.querySelector('.card-action');
	if (copyBtn) {
		copyBtn.addEventListener('click', () => {
			window.__keyCopied = true;
		});
	}
	window.addEventListener('beforeunload', (e) => {
		if (window.__keyCopied === false) {
			e.preventDefault();
			e.returnValue = '';
		}
	});
	if (cardAction) {
		cardAction.addEventListener('click', () => {
			window.__keyCopied = true;
		});
	}
}

// ── API Page ──
function initApiPage() {
	// Placeholder for API page initialization
	// Can be extended for future API page features
}

document.addEventListener('DOMContentLoaded', () => {
	initApiKeyPage();
	initApiPage();
});
