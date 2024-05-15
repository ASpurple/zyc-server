import * as path from "@std/path";
import { v1 } from "@std/uuid";
import { FileModel } from "@zone/file-service";

export type StringKeys<T> = { [K in keyof T]: T[K] extends string ? K : never }[keyof T];

export function setStringValue(field: StringKeys<FileModel>, value: unknown, model: Partial<FileModel>) {
	if (typeof value === "string" && !!value.trim()) model[field] = value;
}

export function toNumber(value?: string) {
	if (!value) return 0;
	const num = parseInt(value);
	return num;
}

export function isInEnum(enumObj: object, value: unknown) {
	const values = Object.values(enumObj).filter((v) => typeof v === "number");
	return values.includes(value);
}

export async function createTempFile() {
	const folder = path.join(Deno.cwd(), "TEMP");
	await Deno.mkdir(folder, { recursive: true });
	const fp = path.join(folder, v1.generate().toString());
	const file = await Deno.create(fp);
	return { file, fp };
}

export function sameString(src: string, dst: string) {
	return src.toLowerCase().includes(dst.toLowerCase());
}
