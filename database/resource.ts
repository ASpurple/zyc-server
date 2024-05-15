import { FileManager } from "@zone/file-service";
import * as path from "@std/path";
import { databaseFolder } from "./mod.ts";

export const fileManager = new FileManager(path.join(databaseFolder, "FILE"));
