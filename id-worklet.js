class SelectURLOperation {
	async run(urls) {
		const hasIdentifier = await sharedStorage.get('has-identifier');

		console.log(`urls = ${JSON.stringify(urls)}`);
		console.log(`hasIdentifier = ${hasIdentifier}`);

		return parseInt(hasIdentifier);
	}
}

register('has-identifier', SelectURLOperation);


class SelectURLOperation2 {
	async run(urls, data) {
		const identifier = await sharedStorage.get(`identifier-${data.index}`);

		console.log(`urls = ${JSON.stringify(urls)}`);
		console.log(`identifier = ${identifier}`);
		console.log(JSON.stringify(this));

		return parseInt(identifier);
	}
}

register('identifier', SelectURLOperation2);
