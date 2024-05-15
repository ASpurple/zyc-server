import { ensureDir } from "@std/fs";
import { Context } from "https://deno.land/x/hono@v3.12.8/mod.ts";
import { failed } from "./mod.ts";
import * as path from "@std/path";
import { databaseFolder } from "../database/mod.ts";
import { decompress } from "@fakoua/zip-ts";
import { success } from "./mod.ts";

export async function webPublish(ctx: Context) {
	const formData = await ctx.req.formData();
	const file = formData.get("file");
	const folderPath = path.join(databaseFolder, "WEB");
	const filePath = path.join(folderPath, "dist.zip");
	if (!file || typeof file === "string") return ctx.json(failed("Invalid file"));

	// 存储文件
	await ensureDir(folderPath);

	await Deno.writeFile(filePath, file.stream(), { create: true, append: false });

	const dirs = Deno.readDir(folderPath);
	for await (const child of dirs) {
		if (child.name === "dist.zip") continue;
		const p = path.join(folderPath, child.name);
		await Deno.remove(p, { recursive: true });
	}

	await decompress(filePath, path.dirname(filePath), { overwrite: true, includeFileName: true });

	await Deno.remove(filePath);

	return ctx.json(success("publish success!"));
}
