import Crudy, { config } from "@allape/gocrud-react";
import { ITag } from "../model/tag.ts";

export const TagCrudy = new Crudy<ITag>(`${config.SERVER_URL}/tag`);
