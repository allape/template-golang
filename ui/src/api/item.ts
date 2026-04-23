import Crudy, { config, get } from "@allape/gocrud-react";
import { IItem, IItemTag } from "../model/item.ts";
import { ITag } from "../model/tag.ts";

export const ItemCrudy = new Crudy<IItem>(`${config.SERVER_URL}/item`);

export const ItemTagCrudy = new Crudy<IItemTag>(
  `${config.SERVER_URL}/item-tag`,
);

export function addTagsToItem(
  tags: ITag["id"][],
  itemId: IItem["id"],
): Promise<IItemTag[]> {
  return get(
    `${config.SERVER_URL}/item-tag/save/itemId?itemId=${itemId}&tagId=${encodeURIComponent(tags.join(","))}`,
    {
      method: "POST",
    },
  );
}
