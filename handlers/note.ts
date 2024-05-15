import * as hono from "https://deno.land/x/hono@v3.12.8/mod.ts";
import { failed, success, getJsonParams } from "./mod.ts";
import * as database from "../database/notes.ts";
import * as typeCheck from "../utils/types.ts";

// 新建/编辑 POST params: database.EditingNote
export async function saveNote(ctx: hono.Context) {
	const data: database.EditingNote | null = await getJsonParams(ctx);
	if (!data) return ctx.json(failed("参数为空"));
	const paramsErr = ctx.json(failed("参数类型错误"));
	if (!data.title || !typeCheck.isString(data.title)) return paramsErr;
	if (!data.content || !typeCheck.isString(data.content)) return paramsErr;
	if (!data.tags || !typeCheck.isStringArray(data.tags)) return paramsErr;
	if (data.id && !typeCheck.isString(data.id)) return paramsErr;
	if (data.image && !typeCheck.isString(data.image)) return paramsErr;
	if (data.summary && !typeCheck.isString(data.summary)) return paramsErr;
	const model: Partial<database.Note> & { title: string; content: string; tags: string[] } = {
		title: data.title,
		content: data.content,
		tags: data.tags,
	};
	if (data.id) model.id = data.id;
	if (data.image) model.image = data.image;
	if (data.summary) model.summary = data.summary;
	const note = await database.saveNote(model);
	if (!note) return ctx.json(failed("保存失败"));
	return ctx.json(success(note));
}

// POST url params: id
export async function deleteNote(ctx: hono.Context) {
	const id = ctx.req.param("id");
	if (!id) return ctx.json(failed("参数错误"));
	await database.deleteNote(id);
	return ctx.json(success(true));
}

// POST url params: id
export async function getNote(ctx: hono.Context) {
	const id = ctx.req.param("id");
	if (!id) return ctx.json(failed("参数错误"));
	const note = await database.getNote(id);
	if (!note) return ctx.json(failed("读取失败"));
	await database.updateNote({ id: note.id, readCount: note.readCount + 1 });
	return ctx.json(success(note));
}

// POST params: { keywords?: string; sortField?: string; sortMode?: SortMode; page?: number; pageSize?: number }
export async function getNotes(ctx: hono.Context) {
	let params = await getJsonParams(ctx);
	if (!params || typeCheck.isArray(params)) params = {};
	const result = await database.queryNotes(params);
	return ctx.json(success(result));
}
