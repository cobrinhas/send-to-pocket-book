export interface Env {
	PROXY_SERVER_URL: string;
	PROXY_SERVER_SECURITY_TOKEN: string;
}

export interface Config {
	proxyServerUrl: string;
	proxyServerSecurityToken: string;
}

export function fromEnv(env: Env): Config {
	return {
		proxyServerUrl: env.PROXY_SERVER_URL,
		proxyServerSecurityToken: env.PROXY_SERVER_SECURITY_TOKEN
	};
}
