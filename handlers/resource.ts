import { FileModel, FilePermission, FileSource, FileStore } from "@zone/file-service";
import * as hono from "https://deno.land/x/hono@v3.12.8/mod.ts";
import { failed, success } from "./mod.ts";
import { isInEnum, sameString, setStringValue } from "../utils/index.ts";
import { fileManager } from "../database/resource.ts";
import { QueryOptions, SortMode } from "@zyc/jsondb";

interface FileQueryOptions extends QueryOptions {
	keywords?: string;
	extension?: string;
	mime?: string;
}

function pickFileModel(form: FormData): Partial<FileModel> {
	const id = form.get("id");
	const name = form.get("name");
	const creatorId = form.get("creatorId");
	const permission = Number(form.get("permission"));
	const downloadable = form.get("downloadable");
	const shared = form.get("shared");
	const description = form.get("description");
	const remark = form.get("remark");
	const downloadCount = Number(form.get("downloadCount"));
	const source = Number(form.get("source"));
	const model: Partial<FileModel> = {};
	setStringValue("id", id, model);
	setStringValue("name", name, model);
	setStringValue("creatorId", creatorId, model);
	setStringValue("description", description, model);
	setStringValue("remark", remark, model);
	if (isInEnum(FilePermission, permission)) model.permission = permission;
	if (isInEnum(FileSource, source)) model.source = source;
	if (!isNaN(downloadCount) && typeof downloadCount === "number") model.downloadCount = downloadCount;
	if (typeof downloadable === "string" && ["true", "false"].includes(downloadable)) {
		model.downloadable = downloadable === "true";
	}
	if (typeof shared === "string" && ["true", "false"].includes(shared)) {
		model.shared = shared === "true";
	}
	return model;
}

function matchingStringFields(model: FileModel, fields: Array<keyof FileModel>, value: string) {
	value = value.toLowerCase();
	for (const i in fields) {
		const f: keyof FileModel = fields[i];
		const v = model[f];
		if (typeof v === "string" && v.toLowerCase().includes(value)) {
			return true;
		}
	}
	return false;
}

function getQueryParams(json: FileQueryOptions | null) {
	if (!json || !(json instanceof Object)) return {};
	const params: FileQueryOptions = {};
	if (json.keywords && typeof json.keywords === "string") params.keywords = json.keywords;
	if (json.extension && typeof json.extension === "string") params.extension = json.extension;
	if (json.mime && typeof json.mime === "string") params.mime = json.mime;
	if (json.sortField && typeof json.sortField === "string") params.sortField = json.sortField;
	if (json.limit && typeof json.limit === "number") params.limit = json.limit;
	if (json.offset && typeof json.offset === "number") params.offset = json.offset;
	if (isInEnum(SortMode, json.sortMode)) params.sortMode = json.sortMode;
	return params;
}

export async function save(ctx: hono.Context) {
	const form = await ctx.req.formData();
	const model = pickFileModel(form);
	const file = form.get("file");
	let fileStream: ReadableStream<Uint8Array> | null = null;
	let fileInfo: { name?: string; size?: number; mime?: string } = {};
	if (file) {
		if (typeof file === "string") {
			fileStream = FileStore.getReadableStreamByText(file);
		} else {
			fileInfo = { name: file.name, size: file.size, mime: file.type };
			fileStream = file.stream();
		}
	}
	if (model.id) {
		const exists = await fileManager.getFileModel(model.id);
		if (!exists) return ctx.json(failed(`file is not exists: ${model.id}`));
		await fileManager.update({ ...fileInfo, ...model, id: model.id }, fileStream);
		return ctx.json(success({ id: model.id }));
	}
	if (!fileStream) return ctx.json(failed("add file: file content is empty"));
	const name = model.name ?? fileInfo.name;
	if (!name) return ctx.json(failed("add file: file name is empty"));
	const fileModel = await fileManager.add(fileStream, { ...fileInfo, ...model, name });
	return ctx.json(success(fileModel));
}

export async function query(ctx: hono.Context) {
	const params = getQueryParams(await ctx.req.json());
	const files = await fileManager.queryFileModels((f) => {
		if (params.keywords && !matchingStringFields(f, ["name", "description"], params.keywords)) return false;
		if (params.extension && !sameString(f.extension, params.extension)) return false;
		if (params.mime && !sameString(f.mime, params.mime)) return false;
		return true;
	}, params);
	return ctx.json(success(files));
}

export async function queryByTypes(ctx: hono.Context) {
	const params = await ctx.req.json();
	if (!params) {
		return ctx.json(failed());
	}
	const queryOptions: QueryOptions = params.queryOptions ?? {};
	if (!(queryOptions instanceof Object)) return ctx.json(failed("type of queryOptions is invalid"));
	const extension: string[] = params.extension ?? [];
	const mime: string[] = params.mime ?? [];
	if (!(extension instanceof Array) || !(mime instanceof Array)) {
		return ctx.json(failed("params field type invalid"));
	}
	const keys = [...extension, ...mime].filter((k) => typeof k === "string");
	const files = await fileManager.queryFileModels((f) => {
		return keys.some((k) => sameString(f.extension, k) || sameString(f.mime, k));
	}, queryOptions);
	return ctx.json(success(files));
}

export async function getFile(ctx: hono.Context) {
	const id = ctx.req.param("id");
	if (!id) return ctx.json(failed());
	const model = await fileManager.getFileModel(id);
	const f = await fileManager.getFileStream(id);
	if (!model || !f) return ctx.json(failed("file not found"));
	ctx.res.headers.set("content-type", model.mime);
	return ctx.body(f);
}

export async function deleteFile(ctx: hono.Context) {
	const id = ctx.req.param("id");
	if (!id) return ctx.json(failed());
	await fileManager.delete(id);
	return ctx.json(success(id));
}
