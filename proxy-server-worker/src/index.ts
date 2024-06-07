import { Env, fromEnv } from './config';

export default {
	async fetch(request: Request, env: Env): Promise<Response> {
		const config = fromEnv(env);

		return fetch(`${config.proxyServerUrl}/send`, {
			body: request.body,
			method: 'POST',
			headers: {
				...request.headers,
				'content-type': 'application/json',
				'x-sec-token': config.proxyServerSecurityToken
			}
		});
	}
};
