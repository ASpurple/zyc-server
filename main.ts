import * as hono from "https://deno.land/x/hono@v3.12.8/mod.ts";
import { cors } from "https://deno.land/x/hono@v3.12.8/middleware.ts";
import { staticHandler } from "./handlers/web.ts";
import * as note from "./handlers/note.ts";
import * as resource from "./handlers/resource.ts";
import { existsSync } from "@std/fs";
import { webPublish } from "./handlers/publish.ts";

const app = new hono.Hono();

app.use("*", cors());

app.post("/publish/web", webPublish);

app.post("/note/save", note.saveNote);
app.post("/note/delete/:id", note.deleteNote);
app.post("/note/detail/:id", note.getNote);
app.post("/note/page", note.getNotes);

app.post("/resource/save", resource.save);
app.post("/resource/query", resource.query);
app.post("/resource/query/types", resource.queryByTypes);
app.get("/resource/file/:id", resource.getFile);
app.post("/resource/delete/:id", resource.deleteFile);

app.notFound(staticHandler);

const port = existsSync("./deno.json") ? 8000 : 80;

Deno.serve({ port }, app.fetch);
