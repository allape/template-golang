import Crudy, { AntdM2MConnectorHandler, config } from "@allape/gocrud-react";
import { IItem, IItemTag } from "common/src/model/item.ts";
import { ITag } from "common/src/model/tag.ts";
import { IItemSearchParams } from "../model/item.ts";
import { TagCrudy } from "./tag.ts";

export const ItemCrudy = new Crudy<IItem, IItemSearchParams>(
  `${config.SERVER_URL}/item`,
);

export const ItemTagHandler = new AntdM2MConnectorHandler<
  IItem,
  ITag,
  IItemTag
>(`${config.SERVER_URL}/item-tag`, ItemCrudy, TagCrudy, "tagId", "itemId");
