import { DataStore, SortMode } from "@zyc/jsondb";
import { FileStore } from "@zone/file-service";
import moment from "https://deno.land/x/deno_ts_moment@0.0.4/mod.ts";
import { DOMParser } from "https://deno.land/x/deno_dom@v0.1.45/deno-dom-wasm.ts";
import { databaseFolder } from "./mod.ts";
import * as path from "@std/path";

export interface Note {
	id: string;
	title: string;
	content: string;
	summary: string;
	tags: string[];
	createTime: string;
	updateTime: string;
	image: string;
	readCount: number;
}

export class EditingNote {
	title: string = "";
	content: string = "";
	tags: string[] = [];
	id?: string;
	image?: string;
	summary?: string;
}

const dataStore = new DataStore<Note>(path.join(databaseFolder, "NOTE.json"));
const noteStore = new FileStore(path.join(databaseFolder, "NOTE"));

async function saveContentFile(content: string, id?: string): Promise<string> {
	let fileId = "";
	if (id) {
		fileId = id;
		await noteStore.update(id, content);
	} else {
		const { id } = await noteStore.save(content);
		fileId = id;
	}
	return fileId;
}

// 新建/编辑
export async function saveNote(data: EditingNote): Promise<Note | null> {
	let exists: Note | null = null;
	if (data.id) {
		exists = await getNote(data.id);
	}
	const fileId = await saveContentFile(data.content, data.id);
	if (!fileId) return null;
	const now = moment().format("YYYY-MM-DD HH:mm:ss");

	let summary = data.summary ? data.summary : "";
	if (!summary) {
		const contentDom = new DOMParser().parseFromString(data.content, "text/html");
		if (contentDom) {
			summary = contentDom.textContent.substring(0, 200);
		}
	}

	const note: Note = {
		id: exists ? exists.id : fileId,
		title: data.title,
		content: "",
		summary,
		tags: data.tags,
		createTime: exists ? exists.createTime : now,
		updateTime: now,
		image: data.image ? data.image : exists ? exists.image : "",
		readCount: exists ? exists.readCount : 0,
	};
	if (data.id) {
		await dataStore.update((it) => it.id === data.id, note);
	} else {
		await dataStore.insert(note);
	}
	return note;
}

export async function updateNote(note: Partial<Note> & { id: string }) {
	if (note.content) {
		await saveContentFile(note.content);
		note.content = "";
	}
	await dataStore.update((it) => it.id === note.id, note);
}

export async function deleteNote(id: string) {
	await dataStore.delete((record) => record.id === id);
	await noteStore.delete(id);
}

export async function getNote(id: string) {
	const note = await dataStore.get((record) => record.id === id);
	if (!note) return null;
	const fileContent = await noteStore.getFileText(note.id);
	note.content = fileContent;
	return note;
}

export async function queryNotes(ops?: {
	keywords?: string;
	sortField?: string;
	sortMode?: SortMode;
	page?: number;
	pageSize?: number;
}) {
	const sortField = ops ? ops.sortField ?? "updateTime" : "updateTime";
	const sortMode = ops ? ops.sortMode ?? SortMode.down : SortMode.down;
	let page = ops ? ops.page ?? 1 : 1;
	const pageSize = ops ? ops.pageSize ?? 10 : 10;
	if (page < 1) page = 1;
	const offset = (page - 1) * pageSize;
	const limit = pageSize;
	const results = await dataStore.query(
		(record) => {
			if (!ops || !ops.keywords) return true;
			const keywords = ops.keywords.toLowerCase();
			return record.title.toLowerCase().includes(keywords) || record.tags.some((t) => t.toLowerCase().includes(keywords!));
		},
		{ sortField, sortMode, offset, limit }
	);
	return results;
}
