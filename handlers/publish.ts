import { ensureDir } from "@std/fs";
import { Context } from "https://deno.land/x/hono@v3.12.8/mod.ts";
import { failed } from "./mod.ts";
import * as path from "@std/path";
import { decompress } from "@fakoua/zip-ts";
import { success } from "./mod.ts";
import { getWebFolder } from "./web.ts";

export async function webPublish(ctx: Context) {
	const auth = ctx.req.header("auth");
	if (!auth || auth !== "zyc") {
		ctx.status(400);
		return ctx.text("No permission");
	}
	const formData = await ctx.req.formData();
	const file = formData.get("file");
	const folderPath = getWebFolder();
	const filePath = path.join(folderPath, "dist.zip");
	if (!file || typeof file === "string") return ctx.json(failed("Invalid file"));

	// 存储文件
	await ensureDir(folderPath);

	await Deno.writeFile(filePath, file.stream(), { create: true, append: false });

	const dirs = Deno.readDir(folderPath);
	const fileName = path.basename(filePath);
	for await (const child of dirs) {
		if (child.name === fileName) continue;
		const p = path.join(folderPath, child.name);
		await Deno.remove(p, { recursive: true });
	}

	await decompress(filePath, path.dirname(filePath));

	await Deno.remove(filePath);

	return ctx.json(success("publish success!"));
}
