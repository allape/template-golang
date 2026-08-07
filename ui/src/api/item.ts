import Crudy, { AntdM2MConnectorHandler, config } from "@allape/gocrud-react";
import { IItem, IItemTag } from "../model/item.ts";
import { ITag } from "../model/tag.ts";
import { TagCrudy } from "./tag.ts";

export const ItemCrudy = new Crudy<IItem>(`${config.SERVER_URL}/item`);

export const ItemTagHandler = new AntdM2MConnectorHandler<
  IItem,
  ITag,
  IItemTag
>(`${config.SERVER_URL}/item-tag`, ItemCrudy, TagCrudy, "tagId", "itemId");
