import * as hono from "https://deno.land/x/hono@v3.12.8/mod.ts";

export enum ResponseCodes {
	success,
	failed,
}

export interface ResponseData<T> {
	code: ResponseCodes;
	message: string;
	data: T;
}

export function success<T>(data: T): ResponseData<T> {
	return {
		code: ResponseCodes.success,
		message: "",
		data,
	};
}

export function failed(message = "请求失败"): ResponseData<null> {
	return {
		code: ResponseCodes.failed,
		message,
		data: null,
	};
}

// deno-lint-ignore no-explicit-any
export async function getJsonParams(ctx: hono.Context): Promise<any> {
	try {
		return await ctx.req.json();
	} catch (_e) {
		return {};
	}
}
