import Crudy, { config } from "@allape/gocrud-react";
import { ITag } from "common/src/model/tag.ts";
import { IItemSearchParams } from "../model/item.ts";

export const TagCrudy = new Crudy<ITag, IItemSearchParams>(
  `${config.SERVER_URL}/tag`,
);
