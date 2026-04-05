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
		.then(() => showToast('Copied to clipboard'))
		.catch(() => {
			const el = document.createElement('textarea');
			el.value = text;
			el.style.cssText = 'position:fixed;opacity:0';
			document.body.appendChild(el);
			el.select();
			document.execCommand('copy');
			el.remove();
			showToast('Copied to clipboard');
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
		navigator.clipboard
			.writeText(key)
			.then(() => {
				if (keyCopiedEl) {
					keyCopiedEl.classList.add('show');
					setTimeout(() => {
						keyCopiedEl.classList.remove('show');
					}, 3000);
				}
				// Mark as copied for beforeunload
				if (window.__keyCopied !== undefined) {
					window.__keyCopied = true;
				}
			})
			.catch((err) => {
				console.error('Copy failed:', err);
				showToast('Failed to copy', 'error');
			});
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
	const apiForm = document.getElementById('api-key-form');
	if (!apiForm) return;

	apiForm.addEventListener('submit', async function (e) {
		e.preventDefault();
		const emailInput = document.getElementById('api-email');
		const email = emailInput?.value.trim();
		if (!email) return;
		const btn = document.getElementById('api-submit-btn');
		const prevText = btn.textContent;
		btn.textContent = 'Sending…';
		btn.disabled = true;

		try {
			const body = new FormData();
			body.append('email', email);

			const res = await fetch('/api/internal/key/request', {
				method: 'POST',
				body,
			});

			if (!res.ok) {
				const err = await res.json().catch(() => ({}));
				throw new Error(err.error || 'Something went wrong');
			}

			if (emailInput) emailInput.closest('.field').style.display = 'none';
			btn.style.display = 'none';
			const sentTo = document.getElementById('api-sent-to');
			if (sentTo) sentTo.textContent = email;
			const success = document.getElementById('api-success');
			if (success) success.style.display = 'flex';
		} catch (err) {
			showToast(err.message, 'error');
		} finally {
			btn.textContent = prevText;
			btn.disabled = false;
		}
	});

	window.copyApiKey = function () {
		const keyValue = document.getElementById('api-key-value');
		const copyBtn = document.getElementById('api-copy-btn');
		if (!keyValue) return;
		const key = keyValue.textContent;
		navigator.clipboard
			.writeText(key)
			.then(() => {
				if (copyBtn) {
					const prev = copyBtn.innerHTML;
					copyBtn.textContent = 'Copied!';
					setTimeout(() => {
						copyBtn.innerHTML = prev;
					}, 2000);
				}
			})
			.catch((err) => {
				console.error('Copy failed:', err);
				showToast('Failed to copy', 'error');
			});
	};

	window.copyCode = function (btn, id) {
		const codeEl = document.getElementById(id);
		if (!codeEl) return;
		const text = codeEl.textContent;
		navigator.clipboard
			.writeText(text)
			.then(() => {
				const prev = btn.textContent;
				btn.textContent = 'Copied!';
				setTimeout(() => {
					btn.textContent = prev;
				}, 2000);
			})
			.catch((err) => {
				console.error('Copy failed:', err);
				showToast('Failed to copy', 'error');
			});
	};
}

document.addEventListener('DOMContentLoaded', () => {
	initApiKeyPage();
	initApiPage();
});
