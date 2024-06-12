import { Env, fromEnv } from './config';

const corsHeaders = {
	"Access-Control-Allow-Origin": "*",
	"Access-Control-Allow-Methods": "GET,HEAD,POST,OPTIONS",
	"Access-Control-Max-Age": "86400",
};

async function handleOptions(request: Request) {
	if (
		request.headers.get("Origin") !== null &&
		request.headers.get("Access-Control-Request-Method") !== null &&
		request.headers.get("Access-Control-Request-Headers") !== null
	) {
		return new Response(null, {
			headers: {
				...corsHeaders,
				"Access-Control-Allow-Headers": request.headers.get(
					"Access-Control-Request-Headers"
				) ?? '',
			},
		});
	} else {
		return new Response(null, {
			headers: {
				Allow: "GET, HEAD, POST, OPTIONS",
			},
		});
	}
}


export default {
	async fetch(request: Request, env: Env): Promise<Response> {
		const config = fromEnv(env);

		if (request.method === "OPTIONS") {
			return handleOptions(request);
		}

		let response = await fetch(`${config.proxyServerUrl}/send`, {
			body: request.body,
			method: 'POST',
			headers: {
				...request.headers,
				'content-type': 'application/json',
				'x-sec-token': config.proxyServerSecurityToken,
				'origin': new URL(`${config.proxyServerUrl}/send`).origin
			}
		});

		response = new Response(response.body, response);

		response.headers.set("Access-Control-Allow-Origin", request.headers.get('origin') ?? '');
		response.headers.append("Vary", "Origin");

		return response;
	}
};
