chrome.tabs.query({ active: true, lastFocusedWindow: true }, (tabs) => {
    const activeTab = tabs.find((x) => x.url);

    document.getElementById('document-title').value = activeTab.title;
    document.getElementById('send-button').addEventListener('click', send);

    const formStatusEl = document.getElementById('form-status');
    const statusEl = document.getElementById('status');

    loadProxySettings();

    function showLoadingSpinner() {
        const spinnerEl = document.getElementById('spinner');

        formStatusEl.classList.remove('hidden');
        formStatusEl.hidden = false;
        spinnerEl.hidden = false;
        statusEl.hidden = true;
    }

    function hideLoadingSpinner() {
        const spinnerEl = document.getElementById('spinner');

        spinnerEl.hidden = true;
    }

    function showStatus(status, tclass) {
        if (tclass) {
            statusEl.classList.length = 0;
            statusEl.classList.add(...['text-sm', 'text-center', tclass]);
        }

        statusEl.hidden = false;
        statusEl.innerText = status;
    }

    async function loadProxySettings() {
        try {
            const storage = await chrome.storage.sync.get([
                'proxy-server-url',
                'pbsync-email',
            ]);

            const proxyServerUrl = storage['proxy-server-url'];
            const pbsyncEmail = storage['pbsync-email'];

            if (proxyServerUrl) {
                document.getElementById('proxy-server-url').value = proxyServerUrl;
            }

            if (pbsyncEmail) {
                document.getElementById('pbsync-email').value = pbsyncEmail;
            }
        } catch {
            console.error('failed to load proxy settings')
        }
    }


    function saveProxySettings(proxyServerUrl, pbsyncEmail) {
        return chrome.storage.sync.set({
            'proxy-server-url': proxyServerUrl,
            'pbsync-email': pbsyncEmail
        });
    }

    async function send() {
        const title = document.getElementById('document-title').value;
        const proxyServerUrl = document.getElementById('proxy-server-url').value;
        const pbsyncEmail = document.getElementById('pbsync-email').value;
        const url = activeTab.url;

        showLoadingSpinner();
        saveProxySettings(proxyServerUrl, pbsyncEmail);

        try {
            const response = await window.fetch(proxyServerUrl, { body: JSON.stringify({ title: title, email: pbsyncEmail, url: url }), method: 'POST', mode: 'no-cors' });

            if (response.status === 202) {
                showStatus('The document has been sent to your pocket book! If this is the first time using the extension, check your inbox to add the proxy e-mail to your whitelist.', 'text-success');
            } else {
                showStatus('The document has been received but no further action has been taken yet.', 'text-warning');
            }
        } catch (e) {
            showStatus('Proxy Server is not reachable.', 'text-error');
        }

        hideLoadingSpinner();
    }
});