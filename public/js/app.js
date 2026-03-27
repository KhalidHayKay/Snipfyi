/* smply.cc — app.js */

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

// ── Shorten form ──
function initShortenForm(formId, resultId, shortUrlBase) {
	const form = document.getElementById(formId);
	if (!form) return;

	form.addEventListener('submit', async (e) => {
		e.preventDefault();
		const urlInput = form.querySelector(
			'input[name="url"], input[name="long_url"]',
		);
		const aliasInput = form.querySelector('input[name="alias"]');
		const btn = form.querySelector('button[type="submit"]');
		const resultEl = document.getElementById(resultId);
		if (!urlInput?.value || !resultEl) return;

		const prev = btn.innerHTML;
		btn.textContent = 'Shortening…';
		btn.disabled = true;

		try {
			const body = new FormData();
			body.append('url', urlInput.value);
			if (aliasInput?.value) {
				body.append('alias', aliasInput.value);
			}

			const res = await fetch('/api/shorten', {
				method: 'POST',
				body,
			});

			if (!res.ok) {
				const err = await res.json().catch(() => ({}));
				throw new Error(err.error || 'Something went wrong');
			}

			const { data } = await res.json();
			const shortUrl = data.short;
			const statsUrl = data.stat;

			renderResult(resultEl, shortUrl, statsUrl);
			resultEl.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
		} catch (err) {
			showToast(err.message, 'error');
		} finally {
			btn.innerHTML = prev;
			btn.disabled = false;
		}
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

// ── Hero result (compact, single column) ──
function renderResult(el, shortUrl, statsUrl) {
	console.log(statsUrl);
	const qrId = 'hero-qr-' + Date.now();
	el.innerHTML = `
    <div class="result-tag">Link ready</div>
    <div class="result-row">
      <div class="result-url"><a href="${shortUrl}" target="_blank" rel="noopener">${shortUrl}</a></div>
      <div class="result-btns">
        <button class="btn btn-primary btn-sm" onclick="copyText('${shortUrl}')">Copy</button>
        <a href="${shortUrl}" class="btn btn-ghost btn-sm" target="_blank" rel="noopener">Open ↗</a>
      </div>
    </div>
    <div class="result-divider"></div>
    <div class="result-bottom">
      <div class="result-qr-wrap">
        <div class="result-qr-box" id="${qrId}"></div>
        <button class="result-qr-dl" onclick="downloadQR('${qrId}', 'smply-qr.png')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          Download PNG
        </button>
      </div>
      <div class="result-stats-col">
        <span class="result-stats-note">Track clicks at</span>
        <a href="${statsUrl}" class="result-stats-link" target="_blank" rel="noopener">
          <svg viewBox="0 0 24 24"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
          ${statsUrl}
        </a>
      </div>
    </div>
  `;
	el.classList.add('show');
	// slight delay lets the DOM paint before QR fills the container
	setTimeout(() => generateQR(qrId, shortUrl), 50);
}

// ── Shorten page result (more room, side-by-side QR) ──
function renderShortenPageResult(el, shortUrl, statsUrl) {
	const qrId = 'sr-qr-' + Date.now();
	el.innerHTML = `
    <div class="sr-tag"><span class="sr-pulse"></span> Link ready</div>
    <div class="sr-main">
      <div class="sr-left">
        <div class="sr-url"><a href="${shortUrl}" target="_blank" rel="noopener">${shortUrl}</a></div>
        <div class="sr-btns">
          <button class="btn btn-primary btn-sm" onclick="copyText('${shortUrl}')">Copy link</button>
          <a href="${shortUrl}" class="btn btn-ghost btn-sm" target="_blank" rel="noopener">Open ↗</a>
        </div>
        <div class="sr-divider"></div>
        <div class="sr-stats">
          <div class="sr-stats-note">View analytics at <strong>${statsUrl}</strong></div>
          <a href="${statsUrl}" class="sr-stats-btn" target="_blank" rel="noopener">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
            /stats
          </a>
        </div>
      </div>
      <div class="sr-qr-col">
        <div class="sr-qr-box" id="${qrId}"></div>
        <button class="sr-qr-dl" onclick="downloadQR('${qrId}', 'smply-qr.png')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          Download QR
        </button>
      </div>
    </div>
  `;
	el.classList.add('show');
	setTimeout(() => generateQR(qrId, shortUrl), 50);
}

// ── Shorten page form (more options) ──
function initAdvancedForm() {
	const form = document.getElementById('shorten-form');
	if (!form) return;

	form.addEventListener('submit', async (e) => {
		e.preventDefault();
		const urlInput = form.querySelector('input[name="long_url"]');
		const aliasInput = form.querySelector('input[name="alias"]');
		const btn = form.querySelector('button[type="submit"]');
		const resultEl = document.getElementById('shorten-result');
		if (!urlInput?.value || !resultEl) return;

		const prev = btn.innerHTML;
		btn.textContent = 'Shortening…';
		btn.disabled = true;

		try {
			const body = new FormData();
			body.append('url', urlInput.value);
			if (aliasInput?.value) {
				body.append('alias', aliasInput.value);
			}

			const res = await fetch('/api/shorten', {
				method: 'POST',
				body,
			});

			if (!res.ok) {
				const err = await res.json().catch(() => ({}));
				throw new Error(err.error || 'Something went wrong');
			}

			const { data } = await res.json();
			const shortUrl = data.short;
			const statsUrl = data.stat;

			renderShortenPageResult(resultEl, shortUrl, statsUrl);
			resultEl.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
		} catch (err) {
			showToast(err.message, 'error');
		} finally {
			btn.innerHTML = prev;
			btn.disabled = false;
		}
	});
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

	apiForm.addEventListener('submit', function (e) {
		e.preventDefault();
		const emailInput = document.getElementById('api-email');
		const email = emailInput?.value.trim();
		if (!email) return;
		const btn = document.getElementById('api-submit-btn');
		const prevText = btn.textContent;
		btn.textContent = 'Sending…';
		btn.disabled = true;
		setTimeout(() => {
			if (emailInput) emailInput.closest('.field').style.display = 'none';
			btn.style.display = 'none';
			const sentTo = document.getElementById('api-sent-to');
			if (sentTo) sentTo.textContent = email;
			const success = document.getElementById('api-success');
			if (success) success.style.display = 'flex';
		}, 1100);
	});

	window.toggleKeyDemo = function () {
		const keyCard = document.getElementById('api-key-card');
		const demoToggle = document.getElementById('demo-toggle');
		if (!keyCard || !demoToggle) return;
		const isVisible = keyCard.style.display === 'block';
		keyCard.style.display = isVisible ? 'none' : 'block';
		demoToggle.textContent = isVisible ? 'Preview key state' : 'Hide key state';
	};

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

// ── Stats Page ──
function initStatsPage() {
	const statsQr = document.getElementById('stats-qr');
	if (!statsQr) return;

	// Generate QR code for stats page short URL
	if (typeof QRCode !== 'undefined') {
		const shortUrl =
			statsQr.closest('.stats-qr-box')?.parentElement?.textContent ||
			document.querySelector('.stats-title em')?.textContent ||
			'';
		if (shortUrl) {
			new QRCode(statsQr, {
				text: shortUrl,
				width: 120,
				height: 120,
				colorDark: '#0d0f12',
				colorLight: '#ffffff',
				correctLevel: QRCode.CorrectLevel.M,
			});
		}
	}
}

// ── Boot ──
document.addEventListener('DOMContentLoaded', () => {
	initShortenForm('hero-form', 'hero-result');
	initAdvancedForm();
	initApiKeyPage();
	initApiPage();
	initStatsPage();
});
