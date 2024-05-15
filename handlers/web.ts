import * as hono from "https://deno.land/x/hono@v3.12.8/mod.ts";
import * as fs from "https://deno.land/std@0.213.0/fs/mod.ts";
import * as path from "https://deno.land/std@0.213.0/path/mod.ts";
import * as mrmime from "https://deno.land/x/mrmime@v2.0.0/mod.ts";

let webFolder = path.join(Deno.cwd(), "WEB");

async function setStaticFolder(folder: string) {
	webFolder = folder;
	if (!(await fs.exists(webFolder))) {
		await Deno.mkdir(webFolder, { recursive: true });
		return;
	}
	const stat = await Deno.stat(webFolder);
	if (!stat.isDirectory) {
		await Deno.remove(webFolder, { recursive: true });
		await Deno.mkdir(webFolder, { recursive: true });
	}
}

setStaticFolder(webFolder);

export async function staticHandler(ctx: hono.Context) {
	let reqPath = ctx.req.path;
	if (reqPath === "/") reqPath = "/index.html";
	let filePath = path.join(webFolder, ...reqPath.split("/"));
	if (!(await fs.exists(filePath))) {
		if (ctx.req.method.toUpperCase() === "GET") {
			filePath = path.join(webFolder, "index.html");
			if (!(await fs.exists(filePath))) return ctx.text("not found", 404);
		} else {
			return ctx.text("not found", 404);
		}
	}
	const mime = mrmime.lookup(filePath) ?? "application/octet-stream";
	const file = await Deno.open(filePath);
	const stat = await file.stat();
	ctx.res.headers.set("Content-Type", mime);
	ctx.res.headers.set("Content-Size", stat.size.toString());

	const st = new ReadableStream<Uint8Array>({
		async pull(controller) {
			for await (const chunk of file.readable) {
				controller.enqueue(chunk);
			}
			controller.close();
		},
	});

	return ctx.body(st);
}
